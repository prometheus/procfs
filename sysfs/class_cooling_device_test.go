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

//go:build linux

package sysfs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestClassCoolingDeviceStats(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	coolingDeviceTest, err := fs.ClassCoolingDeviceStats()
	if err != nil {
		t.Fatal(err)
	}

	classCoolingDeviceStats := []ClassCoolingDeviceStats{
		{
			Name:     "0",
			Type:     "Processor",
			MaxState: 50,
			CurState: 0,
		},
		{
			Name:     "2",
			Type:     "intel_powerclamp",
			MaxState: 27,
			CurState: -1,
		},
	}

	if diff := cmp.Diff(classCoolingDeviceStats, coolingDeviceTest); diff != "" {
		t.Fatalf("unexpected diff (-want +got):\n%s", diff)
	}
}

func TestClassCoolingDeviceStatsNoData(t *testing.T) {
	root := t.TempDir()
	paths := []string{}
	want := []ClassCoolingDeviceStats{}
	for i := range 3 {
		path := filepath.Join(root, "class", "thermal", fmt.Sprintf("cooling_device%d", i))
		if err := os.MkdirAll(path, 0o700); err != nil {
			t.Fatal(err)
		}
		for name, value := range map[string]string{"type": "Processor\n", "max_state": "50\n", "cur_state": "0\n"} {
			if err := os.WriteFile(filepath.Join(path, name), []byte(value), 0o600); err != nil {
				t.Fatal(err)
			}
		}
		paths = append(paths, path)
		want = append(want, ClassCoolingDeviceStats{Name: strconv.Itoa(i), Type: "Processor", MaxState: 50, CurState: 0})
	}
	fs, err := NewFS(root)
	if err != nil {
		t.Fatal(err)
	}
	actual, err := fs.ClassCoolingDeviceStats()
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(want, actual); diff != "" {
		t.Fatalf("real neighboring parsers (-want +got):\n%s", diff)
	}
	noData := []struct {
		name string
		err  error
	}{
		{"direct", syscall.ENODATA},
		{"path", &os.PathError{Op: "read", Path: filepath.Join(paths[1], "cur_state"), Err: syscall.ENODATA}},
		{"wrapped", fmt.Errorf("cooling device read: %w", &os.PathError{Op: "read", Path: filepath.Join(paths[1], "cur_state"), Err: syscall.ENODATA})},
	}
	for _, tc := range noData {
		t.Run(tc.name, func(t *testing.T) {
			seen := []string{}
			actual, err := classCoolingDeviceStats(paths, func(path string) (ClassCoolingDeviceStats, error) {
				seen = append(seen, path)
				if path == paths[1] {
					return ClassCoolingDeviceStats{}, tc.err
				}
				return parseCoolingDeviceStats(path)
			})
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff([]ClassCoolingDeviceStats{want[0], want[2]}, actual); diff != "" {
				t.Fatalf("device-local no-data (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(paths, seen); diff != "" {
				t.Fatalf("must continue after no-data device (-want +got):\n%s", diff)
			}
		})
	}
	for _, failure := range []error{syscall.EIO, fmt.Errorf("read failed: %w", &os.PathError{Op: "read", Path: filepath.Join(paths[1], "cur_state"), Err: syscall.EIO})} {
		actual, err := classCoolingDeviceStats(paths, func(path string) (ClassCoolingDeviceStats, error) {
			if path == paths[1] {
				return ClassCoolingDeviceStats{}, failure
			}
			return parseCoolingDeviceStats(path)
		})
		if !errors.Is(err, syscall.EIO) || len(actual) != 0 {
			t.Fatalf("unexpected I/O failure must propagate: stats=%v err=%v", actual, err)
		}
	}
	if err := os.WriteFile(filepath.Join(paths[1], "cur_state"), []byte("not-a-number\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if actual, err := fs.ClassCoolingDeviceStats(); err == nil || len(actual) != 0 {
		t.Fatalf("real malformed state must propagate: stats=%v err=%v", actual, err)
	}
}
