// Copyright 2021 The Prometheus Authors
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
// +build linux

package sysfs

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
)

const sasPortClassPath = "class/sas_port"

type SASPort struct {
	Name       string   // /sys/class/sas_device/<Name>
	SASPhys    []string // /sys/class/sas_device/<Name>/device/phy-*
	Expanders  []string // /sys/class/sas_port/<Name>/device/expander-*
}

type SASPortClass map[string]SASPort

// SASPortClass parses ports in /sys/class/sas_port.
func (fs FS) SASPortClass() (SASPortClass, error) {
	path := fs.sys.Path(sasPortClassPath)

	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	spc := make(SASPortClass, len(dirs))

	for _, d := range dirs {
		port, err := fs.parseSASPort(d.Name())
		if err != nil {
			return nil, err
		}

		spc[port.Name] = *port
	}

	return spc, nil
}

// Parse a single sas_port.
func (fs FS) parseSASPort(name string) (*SASPort, error) {
	port := SASPort{Name: name}

	portpath := fs.sys.Path(filepath.Join(sasPortClassPath, name, "device"))

	dirs, err := ioutil.ReadDir(portpath)
	if err != nil {
		return nil, err
	}

	phyDevice := regexp.MustCompile(`^phy-[0-9:]+$`)
	expanderDevice := regexp.MustCompile(`^expander-[0-9:]+$`)

	for _, d := range dirs {
		if phyDevice.Match([]byte(d.Name())) {
			port.SASPhys = append(port.SASPhys, d.Name())
		}
		if expanderDevice.Match([]byte(d.Name())) {
			port.Expanders = append(port.Expanders, d.Name())
		}
	}

	return &port, nil
}
