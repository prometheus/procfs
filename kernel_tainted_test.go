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
)

func TestKernelTainted(t *testing.T) {
	fs, err := NewFS(procfsFixtures)
	if err != nil {
		t.Fatalf("failed to access %s: %v", procfsFixtures, err)
	}

	tainted, err := fs.KernelTainted()
	if err != nil {
		t.Fatalf("failed to collect %s/sys/kernel/tainted: %v", procfsFixtures, err)
	}

	// Fixture value is 12288 = bit 12 (O) + bit 13 (E).
	if tainted.Value != 12288 {
		t.Errorf("tainted value: want %d, got %d", 12288, tainted.Value)
	}

	const wantBits = 20
	if len(tainted.Bits) != wantBits {
		t.Fatalf("tainted bits length: want %d, got %d", wantBits, len(tainted.Bits))
	}

	for _, tc := range []struct {
		index int
		flag  string
		set   bool
	}{
		{0, "P", false},
		{12, "O", true},  // Externally-built (out-of-tree) module
		{13, "E", true},  // Unsigned module
		{14, "L", false}, // Soft lockup — must be clear
		{17, "T", false},
		{18, "N", false},
		{19, "J", false},
	} {
		b := tainted.Bits[tc.index]
		if b.Index != tc.index {
			t.Errorf("bit[%d].Index: want %d, got %d", tc.index, tc.index, b.Index)
		}
		if b.Flag != tc.flag {
			t.Errorf("bit[%d].Flag: want %q, got %q", tc.index, tc.flag, b.Flag)
		}
		if b.Set != tc.set {
			t.Errorf("bit[%d] (%s): want Set=%v, got Set=%v", tc.index, tc.flag, tc.set, b.Set)
		}
	}
}

func TestParseTainted(t *testing.T) {
	// Test parsing directly without filesystem access.
	// 16384 = 2^14 → only bit 14 (L, soft lockup) should be set.
	tainted := parseTainted(16384)

	if tainted.Value != 16384 {
		t.Errorf("value: want 16384, got %d", tainted.Value)
	}
	if !tainted.Bits[14].Set {
		t.Errorf("bit 14 (L): want Set=true, got false")
	}
	for i, b := range tainted.Bits {
		if i != 14 && b.Set {
			t.Errorf("bit %d (%s): want Set=false, got true", i, b.Flag)
		}
	}
}
