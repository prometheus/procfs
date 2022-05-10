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

//go:build linux
// +build linux

package sysfs

import (
	"net"
	"reflect"
	"testing"

	"github.com/prometheus/procfs/internal/util"
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

	if len(devices) != 2 {
		t.Errorf("Unexpected number of devices, want %d, have %d", 2, len(devices))
	}
	if !reflect.DeepEqual(devices, []string{"bond0", "eth0"}) {
		t.Errorf("Found unexpected device, want %s, have %s", "eth0", devices[0])
	}
}

func TestNewNetClassDevicesByIface(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fs.NetClassByIface("non-existent")
	if err == nil {
		t.Fatal("expected error, have none")
	}

	device, err := fs.NetClassByIface("eth0")
	if err != nil {
		t.Fatal(err)
	}

	if device.Name != "eth0" {
		t.Errorf("Found unexpected device, want %s, have %s", "eth0", device.Name)
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
	netClass := NetClass{
		"bond0": {
			Address:          "02:02:02:02:02:02",
			AddrAssignType:   util.NewValueParser("3").PInt64(),
			AddrLen:          util.NewValueParser("6").PInt64(),
			Broadcast:        "ff:ff:ff:ff:ff:ff",
			Carrier:          util.NewValueParser("1").PInt64(),
			CarrierChanges:   util.NewValueParser("2").PInt64(),
			CarrierDownCount: util.NewValueParser("1").PInt64(),
			CarrierUpCount:   util.NewValueParser("1").PInt64(),
			DevID:            util.NewValueParser("32").PInt64(),
			Dormant:          util.NewValueParser("1").PInt64(),
			Duplex:           "full",
			Flags:            util.NewValueParser("4867").PInt64(),
			IfAlias:          "",
			IfIndex:          util.NewValueParser("2").PInt64(),
			IfLink:           util.NewValueParser("2").PInt64(),
			LinkMode:         util.NewValueParser("1").PInt64(),
			MTU:              util.NewValueParser("1500").PInt64(),
			Name:             "bond0",
			NameAssignType:   util.NewValueParser("2").PInt64(),
			NetDevGroup:      util.NewValueParser("0").PInt64(),
			OperState:        "up",
			PhysPortID:       "",
			PhysPortName:     "",
			PhysSwitchID:     "",
			Speed:            util.NewValueParser("1000").PInt64(),
			TxQueueLen:       util.NewValueParser("1000").PInt64(),
			Type:             util.NewValueParser("1").PInt64(),
			BondAttrs: &NetClassBondAttrs{
				AdActorKey:                             util.NewValueParser("15").PUInt64(),
				AdActorSysPriority:                     util.NewValueParser("65535").PUInt64(),
				AdActorSystem:                          makeMAC("00:00:00:00:00:00"),
				AdAggregator:                           util.NewValueParser("1").PUInt64(),
				AdNumPorts:                             util.NewValueParser("2").PUInt64(),
				AdPartnerKey:                           util.NewValueParser("1034").PUInt64(),
				AdPartnerMac:                           makeMAC("01:23:45:67:89:AB"),
				AdSelect:                               strPTR("stable 0"),
				AdUserPortKey:                          util.NewValueParser("0").PUInt64(),
				AllDevicesActive:                       util.ParseBool("0"),
				ARPAllTargets:                          strPTR("any 0"),
				ARPInterval:                            util.NewValueParser("0").PInt64(),
				DownDelay:                              util.NewValueParser("200").PInt64(),
				FailoverMac:                            strPTR("none 0"),
				LACPRate:                               strPTR("slow 0"),
				LPInterval:                             util.NewValueParser("1").PInt64(),
				MIIMon:                                 util.NewValueParser("100").PInt64(),
				MIIStatus:                              util.ParseBool("1"),
				MinLinks:                               util.NewValueParser("0").PUInt64(),
				Mode:                                   strPTR("802.3ad 4"),
				NumberGratuitousArp:                    util.NewValueParser("1").PUInt64(),
				NumberUnsolicitedNeighborAdvertisement: util.NewValueParser("1").PUInt64(),
				PacketsPerDevice:                       util.NewValueParser("1").PInt64(),
				PrimaryReselect:                        strPTR("always 0"),
				ResendIgmp:                             util.NewValueParser("1").PInt64(),
				TLBDynamicLB:                           util.NewValueParser("1").PInt64(),
				UpDelay:                                util.NewValueParser("0").PInt64(),
				UseCarrier:                             util.ParseBool("1"),
				TransmitHashPolicy:                     strPTR("layer3+4 1"),
			},
		},
		"eth0": {
			Address:          "01:01:01:01:01:01",
			AddrAssignType:   util.NewValueParser("3").PInt64(),
			AddrLen:          util.NewValueParser("6").PInt64(),
			Broadcast:        "ff:ff:ff:ff:ff:ff",
			Carrier:          util.NewValueParser("1").PInt64(),
			CarrierChanges:   util.NewValueParser("2").PInt64(),
			CarrierDownCount: util.NewValueParser("1").PInt64(),
			CarrierUpCount:   util.NewValueParser("1").PInt64(),
			DevID:            util.NewValueParser("32").PInt64(),
			Dormant:          util.NewValueParser("1").PInt64(),
			Duplex:           "full",
			Flags:            util.NewValueParser("4867").PInt64(),
			IfAlias:          "",
			IfIndex:          util.NewValueParser("2").PInt64(),
			IfLink:           util.NewValueParser("2").PInt64(),
			LinkMode:         util.NewValueParser("1").PInt64(),
			MTU:              util.NewValueParser("1500").PInt64(),
			Name:             "eth0",
			NameAssignType:   util.NewValueParser("2").PInt64(),
			NetDevGroup:      util.NewValueParser("0").PInt64(),
			OperState:        "up",
			PhysPortID:       "",
			PhysPortName:     "",
			PhysSwitchID:     "",
			Speed:            util.NewValueParser("1000").PInt64(),
			TxQueueLen:       util.NewValueParser("1000").PInt64(),
			Type:             util.NewValueParser("1").PInt64(),
			BondDeviceAttrs: &NetClassBondDeviceAttrs{
				AdActorOperationalPortState:   util.NewValueParser("61").PUInt64(),
				AdAggregatorId:                util.NewValueParser("1").PUInt64(),
				AdPartnerOperationalPortState: util.NewValueParser("61").PUInt64(),
				LinkFailureCount:              util.NewValueParser("0").PUInt64(),
				MiiStatus:                     util.ParseBool("1"),
				PermamentHWAddress:            makeMAC("01:01:01:01:01:01"),
				QueueID:                       util.NewValueParser("0").PUInt64(),
			},
		},
	}
	bond0 := netClass["bond0"]
	eth0 := netClass["eth0"]
	netClass["bond0"].BondAttrs.Devices = append(netClass["bond0"].BondAttrs.Devices, &eth0)
	queueIDs := make(map[string]uint64)
	queueIDs["eth0"] = 0

	netClass["bond0"].BondAttrs.DeviceQueueIDs = queueIDs
	netClass["eth0"].BondDeviceAttrs.Controller = &bond0

	if !reflect.DeepEqual(netClass, nc) {
		t.Errorf("Result not correct: want %v, have %v", netClass, nc)
	}
}

func makeMAC(s string) *net.HardwareAddr {
	mac, _ := net.ParseMAC(s)
	return &mac
}

func strPTR(s string) *string {
	return &s
}
