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

type ClassDrmCardPort struct {
	Name    string
	Card    string
	Status  uint64
	Dpms    uint64
	Enabled uint64
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
		portStats.Name = strings.SplitN(filepath.Base(port), "-", 2)[1]
		portStats.Card = strings.SplitN(filepath.Base(port), "-", 2)[0]
		stats = append(stats, portStats)
	}
	return stats, nil
}

func parseClassDrmCardPort(port string) (ClassDrmCardPort, error) {
	portStatusString, err := util.SysReadFile(filepath.Join(port, "status"))
	if err != nil {
		return ClassDrmCardPort{}, err
	}

	portStatus := 0
	if portStatusString == "connected" {
		portStatus = 1
	}

	portDpmsString, err := util.SysReadFile(filepath.Join(port, "dpms"))
	if err != nil {
		return ClassDrmCardPort{}, err
	}

	portDpms := 0
	if portDpmsString == "On" {
		portDpms = 1
	}

	portEnabledString, err := util.SysReadFile(filepath.Join(port, "enabled"))
	if err != nil {
		return ClassDrmCardPort{}, err
	}

	portEnabled := 0
	if portEnabledString == "enabled" {
		portEnabled = 1
	}

	return ClassDrmCardPort{
		Status:  uint64(portStatus),
		Dpms:    uint64(portDpms),
		Enabled: uint64(portEnabled),
	}, nil
}
