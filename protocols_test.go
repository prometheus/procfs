// Copyright 2020 The Prometheus Authors
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

func TestProtocolsParseLine(t *testing.T) {
	rawStr := "TCP       1984  93064  1225378   no     320   yes  kernel      y  y  y  y  y  y  y  y  y  y  y  y  y  n  y  y  y  y  y\n"
	protocols := ProtocolStats{}
	have, err := protocols.parseLine(rawStr)
	if err != nil {
		t.Fatal(err)
	}

	want := ProtocolStatLine{"TCP", 1984, 93064, 1225378, false, 320}
	if want != *have {
		t.Errorf("want %+v have %+v\n", want, have)
	}
}

func TestProtocolsParseProtocols(t *testing.T) {
	fs, err := NewFS(procTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	protocolStats, err := fs.Protocols()
	if err != nil {
		t.Fatal(err)
	}

	lines := map[string]ProtocolStatLine{
		"PACKET":    {"PACKET", 1344, 2, -1, false, 0},
		"PINGv6":    {"PINGv6", 1112, 0, -1, false, 0},
		"RAWv6":     {"RAWv6", 1112, 1, -1, false, 0},
		"UDPLITEv6": {"UDPLITEv6", 1216, 0, 57, false, 0},
		"UDPv6":     {"UDPv6", 1216, 10, 57, false, 0},
		"TCPv6":     {"TCPv6", 2144, 1937, 1225378, false, 320},
		"UNIX":      {"UNIX", 1024, 120, -1, false, 0},
		"UDP-Lite":  {"UDP-Lite", 1024, 0, 57, false, 0},
		"PING":      {"PING", 904, 0, -1, false, 0},
		"RAW":       {"RAW", 912, 0, -1, false, 0},
		"UDP":       {"UDP", 1024, 73, 57, false, 0},
		"TCP":       {"TCP", 1984, 93064, 1225378, true, 320},
		"NETLINK":   {"NETLINK", 1040, 16, -1, false, 0},
	}

	if want, have := len(lines), len(protocolStats); want != have {
		t.Errorf("want %d parsed net/protocols lines, have %d", want, have)
	}
	for _, line := range protocolStats {
		if want, have := lines[line.Name], line; want != have {
			t.Errorf("%s: want %v, have %v", line.Name, want, have)
		}
	}
}
