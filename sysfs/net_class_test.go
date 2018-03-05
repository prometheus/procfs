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
	"reflect"
	"testing"
)

func TestNewNetClass(t *testing.T) {
	fs, err := NewFS("fixtures")
	if err != nil {
		t.Fatal(err)
	}

	nc, err := fs.NewNetClass()
	if err != nil {
		t.Fatal(err)
	}

	netClass := map[string]interfaceNetClass{
		"eth0": {Name: "eth0", AddrAssignType: 3, AddrLen: 6, Address: "01:01:01:01:01:01", Broadcast: "ff:ff:ff:ff:ff:ff", Carrier: 1, CarrierChanges: 2, CarrierUpCount: 1, CarrierDownCount: 1, DevId: 32, Dormant: 1, Duplex: "full", Flags: 4867, IfAlias: "", IfIndex: 2, IfLink: 2, LinkMode: 1, Mtu: 1500, NameAssignType: 2, NetDevGroup: 0, OperState: "up", PhysPortId: "", PhysPortName: "", PhysSwitchId: "", Speed: 1000, TxQueueLen: 1000, Type: 1},
	}

	if want, have := len(netClass), len(nc); want != have {
		t.Errorf("want %d parsed class/net, have %d", want, have)
	}

	fields := []string{"AddrAssignType", "AddrLen", "Address", "Broadcast", "Carrier", "CarrierChanges", "CarrierUpCount", "CarrierDownCount", "DevId", "Dormant", "Duplex", "Flags", "IfAlias", "IfIndex", "IfLink", "LinkMode", "Mtu", "NameAssignType", "NetDevGroup", "OperState", "PhysPortId", "PhysPortName", "PhysSwitchId", "Speed", "TxQueueLen", "Type"}
	for _, interfaceClass := range nc {
		if want, have := netClass[interfaceClass.Name], interfaceClass; want != have {
			t.Errorf("%s: want %v, have %v", interfaceClass.Name, want, have)
		}

		haveElem := reflect.ValueOf(&interfaceClass).Elem()
		wantElem := reflect.ValueOf(netClass["eth0"])

		for _, fieldName := range fields {
			haveValue := haveElem.FieldByName(fieldName)
			wantValue := wantElem.FieldByName(fieldName)

			if want, have := haveValue.Kind(), wantValue.Kind(); want != have {
				t.Errorf("%s Kind: want %v, have %v", fieldName, want, have)
			}
			if haveValue.Kind() == reflect.Uint64 && wantValue.Kind() == reflect.Uint64 {
				if want, have := haveValue.Uint(), wantValue.Uint(); want != have {
					t.Errorf("%s uint64 Value: want %v, have %v", fieldName, want, have)
				}
			}
			if haveValue.Kind() == reflect.String && wantValue.Kind() == reflect.String {
				if want, have := haveValue.String(), wantValue.String(); want != have {
					t.Errorf("%s string Value: want %v, have %v", fieldName, want, have)
				}
			}

		}
	}
}
