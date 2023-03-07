// Copyright 2023 The Prometheus Authors
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
// +build linux

package sysfs

import (
	"testing"
)

func TestMeminfoNUMA(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	meminfo, err := fs.MeminfoNUMA()
	if err != nil {
		t.Fatal(err)
	}

	if want, got := uint64(133), meminfo[1].HugePages_Total; want != got {
		t.Errorf("want meminfo stat HugePages_Total value %d, got %d", want, got)
	}

	if want, got := uint64(134), meminfo[1].HugePages_Free; want != got {
		t.Errorf("want meminfo stat HugePages_Free value %d, got %d", want, got)
	}
}
