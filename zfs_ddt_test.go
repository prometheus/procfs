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
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestZFSDDTStats(t *testing.T) {
	stats, err := getProcFixtures(t).ZFSDDTStats()
	if err != nil {
		t.Fatal(err)
	}

	want := []ZFSDDTStats{
		{
			// Every counter is zero in the fixture, so we don't need to check them all.
			Pool:     "pool1",
			Checksum: "blake3",
		},
		{
			Pool:                 "pool1",
			Checksum:             "sha256",
			Lookup:               800,
			LookupNew:            400,
			LookupExisting:       400,
			LookupLiveHit:        0,
			LookupLiveWait:       0,
			LookupLiveMiss:       800,
			LookupLogHit:         148,
			LookupLogActiveHit:   136,
			LookupLogFlushingHit: 12,
			LookupLogMiss:        652,
			LookupStoredHit:      252,
			LookupStoredMiss:     400,
			LogActiveEntries:     0,
			LogFlushingEntries:   0,
			LogIngestRate:        26, // type 2 (UINT32) rows
			LogFlushRate:         37,
			LogFlushTimeRate:     0,
		},
	}

	if len(stats) != len(want) {
		t.Fatalf("want %d ZFSDDTStats, have %d", len(want), len(stats))
	}

	for i := range want {
		if diff := cmp.Diff(want[i], stats[i]); diff != "" {
			t.Errorf("unexpected ZFSDDTStats (-want +got):\n%s", diff)
		}
	}
}

func TestParseZFSDDTStats(t *testing.T) {
	const header = "103 1 0x01 17 4624 2297405876 621083650087\nname                            type data\n"

	tests := []struct {
		name    string
		input   string
		want    ZFSDDTStats
		wantErr bool
	}{
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
		{
			name:    "missing name/type/data header",
			input:   "lookup                          4    800\n",
			wantErr: true,
		},
		{
			name:  "header only",
			input: header,
			want:  ZFSDDTStats{Pool: "pool1", Checksum: "sha256"},
		},
		{
			name:    "non-numeric value",
			input:   header + "lookup                          4    abc\n",
			wantErr: true,
		},
		{
			name:  "unknown field is ignored",
			input: header + "lookup                          4    800\nfuture_field                    4    5\n",
			want:  ZFSDDTStats{Pool: "pool1", Checksum: "sha256", Lookup: 800},
		},
		{
			name:  "unexpected type is skipped",
			input: header + "lookup                          4    800\nlookup_new                      7    400\n",
			want:  ZFSDDTStats{Pool: "pool1", Checksum: "sha256", Lookup: 800},
		},
		{
			name:  "uint32 row is parsed",
			input: header + "log_ingest_rate                 2    26\n",
			want:  ZFSDDTStats{Pool: "pool1", Checksum: "sha256", LogIngestRate: 26},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseZFSDDTStats("pool1", "sha256", []byte(tt.input))
			if tt.wantErr {
				if !errors.Is(err, ErrFileParse) {
					t.Fatalf("expected ErrFileParse, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("unexpected ZFSDDTStats (-want +got):\n%s", diff)
			}
		})
	}
}
