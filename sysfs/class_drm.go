// Copyright 2018 The Prometheus Authors
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

// +build !windows

package sysfs

import (
	"errors"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/prometheus/procfs/internal/util"
)

// ClassThermalZoneStats contains info from files in /sys/class/thermal/thermal_zone<zone>
// for a single <zone>.
// https://www.kernel.org/doc/Documentation/thermal/sysfs-api.txt
type ClassDrmCardPort struct {
	Name    string
	Status  string
	Dpms    string
	Enabled string
}

func (fs FS) ClassDrmCardPort() ([]ClassDrmCardPort, error) {
	ports, err := filepath.Glob(fs.sys.Path("class/drm/card[0-9]*-*"))
	if err != nil {
		return nil, err
	}

	stats := make([]ClassDrmCardPort, 0, len(ports))
	for _, port := range ports {
		portStats, err := parseClassDrmCardPort(port)
		if err != nil {
			if errors.Is(err, syscall.ENODATA) {
				continue
			}
			return nil, err
		}
		portStats.Name = strings.TrimPrefix(filepath.Base(port), "port")
		stats = append(stats, portStats)
	}
	return stats, nil
}

func parseClassDrmCardPort(port string) (ClassDrmCardPort, error) {
	// Required attributes.
	portStatus, err := util.SysReadFile(filepath.Join(port, "status"))
	if err != nil {
		return ClassDrmCardPort{}, err
	}
	portDpms, err := util.SysReadFile(filepath.Join(port, "dpms"))
	if err != nil {
		return ClassDrmCardPort{}, err
	}
	portEnabled, err := util.SysReadFile(filepath.Join(port, "enabled"))
	if err != nil {
		return ClassDrmCardPort{}, err
	}

	return ClassDrmCardPort{
		Status:  portStatus,
		Dpms:    portDpms,
		Enabled: portEnabled,
	}, nil
}
