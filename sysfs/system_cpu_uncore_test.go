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

package sysfs

import (
	"os"
	"path/filepath"
	"testing"
)

func uncoreByName(ds []SystemUncoreFrequency, name string) *SystemUncoreFrequency {
	for i := range ds {
		if ds[i].Name == name {
			return &ds[i]
		}
	}
	return nil
}

// TestSystemUncoreFrequency uses the committed fixture tree, which represents a
// Granite Rapids machine exposing BOTH directory layouts at once: the legacy
// package_NN_die_MM directories and the newer per-domain uncoreNN directories.
// The parser returns every domain it finds; de-duplication between the two
// layouts is the collector's responsibility (see the node_exporter uncore
// collector), so here we assert that all four domains are parsed with the
// correct per-file values.
func TestSystemUncoreFrequency(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatal(err)
	}
	got, err := fs.SystemUncoreFrequency()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 4 {
		t.Fatalf("expected 4 uncore domains (2 legacy + 2 new), got %d: %+v", len(got), got)
	}

	// New per-domain directories carry current_freq_khz plus package_id/domain_id.
	for _, tc := range []struct {
		name    string
		pkg     uint64
		domain  uint64
		current uint64
	}{
		{"uncore00", 0, 0, 1600000},
		{"uncore01", 1, 4, 2000000},
	} {
		d := uncoreByName(got, tc.name)
		if d == nil {
			t.Errorf("%s not parsed", tc.name)
			continue
		}
		if d.PackageID == nil || *d.PackageID != tc.pkg {
			t.Errorf("%s PackageID = %v, want %d", tc.name, d.PackageID, tc.pkg)
		}
		if d.DomainID == nil || *d.DomainID != tc.domain {
			t.Errorf("%s DomainID = %v, want %d", tc.name, d.DomainID, tc.domain)
		}
		if d.CurrentFreqKhz == nil || *d.CurrentFreqKhz != tc.current {
			t.Errorf("%s CurrentFreqKhz = %v, want %d", tc.name, d.CurrentFreqKhz, tc.current)
		}
	}

	// Legacy directories have min/max/initial but no current_freq_khz and no
	// package_id/domain_id.
	for _, name := range []string{"package_00_die_00", "package_01_die_00"} {
		d := uncoreByName(got, name)
		if d == nil {
			t.Errorf("%s not parsed", name)
			continue
		}
		if d.CurrentFreqKhz != nil {
			t.Errorf("%s: expected nil CurrentFreqKhz, got %v", name, d.CurrentFreqKhz)
		}
		if d.PackageID != nil || d.DomainID != nil {
			t.Errorf("%s: expected nil PackageID/DomainID, got %v/%v", name, d.PackageID, d.DomainID)
		}
		if d.MinFreqKhz == nil || *d.MinFreqKhz != 800000 {
			t.Errorf("%s MinFreqKhz = %v, want 800000", name, d.MinFreqKhz)
		}
		if d.MaxFreqKhz == nil || *d.MaxFreqKhz != 2400000 {
			t.Errorf("%s MaxFreqKhz = %v, want 2400000", name, d.MaxFreqKhz)
		}
	}
}

// TestSystemUncoreFrequencyAbsent verifies that a sys root without an
// intel_uncore_frequency directory yields no domains and no error.
func TestSystemUncoreFrequencyAbsent(t *testing.T) {
	fs, err := NewFS(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	got, err := fs.SystemUncoreFrequency()
	if err != nil {
		t.Fatalf("expected nil error when directory absent, got %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected no domains, got %d", len(got))
	}
}

// TestSystemUncoreFrequencyReadable is a small standalone check that a fully
// readable file is parsed and a missing file yields nil, using a temporary
// tree so it does not depend on the shared fixtures.
func TestSystemUncoreFrequencyReadable(t *testing.T) {
	sys := t.TempDir()
	dir := filepath.Join(sys, "devices/system/cpu/intel_uncore_frequency/uncore00")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "current_freq_khz"), []byte("1600000\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	fs, err := NewFS(sys)
	if err != nil {
		t.Fatal(err)
	}
	got, err := fs.SystemUncoreFrequency()
	if err != nil {
		t.Fatal(err)
	}
	d := uncoreByName(got, "uncore00")
	if d == nil {
		t.Fatal("uncore00 not parsed")
	}
	if d.CurrentFreqKhz == nil || *d.CurrentFreqKhz != 1600000 {
		t.Errorf("CurrentFreqKhz = %v, want 1600000", d.CurrentFreqKhz)
	}
	if d.MinFreqKhz != nil {
		t.Errorf("expected nil MinFreqKhz (file absent), got %v", d.MinFreqKhz)
	}
}
