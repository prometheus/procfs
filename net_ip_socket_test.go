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
	"net"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseNetIPSocketLineParseErrorsIncludeRawValue(t *testing.T) {
	baseFields := []string{"1:", "00000000:0000", "00000000:0000", "07", "00000000:00000001", "0:0", "0", "10", "0", "39309", "2", "000000009bd60d72", "5"}
	tests := []struct {
		name  string
		value string
		want  string
		field int
		isUDP bool
	}{
		{name: "sl", field: 0, value: "invalid-sl:", want: "invalid-sl"},
		{name: "local port", field: 1, value: "00000000:invalid-local-port", want: "invalid-local-port"},
		{name: "remote port", field: 2, value: "00000000:invalid-remote-port", want: "invalid-remote-port"},
		{name: "state", field: 3, value: "invalid-state", want: "invalid-state"},
		{name: "transmit queue", field: 4, value: "invalid-tx-queue:00000001", want: "invalid-tx-queue"},
		{name: "receive queue", field: 4, value: "00000000:invalid-rx-queue", want: "invalid-rx-queue"},
		{name: "UID", field: 7, value: "invalid-uid", want: "invalid-uid"},
		{name: "inode", field: 9, value: "invalid-inode", want: "invalid-inode"},
		{name: "drops", field: 12, value: "invalid-drops", want: "invalid-drops", isUDP: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fields := append([]string(nil), baseFields...)
			fields[test.field] = test.value

			_, err := parseNetIPSocketLine(fields, test.isUDP)
			if err == nil {
				t.Fatal("expected an error, but none occurred")
			}
			if !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error %q does not contain raw value %q", err, test.want)
			}
		})
	}
}

func Test_parseNetIPSocketLine(t *testing.T) {
	tests := []struct {
		fields  []string
		name    string
		want    *netIPSocketLine
		wantErr bool
		isUDP   bool
	}{
		{
			name:   "reading valid lines, no issue should happened",
			fields: []string{"11:", "00000000:0000", "00000000:0000", "0A", "00000017:0000002A", "0:0", "0", "1000", "0", "39309"},
			want: &netIPSocketLine{
				Sl:        11,
				LocalAddr: net.IP{0, 0, 0, 0},
				LocalPort: 0,
				RemAddr:   net.IP{0, 0, 0, 0},
				RemPort:   0,
				St:        10,
				TxQueue:   23,
				RxQueue:   42,
				UID:       1000,
				Inode:     39309,
			},
		},
		{
			name:    "error case - invalid line - number of fields/columns < 10",
			fields:  []string{"1:", "00000000:0000", "00000000:0000", "07", "0:0", "0", "0"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error case - parse sl - not a valid uint",
			fields:  []string{"a:", "00000000:0000", "00000000:0000", "07", "00000000:00000001", "0:0", "0", "0", "0", "39309"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error case - parse local_address - not a valid hex",
			fields:  []string{"1:", "0000000O:0000", "00000000:0000", "07", "00000000:00000001", "0:0", "0", "0", "0", "39309"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error case - parse rem_address - not a valid hex",
			fields:  []string{"1:", "00000000:0000", "0000000O:0000", "07", "00000000:00000001", "0:0", "0", "0", "0", "39309"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error case - cannot parse line - missing colon",
			fields:  []string{"1:", "00000000:0000", "00000000:0000", "07", "0000000000000001", "0:0", "0", "0", "0", "39309"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error case - parse tx_queue - not a valid hex",
			fields:  []string{"1:", "00000000:0000", "00000000:0000", "07", "DEADCODE:00000001", "0:0", "0", "0", "0", "39309"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error case - parse rx_queue - not a valid hex",
			fields:  []string{"1:", "00000000:0000", "00000000:0000", "07", "00000000:FEEDCODE", "0:0", "0", "0", "0", "39309"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error case - parse UID - not a valid uint",
			fields:  []string{"1:", "00000000:0000", "00000000:0000", "07", "00000000:00000001", "0:0", "0", "-10", "0", "39309"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error case - parse Inode - not a valid uint",
			fields:  []string{"1:", "00000000:0000", "00000000:0000", "07", "00000000:00000001", "0:0", "0", "-10", "0", "-39309"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error case - parse Drops - not a valid uint",
			fields:  []string{"1:", "00000000:0000", "00000000:0000", "07", "00000000:00000001", "0:0", "0", "10", "0", "39309", "2", "000000009bd60d72", "-5"},
			want:    nil,
			wantErr: true,
			isUDP:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseNetIPSocketLine(tt.fields, tt.isUDP)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseNetIPSocketLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want == nil && got != nil {
				t.Errorf("parseNetIPSocketLine() = %v, want %v", got, tt.want)
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Fatalf("unexpected diff (-want +got):\n%s", diff)
			}
		})
	}
}
