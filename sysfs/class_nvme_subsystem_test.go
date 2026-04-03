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

func createNVMeSubsystemFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()

	subsysDir := filepath.Join(root, "class", "nvme-subsystem", "nvme-subsys0")

	// Subsystem attributes
	writeFixtureFile(t, subsysDir, "subsysnqn", "nqn.2014-08.org.nvmexpress:uuid:a34c4f3a-0d6f-5cec-dead-beefcafebabe")
	writeFixtureFile(t, subsysDir, "model", "Dell PowerStore")
	writeFixtureFile(t, subsysDir, "serial", "SN12345678")
	writeFixtureFile(t, subsysDir, "iopolicy", "round-robin")

	// Controller nvme0 — live
	ctrl0 := filepath.Join(subsysDir, "nvme0")
	writeFixtureFile(t, ctrl0, "state", "live")
	writeFixtureFile(t, ctrl0, "transport", "fc")
	writeFixtureFile(t, ctrl0, "address", "nn-0x200400a0986b4321:pn-0x210400a0986b4321")

	// Controller nvme1 — live
	ctrl1 := filepath.Join(subsysDir, "nvme1")
	writeFixtureFile(t, ctrl1, "state", "live")
	writeFixtureFile(t, ctrl1, "transport", "fc")
	writeFixtureFile(t, ctrl1, "address", "nn-0x200400a0986b4322:pn-0x210400a0986b4322")

	// Controller nvme2 — dead
	ctrl2 := filepath.Join(subsysDir, "nvme2")
	writeFixtureFile(t, ctrl2, "state", "dead")
	writeFixtureFile(t, ctrl2, "transport", "fc")
	writeFixtureFile(t, ctrl2, "address", "nn-0x200400a0986b4323:pn-0x210400a0986b4323")

	return root
}

func writeFixtureFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestNVMeSubsystemClass(t *testing.T) {
	root := createNVMeSubsystemFixture(t)

	fs, err := NewFS(root)
	if err != nil {
		t.Fatal(err)
	}

	got, err := fs.NVMeSubsystemClass()
	if err != nil {
		t.Fatal(err)
	}

	want := NVMeSubsystemClass{
		{
			Name:     "nvme-subsys0",
			NQN:      "nqn.2014-08.org.nvmexpress:uuid:a34c4f3a-0d6f-5cec-dead-beefcafebabe",
			Model:    "Dell PowerStore",
			Serial:   "SN12345678",
			IOPolicy: "round-robin",
			Controllers: []NVMeSubsystemController{
				{Name: "nvme0", State: "live", Transport: "fc", Address: "nn-0x200400a0986b4321:pn-0x210400a0986b4321"},
				{Name: "nvme1", State: "live", Transport: "fc", Address: "nn-0x200400a0986b4322:pn-0x210400a0986b4322"},
				{Name: "nvme2", State: "dead", Transport: "fc", Address: "nn-0x200400a0986b4323:pn-0x210400a0986b4323"},
			},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected NVMeSubsystemClass (-want +got):\n%s", diff)
	}
}

func TestNVMeSubsystemClassNotPresent(t *testing.T) {
	root := t.TempDir()

	fs, err := NewFS(root)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fs.NVMeSubsystemClass()
	if err == nil {
		t.Fatal("expected error when nvme-subsystem directory does not exist")
	}
}

func TestNVMeSubsystemClassEmpty(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "class", "nvme-subsystem"), 0o755); err != nil {
		t.Fatal(err)
	}

	fs, err := NewFS(root)
	if err != nil {
		t.Fatal(err)
	}

	got, err := fs.NVMeSubsystemClass()
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 0 {
		t.Fatalf("expected 0 subsystems, got %d", len(got))
	}
}

func TestNVMeSubsystemMultipleSubsystems(t *testing.T) {
	root := t.TempDir()

	// Subsystem 0
	subsys0 := filepath.Join(root, "class", "nvme-subsystem", "nvme-subsys0")
	writeFixtureFile(t, subsys0, "subsysnqn", "nqn.target0")
	writeFixtureFile(t, subsys0, "model", "Model0")
	writeFixtureFile(t, subsys0, "serial", "Serial0")
	writeFixtureFile(t, subsys0, "iopolicy", "numa")

	ctrl0 := filepath.Join(subsys0, "nvme0")
	writeFixtureFile(t, ctrl0, "state", "live")
	writeFixtureFile(t, ctrl0, "transport", "tcp")
	writeFixtureFile(t, ctrl0, "address", "traddr=10.0.0.1,trsvcid=4420")

	// Subsystem 1
	subsys1 := filepath.Join(root, "class", "nvme-subsystem", "nvme-subsys1")
	writeFixtureFile(t, subsys1, "subsysnqn", "nqn.target1")
	writeFixtureFile(t, subsys1, "model", "Model1")
	writeFixtureFile(t, subsys1, "serial", "Serial1")
	writeFixtureFile(t, subsys1, "iopolicy", "round-robin")

	ctrl1a := filepath.Join(subsys1, "nvme1")
	writeFixtureFile(t, ctrl1a, "state", "live")
	writeFixtureFile(t, ctrl1a, "transport", "rdma")
	writeFixtureFile(t, ctrl1a, "address", "traddr=10.0.0.2,trsvcid=4420")

	ctrl1b := filepath.Join(subsys1, "nvme2")
	writeFixtureFile(t, ctrl1b, "state", "connecting")
	writeFixtureFile(t, ctrl1b, "transport", "rdma")
	writeFixtureFile(t, ctrl1b, "address", "traddr=10.0.0.3,trsvcid=4420")

	fs, err := NewFS(root)
	if err != nil {
		t.Fatal(err)
	}

	got, err := fs.NVMeSubsystemClass()
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 subsystems, got %d", len(got))
	}

	if got[0].Name != "nvme-subsys0" {
		t.Errorf("expected nvme-subsys0, got %s", got[0].Name)
	}
	if len(got[0].Controllers) != 1 {
		t.Errorf("expected 1 controller for subsys0, got %d", len(got[0].Controllers))
	}

	if got[1].Name != "nvme-subsys1" {
		t.Errorf("expected nvme-subsys1, got %s", got[1].Name)
	}
	if len(got[1].Controllers) != 2 {
		t.Errorf("expected 2 controllers for subsys1, got %d", len(got[1].Controllers))
	}
}
