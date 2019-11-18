// Copyright 2018 The Prometheus Authors
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
	"reflect"
	"testing"
)

func Test_parseNetUDPLine(t *testing.T) {
	type args struct {
		fields []string
	}
	tests := []struct {
		name    string
		args    args
		want    *NetUDPLine
		wantErr bool
	}{
		{
			name: "reading valid lines, no issue should happened",
			args: args{
				fields: []string{"1:", "00000000:0000", "00000000:0000", "07", "00000017:0000002A"},
			},
			want: &NetUDPLine{TxQueue: 23, RxQueue: 42},
		},
		{
			name: "error case - invalid line - number of fields/columns < 5",
			args: args{
				fields: []string{"1:", "00000000:0000", "00000000:0000", "07"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error case - cannot parse line - missing colon",
			args: args{
				fields: []string{"1:", "00000000:0000", "00000000:0000", "07", "0000000000000001"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error case - parse tx_queue - not an valid hex",
			args: args{
				fields: []string{"1:", "00000000:0000", "00000000:0000", "07", "DEADCODE:00000001"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error case - parse rx_queue - not an valid hex",
			args: args{
				fields: []string{"1:", "00000000:0000", "00000000:0000", "07", "00000000:FEEDCODE"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseNetUDPLine(tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseNetUDPLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want == nil && got != nil {
				t.Errorf("parseNetUDPLine() = %v, want %v", got, tt.want)
			}
			if got != nil {
				if (got.RxQueue != tt.want.RxQueue) || (got.TxQueue != tt.want.TxQueue) {
					t.Errorf("parseNetUDPLine() = %#v, want %#v", got, tt.want)
				}
			}
		})
	}
}

func Test_newNetUDP(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    *NetUDP
		wantErr bool
	}{
		{
			name:    "file found, no error should come up",
			args:    args{file: "fixtures/proc/net/udp"},
			want:    &NetUDP{TxQueueLength: 2, RxQueueLength: 2, UsedSockets: 3},
			wantErr: false,
		},
		{
			name:    "error case - file not found",
			args:    args{file: "somewhere over the rainbow"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error case - parse error",
			args:    args{file: "fixtures/proc/net/udp_broken"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newNetUDP(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("newNetUDP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newNetUDP() = %v, want %v", got, tt.want)
			}
		})
	}
}
