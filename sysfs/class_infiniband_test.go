// Copyright 2019 The Prometheus Authors
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

// +build linux

package sysfs

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseSlowRate(t *testing.T) {
	tests := []struct {
		rate string
		want uint64
	}{
		{
			rate: "0 GB/sec",
			want: 0,
		},
		{
			rate: "2.5 Gb/sec (1X SDR)",
			want: 312500000,
		},
		{
			rate: "500 Gb/sec (4X HDR)",
			want: 62500000000,
		},
	}

	for _, tt := range tests {
		rate, err := parseRate(tt.rate)
		if err != nil {
			t.Fatal(err)
		}
		if rate != tt.want {
			t.Errorf("Result for InfiniBand rate not correct: want %v, have %v", tt.want, rate)
		}
	}
}

func TestInfiniBandClass(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	got, err := fs.InfiniBandClass()
	if err != nil {
		t.Fatal(err)
	}

	var (
		hfi1Port1ExcessiveBufferOverrunErrors uint64
		hfi1Port1LinkDowned                   uint64
		hfi1Port1LinkErrorRecovery            uint64
		hfi1Port1LocalLinkIntegrityErrors     uint64
		hfi1Port1PortRcvConstraintErrors      uint64
		hfi1Port1PortRcvData                  uint64 = 1380366808104
		hfi1Port1PortRcvErrors                uint64
		hfi1Port1PortRcvPackets               uint64 = 638036947
		hfi1Port1PortRcvRemotePhysicalErrors  uint64
		hfi1Port1PortRcvSwitchRelayErrors     uint64
		hfi1Port1PortXmitConstraintErrors     uint64
		hfi1Port1PortXmitData                 uint64 = 1094233306172
		hfi1Port1PortXmitDiscards             uint64
		hfi1Port1PortXmitPackets              uint64 = 568318856
		hfi1Port1PortXmitWait                 uint64
		hfi1Port1SymbolError                  uint64
		hfi1Port1VL15Dropped                  uint64

		mlx4Port1ExcessiveBufferOverrunErrors uint64
		mlx4Port1LinkDowned                   uint64
		mlx4Port1LinkErrorRecovery            uint64
		mlx4Port1LocalLinkIntegrityErrors     uint64
		mlx4Port1PortRcvConstraintErrors      uint64
		mlx4Port1PortRcvData                  uint64 = 8884894436
		mlx4Port1PortRcvErrors                uint64
		mlx4Port1PortRcvPackets               uint64 = 87169372
		mlx4Port1PortRcvRemotePhysicalErrors  uint64
		mlx4Port1PortRcvSwitchRelayErrors     uint64
		mlx4Port1PortXmitConstraintErrors     uint64
		mlx4Port1PortXmitData                 uint64 = 106036453180
		mlx4Port1PortXmitDiscards             uint64
		mlx4Port1PortXmitPackets              uint64 = 85734114
		mlx4Port1PortXmitWait                 uint64 = 3599
		mlx4Port1SymbolError                  uint64
		mlx4Port1VL15Dropped                  uint64

		mlx4Port2ExcessiveBufferOverrunErrors uint64
		mlx4Port2LinkDowned                   uint64
		mlx4Port2LinkErrorRecovery            uint64
		mlx4Port2LocalLinkIntegrityErrors     uint64
		mlx4Port2PortRcvConstraintErrors      uint64
		mlx4Port2PortRcvData                  uint64 = 9841747136
		mlx4Port2PortRcvErrors                uint64
		mlx4Port2PortRcvPackets               uint64 = 89332064
		mlx4Port2PortRcvRemotePhysicalErrors  uint64
		mlx4Port2PortRcvSwitchRelayErrors     uint64
		mlx4Port2PortXmitConstraintErrors     uint64
		mlx4Port2PortXmitData                 uint64 = 106161427560
		mlx4Port2PortXmitDiscards             uint64
		mlx4Port2PortXmitPackets              uint64 = 88622850
		mlx4Port2PortXmitWait                 uint64 = 3846
		mlx4Port2SymbolError                  uint64
		mlx4Port2VL15Dropped                  uint64
	)

	want := InfiniBandClass{
		"hfi1_0": InfiniBandDevice{
			Name:            "hfi1_0",
			BoardID:         "HPE 100Gb 1-port OP101 QSFP28 x16 PCIe Gen3 with Intel Omni-Path Adapter",
			FirmwareVersion: "1.27.0",
			HCAType:         "",
			Ports: map[uint]InfiniBandPort{
				1: {
					Name:        "hfi1_0",
					Port:        1,
					State:       "ACTIVE",
					StateID:     4,
					PhysState:   "LinkUp",
					PhysStateID: 5,
					Rate:        12500000000,
					Counters: InfiniBandCounters{
						ExcessiveBufferOverrunErrors: &hfi1Port1ExcessiveBufferOverrunErrors,
						LinkDowned:                   &hfi1Port1LinkDowned,
						LinkErrorRecovery:            &hfi1Port1LinkErrorRecovery,
						LocalLinkIntegrityErrors:     &hfi1Port1LocalLinkIntegrityErrors,
						PortRcvConstraintErrors:      &hfi1Port1PortRcvConstraintErrors,
						PortRcvData:                  &hfi1Port1PortRcvData,
						PortRcvErrors:                &hfi1Port1PortRcvErrors,
						PortRcvPackets:               &hfi1Port1PortRcvPackets,
						PortRcvRemotePhysicalErrors:  &hfi1Port1PortRcvRemotePhysicalErrors,
						PortRcvSwitchRelayErrors:     &hfi1Port1PortRcvSwitchRelayErrors,
						PortXmitConstraintErrors:     &hfi1Port1PortXmitConstraintErrors,
						PortXmitData:                 &hfi1Port1PortXmitData,
						PortXmitDiscards:             &hfi1Port1PortXmitDiscards,
						PortXmitPackets:              &hfi1Port1PortXmitPackets,
						PortXmitWait:                 &hfi1Port1PortXmitWait,
						SymbolError:                  &hfi1Port1SymbolError,
						VL15Dropped:                  &hfi1Port1VL15Dropped,
					},
				},
			},
		},
		"mlx4_0": InfiniBandDevice{
			Name:            "mlx4_0",
			BoardID:         "SM_1141000001000",
			FirmwareVersion: "2.31.5050",
			HCAType:         "MT4099",
			Ports: map[uint]InfiniBandPort{
				1: {
					Name:        "mlx4_0",
					Port:        1,
					State:       "ACTIVE",
					StateID:     4,
					PhysState:   "LinkUp",
					PhysStateID: 5,
					Rate:        5000000000,
					Counters: InfiniBandCounters{
						ExcessiveBufferOverrunErrors: &mlx4Port1ExcessiveBufferOverrunErrors,
						LinkDowned:                   &mlx4Port1LinkDowned,
						LinkErrorRecovery:            &mlx4Port1LinkErrorRecovery,
						LocalLinkIntegrityErrors:     &mlx4Port1LocalLinkIntegrityErrors,
						PortRcvConstraintErrors:      &mlx4Port1PortRcvConstraintErrors,
						PortRcvData:                  &mlx4Port1PortRcvData,
						PortRcvErrors:                &mlx4Port1PortRcvErrors,
						PortRcvPackets:               &mlx4Port1PortRcvPackets,
						PortRcvRemotePhysicalErrors:  &mlx4Port1PortRcvRemotePhysicalErrors,
						PortRcvSwitchRelayErrors:     &mlx4Port1PortRcvSwitchRelayErrors,
						PortXmitConstraintErrors:     &mlx4Port1PortXmitConstraintErrors,
						PortXmitData:                 &mlx4Port1PortXmitData,
						PortXmitDiscards:             &mlx4Port1PortXmitDiscards,
						PortXmitPackets:              &mlx4Port1PortXmitPackets,
						PortXmitWait:                 &mlx4Port1PortXmitWait,
						SymbolError:                  &mlx4Port1SymbolError,
						VL15Dropped:                  &mlx4Port1VL15Dropped,
					},
				},
				2: {
					Name:        "mlx4_0",
					Port:        2,
					State:       "ACTIVE",
					StateID:     4,
					PhysState:   "LinkUp",
					PhysStateID: 5,
					Rate:        5000000000,
					Counters: InfiniBandCounters{
						ExcessiveBufferOverrunErrors: &mlx4Port2ExcessiveBufferOverrunErrors,
						LinkDowned:                   &mlx4Port2LinkDowned,
						LinkErrorRecovery:            &mlx4Port2LinkErrorRecovery,
						LocalLinkIntegrityErrors:     &mlx4Port2LocalLinkIntegrityErrors,
						PortRcvConstraintErrors:      &mlx4Port2PortRcvConstraintErrors,
						PortRcvData:                  &mlx4Port2PortRcvData,
						PortRcvErrors:                &mlx4Port2PortRcvErrors,
						PortRcvPackets:               &mlx4Port2PortRcvPackets,
						PortRcvRemotePhysicalErrors:  &mlx4Port2PortRcvRemotePhysicalErrors,
						PortRcvSwitchRelayErrors:     &mlx4Port2PortRcvSwitchRelayErrors,
						PortXmitConstraintErrors:     &mlx4Port2PortXmitConstraintErrors,
						PortXmitData:                 &mlx4Port2PortXmitData,
						PortXmitDiscards:             &mlx4Port2PortXmitDiscards,
						PortXmitPackets:              &mlx4Port2PortXmitPackets,
						PortXmitWait:                 &mlx4Port2PortXmitWait,
						SymbolError:                  &mlx4Port2SymbolError,
						VL15Dropped:                  &mlx4Port2VL15Dropped,
					},
				},
			},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected InfiniBand class (-want +got):\n%s", diff)
	}
}
