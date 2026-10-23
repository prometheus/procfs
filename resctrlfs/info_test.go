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
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestL3MonInfo(t *testing.T) {
	fs, err := NewFS(resctrlTestFixtures)
	if err != nil {
		t.Fatalf("failed to open filesystem: %v", err)
	}

	got, err := fs.L3MonInfo()
	if err != nil {
		t.Fatalf("failed to read info/L3_MON: %v", err)
	}

	want := L3MonInfo{
		NumRMIDs: 256,
		MonFeatures: []string{
			"llc_occupancy",
			"mbm_total_bytes",
			"mbm_local_bytes",
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("unexpected L3 monitoring info (-want +got):\n%s", diff)
	}
}

func TestL3MonInfoWithoutMonitoring(t *testing.T) {
	fs, err := NewFS(t.TempDir())
	if err != nil {
		t.Fatalf("failed to open filesystem: %v", err)
	}

	_, err = fs.L3MonInfo()
	if err == nil {
		t.Fatal("want L3MonInfo to fail if info/L3_MON is missing")
	}
	if !os.IsNotExist(err) {
		t.Errorf("want a not exist error, got %v", err)
	}
}
