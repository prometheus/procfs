// Copyright 2022 The Prometheus Authors
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
	"github.com/prometheus/procfs/internal/util"
	"io/ioutil"
	"path/filepath"
	"regexp"
)

const sasDeviceClassPath = "class/sas_device"
const sasEndDeviceClassPath = "class/sas_end_device"
const sasExpanderClassPath = "class/sas_expander"

type SASDevice struct {
	Name         string   // /sys/class/sas_device/<Name>
	SASAddress   string   // /sys/class/sas_device/<Name>/sas_address
	SASPhys      []string // /sys/class/sas_device/<Name>/device/phy-*
	SASPorts     []string // /sys/class/sas_device/<Name>/device/ports-*
	BlockDevices []string // /sys/class/sas_device/<Name>/device/target*/*/block/*
}

type SASDeviceClass map[string]SASDevice

// sasDeviceClasses reads all of the SAS devices from a specific set
// of /sys/class/sas*/ entries.  The sas_device, sas_end_device, and
// sas_expander classes are all nearly identical and can be handled by the same basic code.

func (fs FS) parseSASDeviceClass(dir string) (SASDeviceClass, error) {
	path := fs.sys.Path(dir)

	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	sdc := make(SASDeviceClass, len(dirs))

	for _, d := range dirs {
		device, err := fs.parseSASDevice(d.Name())
		if err != nil {
			return nil, err
		}

		sdc[device.Name] = *device
	}

	return sdc, nil
}

// SASDeviceClass parses devices in /sys/class/sas_device.
func (fs FS) SASDeviceClass() (SASDeviceClass, error) {
	return fs.parseSASDeviceClass(sasDeviceClassPath)
}

// SASEndDeviceClass parses devices in /sys/class/sas_end_device.
// This is a subset of sas_device, and excludes expanders.
func (fs FS) SASEndDeviceClass() (SASDeviceClass, error) {
	return fs.parseSASDeviceClass(sasEndDeviceClassPath)
}

// SASExpanderClass parses devices in /sys/class/sas_expander.
// This is a subset of sas_device, but only includes expanders.
func (fs FS) SASExpanderClass() (SASDeviceClass, error) {
	return fs.parseSASDeviceClass(sasExpanderClassPath)
}

// Parse a single sas_device.  This uses /sys/class/sas_device, as
// it's a superset of the other two directories so there's no reason
// to plumb the path through to here.
func (fs FS) parseSASDevice(name string) (*SASDevice, error) {
	device := SASDevice{Name: name}

	devicepath := fs.sys.Path(filepath.Join(sasDeviceClassPath, name, "device"))

	dirs, err := ioutil.ReadDir(devicepath)
	if err != nil {
		return nil, err
	}

	phyDevice := regexp.MustCompile(`^phy-[0-9:]+$`)
	portDevice := regexp.MustCompile(`^port-[0-9:]+$`)

	for _, d := range dirs {
		if phyDevice.MatchString(d.Name()) {
			device.SASPhys = append(device.SASPhys, d.Name())
		}
		if portDevice.MatchString(d.Name()) {
			device.SASPorts = append(device.SASPorts, d.Name())
		}
	}

	address := fs.sys.Path(sasDeviceClassPath, name, "sas_address")
	value, err := util.SysReadFile(address)
	if err != nil {
		return nil, err
	}
	device.SASAddress = value

	device.BlockDevices, err = fs.blockSASDeviceBlockDevices(name)
	if err != nil {
		return nil, err
	}

	return &device, nil
}

// Identify block devices that map to a specific SAS Device
// This info comes from (for example)
// /sys/class/sas_device/end_device-11:2/device/target11:0:0/11:0:0:0/block/sdp
//
// To find that, we have to look in the device directory for target$X
// subdirs, then specific subdirs of $X, then read from directory
// names in the 'block/' subdirectory under that.  This really
// shouldn't be this hard.
func (fs FS) blockSASDeviceBlockDevices(name string) ([]string, error) {
	var devices []string

	devicepath := fs.sys.Path(filepath.Join(sasDeviceClassPath, name, "device"))

	dirs, err := ioutil.ReadDir(devicepath)
	if err != nil {
		return nil, err
	}

	targetDevice := regexp.MustCompile(`^target[0-9:]+$`)
	targetSubDevice := regexp.MustCompile(`[0-9]+:.*`)

	for _, d := range dirs {
		if targetDevice.MatchString(d.Name()) {
			targetdir := d.Name()

			subtargets, err := ioutil.ReadDir(filepath.Join(devicepath, targetdir))
			if err != nil {
				return nil, err
			}

			for _, targetsubdir := range subtargets {

				if !targetSubDevice.MatchString(targetsubdir.Name()) {
					// need to skip 'power', 'subsys', etc.
					continue
				}

				blocks, err := ioutil.ReadDir(filepath.Join(devicepath, targetdir, targetsubdir.Name(), "block"))

				if err != nil {
					return nil, err
				}

				for _, blockdevice := range blocks {
					devices = append(devices, blockdevice.Name())
				}
			}
		}
	}

	return devices, nil
}
