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
)

// NetClassVFPCIAddress returns the PCI address of a Virtual Function (VF)
// for the given Physical Function (PF) network interface by resolving the
// sysfs virtfn symlink at /sys/class/net/<iface>/device/virtfn<vfIndex>.
// Returns the PCI BDF address (e.g. "0000:65:01.0").
func (fs FS) NetClassVFPCIAddress(iface string, vfIndex uint32) (string, error) {
	virtfnPath := fs.sys.Path(netclassPath, iface, "device", fmt.Sprintf("virtfn%d", vfIndex))
	resolved, err := os.Readlink(virtfnPath)
	if err != nil {
		return "", fmt.Errorf("failed to read virtfn symlink %q: %w", virtfnPath, err)
	}
	return filepath.Base(resolved), nil
}
