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

//go:build linux
// +build linux

package sysfs

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"k8s.io/utils/ptr"
)

func TestFibreChannelClass(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	got, err := fs.FibreChannelClass()
	if err != nil {
		t.Fatal(err)
	}

	want := FibreChannelClass{
		"host0": FibreChannelHost{
			Name:             ptr.To("host0"),
			Speed:            ptr.To("16 Gbit"),
			PortState:        ptr.To("Online"),
			PortType:         ptr.To("Point-To-Point (direct nport connection)"),
			PortName:         ptr.To("1000e0071bce95f2"),
			SymbolicName:     ptr.To("Emulex SN1100E2P FV12.4.270.3 DV12.4.0.0. HN:gotest. OS:Linux"),
			NodeName:         ptr.To("2000e0071bce95f2"),
			PortID:           ptr.To("000002"),
			FabricName:       ptr.To("0"),
			DevLossTMO:       ptr.To("30"),
			SupportedClasses: ptr.To("Class 3"),
			SupportedSpeeds:  ptr.To("4 Gbit, 8 Gbit, 16 Gbit"),
			Counters: &FibreChannelCounters{
				DumpedFrames:          ptr.To(^uint64(0)),
				ErrorFrames:           ptr.To(uint64(0)),
				InvalidCRCCount:       ptr.To(uint64(0x2)),
				RXFrames:              ptr.To(uint64(0x3)),
				RXWords:               ptr.To(uint64(0x4)),
				TXFrames:              ptr.To(uint64(0x5)),
				TXWords:               ptr.To(uint64(0x6)),
				SecondsSinceLastReset: ptr.To(uint64(0x7)),
				InvalidTXWordCount:    ptr.To(uint64(0x8)),
				LinkFailureCount:      ptr.To(uint64(0x9)),
				LossOfSyncCount:       ptr.To(uint64(0x10)),
				LossOfSignalCount:     ptr.To(uint64(0x11)),
				NosCount:              ptr.To(uint64(0x12)),
				FCPPacketAborts:       ptr.To(uint64(0x13)),
			},
		},
		"host1": FibreChannelHost{
			Name:      ptr.To("host1"),
			PortState: ptr.To("Online"),
			Counters: &FibreChannelCounters{
				DumpedFrames:          ptr.To(uint64(0)),
				ErrorFrames:           ptr.To(^uint64(0)),
				InvalidCRCCount:       ptr.To(uint64(0x20)),
				RXFrames:              ptr.To(uint64(0x30)),
				RXWords:               ptr.To(uint64(0x40)),
				TXFrames:              ptr.To(uint64(0x50)),
				TXWords:               ptr.To(uint64(0x60)),
				SecondsSinceLastReset: ptr.To(uint64(0x70)),
				InvalidTXWordCount:    ptr.To(uint64(0x80)),
				LinkFailureCount:      ptr.To(uint64(0x90)),
				LossOfSyncCount:       ptr.To(uint64(0x100)),
				LossOfSignalCount:     ptr.To(uint64(0x110)),
				NosCount:              ptr.To(uint64(0x120)),
				FCPPacketAborts:       ptr.To(uint64(0x130)),
			},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected FibreChannel class (-want +got):\n%s", diff)
	}
}
