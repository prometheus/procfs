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

func TestAerCountersByIface(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fs.AerCountersByIface("non-existent")
	if err == nil {
		t.Fatal("expected error, have none")
	}

	device, err := fs.AerCountersByIface("eth0")
	if err != nil {
		t.Fatal(err)
	}

	if device.Name != "eth0" {
		t.Errorf("Found unexpected device, want %s, have %s", "eth0", device.Name)
	}
}

func TestAerCounters(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	ac, _ := fs.AerCounters()
	aerCounters := AllAerCounters{
		"eth0": AerCounters{
			Name: "eth0",
			Correctable: CorrectableAerCounters{
				RxErr:       1,
				BadTLP:      2,
				BadDLLP:     3,
				Rollover:    4,
				Timeout:     5,
				NonFatalErr: 6,
				CorrIntErr:  7,
				HeaderOF:    8,
			},
			Fatal: UncorrectableAerCounters{
				Undefined:        10,
				DLP:              11,
				SDES:             12,
				TLP:              13,
				FCP:              14,
				CmpltTO:          15,
				CmpltAbrt:        16,
				UnxCmplt:         17,
				RxOF:             18,
				MalfTLP:          19,
				ECRC:             20,
				UnsupReq:         21,
				ACSViol:          22,
				UncorrIntErr:     23,
				BlockedTLP:       24,
				AtomicOpBlocked:  25,
				TLPBlockedErr:    26,
				PoisonTLPBlocked: 27,
			},
			NonFatal: UncorrectableAerCounters{
				Undefined:        30,
				DLP:              31,
				SDES:             32,
				TLP:              33,
				FCP:              34,
				CmpltTO:          35,
				CmpltAbrt:        36,
				UnxCmplt:         37,
				RxOF:             38,
				MalfTLP:          39,
				ECRC:             40,
				UnsupReq:         41,
				ACSViol:          42,
				UncorrIntErr:     43,
				BlockedTLP:       44,
				AtomicOpBlocked:  45,
				TLPBlockedErr:    46,
				PoisonTLPBlocked: 47,
			},
		},
	}

	if diff := cmp.Diff(aerCounters, ac); diff != "" {
		t.Fatalf("unexpected diff (-want +got):\n%s", diff)
	}
}
