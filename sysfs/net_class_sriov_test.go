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
	"testing"
)

func TestNetClassVFPCIAddress(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	got, err := fs.NetClassVFPCIAddress("enp3s0f0", 0)
	if err != nil {
		t.Fatal(err)
	}

	if want := "0000:65:01.0"; got != want {
		t.Errorf("NetClassVFPCIAddress() = %q, want %q", got, want)
	}
}

func TestNetClassVFPCIAddressMissing(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fs.NetClassVFPCIAddress("enp3s0f0", 99)
	if err == nil {
		t.Error("expected error for non-existent VF, got nil")
	}
}
