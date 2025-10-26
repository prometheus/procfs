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

//go:build linux
// +build linux

package sysfs

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestClassDRMCardAMDGPUStats(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	drmTest, err := fs.ClassDRMCardAMDGPUStats()
	if err != nil {
		t.Fatal(err)
	}

	classDRMCardStats := []ClassDRMCardAMDGPUStats{
		{
			Name:                          "card0",
			GPUBusyPercent:                4,
			MemoryGTTSize:                 8573157376,
			MemoryGTTUsed:                 144560128,
			MemoryVisibleVRAMSize:         8573157376,
			MemoryVisibleVRAMUsed:         1490378752,
			MemoryVRAMSize:                8573157376,
			MemoryVRAMUsed:                1490378752,
			MemoryVRAMVendor:              "samsung",
			PowerDPMForcePerformanceLevel: "manual",
			UniqueID:                      "0123456789abcdef",
		},
		{
			Name:                  "card1",
			GPUBusyPercent:        0,
			MemoryGTTSize:         0,
			MemoryGTTUsed:         0,
			MemoryVisibleVRAMSize: 0,
			MemoryVisibleVRAMUsed: 0,
			MemoryVRAMSize:        0,
			MemoryVRAMUsed:        0,
		},
	}

	if diff := cmp.Diff(classDRMCardStats, drmTest); diff != "" {
		t.Fatalf("unexpected diff (-want +got):\n%s", diff)
	}
}
