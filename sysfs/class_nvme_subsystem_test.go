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

func TestNVMeSubsystemClass(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
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
			Path:     filepath.Join(sysTestFixtures, "class/nvme-subsystem/nvme-subsys0"),
			NQN:      "nqn.2014-08.org.nvmexpress:uuid:a34c4f3a-0d6f-5cec-dead-beefcafebabe",
			Model:    "Dell PowerStore",
			Serial:   "SN12345678",
			IOPolicy: "round-robin",
			Controllers: []NVMeSubsystemController{
				{Name: "nvme0", State: "live", Transport: "fc", Address: "nn-0x200400a0986b4321:pn-0x210400a0986b4321"},
				{Name: "nvme1", State: "live", Transport: "fc", Address: "nn-0x200400a0986b4322:pn-0x210400a0986b4322"},
				{Name: "nvme2", State: "dead", Transport: "fc", Address: "nn-0x200400a0986b4323:pn-0x210400a0986b4323"},
			},
			Namespaces: []string{"nvme0n1"},
		},
		{
			Name:     "nvme-subsys1",
			Path:     filepath.Join(sysTestFixtures, "class/nvme-subsystem/nvme-subsys1"),
			NQN:      "nqn.2014-08.org.nvmexpress:uuid:b45d5e4b-1e7f-6ded-beef-deadcafe1234",
			Model:    "NetApp ONTAP",
			Serial:   "NTAP98765",
			IOPolicy: "numa",
			Controllers: []NVMeSubsystemController{
				{Name: "nvme3", State: "live", Transport: "tcp", Address: "traddr=10.0.0.1,trsvcid=4420"},
				{Name: "nvme4", State: "connecting", Transport: "rdma", Address: "traddr=10.0.0.2,trsvcid=4420"},
			},
			Namespaces: []string{"nvme3n1", "nvme4n1"},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected NVMeSubsystemClass (-want +got):\n%s", diff)
	}

	// Multipath head diag counters, read on demand.
	for _, tt := range []struct {
		subsys    int
		namespace string
		want      NVMeSubsystemNamespaceDiag
	}{
		{
			subsys:    0,
			namespace: "nvme0n1",
			want: NVMeSubsystemNamespaceDiag{
				IORequeueNoUsablePath: new(uint64(9)),
				IOFailNoAvailablePath: new(uint64(1)),
			},
		},
		{
			// No diag directory, as on kernels older than 7.2.
			subsys:    1,
			namespace: "nvme3n1",
			want:      NVMeSubsystemNamespaceDiag{},
		},
	} {
		diag, err := got[tt.subsys].NamespaceDiag(tt.namespace)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(tt.want, diag); diff != "" {
			t.Errorf("unexpected %s diag (-want +got):\n%s", tt.namespace, diff)
		}
	}

	if _, err := got[0].NamespaceDiag("../nvme-subsys1/nvme3n1"); err == nil {
		t.Error("expected an error for an invalid namespace name")
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
