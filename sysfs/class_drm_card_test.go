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
// +build linux

package sysfs

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestClassDRMCard(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	got, err := fs.DRMCardClass()
	if err != nil {
		t.Fatal(err)
	}

	want := DRMCardClass{
		"card0": DRMCard{
			Name:   "card0",
			Driver: "amdgpu",
			Ports:  map[string]DRMCardPort{},
		},
		"card1": DRMCard{
			Name:   "card1",
			Driver: "i915",
			Ports: map[string]DRMCardPort{
				"card1-DP-1": {
					Name:    "card1-DP-1",
					DPMS:    "Off",
					Enabled: "disabled",
					Status:  "disconnected",
				},
				"card1-DP-5": {
					Name:    "card1-DP-5",
					DPMS:    "On",
					Enabled: "enabled",
					Status:  "connected",
				},
			},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected DRMCard class (-want +got):\n%s", diff)
	}
}
