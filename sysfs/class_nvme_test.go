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
			Namespaces:       []NVMeNamespace{}, // Empty for now since test fixtures don't include namespaces
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
	err := os.MkdirAll(deviceDir, 0755)
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
		err := os.WriteFile(filepath.Join(deviceDir, filename), []byte(content), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Create mock namespace directory and files
	namespaceDir := filepath.Join(deviceDir, "nvme0c0n1")
	err = os.MkdirAll(filepath.Join(namespaceDir, "queue"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	namespaceFiles := map[string]string{
		"nuse":                      "123456",
		"size":                      "1000215216",
		"ana_state":                 "optimized",
		"queue/logical_block_size":  "512",
	}

	for filename, content := range namespaceFiles {
		filePath := filepath.Join(namespaceDir, filename)
		err := os.WriteFile(filePath, []byte(content), 0644)
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
		Cap