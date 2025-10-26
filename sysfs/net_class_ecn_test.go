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
// +build linux

package sysfs

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEcnByIface(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fs.EcnByIface("non-existent")
	if err == nil {
		t.Fatal("expected error, have none")
	}

	device, err := fs.EcnByIface("eth0")
	if err != nil {
		t.Fatal(err)
	}

	if device.Name != "eth0" {
		t.Errorf("Found unexpected device, want %s, have %s", "eth0", device.Name)
	}
}

func TestEcnDevices(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	ed, _ := fs.EcnDevices()
	allEcnDevices := AllEcnIface{
		"eth0": EcnIface{
			Name: "eth0",
			RoceNpEcn: RoceNpEcn{
				Ecn: map[uint8]bool{
					0: true,
					1: true,
					2: false,
					3: true,
					4: true,
					5: false,
					6: true,
					7: true,
				},
				MinTimeBetweenCnps: 4,
				CnpDscp:            48,
				Cnp802pPriority:    6,
			},
			RoceRpEcn: RoceRpEcn{
				Ecn: map[uint8]bool{
					0: true,
					1: true,
					2: false,
					3: true,
					4: true,
					5: false,
					6: true,
					7: false,
				},
				DceTCPG:                 1019,
				DceTCPRtt:               1,
				InitialAlphaValue:       1023,
				RateToSetOnFirstCnp:     10,
				RpgMinDecFac:            50,
				RpgMinRate:              1,
				RpgGd:                   11,
				RateReduceMonitorPeriod: 4,
				ClampTgtRate:            true,
				RpgTimeReset:            300,
				RpgByteReset:            32767,
				RpgThreshold:            1,
				RpgAiRate:               5,
				RpgHaiRate:              50,
			},
		},
	}

	if diff := cmp.Diff(ed, allEcnDevices); diff != "" {
		t.Fatalf("unexpected diff (-want +got):\n%s", diff)
	}
}
