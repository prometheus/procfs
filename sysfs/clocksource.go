// Copyright 2019 The Prometheus Authors
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
	"path/filepath"
	"strings"

	"github.com/prometheus/procfs/internal/util"
)

// Clocksource contains metrics related to the clock source
type Clocksource struct {
	Name      string
	Available []string
	Current   string
}

// NewClocksource returns clocksource information including current and available clocksources.
func NewClocksource() ([]Clocksource, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return []Clocksource{}, err
	}

	return fs.NewClocksource()
}

// NewClocksource returns clocksource information including current and available clocksources.
func (fs FS) NewClocksource() ([]Clocksource, error) {

	clocksourcePaths, err := filepath.Glob(fs.sys.Path("devices/system/clocksource/clocksource[0-9]*"))
	if err != nil {
		return nil, err
	}

	clocksources := make([]Clocksource, len(clocksourcePaths))
	for i, clocksourcePath := range clocksourcePaths {
		clocksourceName := strings.TrimPrefix(filepath.Base(clocksourcePath), "clocksource")

		clocksource, err := parseClocksource(clocksourcePath)
		if err != nil {
			return nil, err
		}
		clocksource.Name = clocksourceName
		clocksources[i] = *clocksource
	}

	return clocksources, nil
}

func parseClocksource(clocksourcePath string) (*Clocksource, error) {

	stringFiles := []string{
		"available_clocksource",
		"current_clocksource",
	}
	stringOut := make([]string, len(stringFiles))
	var err error

	for i, f := range stringFiles {
		stringOut[i], err = util.SysReadFile(filepath.Join(clocksourcePath, f))
		if err != nil {
			return &Clocksource{}, err
		}
	}

	return &Clocksource{
		Available: strings.Fields(stringOut[0]),
		Current:   stringOut[1],
	}, nil
}
