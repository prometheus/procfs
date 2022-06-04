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

const sasHostClassPath = "class/sas_host"

type SASHost struct {
	Name     string           // /sys/class/sas_host/<Name>
	SASPhys []string         // /sys/class/sas_host/<Name>/device/phy-*
	SASPorts []string       // /sys/class/sas_host/<Name>/device/ports-*
}

type SASHostClass map[string]SASHost

// SASHostClass parses host[0-9]+ devices in /sys/class/sas_host.
// This generally only exists so that it can pull in SAS Port and SAS
// PHY entries.
//
// The sas_host class doesn't collect any obvious statistics.  Each
// sas_host contains a scsi_host, which seems to collect a couple
// minor stats (ioc_reset_count and reply_queue_count), but they're
// not worth collecting at this time.
func (fs FS) SASHostClass() (SASHostClass, error) {
	path := fs.sys.Path(sasHostClassPath)

	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	shc := make(SASHostClass, len(dirs))

	for _, d := range dirs {
		host, err := fs.parseSASHost(d.Name())
		if err != nil {
			return nil, err
		}

		shc[host.Name] = *host
	}

	return shc, nil
}

// Parse a single sas_host.
func (fs FS) parseSASHost(name string) (*SASHost, error) {
	//path := fs.sys.Path(sasHostClassPath, name)
	host := SASHost{Name: name}

	devicepath := fs.sys.Path(filepath.Join(sasHostClassPath,name,"device"))

	dirs, err := ioutil.ReadDir(devicepath)
	if err != nil {
		return nil, err
	}

	phyDevice := regexp.MustCompile(`^phy-[0-9:]+$`)
	portDevice := regexp.MustCompile(`^port-[0-9:]+$`)

	for _, d := range dirs {
		if phyDevice.Match([]byte(d.Name())) {
			host.SASPhys = append(host.SASPhys, d.Name())
		}
		if portDevice.Match([]byte(d.Name())) {
			host.SASPorts = append(host.SASPorts, d.Name())
		}
	}

	return &host, nil
}
