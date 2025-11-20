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

	"github.com/google/go-cmp/cmp"
)

func TestNVMeClass(t *testing.T) {
	// Check if test fixtures exist
	if _, err := os.Stat(sysTestFixtures); os.IsNotExist(err) {
		t.Skip("Test fixtures not available, skipping NVMe class tests")
		return
	}

	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	got, err := fs.NVMeClass()
	if err != nil {
		t.Fatal(err)
	}

	want := NVMeClass{
		"nvme0": NVMeDevice{
			Name:             "nvme0",
			FirmwareRevision: "1B2QEXP7",
			Model:            "Samsung SSD 970 PRO 512GB",
			Serial:           "S680HF8N190894I",
			State:            "live",
			ControllerID:     "1997",
			Namespaces: []NVMeNamespace{
				{
					ID:               "0",
					UsedBlocks:       488281250,
					SizeBlocks:       3906250000,
					LogicalBlockSize: 4096,
					ANAState:         "optimized",
					UsedBytes:        2000000000000,
					SizeBytes:        16000000000000,
					CapacityBytes:    16000000000000,
				},
			},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected NVMe class (-want +got):\n%s", diff)
	}
}

func TestNVMeNamespaceParsingWithMockData(t *testing.T) {
	// Create a temporary directory structure for testing namespace parsing
	tempDir := t.TempDir()

	// Create mock NVMe device directory structure
	deviceDir := filepath.Join(tempDir, "class", "nvme", "nvme0")
	err := os.MkdirAll(deviceDir, 0o755)
	if err != nil {
		t.Fatal(err)
	}

	// Create device files
	deviceFiles := map[string]string{
		"firmware_rev": "1B2QEXP7",
		"model":        "Samsung SSD 970 PRO 512GB",
		"serial":       "S680HF8N190894I",
		"state":        "live",
		"cntlid":       "1997",
	}

	for filename, content := range deviceFiles {
		err := os.WriteFile(filepath.Join(deviceDir, filename), []byte(content), 0o644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Create mock namespace directory and files
	namespaceDir := filepath.Join(deviceDir, "nvme0c0n1")
	err = os.MkdirAll(filepath.Join(namespaceDir, "queue"), 0o755)
	if err != nil {
		t.Fatal(err)
	}

	namespaceFiles := map[string]string{
		"nuse":                     "123456",
		"size":                     "1000215216",
		"ana_state":                "optimized",
		"queue/logical_block_size": "512",
	}

	for filename, content := range namespaceFiles {
		filePath := filepath.Join(namespaceDir, filename)
		err := os.WriteFile(filePath, []byte(content), 0o644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Create filesystem and test
	fs, err := NewFS(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	got, err := fs.NVMeClass()
	if err != nil {
		t.Fatal(err)
	}

	// Verify the device was parsed correctly
	if len(got) != 1 {
		t.Fatalf("Expected 1 device, got %d", len(got))
	}

	device, exists := got["nvme0"]
	if !exists {
		t.Fatal("Expected nvme0 device not found")
	}

	// Verify device properties
	if device.Name != "nvme0" {
		t.Errorf("Expected device name nvme0, got %s", device.Name)
	}

	if device.Model != "Samsung SSD 970 PRO 512GB" {
		t.Errorf("Expected model 'Samsung SSD 970 PRO 512GB', got %s", device.Model)
	}

	// Verify namespace was parsed correctly
	if len(device.Namespaces) != 1 {
		t.Fatalf("Expected 1 namespace, got %d", len(device.Namespaces))
	}

	namespace := device.Namespaces[0]
	expectedNamespace := NVMeNamespace{
		ID:               "1",
		UsedBlocks:       123456,
		SizeBlocks:       1000215216,
		LogicalBlockSize: 512,
		ANAState:         "optimized",
		UsedBytes:        123456 * 512,
		SizeBytes:        1000215216 * 512,
		CapacityBytes:    1000215216 * 512,
	}

	if diff := cmp.Diff(expectedNamespace, namespace); diff != "" {
		t.Fatalf("unexpected NVMe namespace (-want +got):\n%s", diff)
	}
}

func TestNVMeMultipleNamespaces(t *testing.T) {
	// Create a temporary directory structure for testing multiple namespaces
	tempDir := t.TempDir()

	// Create mock NVMe device directory structure
	deviceDir := filepath.Join(tempDir, "class", "nvme", "nvme1")
	err := os.MkdirAll(deviceDir, 0o755)
	if err != nil {
		t.Fatal(err)
	}

	// Create device files
	deviceFiles := map[string]string{
		"firmware_rev": "2C3DEXP8",
		"model":        "Test NVMe SSD 1TB",
		"serial":       "TEST123456789",
		"state":        "live",
		"cntlid":       "2048",
	}

	for filename, content := range deviceFiles {
		err := os.WriteFile(filepath.Join(deviceDir, filename), []byte(content), 0o644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Create multiple mock namespace directories
	namespaces := []struct {
		dirName   string
		nsID      string
		nuse      string
		size      string
		anaState  string
		blockSize string
	}{
		{"nvme1c0n1", "1", "100000", "2000000000", "optimized", "4096"},
		{"nvme1c0n2", "2", "50000", "1000000000", "active", "512"},
	}

	for _, ns := range namespaces {
		namespaceDir := filepath.Join(deviceDir, ns.dirName)
		err = os.MkdirAll(filepath.Join(namespaceDir, "queue"), 0o755)
		if err != nil {
			t.Fatal(err)
		}

		namespaceFiles := map[string]string{
			"nuse":                     ns.nuse,
			"size":                     ns.size,
			"ana_state":                ns.anaState,
			"queue/logical_block_size": ns.blockSize,
		}

		for filename, content := range namespaceFiles {
			filePath := filepath.Join(namespaceDir, filename)
			err := os.WriteFile(filePath, []byte(content), 0o644)
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	// Create filesystem and test
	fs, err := NewFS(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	got, err := fs.NVMeClass()
	if err != nil {
		t.Fatal(err)
	}

	// Verify the device was parsed correctly
	if len(got) != 1 {
		t.Fatalf("Expected 1 device, got %d", len(got))
	}

	device, exists := got["nvme1"]
	if !exists {
		t.Fatal("Expected nvme1 device not found")
	}

	// Verify both namespaces were parsed correctly
	if len(device.Namespaces) != 2 {
		t.Fatalf("Expected 2 namespaces, got %d", len(device.Namespaces))
	}

	// Find namespace 1
	var ns1, ns2 *NVMeNamespace
findNamespaces:
	for i := range device.Namespaces {
		switch device.Namespaces[i].ID {
		case "1":
			ns1 = &device.Namespaces[i]
			if ns2 != nil {
				break findNamespaces
			}
		case "2":
			ns2 = &device.Namespaces[i]
			if ns1 != nil {
				break findNamespaces
			}
		}
	}

	if ns1 == nil {
		t.Fatal("Namespace 1 not found")
	}
	if ns2 == nil {
		t.Fatal("Namespace 2 not found")
	}

	// Verify namespace 1
	expectedNS1 := NVMeNamespace{
		ID:               "1",
		UsedBlocks:       100000,
		SizeBlocks:       2000000000,
		LogicalBlockSize: 4096,
		ANAState:         "optimized",
		UsedBytes:        100000 * 4096,
		SizeBytes:        2000000000 * 4096,
		CapacityBytes:    2000000000 * 4096,
	}

	if diff := cmp.Diff(expectedNS1, *ns1); diff != "" {
		t.Errorf("unexpected NVMe namespace 1 (-want +got):\n%s", diff)
	}

	// Verify namespace 2
	expectedNS2 := NVMeNamespace{
		ID:               "2",
		UsedBlocks:       50000,
		SizeBlocks:       1000000000,
		LogicalBlockSize: 512,
		ANAState:         "active",
		UsedBytes:        50000 * 512,
		SizeBytes:        1000000000 * 512,
		CapacityBytes:    1000000000 * 512,
	}

	if diff := cmp.Diff(expectedNS2, *ns2); diff != "" {
		t.Errorf("unexpected NVMe namespace 2 (-want +got):\n%s", diff)
	}
}

func TestNVMeNamespaceMissingFiles(t *testing.T) {
	// Test graceful handling of missing namespace files
	tempDir := t.TempDir()

	// Create mock NVMe device directory structure
	deviceDir := filepath.Join(tempDir, "class", "nvme", "nvme2")
	err := os.MkdirAll(deviceDir, 0o755)
	if err != nil {
		t.Fatal(err)
	}

	// Create device files
	deviceFiles := map[string]string{
		"firmware_rev": "3D4EEXP9",
		"model":        "Incomplete Test SSD",
		"serial":       "INCOMPLETE123",
		"state":        "live",
		"cntlid":       "3072",
	}

	for filename, content := range deviceFiles {
		err := os.WriteFile(filepath.Join(deviceDir, filename), []byte(content), 0o644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Create namespace directory but with missing files
	namespaceDir := filepath.Join(deviceDir, "nvme2c0n1")
	err = os.MkdirAll(filepath.Join(namespaceDir, "queue"), 0o755)
	if err != nil {
		t.Fatal(err)
	}

	// Only create some files, leaving ana_state missing
	namespaceFiles := map[string]string{
		"nuse":                     "75000",
		"size":                     "1500000000",
		"queue/logical_block_size": "4096",
		// ana_state is intentionally missing
	}

	for filename, content := range namespaceFiles {
		filePath := filepath.Join(namespaceDir, filename)
		err := os.WriteFile(filePath, []byte(content), 0o644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Create filesystem and test
	fs, err := NewFS(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	got, err := fs.NVMeClass()
	if err != nil {
		t.Fatal(err)
	}

	// Verify the device was parsed correctly
	if len(got) != 1 {
		t.Fatalf("Expected 1 device, got %d", len(got))
	}

	device, exists := got["nvme2"]
	if !exists {
		t.Fatal("Expected nvme2 device not found")
	}

	// Verify namespace was parsed correctly with default ana_state
	if len(device.Namespaces) != 1 {
		t.Fatalf("Expected 1 namespace, got %d", len(device.Namespaces))
	}

	namespace := device.Namespaces[0]

	// Should have default "unknown" ana_state
	if namespace.ANAState != "unknown" {
		t.Errorf("Expected ana_state 'unknown', got %s", namespace.ANAState)
	}

	// Other values should be parsed correctly
	if namespace.UsedBlocks != 75000 {
		t.Errorf("Expected UsedBlocks 75000, got %d", namespace.UsedBlocks)
	}

	if namespace.SizeBlocks != 1500000000 {
		t.Errorf("Expected SizeBlocks 1500000000, got %d", namespace.SizeBlocks)
	}

	if namespace.LogicalBlockSize != 4096 {
		t.Errorf("Expected LogicalBlockSize 4096, got %d", namespace.LogicalBlockSize)
	}
}
