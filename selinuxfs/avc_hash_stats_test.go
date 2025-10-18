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

//go:build linux && !noselinux
// +build linux,!noselinux

package selinuxfs

import (
	"testing"
)

func TestAVCHashStat(t *testing.T) {
	fs, err := NewFS(selinuxTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	avcHashStats, err := fs.ParseAVCHashStats()
	if err != nil {
		t.Fatal(err)
	}

	if want, got := uint64(503), avcHashStats.Entries; want != got {
		t.Errorf("want avc hash stat entries %v, got %v", want, got)
	}

	if want, got := uint64(512), avcHashStats.BucketsAvailable; want != got {
		t.Errorf("want avc hash stat buckets available %v, got %v", want, got)
	}

	if want, got := uint64(257), avcHashStats.BucketsUsed; want != got {
		t.Errorf("want avc hash stat buckets used %v, got %v", want, got)
	}

	if want, got := uint64(8), avcHashStats.LongestChain; want != got {
		t.Errorf("want avc hash stat longest chain %v, got %v", want, got)
	}
}
