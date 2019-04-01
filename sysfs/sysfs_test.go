// Copyright 2017 The Prometheus Authors
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

package sysfs_test

import (
	"os"
	"testing"
)

const (
	sysTestFixtures = "../fixtures/sys"
)

func TestSysFSFixturesDir(t *testing.T) {
	info, err := os.Stat(sysTestFixtures)
	if err != nil {
		t.Errorf("could not read %s: %s", sysTestFixtures, err)
	}
	if !info.IsDir() {
		t.Errorf("mount point %s is not a directory", sysTestFixtures)
	}
}
