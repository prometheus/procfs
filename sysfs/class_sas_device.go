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
	"github.com/prometheus/procfs/internal/util"
	"io/ioutil"
	"path/filepath"
	"regexp"
)

const sasDeviceClassPath = "class/sas_device"
const sasEndDeviceClassPath = "class/sas_end_device"
const sasExpanderClassPath = "class/sas_expander"

type SASDevice struct {
	Name       string   // /sys/class/sas_device/<Name>
	SASAddress string   // /sys/class/sas_device/<Name>/sas_address
	SASPhys    []string // /sys/class/sas_device/<Name>/device/phy-*
	SASPorts   []string // /sys/class/sas_device/<Name>/device/ports-*
	BlockDevices []string  // /sys/class/sas_device/<Name>/device/target*/*/block/*
}

type SASDeviceClass map[string]SASDevice

// SASDeviceClass parses devices in /sys/class/sas_device.
func (fs FS) SASDeviceClass() (SASDeviceClass, error) {
	path := fs.sys.Path(sasDeviceClassPath)

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

// SASEndDeviceClass parses devices in /sys/class/sas_end_device.
// This is *almost* identical to sas_device, just with a different
// base directory.  The major difference is that end_devices don't
// include expanders and other infrastructure devices.
func (fs FS) SASEndDeviceClass() (SASDeviceClass, error) {
	path := fs.sys.Path(sasEndDeviceClassPath)

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

// SASExpanderClass parses devices in /sys/class/sas_expander.
// This is *almost* identical to sas_device, but only includes expanders.
func (fs FS) SASExpanderClass() (SASDeviceClass, error) {
	path := fs.sys.Path(sasExpanderClassPath)

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

// Parse a single sas_device.
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
		if phyDevice.Match([]byte(d.Name())) {
			device.SASPhys = append(device.SASPhys, d.Name())
		}
		if portDevice.Match([]byte(d.Name())) {
			device.SASPorts = append(device.SASPorts, d.Name())
		}
	}

	address := fs.sys.Path(sasDeviceClassPath, name, "sas_address")
	value, err := util.SysReadFile(address)

	if err != nil {
		return &device, err
	} else {
		device.SASAddress = value
	}

	device.BlockDevices, err = fs.blockSASDeviceBlockDevices(name)
	if err != nil {
		return &device, err
	}

	return &device, nil
}

// Identify block devices that map to a specific SAS Device
// This info comes from (for example) /sys/class/sas_device/end_device-11:2/device/target11:0:0/11:0:0:0/block/sdp
//
// To find that, we have to look in the device directory for target$X
// subdirs, then a subdir of $X, then read from directory names in the
// 'block/' subdirectory under that.
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
