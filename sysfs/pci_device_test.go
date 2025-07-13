// Copyright 2025 The Prometheus Authors
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
)

func TestPciDevices(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	got, err := fs.PciDevices()
	if err != nil {
		t.Fatal(err)
	}

	var (
		LinkSpeed8GTs = 8.0
		LinkWidth4    = 4.0
		LinkWidth8    = 8.0
	)
	want := PciDevices{
		"0000:00:02:1": PciDevice{
			Location: PciDeviceLocation{
				Segment:  0,
				Bus:      0,
				Device:   2,
				Function: 1,
			},
			ParentLocation: nil,

			Class:           0x060400,
			Vendor:          0x1022,
			Device:          0x1634,
			SubsystemVendor: 0x17aa,
			SubsystemDevice: 0x5095,
			Revision:        0x00,

			MaxLinkSpeed:     &LinkSpeed8GTs,
			MaxLinkWidth:     &LinkWidth8,
			CurrentLinkSpeed: &LinkSpeed8GTs,
			CurrentLinkWidth: &LinkWidth4,
		},
		"0000:01:00:0": PciDevice{
			Location: PciDeviceLocation{
				Segment:  0,
				Bus:      1,
				Device:   0,
				Function: 0,
			},
			ParentLocation: &PciDeviceLocation{
				Segment:  0,
				Bus:      0,
				Device:   2,
				Function: 1,
			},

			Class:           0x010802,
			Vendor:          0xc0a9,
			Device:          0x540a,
			SubsystemVendor: 0xc0a9,
			SubsystemDevice: 0x5021,
			Revision:        0x01,

			MaxLinkSpeed:     &LinkSpeed8GTs,
			MaxLinkWidth:     &LinkWidth4,
			CurrentLinkSpeed: &LinkSpeed8GTs,
			CurrentLinkWidth: &LinkWidth4,
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected PciDevices (-want +got):\n%s", diff)
	}
}

func TestParseDeviceLocation(t *testing.T) {
	got, err := parsePciDeviceLocation("0001:9b:0c.0")
	if err != nil {
		t.Fatal(err)
	}

	want := &PciDeviceLocation{
		Segment:  1,
		Bus:      0x9b,
		Device:   0xc,
		Function: 0,
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected location (-want +got):\n%s", diff)
	}
}
