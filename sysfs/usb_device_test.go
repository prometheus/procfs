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

	"github.com/google/go-cmp/cmp"
)

func TestUSBDevices(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	got, err := fs.USBDevices()
	if err != nil {
		t.Fatal(err)
	}

	boolPtr := func(v bool) *bool { return &v }
	removablePtr := func(v USBRemovable) *USBRemovable { return &v }
	int64Ptr := func(v int64) *int64 { return &v }

	want := []USBDevice{
		{
			Name: "1-2", BusNum: "1", DevNum: "14",
			IDVendor: "0658", IDProduct: "0200",
			DeviceClass: "02", DeviceSubClass: "00", DeviceProtocol: "00", BCDDevice: "0000",
			BCDUSB: "2.00", Speed: 12, NumConfigurations: 1, NumInterfaces: 2, BMaxPower: 100,
			Authorized: boolPtr(true), Removable: removablePtr(USBRemovableRemovable), URBNum: 9,
			PowerControl: "on", RuntimeStatus: "active", AutosuspendDelayMillis: int64Ptr(2000),
			Interfaces: []USBInterface{
				{Configuration: "1", Interface: "0", Class: "02", SubClass: "02", Protocol: "01", Driver: "cdc_acm"},
				{Configuration: "1", Interface: "1", Class: "0a", SubClass: "00", Protocol: "00", Driver: "cdc_acm"},
			},
		},
		{
			Name: "1-3", BusNum: "1", DevNum: "2",
			IDVendor: "058f", IDProduct: "9540",
			Manufacturer: "Generic", Product: "EMV Smartcard Reader",
			DeviceClass: "00", DeviceSubClass: "00", DeviceProtocol: "00", BCDDevice: "0120",
			BCDUSB: "2.01", Speed: 12, NumConfigurations: 1, NumInterfaces: 1, BMaxPower: 50,
			Authorized: boolPtr(true), Removable: removablePtr(USBRemovableFixed), URBNum: 2977,
			PowerControl: "auto", RuntimeStatus: "suspended", AutosuspendDelayMillis: int64Ptr(2000),
			Interfaces: []USBInterface{
				{Configuration: "1", Interface: "0", Class: "0b", SubClass: "00", Protocol: "00"},
			},
		},
		{
			// 1-4 only has an idVendor attribute: the optional attributes must
			// stay nil instead of being reported as zero values.
			Name: "1-4", IDVendor: "1234",
		},
		{
			Name: "4-2", BusNum: "4", DevNum: "2",
			IDVendor: "1050", IDProduct: "0407",
			Manufacturer: "Yubico", Product: "Yubikey 4 OTP+U2F+CCID",
			DeviceClass: "00", DeviceSubClass: "00", DeviceProtocol: "00", BCDDevice: "0426",
			BCDUSB: "2.00", Speed: 12, NumConfigurations: 1, NumInterfaces: 3, BMaxPower: 30,
			Authorized: boolPtr(true), Removable: removablePtr(USBRemovableRemovable), URBNum: 13724,
			PowerControl: "auto", RuntimeStatus: "active", AutosuspendDelayMillis: int64Ptr(2000),
			Interfaces: []USBInterface{
				{Configuration: "1", Interface: "0", Class: "03", SubClass: "01", Protocol: "01", Driver: "usbhid"},
				{Configuration: "1", Interface: "1", Class: "03", SubClass: "00", Protocol: "00", Driver: "usbhid"},
				{Configuration: "1", Interface: "2", Class: "0b", SubClass: "00", Protocol: "00", Driver: "usbfs"},
			},
		},
		{
			Name: "4-3", BusNum: "4", DevNum: "3",
			IDVendor: "06cb", IDProduct: "00bd",
			Serial:      "7ba268d87336",
			DeviceClass: "ff", DeviceSubClass: "10", DeviceProtocol: "ff", BCDDevice: "0000",
			BCDUSB: "2.00", Speed: 12, NumConfigurations: 1, NumInterfaces: 1, BMaxPower: 100,
			Authorized: boolPtr(true), Removable: removablePtr(USBRemovableFixed), URBNum: 35,
			PowerControl: "auto", RuntimeStatus: "suspended", AutosuspendDelayMillis: int64Ptr(2000),
			Interfaces: []USBInterface{
				{Configuration: "1", Interface: "0", Class: "ff", SubClass: "00", Protocol: "00"},
			},
		},
		{
			Name: "4-4", BusNum: "4", DevNum: "4",
			IDVendor: "8087", IDProduct: "0029",
			DeviceClass: "e0", DeviceSubClass: "01", DeviceProtocol: "01", BCDDevice: "0001",
			BCDUSB: "2.01", Speed: 12, NumConfigurations: 1, NumInterfaces: 2, BMaxPower: 100,
			Authorized: boolPtr(true), Removable: removablePtr(USBRemovableFixed), URBNum: 33505,
			PowerControl: "auto", RuntimeStatus: "suspended", AutosuspendDelayMillis: int64Ptr(2000),
			Interfaces: []USBInterface{
				{Configuration: "1", Interface: "0", Class: "e0", SubClass: "01", Protocol: "01", Driver: "btusb"},
				{Configuration: "1", Interface: "1", Class: "e0", SubClass: "01", Protocol: "01", Driver: "btusb"},
			},
		},
		{
			Name: "6-2", BusNum: "6", DevNum: "2",
			IDVendor: "04f2", IDProduct: "b6cb",
			Manufacturer: "Chicony Electronics Co.,Ltd.", Product: "Integrated Camera", Serial: "0001",
			DeviceClass: "ef", DeviceSubClass: "02", DeviceProtocol: "01", BCDDevice: "5818",
			BCDUSB: "2.01", Speed: 480, NumConfigurations: 1, NumInterfaces: 5, BMaxPower: 500,
			Authorized: boolPtr(true), Removable: removablePtr(USBRemovableFixed), URBNum: 804,
			PowerControl: "auto", RuntimeStatus: "suspended", AutosuspendDelayMillis: int64Ptr(2000),
			Interfaces: []USBInterface{
				{Configuration: "1", Interface: "0", Class: "0e", SubClass: "01", Protocol: "00", Driver: "uvcvideo"},
				{Configuration: "1", Interface: "1", Class: "0e", SubClass: "02", Protocol: "00", Driver: "uvcvideo"},
				{Configuration: "1", Interface: "2", Class: "0e", SubClass: "01", Protocol: "01", Driver: "uvcvideo"},
				{Configuration: "1", Interface: "3", Class: "0e", SubClass: "02", Protocol: "01", Driver: "uvcvideo"},
				{Configuration: "1", Interface: "4", Class: "fe", SubClass: "01", Protocol: "01"},
			},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected USBDevices (-want +got):\n%s", diff)
	}
}

func TestUSBControllers(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	got, err := fs.USBControllers()
	if err != nil {
		t.Fatal(err)
	}

	want := []USBController{
		{BusNum: "1", BCDUSB: "2.00", Speed: 480, MaxChild: 4},
		{BusNum: "2", BCDUSB: "3.10", Speed: 10000, MaxChild: 2},
		{BusNum: "3", BCDUSB: "2.00", Speed: 480, MaxChild: 1},
		{BusNum: "4", BCDUSB: "2.00", Speed: 480, MaxChild: 4},
		{BusNum: "5", BCDUSB: "3.10", Speed: 10000, MaxChild: 2},
		{BusNum: "6", BCDUSB: "2.00", Speed: 480, MaxChild: 2},
		{BusNum: "7", BCDUSB: "3.00", Speed: 5000, MaxChild: 2},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected USBControllers (-want +got):\n%s", diff)
	}
}
