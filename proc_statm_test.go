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
	"math"
	"os"
	"testing"
)

func TestProcStatm(t *testing.T) {
	statm, err := testProcStatm(26231)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name string
		want uint64
		have uint64
	}{
		{name: "Pid", want: 26231, have: uint64(statm.PID)},
		{name: "Size", want: 149919, have: statm.Size},
		{name: "Resident", want: 12547, have: statm.Resident},
		{name: "Shared", want: 18446744073709551615, have: statm.Shared},
		{name: "Text", want: 19864, have: statm.Text},
		{name: "Lib", want: 0, have: statm.Lib},
		{name: "Data", want: 14531, have: statm.Data},
		{name: "Dt", want: 0, have: statm.Dt},
	} {
		if test.want != test.have {
			t.Errorf("want %s %d, have %d", test.name, test.want, test.have)
		}
	}
}

func TestProcStatmLimits(t *testing.T) {
	statm, err := testProcStatm(26231)
	if err != nil {
		t.Fatal(err)
	}

	// max values of statm int fields
	for _, test := range []struct {
		name string
		want uint64
		have uint64
	}{
		{name: "number of resident shared pages in process", want: math.MaxUint64, have: statm.Shared},
		{name: "number of dirty pages in process", want: 0, have: statm.Dt},
	} {
		if test.want != test.have {
			t.Errorf("want %s %d, have %d", test.name, test.want, test.have)
		}
	}
}

func TestSizeBytes(t *testing.T) {
	statm, err := testProcStatm(26231)
	if err != nil {
		t.Fatal(err)
	}

	if want, have := statm.Size*uint64(os.Getpagesize()), statm.SizeBytes(); want != have {
		t.Errorf("want total program memory %d, have %d", want, have)
	}
}

func TestResidentBytes(t *testing.T) {
	statm, err := testProcStatm(26231)
	if err != nil {
		t.Fatal(err)
	}

	if want, have := statm.Resident*uint64(os.Getpagesize()), statm.ResidentBytes(); want != have {
		t.Errorf("want resident memory %d, have %d", want, have)
	}
}

func TestSHRBytes(t *testing.T) {
	statm, err := testProcStatm(26231)
	if err != nil {
		t.Fatal(err)
	}

	if want, have := statm.Shared*uint64(os.Getpagesize()), statm.SHRBytes(); want != have {
		t.Errorf("want share memory %d, have %d", want, have)
	}
}

func TestTextBytes(t *testing.T) {
	statm, err := testProcStatm(26231)
	if err != nil {
		t.Fatal(err)
	}

	if want, have := statm.Text*uint64(os.Getpagesize()), statm.TextBytes(); want != have {
		t.Errorf("want text (code) size %d, have %d", want, have)
	}
}

func TestDataBytes(t *testing.T) {
	statm, err := testProcStatm(26231)
	if err != nil {
		t.Fatal(err)
	}

	if want, have := statm.Data*uint64(os.Getpagesize()), statm.DataBytes(); want != have {
		t.Errorf("want data + stack size %d, have %d", want, have)
	}
}

func testProcStatm(pid int) (ProcStatm, error) {
	fs, err := NewFS(procTestFixtures)
	if err != nil {
		return ProcStatm{}, err
	}
	p, err := fs.Proc(pid)
	if err != nil {
		return ProcStatm{}, err
	}

	return p.Statm()
}
