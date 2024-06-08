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
	"reflect"
	"testing"
)

func TestClassDRMCard(t *testing.T) {
	fs, err := NewFS(sysTestFixtures)
	if err != nil {
		t.Fatal(err)
	}

	drmCardTest, err := fs.ClassDrmCard()
	if err != nil {
		t.Fatal(err)
	}

	classDrmCard := []ClassDrmCard{
		{
			Name:   "card0",
			Enable: 1,
			Driver: "amdgpu",
		},
		{
			Name:   "card1",
			Enable: 1,
			Driver: "i915",
		},
	}

	if !reflect.DeepEqual(classDrmCard, drmCardTest) {
		t.Errorf("Result not correct: want %v, have %v", classDrmCard, drmCardTest)
	}
}
