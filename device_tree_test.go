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

//go:build !windows

package procfs

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestDeviceTreeSerialNumber(t *testing.T) {
	for _, tt := range []struct {
		name string
		data string
		want string
	}{
		{name: "terminated", data: "000000001234abcd\x00", want: "000000001234abcd"},
		{name: "unterminated", data: "000000001234abcd", want: "000000001234abcd"},
		{name: "whitespace", data: " serial number \x00", want: " serial number "},
		{name: "empty", data: "\x00", want: ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			procDir := filepath.Join(root, "proc")
			treeDir := filepath.Join(root, "sys", "firmware", "devicetree", "base")
			for _, dir := range []string{procDir, treeDir} {
				if err := os.MkdirAll(dir, 0o755); err != nil {
					t.Fatal(err)
				}
			}
			if err := os.Symlink(treeDir, filepath.Join(procDir, "device-tree")); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(treeDir, "serial-number"), []byte(tt.data), 0o644); err != nil {
				t.Fatal(err)
			}
			pfs, err := NewFS(procDir)
			if err != nil {
				t.Fatal(err)
			}
			got, err := pfs.DeviceTreeSerialNumber()
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("want serial number %q, got %q", tt.want, got)
			}
		})
	}
}

func TestDeviceTreeSerialNumberMissing(t *testing.T) {
	for _, name := range []string{"missing device tree", "missing serial number", "dangling symlink"} {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			treePath := filepath.Join(root, "device-tree")
			switch name {
			case "missing serial number":
				if err := os.Mkdir(treePath, 0o755); err != nil {
					t.Fatal(err)
				}
			case "dangling symlink":
				if err := os.Symlink(filepath.Join(root, "absent"), treePath); err != nil {
					t.Fatal(err)
				}
			}
			pfs, err := NewFS(root)
			if err != nil {
				t.Fatal(err)
			}
			got, err := pfs.DeviceTreeSerialNumber()
			if !errors.Is(err, fs.ErrNotExist) {
				t.Fatalf("want missing-file error, got %v", err)
			}
			if got != "" {
				t.Errorf("want no serial number, got %q", got)
			}
		})
	}
}
