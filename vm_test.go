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

//go:build !windows

package procfs

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestVM(t *testing.T) {
	fs, err := NewFS(procTestFixtures)
	if err != nil {
		t.Fatal(err)
	}
	got, err := fs.VM()
	if err != nil {
		t.Fatal(err)
	}
	zeroPointer := new(int64(0))
	lowmemreserveratio := []*int64{new(int64(256)), new(int64(256)), new(int64(32)), zeroPointer, zeroPointer}
	want := &VM{
		AdminReserveKbytes:        new(int64(8192)),
		BlockDump:                 zeroPointer,
		CompactUnevictableAllowed: new(int64(1)),
		DirtyBackgroundBytes:      zeroPointer,
		DirtyBackgroundRatio:      new(int64(10)),
		DirtyBytes:                zeroPointer,
		DirtyExpireCentisecs:      new(int64(3000)),
		DirtyRatio:                new(int64(20)),
		DirtytimeExpireSeconds:    new(int64(43200)),
		DirtyWritebackCentisecs:   new(int64(500)),
		DropCaches:                zeroPointer,
		ExtfragThreshold:          new(int64(500)),
		HugetlbShmGroup:           zeroPointer,
		LaptopMode:                new(int64(5)),
		LegacyVaLayout:            zeroPointer,
		LowmemReserveRatio:        lowmemreserveratio,
		MaxMapCount:               new(int64(65530)),
		MemoryFailureEarlyKill:    zeroPointer,
		MemoryFailureRecovery:     new(int64(1)),
		MinFreeKbytes:             new(int64(67584)),
		MinSlabRatio:              new(int64(5)),
		MinUnmappedRatio:          new(int64(1)),
		MmapMinAddr:               new(int64(65536)),
		NumaStat:                  new(int64(1)),
		NumaZonelistOrder:         "Node",
		NrHugepages:               zeroPointer,
		NrHugepagesMempolicy:      zeroPointer,
		NrOvercommitHugepages:     zeroPointer,
		OomDumpTasks:              new(int64(1)),
		OomKillAllocatingTask:     zeroPointer,
		OvercommitKbytes:          zeroPointer,
		OvercommitMemory:          zeroPointer,
		OvercommitRatio:           new(int64(50)),
		PageCluster:               new(int64(3)),
		PanicOnOom:                zeroPointer,
		PercpuPagelistFraction:    zeroPointer,
		StatInterval:              new(int64(1)),
		Swappiness:                new(int64(60)),
		UserReserveKbytes:         new(int64(131072)),
		VfsCachePressure:          new(int64(100)),
		WatermarkBoostFactor:      new(int64(15000)),
		WatermarkScaleFactor:      new(int64(10)),
		ZoneReclaimMode:           zeroPointer,
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected power supply class (-want +got):\n%s", diff)
	}
}
