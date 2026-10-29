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
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestProcEnviron(t *testing.T) {
	p, err := getProcFixtures(t).Proc(26231)
	if err != nil {
		t.Fatal(err)
	}

	environments, err := p.Environ()
	if err != nil {
		t.Fatal(err)
	}

	expectedEnvironments := []string{
		"PATH=/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		"HOSTNAME=cd24e11f73a5",
		"TERM=xterm",
		"GOLANG_VERSION=1.12.5",
		"GOPATH=/go",
		"HOME=/root",
	}

	if want, have := len(expectedEnvironments), len(environments); want != have {
		t.Errorf("want %d parsed environments, have %d", want, have)
	}

	for i, environment := range environments {
		if want, have := expectedEnvironments[i], environment; want != have {
			t.Errorf("%d: want %v, have %v", i, want, have)
		}
	}
}

func TestProcEnvironTrailingNull(t *testing.T) {
	for _, tt := range []struct {
		name string
		data string
		want []string
	}{
		{name: "empty", data: "", want: []string{}},
		{name: "single terminated", data: "A=1\x00", want: []string{"A=1"}},
		{name: "single unterminated", data: "A=1", want: []string{"A=1"}},
		{name: "multiple terminated", data: "A=1\x00B=2\x00", want: []string{"A=1", "B=2"}},
		{name: "multiple unterminated", data: "A=1\x00B=2", want: []string{"A=1", "B=2"}},
		{name: "empty value", data: "A=1\x00B=", want: []string{"A=1", "B="}},
		{name: "empty entry", data: "A=1\x00\x00B=2\x00", want: []string{"A=1", "", "B=2"}},
		{name: "trailing empty entry", data: "A=1\x00\x00", want: []string{"A=1", ""}},
		{name: "only nulls", data: "\x00\x00", want: []string{"", ""}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			procDir := filepath.Join(dir, "1")
			if err := os.Mkdir(procDir, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(procDir, "environ"), []byte(tt.data), 0o644); err != nil {
				t.Fatal(err)
			}
			fs, err := NewFS(dir)
			if err != nil {
				t.Fatal(err)
			}
			p, err := fs.Proc(1)
			if err != nil {
				t.Fatal(err)
			}
			got, err := p.Environ()
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("unexpected environments (-want +got):\n%s", diff)
			}
		})
	}
}
