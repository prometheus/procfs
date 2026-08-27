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

package resctrlfs

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func uint64Ptr(v uint64) *uint64 {
	return &v
}

func TestMonData(t *testing.T) {
	fs, err := NewFS(resctrlTestFixtures)
	if err != nil {
		t.Fatalf("failed to open filesystem: %v", err)
	}

	got, err := fs.MonData()
	if err != nil {
		t.Fatalf("failed to read mon_data: %v", err)
	}

	want := []MonData{
		{
			Resource:      "L3",
			ID:            "0",
			LLCOccupancy:  uint64Ptr(44040192),
			MBMTotalBytes: uint64Ptr(315294410752),
			MBMLocalBytes: uint64Ptr(210196273664),
		},
		{
			Resource:      "L3",
			ID:            "1",
			LLCOccupancy:  uint64Ptr(8388608),
			MBMTotalBytes: uint64Ptr(105098136832),
			// mbm_local_bytes reads "Unavailable" in the fixture.
			MBMLocalBytes: nil,
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("unexpected mon_data (-want +got):\n%s", diff)
	}
}

func TestMonDataWithoutMonData(t *testing.T) {
	fs, err := NewFS(t.TempDir())
	if err != nil {
		t.Fatalf("failed to open filesystem: %v", err)
	}

	if _, err := fs.MonData(); err == nil {
		t.Error("want MonData to fail if the mon_data directory is missing")
	}
}

func TestParseMonDirName(t *testing.T) {
	for _, test := range []struct {
		name         string
		wantResource string
		wantID       string
		wantOK       bool
	}{
		{name: "mon_L3_00", wantResource: "L3", wantID: "0", wantOK: true},
		{name: "mon_L3_01", wantResource: "L3", wantID: "1", wantOK: true},
		{name: "mon_L3_12", wantResource: "L3", wantID: "12", wantOK: true},
		{name: "mon_MB_0", wantResource: "MB", wantID: "0", wantOK: true},
		{name: "mon_L3_MON_3", wantResource: "L3_MON", wantID: "3", wantOK: true},
		{name: "mon_L3"},
		{name: "mon_L3_"},
		{name: "mon__0"},
		{name: "mon_L3_xy"},
		{name: "L3_00"},
		{name: "mon_"},
		{name: ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			resource, id, ok := parseMonDirName(test.name)
			if ok != test.wantOK {
				t.Fatalf("want ok %v, got %v", test.wantOK, ok)
			}
			if resource != test.wantResource {
				t.Errorf("want resource %q, got %q", test.wantResource, resource)
			}
			if id != test.wantID {
				t.Errorf("want id %q, got %q", test.wantID, id)
			}
		})
	}
}
