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

package procfs

import (
	"fmt"
	"testing"
)

func TestNetDevSNMP6(t *testing.T) {
	fs, err := NewFS(procTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	netDevSNMP6, err := fs.NetDevSNMP6()
	if err != nil {
		t.Fatal(err)
	}

	if err := validateNetDevSNMP6(netDevSNMP6); err != nil {
		t.Error(err.Error())
	}
}

func TestProcNetDevSNMP6(t *testing.T) {
	p, err := getProcFixtures(t).Proc(26231)
	if err != nil {
		t.Fatal(err)
	}

	procNetDevSNMP6, err := p.NetDevSNMP6()
	if err != nil {
		t.Fatal(err)
	}

	if err := validateNetDevSNMP6(procNetDevSNMP6); err != nil {
		t.Error(err.Error())
	}
}

func validateNetDevSNMP6(have NetDevSNMP6) error {
	var wantNetDevSNMP6 = map[string]map[string]uint64{
		"eth0": {
			"ifIndex":      1,
			"Ip6InOctets":  14064059261,
			"Ip6OutOctets": 811213622,
			"Icmp6InMsgs":  53293,
			"Icmp6OutMsgs": 20400,
		},
		"eth1": {
			"ifIndex":      2,
			"Ip6InOctets":  303177290674,
			"Ip6OutOctets": 29245052746,
			"Icmp6InMsgs":  37911,
			"Icmp6OutMsgs": 114015,
		},
	}

	for wantIface, wantData := range wantNetDevSNMP6 {
		if haveData, ok := have[wantIface]; ok {
			for wantStat, wantVal := range wantData {
				if haveVal, ok := haveData[wantStat]; !ok {
					return fmt.Errorf("stat %s missing from %s test data", wantStat, wantIface)
				} else if wantVal != haveVal {
					return fmt.Errorf("%s - %s: want %d, have %d", wantIface, wantStat, wantVal, haveVal)
				}
			}
		} else {
			return fmt.Errorf("%s not found in test data", wantIface)
		}
	}

	return nil
}
