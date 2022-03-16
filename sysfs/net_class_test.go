// Copyright The Prometheus Authors
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

package sysfs

import (
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/prometheus/procfs/internal/parsers"
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

	found := false
	for _, d := range devices {
		if d == "eth0" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected device eth0 not found in %v", devices)
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
		"enp3s0f0": {Name: "enp3s0f0"},
		"bond0": {
			Address:          "02:02:02:02:02:02",
			AddrAssignType:   parsers.NewValueParser("3").PInt64(),
			AddrLen:          parsers.NewValueParser("6").PInt64(),
			Broadcast:        "ff:ff:ff:ff:ff:ff",
			Carrier:          parsers.NewValueParser("1").PInt64(),
			CarrierChanges:   parsers.NewValueParser("2").PInt64(),
			CarrierDownCount: parsers.NewValueParser("1").PInt64(),
			CarrierUpCount:   parsers.NewValueParser("1").PInt64(),
			DevID:            parsers.NewValueParser("32").PInt64(),
			Dormant:          parsers.NewValueParser("1").PInt64(),
			Duplex:           "full",
			Flags:            parsers.NewValueParser("4867").PInt64(),
			IfAlias:          "",
			IfIndex:          parsers.NewValueParser("2").PInt64(),
			IfLink:           parsers.NewValueParser("2").PInt64(),
			LinkMode:         parsers.NewValueParser("1").PInt64(),
			MTU:              parsers.NewValueParser("1500").PInt64(),
			Name:             "bond0",
			NameAssignType:   parsers.NewValueParser("2").PInt64(),
			NetDevGroup:      parsers.NewValueParser("0").PInt64(),
			OperState:        "up",
			PhysPortID:       "",
			PhysPortName:     "",
			PhysSwitchID:     "",
			Speed:            parsers.NewValueParser("1000").PInt64(),
			TxQueueLen:       parsers.NewValueParser("1000").PInt64(),
			Type:             parsers.NewValueParser("1").PInt64(),
			BondAttrs: &NetClassBondAttrs{
				AdActorKey:                             parsers.NewValueParser("15").PUInt64(),
				AdActorSysPriority:                     parsers.NewValueParser("65535").PUInt64(),
				AdActorSystem:                          makeMAC("00:00:00:00:00:00"),
				AdAggregator:                           parsers.NewValueParser("1").PUInt64(),
				AdNumPorts:                             parsers.NewValueParser("2").PUInt64(),
				AdPartnerKey:                           parsers.NewValueParser("1034").PUInt64(),
				AdPartnerMac:                           makeMAC("01:23:45:67:89:AB"),
				AdSelect:                               strPTR("stable"),
				AdSelectID:                             parsers.NewValueParser("0").PUInt64(),
				AdUserPortKey:                          parsers.NewValueParser("0").PUInt64(),
				AllDevicesActive:                       parsers.ParseBool("0"),
				ARPAllTargets:                          strPTR("any"),
				ARPAllTargetsID:                        parsers.NewValueParser("0").PUInt64(),
				ARPInterval:                            parsers.NewValueParser("0").PInt64(),
				ARPValidate:                            strPTR("none"),
				ARPValidateID:                          parsers.NewValueParser("0").PUInt64(),
				DownDelay:                              parsers.NewValueParser("200").PInt64(),
				FailoverMac:                            strPTR("none"),
				FailoverMacID:                          parsers.NewValueParser("0").PUInt64(),
				LACPRate:                               strPTR("slow"),
				LACPRateID:                             parsers.NewValueParser("0").PUInt64(),
				LPInterval:                             parsers.NewValueParser("1").PInt64(),
				MIIMon:                                 parsers.NewValueParser("100").PInt64(),
				MIIStatus:                              parsers.ParseBool("1"),
				MinLinks:                               parsers.NewValueParser("0").PUInt64(),
				Mode:                                   strPTR("802.3ad"),
				ModeID:                                 parsers.NewValueParser("4").PUInt64(),
				NumberGratuitousArp:                    parsers.NewValueParser("1").PUInt64(),
				NumberUnsolicitedNeighborAdvertisement: parsers.NewValueParser("1").PUInt64(),
				PacketsPerDevice:                       parsers.NewValueParser("1").PInt64(),
				PrimaryReselect:                        strPTR("always"),
				PrimaryReselectID:                      parsers.NewValueParser("0").PUInt64(),
				ResendIgmp:                             parsers.NewValueParser("1").PInt64(),
				TLBDynamicLB:                           parsers.NewValueParser("1").PInt64(),
				UpDelay:                                parsers.NewValueParser("0").PInt64(),
				UseCarrier:                             parsers.ParseBool("1"),
				TransmitHashPolicy:                     strPTR("layer3+4"),
				TransmitHashPolicyID:                   parsers.NewValueParser("1").PUInt64(),
			},
		},
		"eth0": {
			Address:          "01:01:01:01:01:01",
			AddrAssignType:   parsers.NewValueParser("3").PInt64(),
			AddrLen:          parsers.NewValueParser("6").PInt64(),
			Broadcast:        "ff:ff:ff:ff:ff:ff",
			Carrier:          parsers.NewValueParser("1").PInt64(),
			CarrierChanges:   parsers.NewValueParser("2").PInt64(),
			CarrierDownCount: parsers.NewValueParser("1").PInt64(),
			CarrierUpCount:   parsers.NewValueParser("1").PInt64(),
			DevID:            parsers.NewValueParser("32").PInt64(),
			Dormant:          parsers.NewValueParser("1").PInt64(),
			Duplex:           "full",
			Flags:            parsers.NewValueParser("4867").PInt64(),
			IfAlias:          "",
			IfIndex:          parsers.NewValueParser("2").PInt64(),
			IfLink:           parsers.NewValueParser("2").PInt64(),
			LinkMode:         parsers.NewValueParser("1").PInt64(),
			MTU:              parsers.NewValueParser("1500").PInt64(),
			Name:             "eth0",
			NameAssignType:   parsers.NewValueParser("2").PInt64(),
			NetDevGroup:      parsers.NewValueParser("0").PInt64(),
			OperState:        "up",
			PhysPortID:       "",
			PhysPortName:     "",
			PhysSwitchID:     "",
			Speed:            parsers.NewValueParser("1000").PInt64(),
			TxQueueLen:       parsers.NewValueParser("1000").PInt64(),
			Type:             parsers.NewValueParser("1").PInt64(),
			BondDeviceAttrs: &NetClassBondDeviceAttrs{
				AdActorOperationalPortState:   parsers.NewValueParser("61").PUInt64(),
				AdAggregatorID:                parsers.NewValueParser("1").PUInt64(),
				AdPartnerOperationalPortState: parsers.NewValueParser("61").PUInt64(),
				LinkFailureCount:              parsers.NewValueParser("0").PUInt64(),
				MiiStatus:                     parsers.ParseBool("1"),
				PermamentHWAddress:            makeMAC("01:01:01:01:01:01"),
				QueueID:                       parsers.NewValueParser("0").PUInt64(),
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

	if diff := cmp.Diff(netClass, nc); diff != "" {
		t.Fatalf("unexpected diff (-want +got):\n%s", diff)
	}
}

func makeMAC(s string) *net.HardwareAddr {
	mac, _ := net.ParseMAC(s)
	return &mac
}

func strPTR(s string) *string {
	return &s
}
