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

package sysfs

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

// NetClassIface contains info from files in /sys/class/net/<iface>
// for single interface (iface).
type NetClassIface struct {
	Name             string // Interface name
	AddrAssignType   uint64 `fileName:"addr_assign_type"`   // /sys/class/net/<iface>/addr_assign_type
	AddrLen          uint64 `fileName:"addr_len"`           // /sys/class/net/<iface>/addr_len
	Address          string `fileName:"address"`            // /sys/class/net/<iface>/address
	Broadcast        string `fileName:"broadcast"`          // /sys/class/net/<iface>/broadcast
	Carrier          uint64 `fileName:"carrier"`            // /sys/class/net/<iface>/carrier
	CarrierChanges   uint64 `fileName:"carrier_changes"`    // /sys/class/net/<iface>/carrier_changes
	CarrierUpCount   uint64 `fileName:"carrier_up_count"`   // /sys/class/net/<iface>/carrier_up_count
	CarrierDownCount uint64 `fileName:"carrier_down_count"` // /sys/class/net/<iface>/carrier_down_count
	DevID            uint64 `fileName:"dev_id"`             // /sys/class/net/<iface>/dev_id
	Dormant          uint64 `fileName:"dormant"`            // /sys/class/net/<iface>/dormant
	Duplex           string `fileName:"duplex"`             // /sys/class/net/<iface>/duplex
	Flags            uint64 `fileName:"flags"`              // /sys/class/net/<iface>/flags
	IfAlias          string `fileName:"ifalias"`            // /sys/class/net/<iface>/ifalias
	IfIndex          uint64 `fileName:"ifindex"`            // /sys/class/net/<iface>/ifindex
	IfLink           uint64 `fileName:"iflink"`             // /sys/class/net/<iface>/iflink
	LinkMode         uint64 `fileName:"link_mode"`          // /sys/class/net/<iface>/link_mode
	MTU              uint64 `fileName:"mtu"`                // /sys/class/net/<iface>/mtu
	NameAssignType   uint64 `fileName:"name_assign_type"`   // /sys/class/net/<iface>/name_assign_type
	NetDevGroup      uint64 `fileName:"netdev_group"`       // /sys/class/net/<iface>/netdev_group
	OperState        string `fileName:"operstate"`          // /sys/class/net/<iface>/operstate
	PhysPortID       string `fileName:"phys_port_id"`       // /sys/class/net/<iface>/phys_port_id
	PhysPortName     string `fileName:"phys_port_name"`     // /sys/class/net/<iface>/phys_port_name
	PhysSwitchID     string `fileName:"phys_switch_id"`     // /sys/class/net/<iface>/phys_switch_id
	Speed            uint64 `fileName:"speed"`              // /sys/class/net/<iface>/speed
	TxQueueLen       uint64 `fileName:"tx_queue_len"`       // /sys/class/net/<iface>/tx_queue_len
	Type             uint64 `fileName:"type"`               // /sys/class/net/<iface>/type
}

// NetClass is collection of info for every interface (iface) in /sys/class/net. The map keys
// are interface (iface) names.
type NetClass map[string]NetClassIface

// NewNetClass returns info for all net interfaces (iface) read from /sys/class/net/<iface>.
func NewNetClass() (NetClass, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return nil, err
	}

	return fs.NewNetClass()
}

// NewNetClass returns info for all net interfaces (iface) read from /sys/class/net/<iface>.
func (fs FS) NewNetClass() (NetClass, error) {
	path := fs.Path("class/net")

	devices, err := ioutil.ReadDir(path)
	if err != nil {
		return NetClass{}, fmt.Errorf("cannot access %s dir %s", path, err)
	}

	netClass := NetClass{}
	for _, deviceDir := range devices {
		interfaceClass, err := netClass.parseNetClassIface(path + "/" + deviceDir.Name())
		if err != nil {
			continue
		}
		interfaceClass.Name = deviceDir.Name()
		netClass[deviceDir.Name()] = *interfaceClass
	}
	return netClass, nil
}

// parseNetClassIface scans predefined files in /sys/class/net/<iface>
// directory and gets their contents.
func (nc NetClass) parseNetClassIface(devicePath string) (*NetClassIface, error) {
	interfaceClass := NetClassIface{}
	interfaceElem := reflect.ValueOf(&interfaceClass).Elem()
	interfaceType := reflect.TypeOf(interfaceClass)

	for i := 0; i < interfaceElem.NumField(); i++ {
		fieldType := interfaceType.Field(i)
		fieldValue := interfaceElem.Field(i)

		if fieldType.Tag.Get("fileName") == "" {
			continue
		}

		fileContents, err := ioutil.ReadFile(devicePath + "/" + fieldType.Tag.Get("fileName"))
		if err != nil {
			continue
		}
		value := strings.TrimSpace(string(fileContents))

		switch fieldValue.Kind() {
		case reflect.Uint64:
			if strings.HasPrefix(value, "0x") {
				intValue, err := strconv.ParseUint(value[2:], 16, 64)
				if err != nil {
					continue
				}
				fieldValue.SetUint(intValue)
			} else {
				intValue, err := strconv.ParseUint(value, 10, 64)
				if err != nil {
					continue
				}
				fieldValue.SetUint(intValue)
			}
		case reflect.String:
			fieldValue.SetString(value)
		default:
			return nil, fmt.Errorf("unhandled type %q", fieldValue.Kind())
		}
	}

	return &interfaceClass, nil
}
