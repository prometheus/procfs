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

// Package resctrlfs provides access to the monitoring data of the resctrl
// filesystem.
//
// resctrl is the kernel interface to the cache and memory bandwidth resource
// control of the CPU: Intel calls it RDT (Resource Director Technology), AMD
// calls it PQoS (Platform Quality of Service), Arm calls it MPAM (Memory System
// Resource Partitioning and Monitoring). All vendors are served by the same
// kernel driver, so the file layout is identical.
//
// The filesystem is not mounted by default. Mount it with:
//
//	mount -t resctrl resctrl /sys/fs/resctrl
//
// This package only reads monitoring data. It does not configure allocation.
//
// The filesystem is documented in the kernel tree:
//   - https://docs.kernel.org/filesystems/resctrl.html
//   - https://docs.kernel.org/arch/arm64/mpam.html
package resctrlfs

import (
	"github.com/prometheus/procfs/internal/fs"
)

// FS represents the pseudo-filesystem resctrl, which provides an interface to
// the cache and memory bandwidth monitoring of the CPU.
type FS struct {
	resctrl fs.FS
}

// DefaultMountPoint is the common mount point of the resctrl filesystem.
const DefaultMountPoint = fs.DefaultResctrlMountPoint

// NewDefaultFS returns a new FS mounted under the default mountPoint. It will error
// if the mount point can't be read.
func NewDefaultFS() (FS, error) {
	return NewFS(DefaultMountPoint)
}

// NewFS returns a new FS mounted under the given mountPoint. It will error
// if the mount point can't be read.
func NewFS(mountPoint string) (FS, error) {
	fs, err := fs.NewFS(mountPoint)
	if err != nil {
		return FS{}, err
	}
	return FS{fs}, nil
}
