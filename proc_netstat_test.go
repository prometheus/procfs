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
	"reflect"
	"strings"
	"testing"
)

func TestProcNetstat(t *testing.T) {
	p, err := getProcFixtures(t).Proc(26231)
	if err != nil {
		t.Fatal(err)
	}

	procNetstat, err := p.Netstat()
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name string
		want float64
		have float64
	}{
		{name: "pid", want: 26231, have: float64(procNetstat.PID)},
		{name: "TcpExt:SyncookiesSent", want: 0, have: *procNetstat.SyncookiesSent},
		{name: "TcpExt:EmbryonicRsts", want: 1, have: *procNetstat.EmbryonicRsts},
		{name: "TcpExt:TW", want: 83, have: *procNetstat.TW},
		{name: "TcpExt:PAWSEstab", want: 3640, have: *procNetstat.PAWSEstab},

		{name: "IpExt:InNoRoutes", want: 0, have: *procNetstat.InNoRoutes},
		{name: "IpExt:InMcastPkts", want: 208, have: *procNetstat.InMcastPkts},
		{name: "IpExt:OutMcastPkts", want: 214, have: *procNetstat.OutMcastPkts},
	} {
		if test.want != test.have {
			t.Errorf("want %s %f, have %f", test.name, test.want, test.have)
		}
	}
}

// TestParseProcNetstatLegacyFields covers TcpExt counters that are no longer
// part of the kernel's TcpExt table but were still printed by older kernels.
func TestParseProcNetstatLegacyFields(t *testing.T) {
	payload := `TcpExt: SyncookiesSent PAWSPassive TCPPrequeued TCPDirectCopyFromBacklog TCPDirectCopyFromPrequeue TCPPrequeueDropped TCPLoss TCPFACKReorder TCPForwardRetrans TCPHPHitsToUser TCPSchedulerFailed
TcpExt: 1 2 3 4 5 6 7 8 9 10 11
IpExt: InNoRoutes OutOctets
IpExt: 12 13`

	procNetstat, err := parseProcNetstat(strings.NewReader(payload), "net/netstat")
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name string
		want float64
		have float64
	}{
		{name: "TcpExt:SyncookiesSent", want: 1, have: *procNetstat.SyncookiesSent},
		{name: "TcpExt:PAWSPassive", want: 2, have: *procNetstat.PAWSPassive},
		{name: "TcpExt:TCPPrequeued", want: 3, have: *procNetstat.TCPPrequeued},
		{name: "TcpExt:TCPDirectCopyFromBacklog", want: 4, have: *procNetstat.TCPDirectCopyFromBacklog},
		{name: "TcpExt:TCPDirectCopyFromPrequeue", want: 5, have: *procNetstat.TCPDirectCopyFromPrequeue},
		{name: "TcpExt:TCPPrequeueDropped", want: 6, have: *procNetstat.TCPPrequeueDropped},
		{name: "TcpExt:TCPLoss", want: 7, have: *procNetstat.TCPLoss},
		{name: "TcpExt:TCPFACKReorder", want: 8, have: *procNetstat.TCPFACKReorder},
		{name: "TcpExt:TCPForwardRetrans", want: 9, have: *procNetstat.TCPForwardRetrans},
		{name: "TcpExt:TCPHPHitsToUser", want: 10, have: *procNetstat.TCPHPHitsToUser},
		{name: "TcpExt:TCPSchedulerFailed", want: 11, have: *procNetstat.TCPSchedulerFailed},
		{name: "IpExt:InNoRoutes", want: 12, have: *procNetstat.InNoRoutes},
		{name: "IpExt:OutOctets", want: 13, have: *procNetstat.OutOctets},
	} {
		if test.want != test.have {
			t.Errorf("want %s %f, have %f", test.name, test.want, test.have)
		}
	}

	modern := `TcpExt: SyncookiesSent
TcpExt: 1
IpExt: InNoRoutes
IpExt: 12`

	pn, err := parseProcNetstat(strings.NewReader(modern), "net/netstat")
	if err != nil {
		t.Fatal(err)
	}

	if pn.SyncookiesSent == nil || *pn.SyncookiesSent != 1 {
		t.Error("want TcpExt:SyncookiesSent 1")
	}
	if pn.TCPPrequeued != nil {
		t.Error("TCPPrequeued should be nil when the kernel does not report it")
	}
}

// TestParseProcNetstatAllFields makes sure every field modelled in TcpExt and
// IpExt is filled by the parser. A field added to a struct without a matching
// case in parseProcNetstat would otherwise silently stay nil.
func TestParseProcNetstatAllFields(t *testing.T) {
	var payload strings.Builder
	for _, s := range []any{TcpExt{}, IpExt{}} {
		typ := reflect.TypeOf(s)
		names, values := []string{typ.Name() + ":"}, []string{typ.Name() + ":"}
		for i := 0; i < typ.NumField(); i++ {
			names = append(names, typ.Field(i).Name)
			values = append(values, fmt.Sprint(i+1))
		}
		fmt.Fprintf(&payload, "%s\n%s\n", strings.Join(names, " "), strings.Join(values, " "))
	}

	procNetstat, err := parseProcNetstat(strings.NewReader(payload.String()), "net/netstat")
	if err != nil {
		t.Fatal(err)
	}

	for _, s := range []any{procNetstat.TcpExt, procNetstat.IpExt} {
		v := reflect.ValueOf(s)
		for i := 0; i < v.NumField(); i++ {
			name := v.Type().Name() + ":" + v.Type().Field(i).Name
			f, ok := v.Field(i).Interface().(*float64)
			switch {
			case !ok:
				t.Errorf("%s: unexpected field type %s", name, v.Field(i).Type())
			case f == nil:
				t.Errorf("%s: not parsed", name)
			case *f != float64(i+1):
				t.Errorf("%s: want %d, have %f", name, i+1, *f)
			}
		}
	}
}
