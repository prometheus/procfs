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

//go:build !windows

package procfs

import (
	"strings"

	"github.com/prometheus/procfs/internal/parsers"
)

// DeviceTreeSerialNumber reads the system serial number from
// /proc/device-tree/serial-number, following the device-tree symlink.
// The terminating NUL byte is removed. An error is returned if the property
// cannot be read, including on systems without a device tree.
// See https://www.kernel.org/doc/Documentation/ABI/testing/sysfs-firmware-ofw.
func (fs FS) DeviceTreeSerialNumber() (string, error) {
	data, err := parsers.ReadFileNoStat(fs.proc.Path("device-tree", "serial-number"))
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(data), "\x00"), nil
}
