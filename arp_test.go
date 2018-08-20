// Copyright 2017 The Prometheus Authors
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
	"testing"
)

func TestARP(t *testing.T) {
	arpFile, err := FS("fixtures/arp/valid").GatherARPEntries()
	if err != nil {
		t.Fatal(err)
	}

	if want, got := "192.168.224.254", arpFile[0].IPAddr.String(); want != got {
		t.Errorf("want 192.168.224.254, got %s", got)
	}
}
