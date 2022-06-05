// Copyright 2021 The Prometheus Authors
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
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestSASPortClass(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	got, err := fs.SASPortClass()
	if err != nil {
		t.Fatal(err)
	}

	want := SASPortClass{
		"port-11:0": {
			Name:      "port-11:0",
			SASPhys:   []string{"phy-11:10", "phy-11:11", "phy-11:8", "phy-11:9"},
			Expanders: []string{"expander-11:0"},
		},
		"port-11:0:0": {Name: "port-11:0:0", SASPhys: []string{"phy-11:0:2"}},
		"port-11:0:1": {Name: "port-11:0:1", SASPhys: []string{"phy-11:0:4"}},
		"port-11:0:2": {Name: "port-11:0:2", SASPhys: []string{"phy-11:0:6"}},
		"port-11:1": {
			Name:      "port-11:1",
			SASPhys:   []string{"phy-11:12", "phy-11:13", "phy-11:14", "phy-11:15"},
			Expanders: []string{"expander-11:1"},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected SASDevice class (-want +got):\n%s", diff)
	}
}
