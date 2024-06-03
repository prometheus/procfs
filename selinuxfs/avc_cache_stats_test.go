// Copyright 2023 The Prometheus Authors
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

//go:build linux && !noselinux
// +build linux,!noselinux

package selinuxfs

import (
	"testing"
)

func TestAVCStats(t *testing.T) {
	fs, err := NewFS(selinuxTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	avcStats, err := fs.ParseAVCStats()
	if err != nil {
		t.Fatal(err)
	}

	if want, got := uint64(91590784), avcStats.Lookups; want != got {
		t.Errorf("want avcstat lookups %v, got %v", want, got)
	}

	if want, got := uint64(91569452), avcStats.Hits; want != got {
		t.Errorf("want avcstat hits %v, got %v", want, got)
	}

	if want, got := uint64(21332), avcStats.Misses; want != got {
		t.Errorf("want avcstat misses %v, got %v", want, got)
	}

	if want, got := uint64(21332), avcStats.Allocations; want != got {
		t.Errorf("want avcstat allocations %v, got %v", want, got)
	}

	if want, got := uint64(20400), avcStats.Reclaims; want != got {
		t.Errorf("want avcstat reclaims %v, got %v", want, got)
	}

	if want, got := uint64(20826), avcStats.Frees; want != got {
		t.Errorf("want avcstat frees %v, got %v", want, got)
	}
}
