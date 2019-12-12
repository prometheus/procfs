// Copyright 2019 The Prometheus Authors
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
	"reflect"
	"testing"
)

func TestSwaps(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		swap    *Swap
		invalid bool
	}{
		{
			name:    "device-mapper volume",
			s:       "/dev/dm-2                               partition       131068  1024    -2",
			invalid: false,
			swap: &Swap{
				Filename: "/dev/dm-2",
				Type:     "partition",
				Size:     131068,
				Used:     1024,
				Priority: -2,
			},
		},
		{
			name:    "Swap file",
			s:       "/foo                                    file            1048572 0       -3",
			invalid: false,
			swap: &Swap{
				Filename: "/foo",
				Type:     "file",
				Size:     1048572,
				Used:     0,
				Priority: -3,
			},
		},
		{
			name:    "Invalid number",
			s:       "/dev/sda2                               partition       hello   world   -2",
			invalid: true,
		},
		{
			name:    "Not enough fields",
			s:       "/dev/dm-2                               partition       131068  1024",
			invalid: true,
		},
	}

	for i, test := range tests {
		t.Logf("[%02d] test %q", i, test.name)

		swap, err := parseSwapString(test.s)

		if test.invalid && err == nil {
			t.Error("unexpected success")
		}
		if !test.invalid && err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(test.swap, swap) {
			t.Errorf("swap:\nwant:\n%+v\nhave:\n%+v", test.swap, swap)
		}
	}
}
