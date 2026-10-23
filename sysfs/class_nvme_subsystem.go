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
	"strings"

	"github.com/prometheus/procfs/internal/parsers"
)

const nvmeSubsystemClassPath = "class/nvme-subsystem"

var (
	nvmeSubsystemControllerRE = regexp.MustCompile(`^nvme\d+$`)
	nvmeNamespaceRE           = regexp.MustCompile(`^nvme\d+n\d+$`)
)

// NVMeSubsystem contains info from /sys/class/nvme-subsystem/<subsys>/.
type NVMeSubsystem struct {
	// Name is the subsystem directory name, e.g. "nvme-subsys0".
	Name string
	// Path is the subsystem's sysfs directory,
	// /sys/class/nvme-subsystem/<Name>.
	Path string
	// NQN is the NVMe Qualified Name from subsysnqn.
	NQN string
	// Model is the subsystem model string.
	Model string
	// Serial is the subsystem serial number.
	Serial string
	// IOPolicy is the multipath I/O policy, e.g. "numa", "round-robin".
	IOPolicy string
	// Controllers lists the NVMe controllers under this subsystem.
	Controllers []NVMeSubsystemController
	// Namespaces lists the NVMe namespace block devices (e.g. "nvme0n1")
	// that belong to this subsystem.
	Namespaces []string
}

// NVMeSubsystemNamespaceDiag holds the diagnostic counters of a multipath
// namespace head from /sys/class/nvme-subsystem/<subsys>/<namespace>/diag.
// A nil field means the kernel does not expose that counter. Use
// NVMeSubsystem.NamespaceDiag to read them.
type NVMeSubsystemNamespaceDiag struct {
	// IORequeueNoUsablePath counts I/O requeued because no path was usable
	// (diag/io_requeue_no_usable_path_count).
	IORequeueNoUsablePath *uint64
	// IOFailNoAvailablePath counts I/O failed because no path was available
	// (diag/io_fail_no_available_path_count).
	IOFailNoAvailablePath *uint64
}

// NVMeSubsystemController contains info about a single NVMe controller
// within an NVMe subsystem.
type NVMeSubsystemController struct {
	// Name is the controller directory name, e.g. "nvme0".
	Name string
	// State is the controller state, e.g. "live", "connecting", "dead".
	State string
	// Transport is the transport type, e.g. "tcp", "fc", "rdma".
	Transport string
	// Address is the controller address string.
	Address string
}

// NVMeSubsystemClass is a collection of NVMe subsystems from
// /sys/class/nvme-subsystem.
type NVMeSubsystemClass []NVMeSubsystem

// NVMeSubsystemClass returns info for all NVMe subsystems read from
// /sys/class/nvme-subsystem.
func (fs FS) NVMeSubsystemClass() (NVMeSubsystemClass, error) {
	path := fs.sys.Path(nvmeSubsystemClassPath)

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var subsystems NVMeSubsystemClass
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), "nvme-subsys") {
			continue
		}
		subsys, err := fs.parseNVMeSubsystem(entry.Name())
		if err != nil {
			return nil, err
		}
		subsystems = append(subsystems, *subsys)
	}

	return subsystems, nil
}

func (fs FS) parseNVMeSubsystem(name string) (*NVMeSubsystem, error) {
	path := fs.sys.Path(nvmeSubsystemClassPath, name)
	subsys := &NVMeSubsystem{Name: name, Path: path}

	for _, attr := range [...]struct {
		file string
		dest *string
	}{
		{"subsysnqn", &subsys.NQN},
		{"model", &subsys.Model},
		{"serial", &subsys.Serial},
		{"iopolicy", &subsys.IOPolicy},
	} {
		val, err := parsers.SysReadFile(fs.sys.Path(nvmeSubsystemClassPath, name, attr.file))
		if err != nil {
			return nil, fmt.Errorf("failed to read %s for %s: %w", attr.file, name, err)
		}
		*attr.dest = val
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list controllers for %s: %w", name, err)
	}

	seen := make(map[string]bool)
	for _, entry := range entries {
		eName := entry.Name()
		if nvmeSubsystemControllerRE.MatchString(eName) {
			ctrl, err := fs.parseNVMeSubsystemController(name, eName)
			if err != nil {
				return nil, err
			}
			subsys.Controllers = append(subsys.Controllers, *ctrl)
		}
		if nvmeNamespaceRE.MatchString(eName) && !seen[eName] {
			seen[eName] = true
			subsys.Namespaces = append(subsys.Namespaces, eName)
		}
	}

	return subsys, nil
}

func (fs FS) parseNVMeSubsystemController(subsysName, ctrlName string) (*NVMeSubsystemController, error) {
	ctrl := &NVMeSubsystemController{Name: ctrlName}

	for _, attr := range [...]struct {
		file string
		dest *string
	}{
		{"state", &ctrl.State},
		{"transport", &ctrl.Transport},
		{"address", &ctrl.Address},
	} {
		val, err := parsers.SysReadFile(fs.sys.Path(nvmeSubsystemClassPath, subsysName, ctrlName, attr.file))
		if err != nil {
			return nil, fmt.Errorf("failed to read %s for %s/%s: %w", attr.file, subsysName, ctrlName, err)
		}
		*attr.dest = val
	}

	return ctrl, nil
}

// NamespaceDiag reads the diagnostic counters of one of the subsystem's
// multipath namespace head devices, given by name as in Namespaces (e.g.
// nvme0n1), from /sys/class/nvme-subsystem/<subsys>/<namespace>/diag. They
// are read on demand rather than by NVMeSubsystemClass, so callers that do
// not need them avoid the extra file reads.
func (s NVMeSubsystem) NamespaceDiag(namespace string) (NVMeSubsystemNamespaceDiag, error) {
	var (
		diag NVMeSubsystemNamespaceDiag
		err  error
	)
	if !nvmeNamespaceRE.MatchString(namespace) {
		return diag, fmt.Errorf("invalid NVMe namespace name %q", namespace)
	}
	dir := filepath.Join(s.Path, namespace, "diag")
	if diag.IORequeueNoUsablePath, err = readOptionalUint(filepath.Join(dir, "io_requeue_no_usable_path_count")); err != nil {
		return diag, err
	}
	if diag.IOFailNoAvailablePath, err = readOptionalUint(filepath.Join(dir, "io_fail_no_available_path_count")); err != nil {
		return diag, err
	}
	return diag, nil
}
