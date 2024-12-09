// Copyright (c) 2024 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.

package sysfs

import (
	"reflect"
	"testing"
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
			Fatal: NonCorrectableAerCounters{
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
			NonFatal: NonCorrectableAerCounters{
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

	if !reflect.DeepEqual(aerCounters, ac) {
		t.Errorf("Result not correct: want %v, have %v", aerCounters, ac)
	}
}
