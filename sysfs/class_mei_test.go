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

func TestMEIClass(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	got, err := fs.MEIClass()
	if err != nil {
		t.Fatal(err)
	}

	dev0 := "244:0"
	dev1 := "243:0"
	devState := "ENABLED"
	fwStatus0 := "90000245\n00110500\n00000020\n00000000\n02F41F03\n40000000"
	fwStatus1 := "90000245\n00110500\n00000020\n00000004\n02F61F03\n40000000"
	fwVer0 := "0:18.0.5.2098\n0:18.0.5.2098\n0:18.0.5.2098"
	fwVer1 := "0:18.1.18.2595\n0:18.1.18.2595\n0:18.1.18.2599"
	hbmVer := "2.2"
	hbmVerDrv := "2.2"
	kind := "mei"
	trc0 := "00000889"
	trc1 := "00000879"
	txQueueLimit := "50"

	want := &MEIClass{
		"mei0": MEIDev{
			Dev:           &dev0,
			DevState:      &devState,
			FWStatus:      &fwStatus0,
			FWVersion:     &fwVer0,
			HBMVersion:    &hbmVer,
			HBMVersionDrv: &hbmVerDrv,
			Kind:          &kind,
			Trc:           &trc0,
			TxQueueLimit:  &txQueueLimit,
		},
		"mei1": MEIDev{
			Dev:           &dev1,
			DevState:      &devState,
			FWStatus:      &fwStatus1,
			FWVersion:     &fwVer1,
			HBMVersion:    &hbmVer,
			HBMVersionDrv: &hbmVerDrv,
			Kind:          &kind,
			Trc:           &trc1,
			TxQueueLimit:  &txQueueLimit,
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected MEI class (-want +got):\n%s", diff)
	}
}
