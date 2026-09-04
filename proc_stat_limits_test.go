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
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

const procStatLimitsTemplate = "26231 (vim) R 5392 7446 5392 34835 7446 4218880 32533 309516 26 82 1677 44 %d %d 20 0 1 0 82375 56274944 1981 18446744073709551615 4194304 6294284 140736914091744 140736914087944 139965136429984 0 0 12288 1870679807 0 0 0 17 0 0 0 31 0 0 8391624 8481048 16420864 140736914093252 140736914093279 140736914093279 140736914096107 0\n"

func testProcStatLimits(t *testing.T, minInt, maxInt int) {
	t.Helper()

	p := writeProcStatLimitFixture(t, fmt.Sprintf(procStatLimitsTemplate, minInt, maxInt))

	s, err := p.Stat()
	if err != nil {
		t.Fatalf("want no error, have %s", err)
	}

	for _, test := range []struct {
		name string
		want int
		have int
	}{
		{name: "waited for children user time", want: minInt, have: s.CUTime},
		{name: "waited for children system time", want: maxInt, have: s.CSTime},
	} {
		if test.want != test.have {
			t.Errorf("want %s %d, have %d", test.name, test.want, test.have)
		}
	}
}

func writeProcStatLimitFixture(t *testing.T, stat string) Proc {
	t.Helper()

	const pid = 26231

	root := t.TempDir()
	procDir := filepath.Join(root, strconv.Itoa(pid))
	if err := os.MkdirAll(procDir, 0o755); err != nil {
		t.Fatalf("create proc fixture dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(procDir, "stat"), []byte(stat), 0o644); err != nil {
		t.Fatalf("write stat fixture: %v", err)
	}

	fs, err := NewFS(root)
	if err != nil {
		t.Fatalf("create pseudo fs: %v", err)
	}
	p, err := fs.Proc(pid)
	if err != nil {
		t.Fatalf("open proc fixture: %v", err)
	}

	return p
}
