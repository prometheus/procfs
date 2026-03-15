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

package blockdevice

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDMMultipathDevices(t *testing.T) {
	blockdevice, err := NewFS(procfsFixtures, sysfsFixtures)
	if err != nil {
		t.Fatalf("failed to access blockdevice fs: %v", err)
	}

	devices, err := blockdevice.DMMultipathDevices()
	if err != nil {
		t.Fatal(err)
	}

	expected := []DMMultipathDevice{
		{
			Name:      "mpathA",
			SysfsName: "dm-1",
			UUID:      "mpath-360000000000001",
			Suspended: false,
			SizeBytes: 104857600 * 512,
			Paths: []DMMultipathPath{
				{Device: "sdb", State: "running"},
				{Device: "sdc", State: "offline"},
			},
		},
		{
			Name:      "mpathB",
			SysfsName: "dm-2",
			UUID:      "mpath-360000000000002",
			Suspended: true,
			SizeBytes: 209715200 * 512,
			Paths: []DMMultipathPath{
				{Device: "sdd", State: "running"},
				{Device: "sde", State: "running"},
			},
		},
	}

	if diff := cmp.Diff(expected, devices); diff != "" {
		t.Fatalf("unexpected DMMultipathDevices (-want +got):\n%s", diff)
	}
}

func TestDMMultipathDevicesFiltersNonMultipath(t *testing.T) {
	blockdevice, err := NewFS(procfsFixtures, sysfsFixtures)
	if err != nil {
		t.Fatalf("failed to access blockdevice fs: %v", err)
	}

	devices, err := blockdevice.DMMultipathDevices()
	if err != nil {
		t.Fatal(err)
	}

	for _, dev := range devices {
		if dev.SysfsName == "dm-0" {
			t.Error("dm-0 (LVM device) should have been filtered out")
		}
	}
}
