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

	dev := "244:0"
	devState := "ENABLED"
	fwStatus := "90000245\n00110500\n00000020\n00000000\n02F41F03\n40000000"
	fwVer := "0:18.0.5.2098\n0:18.0.5.2098\n0:18.0.5.2098"
	hbmVer := "2.2"
	hbmVerDrv := "2.2"
	kind := "mei"
	trc := "00000889"
	txQueueLimit := "50"

	want := &MEIClass{
		Dev:           &dev,
		DevState:      &devState,
		FWStatus:      &fwStatus,
		FWVersion:     &fwVer,
		HBMVersion:    &hbmVer,
		HBMVersionDrv: &hbmVerDrv,
		Kind:          &kind,
		Trc:           &trc,
		TxQueueLimit:  &txQueueLimit,
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected MEI class (-want +got):\n%s", diff)
	}
}
