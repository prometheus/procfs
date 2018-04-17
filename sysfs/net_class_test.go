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

	var (
		AddrAssignType   int64 = 3
		AddrLen          int64 = 6
		Carrier          int64 = 1
		CarrierChanges   int64 = 2
		CarrierDownCount int64 = 1
		CarrierUpCount   int64 = 1
		DevID            int64 = 32
		Dormant          int64 = 1
		Flags            int64 = 4867
		IfIndex          int64 = 2
		IfLink           int64 = 2
		LinkMode         int64 = 1
		MTU              int64 = 1500
		NameAssignType   int64 = 2
		NetDevGroup      int64 = 0
		Speed            int64 = 1000
		TxQueueLen       int64 = 1000
		Type             int64 = 1
	)

	netClass := NetClass{
		"eth0": {
			Address:          "01:01:01:01:01:01",
			AddrAssignType:   &AddrAssignType,
			AddrLen:          &AddrLen,
			Broadcast:        "ff:ff:ff:ff:ff:ff",
			Carrier:          &Carrier,
			CarrierChanges:   &CarrierChanges,
			CarrierDownCount: &CarrierDownCount,
			CarrierUpCount:   &CarrierUpCount,
			DevID:            &DevID,
			Dormant:          &Dormant,
			Duplex:           "full",
			Flags:            &Flags,
			IfAlias:          "",
			IfIndex:          &IfIndex,
			IfLink:           &IfLink,
			LinkMode:         &LinkMode,
			MTU:              &MTU,
			Name:             "eth0",
			NameAssignType:   &NameAssignType,
			NetDevGroup:      &NetDevGroup,
			OperState:        "up",
			PhysPortID:       "",
			PhysPortName:     "",
			PhysSwitchID:     "",
			Speed:            &Speed,
			TxQueueLen:       &TxQueueLen,
			Type:             &Type,
		},
	}

	if !reflect.DeepEqual(netClass, nc) {
		t.Errorf("Result not correct: want %v, have %v", netClass, nc)
	}
}
