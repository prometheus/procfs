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

func TestNetClassPCIDevice(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	dev, err := fs.NetClassPCIDevice("enp3s0f0")
	if err != nil {
		t.Fatal(err)
	}

	if want := 0xa2; dev.Location.Bus != want {
		t.Errorf("NetClassPCIDevice() Bus = %#x, want %#x", dev.Location.Bus, want)
	}

	if dev.NumaNode == nil {
		t.Fatal("expected NumaNode to be set")
	}
	if want := int32(1); *dev.NumaNode != want {
		t.Errorf("NumaNode = %d, want %d", *dev.NumaNode, want)
	}
}

func TestNetClassPCIDeviceMissing(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fs.NetClassPCIDevice("nonexistent")
	if err == nil {
		t.Error("expected error for non-existent interface, got nil")
	}
}

func TestPciDeviceVFAddress(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	dev, err := fs.NetClassPCIDevice("enp3s0f0")
	if err != nil {
		t.Fatal(err)
	}

	got, err := fs.PciDeviceVFAddress(dev, 0)
	if err != nil {
		t.Fatal(err)
	}

	if want := "0000:a2:01.0"; got != want {
		t.Errorf("PciDeviceVFAddress() = %q, want %q", got, want)
	}
}

func TestPciDeviceVFAddressMissing(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	dev, err := fs.NetClassPCIDevice("enp3s0f0")
	if err != nil {
		t.Fatal(err)
	}

	_, err = fs.PciDeviceVFAddress(dev, 99)
	if err == nil {
		t.Error("expected error for non-existent VF, got nil")
	}
}
