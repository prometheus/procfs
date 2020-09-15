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

package procfs

import (
	"reflect"
	"testing"
)

func TestFS_MDStat(t *testing.T) {
	fs := getProcFixtures(t)
	mdStats, err := fs.MDStat()

	if err != nil {
		t.Fatalf("parsing of reference-file failed entirely: %s", err)
	}

	refs := map[string]MDStat{
		"md127": {
			Name:                 "md127",
			ActivityState:        "active",
			Personality:          "raid1",
			DisksActive:          2,
			DisksTotal:           2,
			DisksFailed:          0,
			DisksSpare:           0,
			BlocksTotal:          312319552,
			BlocksSynced:         312319552,
			PercentSynced:        100,
			RemainingSyncMinutes: 0,
			AssignedDevices: []MDAssignedDevice{
				{
					Name:  "sdi2",
					Role:  0,
					State: "active",
				},
				{
					Name:  "sdj2",
					Role:  1,
					State: "active",
				},
			},
		},
		"md0": {
			Name:                 "md0",
			ActivityState:        "active",
			Personality:          "raid1",
			DisksActive:          2,
			DisksTotal:           2,
			DisksFailed:          0,
			DisksSpare:           0,
			BlocksTotal:          248896,
			BlocksSynced:         248896,
			PercentSynced:        100,
			RemainingSyncMinutes: 0,
			AssignedDevices: []MDAssignedDevice{
				{
					Name:  "sdi1",
					Role:  0,
					State: "active",
				},
				{
					Name:  "sdj1",
					Role:  1,
					State: "active",
				},
			},
		},
		"md4": {
			Name:                 "md4",
			ActivityState:        "inactive",
			Personality:          "raid1",
			DisksActive:          0,
			DisksTotal:           0,
			DisksFailed:          1,
			DisksSpare:           1,
			BlocksTotal:          4883648,
			BlocksSynced:         4883648,
			PercentSynced:        100,
			RemainingSyncMinutes: 0,
			AssignedDevices: []MDAssignedDevice{
				{
					Name:  "sda3",
					Role:  0,
					State: "failed",
				},
				{
					Name:  "sdb3",
					Role:  1,
					State: "spare",
				},
			},
		},
		"md6": {
			Name:                 "md6",
			ActivityState:        "recovering",
			Personality:          "raid1",
			DisksActive:          1,
			DisksTotal:           2,
			DisksFailed:          1,
			DisksSpare:           1,
			BlocksTotal:          195310144,
			BlocksSynced:         16775552,
			PercentSynced:        8.5,
			RemainingSyncMinutes: 17.0,
			AssignedDevices: []MDAssignedDevice{
				{
					Name:  "sdb2",
					Role:  2,
					State: "failed",
				},
				{
					Name:  "sdc",
					Role:  1,
					State: "spare",
				},
				{
					Name:  "sda2",
					Role:  0,
					State: "active",
				},
			},
		},
		"md3": {
			Name:                 "md3",
			ActivityState:        "active",
			Personality:          "raid6",
			DisksActive:          8,
			DisksTotal:           8,
			DisksFailed:          0,
			DisksSpare:           2,
			BlocksTotal:          5853468288,
			BlocksSynced:         5853468288,
			PercentSynced:        100,
			RemainingSyncMinutes: 0,
			AssignedDevices: []MDAssignedDevice{
				{
					Name:  "sda1",
					Role:  8,
					State: "active",
				},
				{
					Name:  "sdh1",
					Role:  7,
					State: "active",
				},
				{
					Name:  "sdg1",
					Role:  6,
					State: "active",
				},
				{
					Name:  "sdf1",
					Role:  5,
					State: "active",
				},
				{
					Name:  "sde1",
					Role:  11,
					State: "active",
				},
				{
					Name:  "sdd1",
					Role:  3,
					State: "active",
				},
				{
					Name:  "sdc1",
					Role:  10,
					State: "active",
				},
				{
					Name:  "sdb1",
					Role:  9,
					State: "active",
				},
				{
					Name:  "sdd1",
					Role:  10,
					State: "spare",
				},
				{
					Name:  "sdd2",
					Role:  11,
					State: "spare",
				},
			},
		},
		"md8": {
			Name:                 "md8",
			ActivityState:        "resyncing",
			Personality:          "raid1",
			DisksActive:          2,
			DisksTotal:           2,
			DisksFailed:          0,
			DisksSpare:           2,
			BlocksTotal:          195310144,
			BlocksSynced:         16775552,
			PercentSynced:        8.5,
			RemainingSyncMinutes: 17.0,
			AssignedDevices: []MDAssignedDevice{
				{
					Name:  "sdb1",
					Role:  1,
					State: "active",
				},
				{
					Name:  "sda1",
					Role:  0,
					State: "active",
				},
				{
					Name:  "sdc",
					Role:  2,
					State: "spare",
				},
				{
					Name:  "sde",
					Role:  3,
					State: "spare",
				},
			},
		},
		"md7": {
			Name:                 "md7",
			ActivityState:        "active",
			Personality:          "raid6",
			DisksActive:          3,
			DisksTotal:           4,
			DisksFailed:          1,
			DisksSpare:           0,
			BlocksTotal:          7813735424,
			BlocksSynced:         7813735424,
			PercentSynced:        100,
			RemainingSyncMinutes: 0,
			AssignedDevices: []MDAssignedDevice{
				{
					Name:  "sdb1",
					Role:  0,
					State: "active",
				},
				{
					Name:  "sde1",
					Role:  3,
					State: "active",
				},
				{
					Name:  "sdd1",
					Role:  2,
					State: "active",
				},
				{
					Name:  "sdc1",
					Role:  1,
					State: "failed",
				},
			},
		},
		"md9": {
			Name:                 "md9",
			ActivityState:        "resyncing",
			Personality:          "raid1",
			DisksActive:          4,
			DisksTotal:           4,
			DisksSpare:           1,
			DisksFailed:          2,
			BlocksTotal:          523968,
			BlocksSynced:         0,
			PercentSynced:        0,
			RemainingSyncMinutes: 0,
			AssignedDevices: []MDAssignedDevice{
				{
					Name:  "sdc2",
					Role:  2,
					State: "active",
				},
				{
					Name:  "sdd2",
					Role:  3,
					State: "active",
				},
				{
					Name:  "sdb2",
					Role:  1,
					State: "active",
				},
				{
					Name:  "sda2",
					Role:  0,
					State: "active",
				},
				{
					Name:  "sde",
					Role:  4,
					State: "failed",
				},
				{
					Name:  "sdf",
					Role:  5,
					State: "failed",
				},
				{
					Name:  "sdg",
					Role:  6,
					State: "spare",
				},
			},
		},
		"md10": {
			Name:                 "md10",
			ActivityState:        "active",
			Personality:          "raid0",
			DisksActive:          2,
			DisksTotal:           2,
			DisksFailed:          0,
			DisksSpare:           0,
			BlocksTotal:          314159265,
			BlocksSynced:         314159265,
			PercentSynced:        100,
			RemainingSyncMinutes: 0,
			AssignedDevices: []MDAssignedDevice{
				{
					Name:  "sda1",
					Role:  0,
					State: "active",
				},
				{
					Name:  "sdb1",
					Role:  1,
					State: "active",
				},
			},
		},
		"md11": {
			Name:                 "md11",
			ActivityState:        "resyncing",
			Personality:          "raid1",
			DisksActive:          2,
			DisksTotal:           2,
			DisksFailed:          1,
			DisksSpare:           2,
			BlocksTotal:          4190208,
			BlocksSynced:         0,
			PercentSynced:        0,
			RemainingSyncMinutes: 0,
			AssignedDevices: []MDAssignedDevice{
				{
					Name:  "sdb2",
					Role:  0,
					State: "active",
				},
				{
					Name:  "sdc2",
					Role:  1,
					State: "active",
				},
				{
					Name:  "sdc3",
					Role:  2,
					State: "failed",
				},
				{
					Name:  "hdb",
					Role:  4,
					State: "spare",
				},
				{
					Name:  "ssdc2",
					Role:  3,
					State: "spare",
				},
			},
		},
		"md12": {
			Name:                 "md12",
			ActivityState:        "active",
			Personality:          "raid0",
			DisksActive:          2,
			DisksTotal:           2,
			DisksSpare:           0,
			DisksFailed:          0,
			BlocksTotal:          3886394368,
			BlocksSynced:         3886394368,
			PercentSynced:        100,
			RemainingSyncMinutes: 0,
			AssignedDevices: []MDAssignedDevice{
				{
					Name:  "sdc2",
					Role:  0,
					State: "active",
				},
				{
					Name:  "sdd2",
					Role:  1,
					State: "active",
				},
			},
		},
		"md120": {
			Name:                 "md120",
			ActivityState:        "active",
			Personality:          "linear",
			DisksActive:          2,
			DisksTotal:           2,
			DisksFailed:          0,
			DisksSpare:           0,
			BlocksTotal:          2095104,
			BlocksSynced:         2095104,
			PercentSynced:        100,
			RemainingSyncMinutes: 0,
			AssignedDevices: []MDAssignedDevice{
				{
					Name:  "sda1",
					Role:  1,
					State: "active",
				},
				{
					Name:  "sdb1",
					Role:  0,
					State: "active",
				},
			},
		},
		"md126": {
			Name:                 "md126",
			ActivityState:        "active",
			Personality:          "raid0",
			DisksActive:          2,
			DisksTotal:           2,
			DisksFailed:          0,
			DisksSpare:           0,
			BlocksTotal:          1855870976,
			BlocksSynced:         1855870976,
			PercentSynced:        100,
			RemainingSyncMinutes: 0,
			AssignedDevices: []MDAssignedDevice{
				{
					Name:  "sdb",
					Role:  1,
					State: "active",
				},
				{
					Name:  "sdc",
					Role:  0,
					State: "active",
				},
			},
		},
		"md219": {
			Name:                 "md219",
			ActivityState:        "inactive",
			DisksTotal:           0,
			DisksFailed:          0,
			DisksActive:          0,
			DisksSpare:           3,
			BlocksTotal:          7932,
			BlocksSynced:         7932,
			PercentSynced:        100,
			RemainingSyncMinutes: 0,
			AssignedDevices: []MDAssignedDevice{
				{
					Name:  "sdb",
					Role:  2,
					State: "spare",
				},
				{
					Name:  "sdc",
					Role:  1,
					State: "spare",
				},
				{
					Name:  "sda",
					Role:  0,
					State: "spare",
				},
			},
		},
		"md00": {
			Name:                 "md00",
			ActivityState:        "active",
			Personality:          "raid0",
			DisksActive:          1,
			DisksTotal:           1,
			DisksFailed:          0,
			DisksSpare:           0,
			BlocksTotal:          4186624,
			BlocksSynced:         4186624,
			PercentSynced:        100,
			RemainingSyncMinutes: 0,
			AssignedDevices: []MDAssignedDevice{
				{
					Name:  "xvdb",
					Role:  0,
					State: "active",
				},
			},
		},
		"md101": {
			Name:                 "md101",
			ActivityState:        "active",
			Personality:          "raid0",
			DisksActive:          3,
			DisksTotal:           3,
			DisksFailed:          0,
			DisksSpare:           0,
			BlocksTotal:          322560,
			BlocksSynced:         322560,
			PercentSynced:        100,
			RemainingSyncMinutes: 0,
			AssignedDevices: []MDAssignedDevice{
				{
					Name:  "sdb",
					Role:  2,
					State: "active",
				},
				{
					Name:  "sdd",
					Role:  1,
					State: "active",
				},
				{
					Name:  "sdc",
					Role:  0,
					State: "active",
				},
			},
		},
		"md201": {
			Name:                 "md201",
			ActivityState:        "checking",
			Personality:          "raid1",
			DisksActive:          2,
			DisksTotal:           2,
			DisksFailed:          0,
			DisksSpare:           0,
			BlocksTotal:          1993728,
			BlocksSynced:         114176,
			PercentSynced:        5.7,
			RemainingSyncMinutes: 0.2,
			AssignedDevices: []MDAssignedDevice{
				{
					Name:  "sda3",
					Role:  0,
					State: "active",
				},
				{
					Name:  "sdb3",
					Role:  1,
					State: "active",
				},
			},
		},
	}

	if want, have := len(refs), len(mdStats); want != have {
		t.Errorf("want %d parsed md-devices, have %d", want, have)
	}
	for _, md := range mdStats {
		if want, have := refs[md.Name], md; !reflect.DeepEqual(want, have) {
			t.Errorf("%s: want %v, have %v", md.Name, want, have)
		}
	}

}

func TestInvalidMdstat(t *testing.T) {
	invalidMount := []byte(`
Personalities : [invalid]
md3 : invalid
      314159265 blocks 64k chunks

unused devices: <none>
`)

	_, err := parseMDStat(invalidMount)
	if err == nil {
		t.Fatalf("parsing of invalid reference file did not find any errors")
	}
}
