// Copyright 2020 The Prometheus Authors
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
)

func TestLnstat(t *testing.T) {
	const (
		filesCount             = 2
		CPUsCount              = 2
		arpCacheMetricsCount   = 13
		ndiscCacheMetricsCount = 13
	)

	fs, err := NewFS(procTestFixtures)
	if err != nil {
		t.Fatalf("failed to open procfs: %v", err)
	}

	lnstats, err := fs.Lnstat()
	if err != nil {
		t.Fatalf("Lnstat() error: %s", err)
	}

	if len(lnstats) != filesCount {
		t.Fatalf("unexpected number of files parsed %d, expected %d", len(lnstats), filesCount)
	}

	expectedStats := [2]Lnstats{
		{
			Filename: "arp_cache",
			Stats:    make(map[string][]uint64),
		},
		{
			Filename: "ndisc_cache",
			Stats:    make(map[string][]uint64),
		},
	}

	for _, expected := range expectedStats {
		if expected.Filename == "arp_cache" {
			expected.Stats["entries"] = []uint64{20, 20}
			expected.Stats["allocs"] = []uint64{1, 13}
			expected.Stats["destroys"] = []uint64{2, 14}
			expected.Stats["hash_grows"] = []uint64{3, 15}
			expected.Stats["lookups"] = []uint64{4, 16}
			expected.Stats["hits"] = []uint64{5, 17}
			expected.Stats["res_failed"] = []uint64{6, 18}
			expected.Stats["rcv_probes_mcast"] = []uint64{7, 19}
			expected.Stats["rcv_probes_ucast"] = []uint64{8, 20}
			expected.Stats["periodic_gc_runs"] = []uint64{9, 21}
			expected.Stats["forced_gc_runs"] = []uint64{10, 22}
			expected.Stats["unresolved_discards"] = []uint64{11, 23}
			expected.Stats["table_fulls"] = []uint64{12, 24}
		}
		if expected.Filename == "ndisc_cache" {
			expected.Stats["entries"] = []uint64{36, 36}
			expected.Stats["allocs"] = []uint64{240, 252}
			expected.Stats["destroys"] = []uint64{241, 253}
			expected.Stats["hash_grows"] = []uint64{242, 254}
			expected.Stats["lookups"] = []uint64{243, 255}
			expected.Stats["hits"] = []uint64{244, 256}
			expected.Stats["res_failed"] = []uint64{245, 257}
			expected.Stats["rcv_probes_mcast"] = []uint64{246, 258}
			expected.Stats["rcv_probes_ucast"] = []uint64{247, 259}
			expected.Stats["periodic_gc_runs"] = []uint64{248, 260}
			expected.Stats["forced_gc_runs"] = []uint64{249, 261}
			expected.Stats["unresolved_discards"] = []uint64{250, 262}
			expected.Stats["table_fulls"] = []uint64{251, 263}
		}
	}

	for _, lnstatFile := range lnstats {
		if lnstatFile.Filename == "arp_cache" && len(lnstatFile.Stats) != arpCacheMetricsCount {
			t.Fatalf("unexpected arp_cache metrics count %d, expected %d", len(lnstatFile.Stats), arpCacheMetricsCount)
		}
		if lnstatFile.Filename == "ndisc_cache" && len(lnstatFile.Stats) != ndiscCacheMetricsCount {
			t.Fatalf("unexpected ndisc_cache metrics count %d, expected %d", len(lnstatFile.Stats), ndiscCacheMetricsCount)
		}
		for _, expected := range expectedStats {
			for header, stats := range lnstatFile.Stats {
				if header == "" {
					t.Fatalf("Found empty metric name")
				}
				if len(stats) != CPUsCount {
					t.Fatalf("Lnstat parsed %d lines with metrics, expected %d", len(stats), CPUsCount)
				}
				if lnstatFile.Filename == expected.Filename {
					if expected.Stats[header] == nil {
						t.Fatalf("unexpected metric header: %s", header)
					}
					for cpu, value := range lnstatFile.Stats[header] {
						if expected.Stats[header][cpu] != value {
							t.Fatalf("unexpected value for %s for cpu %d in %s: %d, expected %d", header, cpu, lnstatFile.Filename, value, expected.Stats[header][cpu])
						}
					}
				}
			}
		}
	}
}
