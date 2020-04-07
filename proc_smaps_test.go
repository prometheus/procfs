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

// +build !windows

package procfs

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/sys/unix"
)

func TestProcSmaps(t *testing.T) {
	p, err := getProcFixtures(t).Proc(26231)
	if err != nil {
		t.Fatal(err)
	}

	s, err := p.ProcSMaps()
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name string
		want int
		have int
	}{
		{name: "length", want: 13, have: len(s)},
		{name: "SizeSum", want: 154988 * 1024, have: int(s.SizeSum())},
		{name: "SizeSum", want: 11268 * 1024, have: int(s.RssSum())},
		{name: "SizeSum", want: 2852 * 1024, have: int(s.PssSum())},
		{name: "SizeSum", want: 5268 * 1024, have: int(s.SharedCleanSum())},
		{name: "SizeSum", want: 4716 * 1024, have: int(s.SharedDirtySum())},
		{name: "SizeSum", want: 0 * 1024, have: int(s.PrivateCleanSum())},
		{name: "SizeSum", want: 1284 * 1024, have: int(s.PrivateDirtySum())},
		{name: "SizeSum", want: 10300 * 1024, have: int(s.ReferencedSum())},
		{name: "SizeSum", want: 1608 * 1024, have: int(s.AnonymousSum())},
		{name: "SizeSum", want: 8 * 1024, have: int(s.SwapSum())},
		{name: "SizeSum", want: 0 * 1024, have: int(s.SwapPssSum())},
	} {
		if test.want != test.have {
			t.Errorf("want %s %d, have %d", test.name, test.want, test.have)
		}
	}

	wantMap := &ProcMap{
		StartAddr: 0x561c3902b000,
		EndAddr:   0x561c39757000,
		Perms:     &ProcMapPermissions{true, false, true, false, true},
		Offset:    0,
		Dev:       unix.Mkdev(0xfd, 0x01),
		Inode:     664619,
		Pathname:  "/usr/lib/postgresql/11/bin/postgres",
	}

	if len(s) > 0 {
		if diff := cmp.Diff(wantMap, s[0].ProcMap); diff != "" {
			t.Fatalf("unexpected proc/map entry (-want +got):\n%s", diff)
		}
	}
}
