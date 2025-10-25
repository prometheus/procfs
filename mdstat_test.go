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

package procfs

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFS_MDStat(t *testing.T) {
	fs := getProcFixtures(t)
	mdStats, err := fs.MDStat()

	if err != nil {
		t.Fatalf("parsing of reference-file failed entirely: %s", err)
	}
	// TODO: Test cases to capture in future:
	// WriteMostly devices
	// Journal devices
	// Replacement devices
	// Global hotspares

	refs := map[string]MDStat{
		"md127": {
			Name:                   "md127",
			Type:                   "raid1",
			ActivityState:          "active",
			DisksActive:            2,
			DisksTotal:             2,
			DisksFailed:            0,
			DisksDown:              0,
			DisksSpare:             0,
			BlocksTotal:            312319552,
			BlocksSynced:           312319552,
			BlocksToBeSynced:       312319552,
			BlocksSyncedPct:        0,
			BlocksSyncedFinishTime: 0,
			BlocksSyncedSpeed:      0,
			Devices:                []MDStatComponent{{Name: "sdi2", DescriptorIndex: 0}, {Name: "sdj2", DescriptorIndex: 1}}},
		"md0": {
			Name:                   "md0",
			Type:                   "raid1",
			ActivityState:          "active",
			DisksActive:            2,
			DisksTotal:             2,
			DisksFailed:            0,
			DisksDown:              0,
			DisksSpare:             0,
			BlocksTotal:            248896,
			BlocksSynced:           248896,
			BlocksToBeSynced:       248896,
			BlocksSyncedPct:        0,
			BlocksSyncedFinishTime: 0,
			BlocksSyncedSpeed:      0,
			Devices:                []MDStatComponent{{Name: "sdi1", DescriptorIndex: 0}, {Name: "sdj1", DescriptorIndex: 1}}},
		"md4": {
			Name:                   "md4",
			Type:                   "raid1",
			ActivityState:          "inactive",
			DisksActive:            0,
			DisksTotal:             0,
			DisksFailed:            1,
			DisksDown:              0,
			DisksSpare:             1,
			BlocksTotal:            4883648,
			BlocksSynced:           4883648,
			BlocksToBeSynced:       4883648,
			BlocksSyncedPct:        0,
			BlocksSyncedFinishTime: 0,
			BlocksSyncedSpeed:      0,
			Devices:                []MDStatComponent{{Name: "sda3", Faulty: true, DescriptorIndex: 0}, {Name: "sdb3", Spare: true, DescriptorIndex: 1}}},
		"md6": {
			Name:                   "md6",
			Type:                   "raid1",
			ActivityState:          "recovering",
			DisksActive:            1,
			DisksTotal:             2,
			DisksFailed:            1,
			DisksDown:              1,
			DisksSpare:             1,
			BlocksTotal:            195310144,
			BlocksSynced:           16775552,
			BlocksToBeSynced:       195310144,
			BlocksSyncedPct:        8.5,
			BlocksSyncedFinishTime: 17,
			BlocksSyncedSpeed:      259783,
			Devices:                []MDStatComponent{{Name: "sdb2", DescriptorIndex: 2, Faulty: true}, {Name: "sdc", DescriptorIndex: 1, Spare: true}, {Name: "sda2", DescriptorIndex: 0}}},
		"md3": {
			Name:                   "md3",
			Type:                   "raid6",
			ActivityState:          "active",
			DisksActive:            8,
			DisksTotal:             8,
			DisksFailed:            0,
			DisksDown:              0,
			DisksSpare:             2,
			BlocksTotal:            5853468288,
			BlocksSynced:           5853468288,
			BlocksToBeSynced:       5853468288,
			BlocksSyncedPct:        0,
			BlocksSyncedFinishTime: 0,
			BlocksSyncedSpeed:      0,
			Devices:                []MDStatComponent{{Name: "sda1", DescriptorIndex: 8}, {Name: "sdh1", DescriptorIndex: 7}, {Name: "sdg1", DescriptorIndex: 6}, {Name: "sdf1", DescriptorIndex: 5}, {Name: "sde1", DescriptorIndex: 11}, {Name: "sdd1", DescriptorIndex: 3}, {Name: "sdc1", DescriptorIndex: 10}, {Name: "sdb1", DescriptorIndex: 9}, {Name: "sdd1", DescriptorIndex: 10, Spare: true}, {Name: "sdd2", DescriptorIndex: 11, Spare: true}}},
		"md8": {
			Name:                   "md8",
			Type:                   "raid1",
			ActivityState:          "resyncing",
			DisksActive:            2,
			DisksTotal:             2,
			DisksFailed:            0,
			DisksDown:              0,
			DisksSpare:             2,
			BlocksTotal:            195310144,
			BlocksSynced:           16775552,
			BlocksToBeSynced:       195310144,
			BlocksSyncedPct:        8.5,
			BlocksSyncedFinishTime: 17,
			BlocksSyncedSpeed:      259783,
			Devices:                []MDStatComponent{{Name: "sdb1", DescriptorIndex: 1}, {Name: "sda1", DescriptorIndex: 0}, {Name: "sdc", DescriptorIndex: 2, Spare: true}, {Name: "sde", DescriptorIndex: 3, Spare: true}}},
		"md7": {
			Name:                   "md7",
			Type:                   "raid6",
			ActivityState:          "active",
			DisksActive:            3,
			DisksTotal:             4,
			DisksFailed:            1,
			DisksDown:              1,
			DisksSpare:             0,
			BlocksTotal:            7813735424,
			BlocksSynced:           7813735424,
			BlocksToBeSynced:       7813735424,
			BlocksSyncedPct:        0,
			BlocksSyncedFinishTime: 0,
			BlocksSyncedSpeed:      0,
			Devices:                []MDStatComponent{{Name: "sdb1", DescriptorIndex: 0}, {Name: "sde1", DescriptorIndex: 3}, {Name: "sdd1", DescriptorIndex: 2}, {Name: "sdc1", DescriptorIndex: 1, Faulty: true}}},
		"md9": {
			Name:                   "md9",
			Type:                   "raid1",
			ActivityState:          "resyncing",
			DisksActive:            4,
			DisksTotal:             4,
			DisksSpare:             1,
			DisksDown:              0,
			DisksFailed:            2,
			BlocksTotal:            523968,
			BlocksSynced:           0,
			BlocksToBeSynced:       523968,
			BlocksSyncedPct:        0,
			BlocksSyncedFinishTime: 0,
			BlocksSyncedSpeed:      0,
			Devices:                []MDStatComponent{{Name: "sdc2", DescriptorIndex: 2}, {Name: "sdd2", DescriptorIndex: 3}, {Name: "sdb2", DescriptorIndex: 1}, {Name: "sda2", DescriptorIndex: 0}, {Name: "sde", DescriptorIndex: 4, Faulty: true}, {Name: "sdf", DescriptorIndex: 5, Faulty: true}, {Name: "sdg", DescriptorIndex: 6, Spare: true}}},
		"md10": {
			Name:                   "md10",
			Type:                   "raid0",
			ActivityState:          "active",
			DisksActive:            2,
			DisksTotal:             2,
			DisksFailed:            0,
			DisksDown:              0,
			DisksSpare:             0,
			BlocksTotal:            314159265,
			BlocksSynced:           314159265,
			BlocksToBeSynced:       314159265,
			BlocksSyncedPct:        0,
			BlocksSyncedFinishTime: 0,
			BlocksSyncedSpeed:      0,
			Devices:                []MDStatComponent{{Name: "sda1", DescriptorIndex: 0}, {Name: "sdb1", DescriptorIndex: 1}}},
		"md11": {
			Name:                   "md11",
			Type:                   "raid1",
			ActivityState:          "resyncing",
			DisksActive:            2,
			DisksTotal:             2,
			DisksFailed:            1,
			DisksDown:              0,
			DisksSpare:             2,
			BlocksTotal:            4190208,
			BlocksSynced:           0,
			BlocksToBeSynced:       4190208,
			BlocksSyncedPct:        0,
			BlocksSyncedFinishTime: 0,
			BlocksSyncedSpeed:      0,
			Devices:                []MDStatComponent{{Name: "sdb2", DescriptorIndex: 0}, {Name: "sdc2", DescriptorIndex: 1}, {Name: "sdc3", DescriptorIndex: 2, Faulty: true}, {Name: "hda", DescriptorIndex: 4, Spare: true}, {Name: "ssdc2", DescriptorIndex: 3, Spare: true}}},
		"md12": {
			Name:                   "md12",
			Type:                   "raid0",
			ActivityState:          "active",
			DisksActive:            2,
			DisksTotal:             2,
			DisksSpare:             0,
			DisksDown:              0,
			DisksFailed:            0,
			BlocksTotal:            3886394368,
			BlocksSynced:           3886394368,
			BlocksToBeSynced:       3886394368,
			BlocksSyncedPct:        0,
			BlocksSyncedFinishTime: 0,
			BlocksSyncedSpeed:      0,
			Devices:                []MDStatComponent{{Name: "sdc2", DescriptorIndex: 0}, {Name: "sdd2", DescriptorIndex: 1}}},
		"md120": {
			Name:                   "md120",
			Type:                   "linear",
			ActivityState:          "active",
			DisksActive:            2,
			DisksTotal:             2,
			DisksFailed:            0,
			DisksDown:              0,
			DisksSpare:             0,
			BlocksTotal:            2095104,
			BlocksSynced:           2095104,
			BlocksToBeSynced:       2095104,
			BlocksSyncedPct:        0,
			BlocksSyncedFinishTime: 0,
			BlocksSyncedSpeed:      0,
			Devices:                []MDStatComponent{{Name: "sda1", DescriptorIndex: 1}, {Name: "sdb1", DescriptorIndex: 0}}},
		"md126": {
			Name:                   "md126",
			Type:                   "raid0",
			ActivityState:          "active",
			DisksActive:            2,
			DisksTotal:             2,
			DisksFailed:            0,
			DisksDown:              0,
			DisksSpare:             0,
			BlocksTotal:            1855870976,
			BlocksSynced:           1855870976,
			BlocksToBeSynced:       1855870976,
			BlocksSyncedPct:        0,
			BlocksSyncedFinishTime: 0,
			BlocksSyncedSpeed:      0,
			Devices:                []MDStatComponent{{Name: "sdb", DescriptorIndex: 1}, {Name: "sdc", DescriptorIndex: 0}}},
		"md219": {
			Name:                   "md219",
			Type:                   "unknown",
			ActivityState:          "inactive",
			DisksTotal:             0,
			DisksFailed:            0,
			DisksActive:            0,
			DisksDown:              0,
			DisksSpare:             3,
			BlocksTotal:            7932,
			BlocksSynced:           7932,
			BlocksToBeSynced:       7932,
			BlocksSyncedPct:        0,
			BlocksSyncedFinishTime: 0,
			BlocksSyncedSpeed:      0,
			Devices:                []MDStatComponent{{Name: "sdb", DescriptorIndex: 2, Spare: true}, {Name: "sdc", DescriptorIndex: 1, Spare: true}, {Name: "sda", DescriptorIndex: 0, Spare: true}}},
		"md00": {
			Name:                   "md00",
			Type:                   "raid0",
			ActivityState:          "active",
			DisksActive:            1,
			DisksTotal:             1,
			DisksFailed:            0,
			DisksDown:              0,
			DisksSpare:             0,
			BlocksTotal:            4186624,
			BlocksSynced:           4186624,
			BlocksToBeSynced:       4186624,
			BlocksSyncedPct:        0,
			BlocksSyncedFinishTime: 0,
			BlocksSyncedSpeed:      0,
			Devices:                []MDStatComponent{{Name: "xvdb", DescriptorIndex: 0}}},
		"md101": {
			Name:                   "md101",
			Type:                   "raid0",
			ActivityState:          "active",
			DisksActive:            3,
			DisksTotal:             3,
			DisksFailed:            0,
			DisksDown:              0,
			DisksSpare:             0,
			BlocksTotal:            322560,
			BlocksSynced:           322560,
			BlocksToBeSynced:       322560,
			BlocksSyncedPct:        0,
			BlocksSyncedFinishTime: 0,
			BlocksSyncedSpeed:      0,
			Devices:                []MDStatComponent{{Name: "sdb", DescriptorIndex: 2}, {Name: "sdd", DescriptorIndex: 1}, {Name: "sdc", DescriptorIndex: 0}}},
		"md201": {
			Name:                   "md201",
			Type:                   "raid1",
			ActivityState:          "checking",
			DisksActive:            2,
			DisksTotal:             2,
			DisksFailed:            0,
			DisksDown:              0,
			DisksSpare:             0,
			BlocksTotal:            1993728,
			BlocksSynced:           114176,
			BlocksToBeSynced:       1993728,
			BlocksSyncedPct:        5.7,
			BlocksSyncedFinishTime: 0.2,
			BlocksSyncedSpeed:      114176,
			Devices:                []MDStatComponent{{Name: "sda3", DescriptorIndex: 0}, {Name: "sdb3", DescriptorIndex: 1}}},
		"md42": {
			Name:                   "md42",
			Type:                   "raid5",
			ActivityState:          "reshaping",
			DisksActive:            2,
			DisksTotal:             3,
			DisksFailed:            0,
			DisksDown:              1,
			DisksSpare:             1,
			BlocksTotal:            1953381440,
			BlocksSynced:           1096879076,
			BlocksToBeSynced:       1953381440,
			BlocksSyncedPct:        56.1,
			BlocksSyncedFinishTime: 1868.1,
			BlocksSyncedSpeed:      7640,
			Devices:                []MDStatComponent{{Name: "sda1", DescriptorIndex: 3, Spare: true}, {Name: "sdd1", DescriptorIndex: 0}, {Name: "sde1", DescriptorIndex: 1}}},
	}

	if want, have := len(refs), len(mdStats); want != have {
		t.Errorf("want %d parsed md-devices, have %d", want, have)
	}
	for _, md := range mdStats {
		if want, have := refs[md.Name], md; !cmp.Equal(want, have) {
			t.Errorf("%s: want %v, have %v", md.Name, want, have)
		}
	}

}

func TestInvalidMdstat(t *testing.T) {
	invalidMount := [][]byte{
		// Test invalid Personality and format
		[]byte(`
Personalities : [invalid]
md3 : invalid
      314159265 blocks 64k chunks

unused devices: <none>
`),
		// Test extra blank line
		[]byte(`
md12 : active raid0 sdc2[0] sdd2[1]

      3886394368 blocks super 1.2 512k chunks
`),
		// test for impossible component state
		[]byte(`
md127 : active raid1 sdi2[0] sdj2[1](Z)
      312319552 blocks [2/2] [UU]
`),
		// test for malformed component state
		[]byte(`
md127 : active raid1 sdi2[0] sdj2[X]
      312319552 blocks [2/2] [UU]
`),
	}

	for _, invalid := range invalidMount {
		_, err := parseMDStat(invalid)
		if err == nil {
			t.Fatalf("parsing of invalid reference file did not find any errors")
		}
	}
}
