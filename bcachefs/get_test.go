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

package bcachefs

import "testing"

func TestParseHumanReadableBytes(t *testing.T) {
	tests := []struct {
		in   string
		want uint64
	}{
		{"542k", 555008},
		{"322M", 337641472},
		{"1.5G", 1610612736},
		{"112k", 114688},
		{"405", 405},
	}

	for _, tt := range tests {
		got, err := parseHumanReadableBytes(tt.in)
		if err != nil {
			t.Fatalf("parseHumanReadableBytes(%q) error: %v", tt.in, err)
		}
		if got != tt.want {
			t.Fatalf("parseHumanReadableBytes(%q) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestFSBcachefsStats(t *testing.T) {
	uuid := "deadbeef-1234-5678-9012-abcdefabcdef"
	fs, err := NewFS("testdata/fixtures/sys")
	if err != nil {
		t.Fatalf("NewFS error: %v", err)
	}
	stats, err := fs.Stats()
	if err != nil {
		t.Fatalf("Stats error: %v", err)
	}
	if len(stats) != 1 {
		t.Fatalf("expected 1 filesystem, got %d", len(stats))
	}

	got := stats[0]
	if got.UUID != uuid {
		t.Fatalf("unexpected uuid: %q", got.UUID)
	}
	if got.BtreeCacheSizeBytes != 1610612736 {
		t.Fatalf("unexpected btree cache size: %d", got.BtreeCacheSizeBytes)
	}
	if got.Compression["lz4"].CompressedBytes == 0 {
		t.Fatalf("missing compression stats for lz4")
	}
	if got.Errors["btree_node_read_err"].Count != 5 {
		t.Fatalf("unexpected error count: %d", got.Errors["btree_node_read_err"].Count)
	}
	if got.Counters["btree_node_read"].SinceFilesystemCreation != 67890 {
		t.Fatalf("unexpected counter value: %d", got.Counters["btree_node_read"].SinceFilesystemCreation)
	}
	if got.BtreeWrites["initial"].Count != 19088 {
		t.Fatalf("unexpected btree write count: %d", got.BtreeWrites["initial"].Count)
	}
	dev := got.Devices["1"]
	if dev == nil {
		t.Fatalf("missing device stats")
	}
	if dev.Buckets != 4096 {
		t.Fatalf("unexpected nbuckets: %d", dev.Buckets)
	}
	if dev.IOErrors["read"] != 197346 {
		t.Fatalf("unexpected io_errors read: %d", dev.IOErrors["read"])
	}
	if dev.IODone["read"]["btree"] != 4411097088 {
		t.Fatalf("unexpected io_done read btree: %d", dev.IODone["read"]["btree"])
	}
}
