// Copyright 2019 The Prometheus Authors
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

package procfs

import (
	"net"
	"testing"
)

func TestARP(t *testing.T) {
	fs, err := NewFS(procTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	arpFile, err := fs.GatherARPEntries()
	if err != nil {
		t.Fatal(err)
	}

	if want, got := "192.168.224.1", arpFile[0].IPAddr.String(); want != got {
		t.Errorf("want 192.168.224.1, got %s", got)
	}

	hwAddr, err := net.ParseMAC("00:50:56:c0:00:08")
	if err != nil {
		t.Error(err)
	}

	if want, got := arpFile[0].HWAddr.String(), hwAddr.String(); want != got {
		t.Errorf("want 00:50:56:c0:00:08, got %s", got)
	}

	if want, got := "ens33", arpFile[0].Device; want != got {
		t.Errorf("want ens33, got %s", got)
	}

	if want, got := "10.50.0.1", arpFile[1].IPAddr.String(); want != got {
		t.Errorf("want 10.50.0.1, got %s", got)
	}

	hwAddr, err = net.ParseMAC("00:00:00:00:00:00")
	if err != nil {
		t.Error(err)
	}

	if want, got := arpFile[1].HWAddr.String(), hwAddr.String(); want != got {
		t.Errorf("want 00:00:00:00:00:00, got %s", got)
	}

	if want, got := "eth0", arpFile[1].Device; want != got {
		t.Errorf("want eth0, got %s", got)
	}
}
