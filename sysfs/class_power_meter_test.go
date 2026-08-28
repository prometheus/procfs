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
	"sort"
	"testing"
)

func TestGetPowerMeters(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatalf("failed to open filesystem: %v", err)
	}

	meters, err := GetPowerMeters(fs)
	if err != nil {
		t.Fatalf("failed to get power meters: %v", err)
	}
	if len(meters) != 2 {
		t.Fatalf("expected 2 power meters, got %d", len(meters))
	}

	// Sort by name for deterministic ordering (glob order varies by OS).
	sort.Slice(meters, func(i, j int) bool { return meters[i].Name < meters[j].Name })

	// --- ACPI000D:00: full-field meter ---
	m0 := meters[0]
	if m0.Name != "ACPI000D:00" {
		t.Fatalf("expected Name=ACPI000D:00, got %q", m0.Name)
	}

	assertPInt64(t, "ACPI000D:00.Average", m0.Average, 15000000)
	assertPInt64(t, "ACPI000D:00.AverageMin", m0.AverageMin, 0)
	assertPInt64(t, "ACPI000D:00.AverageMax", m0.AverageMax, 60000000)
	assertPInt64(t, "ACPI000D:00.AverageInterval", m0.AverageInterval, 1000)
	assertPInt64(t, "ACPI000D:00.AverageIntervalMin", m0.AverageIntervalMin, 100)
	assertPInt64(t, "ACPI000D:00.AverageIntervalMax", m0.AverageIntervalMax, 10000)
	assertPInt64(t, "ACPI000D:00.Alarm", m0.Alarm, 0)
	assertPInt64(t, "ACPI000D:00.Cap", m0.Cap, 25000000)
	assertPInt64(t, "ACPI000D:00.CapMin", m0.CapMin, 1000000)
	assertPInt64(t, "ACPI000D:00.CapMax", m0.CapMax, 100000000)
	assertPInt64(t, "ACPI000D:00.CapHyst", m0.CapHyst, 500000)
	assertPInt64(t, "ACPI000D:00.IsBattery", m0.IsBattery, 0)

	if m0.Accuracy != "1.50%" {
		t.Errorf("expected Accuracy=%q, got %q", "1.50%", m0.Accuracy)
	}
	if m0.ModelNumber != "ACME PM01" {
		t.Errorf("expected ModelNumber=%q, got %q", "ACME PM01", m0.ModelNumber)
	}
	if m0.SerialNumber != "SN12345" {
		t.Errorf("expected SerialNumber=%q, got %q", "SN12345", m0.SerialNumber)
	}
	if m0.OEMInfo != "ACME Corp" {
		t.Errorf("expected OEMInfo=%q, got %q", "ACME Corp", m0.OEMInfo)
	}

	sort.Strings(m0.Measures)
	if len(m0.Measures) != 2 || m0.Measures[0] != "LNXCPU:00" || m0.Measures[1] != "LNXMEM:00" {
		t.Errorf("expected Measures=[LNXCPU:00, LNXMEM:00], got %v", m0.Measures)
	}

	// --- ACPI000D:01: minimal-field meter (verify optional fields are nil) ---
	m1 := meters[1]
	if m1.Name != "ACPI000D:01" {
		t.Fatalf("expected Name=ACPI000D:01, got %q", m1.Name)
	}

	assertPInt64(t, "ACPI000D:01.Average", m1.Average, 5000000)
	assertPInt64(t, "ACPI000D:01.AverageInterval", m1.AverageInterval, 500)
	assertPInt64(t, "ACPI000D:01.Alarm", m1.Alarm, 0)

	assertNilPInt64(t, "ACPI000D:01.AverageMin", m1.AverageMin)
	assertNilPInt64(t, "ACPI000D:01.AverageMax", m1.AverageMax)
	assertNilPInt64(t, "ACPI000D:01.AverageIntervalMin", m1.AverageIntervalMin)
	assertNilPInt64(t, "ACPI000D:01.AverageIntervalMax", m1.AverageIntervalMax)
	assertNilPInt64(t, "ACPI000D:01.Cap", m1.Cap)
	assertNilPInt64(t, "ACPI000D:01.CapMin", m1.CapMin)
	assertNilPInt64(t, "ACPI000D:01.CapMax", m1.CapMax)
	assertNilPInt64(t, "ACPI000D:01.CapHyst", m1.CapHyst)
	assertNilPInt64(t, "ACPI000D:01.IsBattery", m1.IsBattery)

	if m1.Accuracy != "" {
		t.Errorf("expected Accuracy=\"\" for ACPI000D:01, got %q", m1.Accuracy)
	}
	if m1.ModelNumber != "" {
		t.Errorf("expected ModelNumber=\"\", got %q", m1.ModelNumber)
	}
	if m1.SerialNumber != "" {
		t.Errorf("expected SerialNumber=\"\", got %q", m1.SerialNumber)
	}
	if m1.OEMInfo != "" {
		t.Errorf("expected OEMInfo=\"\", got %q", m1.OEMInfo)
	}
	if len(m1.Measures) != 0 {
		t.Errorf("expected empty Measures for ACPI000D:01, got %v", m1.Measures)
	}
}

func TestGetPowerMeters_NoDevices(t *testing.T) {
	fs, err := NewFS(t.TempDir())
	if err != nil {
		t.Fatalf("failed to open filesystem: %v", err)
	}

	meters, err := GetPowerMeters(fs)
	if !os.IsNotExist(err) {
		t.Fatalf("expected os.ErrNotExist when no devices exist, got %v", err)
	}
	if len(meters) != 0 {
		t.Fatalf("expected empty meters, got %v", meters)
	}
}

func TestGetPowerMeters_SingleMeter(t *testing.T) {
	// Verify parsing works when only a subset of attributes is present.
	tmp := t.TempDir()
	meterDir := filepath.Join(tmp, "bus", "acpi", "drivers", "power_meter", "ACPI000D:00")
	if err := os.MkdirAll(filepath.Join(meterDir, "measures"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(meterDir, "power1_average"), []byte("1000"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(meterDir, "power1_alarm"), []byte("1"), 0o644); err != nil {
		t.Fatal(err)
	}

	fs, err := NewFS(tmp)
	if err != nil {
		t.Fatal(err)
	}

	meters, err := GetPowerMeters(fs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(meters) != 1 {
		t.Fatalf("expected 1 meter, got %d", len(meters))
	}
	if meters[0].Name != "ACPI000D:00" {
		t.Errorf("expected Name=ACPI000D:00, got %q", meters[0].Name)
	}
	assertPInt64(t, "Average", meters[0].Average, 1000)
	assertPInt64(t, "Alarm", meters[0].Alarm, 1)
	assertNilPInt64(t, "Cap", meters[0].Cap)
}

// assertPInt64 checks that a *int64 field is non-nil and holds the expected value.
func assertPInt64(t *testing.T, field string, got *int64, want int64) {
	t.Helper()
	if got == nil {
		t.Errorf("%s: expected %d, got nil", field, want)
		return
	}
	if *got != want {
		t.Errorf("%s: expected %d, got %d", field, want, *got)
	}
}

// assertNilPInt64 checks that a *int64 field is nil.
func assertNilPInt64(t *testing.T, field string, got *int64) {
	t.Helper()
	if got != nil {
		t.Errorf("%s: expected nil, got %d", field, *got)
	}
}
