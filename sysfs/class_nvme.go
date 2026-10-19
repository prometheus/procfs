// Copyright The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build linux

package sysfs

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/prometheus/procfs/internal/parsers"
)

const (
	nvmeClassPath = "class/nvme"

	// The block layer reports the size attribute in 512-byte sectors.
	blockSectorSize = 512
)

// nvmeNamespacePattern matches the namespace block devices below a controller:
// nvme<ctrl>n<inst> for private namespaces and, when native multipath is
// enabled for a shared namespace, the hidden per-path device
// nvme<subsys>c<ctrl>n<inst>.
var nvmeNamespacePattern = regexp.MustCompile(`^nvme\d+(?:c\d+)?n(\d+)$`)

// NVMeNamespace contains info from files in /sys/class/nvme/<device>/<namespace>.
type NVMeNamespace struct {
	Name             string // block device name, e.g. nvme0n1 or nvme0c1n1
	ID               string // namespace ID from the nsid file, or from the directory name if absent
	UsedBlocks       uint64 // from nuse file (blocks used)
	SizeBlocks       uint64 // from size file (total size in 512-byte sectors, as for every block device)
	LogicalBlockSize uint64 // from queue/logical_block_size file
	ANAState         string // from ana_state file
	UsedBytes        uint64 // calculated: UsedBlocks * LogicalBlockSize
	SizeBytes        uint64 // calculated: SizeBlocks * 512
	CapacityBytes    uint64 // calculated: SizeBlocks * 512
}

// NVMeDevice contains info from files in /sys/class/nvme for a single NVMe device.
type NVMeDevice struct {
	Name             string
	Serial           string          // /sys/class/nvme/<Name>/serial
	Model            string          // /sys/class/nvme/<Name>/model
	State            string          // /sys/class/nvme/<Name>/state
	FirmwareRevision string          // /sys/class/nvme/<Name>/firmware_rev
	ControllerID     string          // /sys/class/nvme/<Name>/cntlid
	Namespaces       []NVMeNamespace // NVMe namespaces for this device
}

// NVMeClass is a collection of every NVMe device in /sys/class/nvme.
//
// The map keys are the names of the NVMe devices.
type NVMeClass map[string]NVMeDevice

// NVMeClass returns info for all NVMe devices read from /sys/class/nvme.
func (fs FS) NVMeClass() (NVMeClass, error) {
	path := fs.sys.Path(nvmeClassPath)

	dirs, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list NVMe devices at %q: %w", path, err)
	}

	nc := make(NVMeClass, len(dirs))
	for _, d := range dirs {
		device, err := fs.parseNVMeDevice(d.Name())
		if err != nil {
			return nil, err
		}

		nc[device.Name] = *device
	}

	return nc, nil
}

// Parse one NVMe device.
func (fs FS) parseNVMeDevice(name string) (*NVMeDevice, error) {
	path := fs.sys.Path(nvmeClassPath, name)
	device := NVMeDevice{Name: name}

	// Parse device-level attributes
	for _, f := range [...]string{"firmware_rev", "model", "serial", "state", "cntlid"} {
		name := filepath.Join(path, f)
		value, err := parsers.SysReadFile(name)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %q: %w", name, err)
		}

		switch f {
		case "firmware_rev":
			device.FirmwareRevision = value
		case "model":
			device.Model = value
		case "serial":
			device.Serial = value
		case "state":
			device.State = value
		case "cntlid":
			device.ControllerID = value
		}
	}

	// Parse namespaces - read directory and filter using regex
	dirs, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list NVMe namespaces at %q: %w", path, err)
	}

	var namespaces []NVMeNamespace

	for _, d := range dirs {
		// Use regex to identify namespace directories and extract namespace ID
		match := nvmeNamespacePattern.FindStringSubmatch(d.Name())
		if len(match) < 2 {
			// Skip if not a namespace directory
			continue
		}
		nsid := match[1]
		namespacePath := filepath.Join(path, d.Name())

		namespace := NVMeNamespace{
			Name:     d.Name(),
			ID:       nsid,
			ANAState: "unknown", // Default value
		}

		// The number in the device name is the kernel's namespace head
		// instance, which only matches the NSID by convention. Prefer the
		// real NSID when the kernel exposes it.
		if value, err := parsers.SysReadFile(filepath.Join(namespacePath, "nsid")); err == nil {
			namespace.ID = value
		}

		// Parse namespace attributes using the same approach as device attributes
		for _, f := range [...]string{"nuse", "size", "queue/logical_block_size", "ana_state"} {
			filePath := filepath.Join(namespacePath, f)
			value, err := parsers.SysReadFile(filePath)
			if err != nil {
				if f == "ana_state" {
					// ana_state may not exist, skip silently
					continue
				}
				return nil, fmt.Errorf("failed to read file %q: %w", filePath, err)
			}

			switch f {
			case "nuse":
				if val, parseErr := strconv.ParseUint(value, 10, 64); parseErr == nil {
					namespace.UsedBlocks = val
				}
			case "size":
				if val, parseErr := strconv.ParseUint(value, 10, 64); parseErr == nil {
					namespace.SizeBlocks = val
				}
			case "queue/logical_block_size":
				if val, parseErr := strconv.ParseUint(value, 10, 64); parseErr == nil {
					namespace.LogicalBlockSize = val
				}
			case "ana_state":
				namespace.ANAState = value
			}
		}

		// Calculate derived values. nuse is reported in logical blocks, but
		// the block layer always reports size in 512-byte sectors,
		// independent of the namespace's logical block size.
		namespace.UsedBytes = namespace.UsedBlocks * namespace.LogicalBlockSize
		namespace.SizeBytes = namespace.SizeBlocks * blockSectorSize
		namespace.CapacityBytes = namespace.SizeBlocks * blockSectorSize

		namespaces = append(namespaces, namespace)
	}

	device.Namespaces = namespaces

	return &device, nil
}
