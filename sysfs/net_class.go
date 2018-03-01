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
	"io/ioutil"
	"fmt"
	"reflect"
	"strings"
	"strconv"
)

type interfaceNetClass struct {
	Name             string `json:"name"`               // Interface name
	AddrAssignType   uint64 `json:"addr_assign_type"`   // /sys/class/net/<iface>/addr_assign_type
	AddrLen          uint64 `json:"addr_len"`           // /sys/class/net/<iface>/addr_len
	Address          string `json:"address"`            // /sys/class/net/<iface>/address
	Broadcast        string `json:"broadcast"`          // /sys/class/net/<iface>/broadcast
	Carrier          uint64 `json:"carrier"`            // /sys/class/net/<iface>/carrier
	CarrierChanges   uint64 `json:"carrier_changes"`    // /sys/class/net/<iface>/carrier_changes
	CarrierUpCount   uint64 `json:"carrier_up_count"`   // /sys/class/net/<iface>/carrier_up_count
	CarrierDownCount uint64 `json:"carrier_down_count"` // /sys/class/net/<iface>/carrier_down_count
	DevId            uint64 `json:"dev_id"`             // /sys/class/net/<iface>/dev_id
	Dormant          uint64 `json:"dormant"`            // /sys/class/net/<iface>/dormant
	Duplex           string `json:"duplex"`             // /sys/class/net/<iface>/duplex
	Flags            uint64 `json:"flags"`              // /sys/class/net/<iface>/flags
	IfAlias          string `json:"ifalias"`            // /sys/class/net/<iface>/ifalias
	IfIndex          uint64 `json:"ifindex"`            // /sys/class/net/<iface>/ifindex
	IfLink           uint64 `json:"iflink"`             // /sys/class/net/<iface>/iflink
	LinkMode         uint64 `json:"link_mode"`          // /sys/class/net/<iface>/link_mode
	Mtu              uint64 `json:"mtu"`                // /sys/class/net/<iface>/link_mode
	NameAssignType   string `json:"name_assign_type"`   // /sys/class/net/<iface>/name_assign_type
	NetDevGroup      uint64 `json:"netdev_group"`       // /sys/class/net/<iface>/netdev_group
	OperState        string `json:"operstate"`          // /sys/class/net/<iface>/operstate
	PhysPortId       string `json:"phys_port_id"`       // /sys/class/net/<iface>/phys_port_id
	PhysPortName     string `json:"phys_port_name"`     // /sys/class/net/<iface>/phys_port_name
	PhysSwitchId     uint64 `json:"phys_switch_id"`     // /sys/class/net/<iface>/phys_switch_id
	Speed            uint64 `json:"speed"`              // /sys/class/net/<iface>/speed
	TxQueueLen       uint64 `json:"tx_queue_len"`       // /sys/class/net/<iface>/tx_queue_len
	Type             uint64 `json:"type"`               // /sys/class/net/<iface>/type
}

type NetClass map[string]interfaceNetClass

// NewNetDev returns kernel/system statistics read from /proc/net/dev.
func NewNetClass() (NetClass, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return nil, err
	}

	return fs.NewNetClass()
}

func (fs FS) NewNetClass() (NetClass, error) {
	return newNetClass(fs.Path("class/net"))
}

func newNetClass(path string) (NetClass, error) {
	devices, err := ioutil.ReadDir(path)
	if err != nil {
		return NetClass{}, fmt.Errorf("cannot access %s dir %s", path, err)
	}

	netclass := NetClass{}
	for _, deviceDir := range devices {
		interfaceClass, err := netclass.parseInterfaceNetClass(path + "/" + deviceDir.Name())
		if err != nil {
			//log?
			continue
		}
		interfaceClass.Name = deviceDir.Name()
		netclass[deviceDir.Name()] = *interfaceClass
	}
	return netclass, nil
}

func (nc NetClass) parseInterfaceNetClass(devicePath string) (*interfaceNetClass, error) {
	fields := []string{"AddrAssignType", "AddrLen", "Address", "Broadcast", "Carrier", "CarrierChanges", "CarrierUpCount", "CarrierDownCount", "DevId", "Dormant", "Duplex", "Flags", "IfAlias", "IfIndex", "IfLink", "LinkMode", "Mtu", "NameAssignType", "NetDevGroup", "OperState", "PhysPortId", "PhysPortName", "PhysSwitchId", "Speed", "TxQueueLen", "Type"}
	interfaceClass := interfaceNetClass{}
	interfaceElem := reflect.ValueOf(&interfaceClass).Elem()
	interfaceType := reflect.TypeOf(interfaceClass)

	for _, fieldName := range fields {
		fieldType, found := interfaceType.FieldByName(fieldName)
		if !found {
			continue
		}
		fieldValue := interfaceElem.FieldByName(fieldName)



		fileContents, err := ioutil.ReadFile(devicePath + "/" + fieldType.Tag.Get("json"))
		if err != nil {
			continue
		}
		value := strings.TrimSpace(string(fileContents))

		if fieldValue.Kind() == reflect.Uint64 {
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
		} else if fieldValue.Kind() == reflect.String {
			fieldValue.SetString(value)
		}
	}

	return &interfaceClass, nil
}