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

package procfs

import (
	"reflect"
	"testing"
)

func TestMeminfo(t *testing.T) {
	expected := Meminfo{
		MemTotal:          newuint64(15666184),
		MemFree:           newuint64(440324),
		Buffers:           newuint64(1020128),
		Cached:            newuint64(12007640),
		SwapCached:        newuint64(0),
		Active:            newuint64(6761276),
		Inactive:          newuint64(6532708),
		ActiveAnon:        newuint64(267256),
		InactiveAnon:      newuint64(268),
		ActiveFile:        newuint64(6494020),
		InactiveFile:      newuint64(6532440),
		Unevictable:       newuint64(0),
		Mlocked:           newuint64(0),
		SwapTotal:         newuint64(0),
		SwapFree:          newuint64(0),
		Dirty:             newuint64(768),
		Writeback:         newuint64(0),
		AnonPages:         newuint64(266216),
		Mapped:            newuint64(44204),
		Shmem:             newuint64(1308),
		Slab:              newuint64(1807264),
		SReclaimable:      newuint64(1738124),
		SUnreclaim:        newuint64(69140),
		KernelStack:       newuint64(1616),
		PageTables:        newuint64(5288),
		NFSUnstable:       newuint64(0),
		Bounce:            newuint64(0),
		WritebackTmp:      newuint64(0),
		CommitLimit:       newuint64(7833092),
		CommittedAS:       newuint64(530844),
		VmallocTotal:      newuint64(34359738367),
		VmallocUsed:       newuint64(36596),
		VmallocChunk:      newuint64(34359637840),
		HardwareCorrupted: newuint64(0),
		AnonHugePages:     newuint64(12288),
		HugePagesTotal:    newuint64(0),
		HugePagesFree:     newuint64(0),
		HugePagesRsvd:     newuint64(0),
		HugePagesSurp:     newuint64(0),
		Hugepagesize:      newuint64(2048),
		DirectMap4k:       newuint64(91136),
		DirectMap2M:       newuint64(16039936),
	}

	have, err := getProcFixtures(t).Meminfo()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(have, expected) {
		t.Logf("have: %+v", have)
		t.Logf("expected: %+v", expected)
		t.Errorf("structs are not equal")
	}
}
