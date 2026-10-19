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

func TestMeminfo(t *testing.T) {
	want := Meminfo{
		MemTotal:          new(uint64(15666184)),
		MemFree:           new(uint64(440324)),
		Buffers:           new(uint64(1020128)),
		Cached:            new(uint64(12007640)),
		SwapCached:        new(uint64(0)),
		Active:            new(uint64(6761276)),
		Inactive:          new(uint64(6532708)),
		ActiveAnon:        new(uint64(267256)),
		InactiveAnon:      new(uint64(268)),
		ActiveFile:        new(uint64(6494020)),
		InactiveFile:      new(uint64(6532440)),
		Unevictable:       new(uint64(0)),
		Mlocked:           new(uint64(0)),
		SwapTotal:         new(uint64(0)),
		SwapFree:          new(uint64(0)),
		Zswap:             new(uint64(22414)),
		Zswapped:          new(uint64(10502)),
		Dirty:             new(uint64(768)),
		Writeback:         new(uint64(0)),
		AnonPages:         new(uint64(266216)),
		Mapped:            new(uint64(44204)),
		Shmem:             new(uint64(1308)),
		Slab:              new(uint64(1807264)),
		SReclaimable:      new(uint64(1738124)),
		SUnreclaim:        new(uint64(69140)),
		KernelStack:       new(uint64(1616)),
		PageTables:        new(uint64(5288)),
		NFSUnstable:       new(uint64(0)),
		Bounce:            new(uint64(0)),
		WritebackTmp:      new(uint64(0)),
		CommitLimit:       new(uint64(7833092)),
		CommittedAS:       new(uint64(530844)),
		VmallocTotal:      new(uint64(34359738367)),
		VmallocUsed:       new(uint64(36596)),
		VmallocChunk:      new(uint64(34359637840)),
		Percpu:            new(uint64(26176)),
		HardwareCorrupted: new(uint64(0)),
		AnonHugePages:     new(uint64(12288)),
		HugePagesTotal:    new(uint64(0)),
		HugePagesFree:     new(uint64(0)),
		HugePagesRsvd:     new(uint64(0)),
		HugePagesSurp:     new(uint64(0)),
		Hugepagesize:      new(uint64(2048)),
		DirectMap4k:       new(uint64(91136)),
		DirectMap2M:       new(uint64(16039936)),

		MemTotalBytes:          new(uint64(16042172416)),
		MemFreeBytes:           new(uint64(450891776)),
		BuffersBytes:           new(uint64(1044611072)),
		CachedBytes:            new(uint64(12295823360)),
		SwapCachedBytes:        new(uint64(0)),
		ActiveBytes:            new(uint64(6923546624)),
		InactiveBytes:          new(uint64(6689492992)),
		ActiveAnonBytes:        new(uint64(273670144)),
		InactiveAnonBytes:      new(uint64(274432)),
		ActiveFileBytes:        new(uint64(6649876480)),
		InactiveFileBytes:      new(uint64(6689218560)),
		UnevictableBytes:       new(uint64(0)),
		MlockedBytes:           new(uint64(0)),
		SwapTotalBytes:         new(uint64(0)),
		SwapFreeBytes:          new(uint64(0)),
		ZswapBytes:             new(uint64(22951936)),
		ZswappedBytes:          new(uint64(10754048)),
		DirtyBytes:             new(uint64(786432)),
		WritebackBytes:         new(uint64(0)),
		AnonPagesBytes:         new(uint64(272605184)),
		MappedBytes:            new(uint64(45264896)),
		ShmemBytes:             new(uint64(1339392)),
		SlabBytes:              new(uint64(1850638336)),
		SReclaimableBytes:      new(uint64(1779838976)),
		SUnreclaimBytes:        new(uint64(70799360)),
		KernelStackBytes:       new(uint64(1654784)),
		PageTablesBytes:        new(uint64(5414912)),
		NFSUnstableBytes:       new(uint64(0)),
		BounceBytes:            new(uint64(0)),
		WritebackTmpBytes:      new(uint64(0)),
		CommitLimitBytes:       new(uint64(8021086208)),
		CommittedASBytes:       new(uint64(543584256)),
		VmallocTotalBytes:      new(uint64(35184372087808)),
		VmallocUsedBytes:       new(uint64(37474304)),
		VmallocChunkBytes:      new(uint64(35184269148160)),
		PercpuBytes:            new(uint64(26804224)),
		HardwareCorruptedBytes: new(uint64(0)),
		AnonHugePagesBytes:     new(uint64(12582912)),
		HugepagesizeBytes:      new(uint64(2097152)),
		DirectMap4kBytes:       new(uint64(93323264)),
		DirectMap2MBytes:       new(uint64(16424894464)),
	}

	got, err := getProcFixtures(t).Meminfo()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected meminfo entry (-want +got):\n%s", diff)
	}
}
