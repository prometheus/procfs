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
	"reflect"
	"testing"
)

func TestNewNetClassDevices(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	devices, err := fs.NetClassDevices()
	if err != nil {
		t.Fatal(err)
	}

	if len(devices) != 1 {
		t.Errorf("Unexpected number of devices, want %d, have %d", 1, len(devices))
	}
	if devices[0] != "eth0" {
		t.Errorf("Found unexpected device, want %s, have %s", "eth0", devices[0])
	}
}

func TestNetClass(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	nc, err := fs.NetClass()
	if err != nil {
		t.Fatal(err)
	}

	var (
		addrAssignType   int64 = 3
		addrLen          int64 = 6
		carrier          int64 = 1
		carrierChanges   int64 = 2
		carrierDownCount int64 = 1
		carrierUpCount   int64 = 1
		devID            int64 = 32
		dormant          int64 = 1
		flags            int64 = 4867
		ifIndex          int64 = 2
		ifLink           int64 = 2
		linkMode         int64 = 1
		mtu              int64 = 1500
		nameAssignType   int64 = 2
		netDevGroup      int64
		speed            int64 = 1000
		txQueueLen       int64 = 1000
		netType          int64 = 1
		collisions       int64
		multicast        int64 = 685
		rxBytes          int64 = 67684
		rxCompressed     int64
		rxCrcErrors      int64
		rxDropped        int64
		rxErrrors        int64
		rxFifoErrors     int64
		rxFrameErrors    int64
		rxLengthErrors   int64
		rxMissedErrors   int64
		rxNoHandler      int64
		rxOverErrors     int64
		rxPackets        int64 = 869
		txAbortedErrors  int64
		txBytes          int64 = 2113
		txCarrierErrors  int64
		txCompressed     int64
		txDropped        int64
		txErrors         int64
		txFifoErrors     int64
		txHeartbtErrors  int64
		txPackets        int64 = 11
		txWindowErrors   int64
	)

	netClass := NetClass{
		"eth0": {
			Address:          "01:01:01:01:01:01",
			AddrAssignType:   &addrAssignType,
			AddrLen:          &addrLen,
			Broadcast:        "ff:ff:ff:ff:ff:ff",
			Carrier:          &carrier,
			CarrierChanges:   &carrierChanges,
			CarrierDownCount: &carrierDownCount,
			CarrierUpCount:   &carrierUpCount,
			DevID:            &devID,
			Dormant:          &dormant,
			Duplex:           "full",
			Flags:            &flags,
			IfAlias:          "",
			IfIndex:          &ifIndex,
			IfLink:           &ifLink,
			LinkMode:         &linkMode,
			MTU:              &mtu,
			Name:             "eth0",
			NameAssignType:   &nameAssignType,
			NetDevGroup:      &netDevGroup,
			OperState:        "up",
			PhysPortID:       "",
			PhysPortName:     "",
			PhysSwitchID:     "",
			Speed:            &speed,
			TxQueueLen:       &txQueueLen,
			Type:             &netType,
			Collisions:       &collisions,
			Multicast:        &multicast,
			RxBytes:          &rxBytes,
			RxCompressed:     &rxCompressed,
			RxCrcErrors:      &rxCrcErrors,
			RxDropped:        &rxDropped,
			RxErrors:         &rxErrrors,
			RxFifoErrors:     &rxFifoErrors,
			RxFrameErrors:    &rxFrameErrors,
			RxLengthErrors:   &rxLengthErrors,
			RxMissedErrors:   &rxMissedErrors,
			RxNoHandler:      &rxNoHandler,
			RxOverErrors:     &rxOverErrors,
			RxPackets:        &rxPackets,
			TxAbortedErrors:  &txAbortedErrors,
			TxBytes:          &txBytes,
			TxCarrierErrors:  &txCarrierErrors,
			TxCompressed:     &txCompressed,
			TxDropped:        &txDropped,
			TxErrors:         &txErrors,
			TxFifoErrors:     &txFifoErrors,
			TxHeartbtErrors:  &txHeartbtErrors,
			TxPackets:        &txPackets,
			TxWindowErrors:   &txWindowErrors,
		},
	}

	if !reflect.DeepEqual(netClass, nc) {
		t.Errorf("Result not correct: want %v, have %v", netClass, nc)
	}
}
