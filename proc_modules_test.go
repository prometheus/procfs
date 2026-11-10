// Copyright 2022 The Prometheus Authors
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

func TestProcModules(t *testing.T) {
	fs := getProcFixtures(t)

	modules, err := fs.Modules()
	if err != nil {
		t.Fatal(err)
	}

	if want, have := 89, len(modules); want != have {
		t.Errorf("want length %d, have %d", want, have)
	}

	for _, test := range []struct {
		name  string
		index int
		want  Module
	}{
		{
			name:  "no dependencies and taints",
			index: 0,
			want: Module{
				Name:         "nft_counter",
				Size:         16384,
				Instances:    4,
				Dependencies: []string{},
				State:        "Live",
				Offset:       0,
				Taints:       []string{},
			},
		},
		{
			name:  "have dependencies with no taints",
			index: 11,
			want: Module{
				Name:         "nf_tables",
				Size:         245760,
				Instances:    19,
				Dependencies: []string{"nft_counter", "nft_chain_nat", "nft_compat"},
				State:        "Live",
				Offset:       0,
				Taints:       []string{},
			},
		},
		{
			name:  "have multiple taints with multiple dependencies",
			index: 83,
			want: Module{
				Name:         "drm",
				Size:         622592,
				Instances:    3,
				Dependencies: []string{"virtio_gpu", "drm_kms_helper"},
				State:        "Live",
				Offset:       0,
				Taints:       []string{"P", "O", "E"},
			},
		},
		{
			name:  "have single taint with single dependency",
			index: 88,
			want: Module{
				Name:         "failover",
				Size:         16384,
				Instances:    1,
				Dependencies: []string{"net_failover"},
				State:        "Live",
				Offset:       0,
				Taints:       []string{"P"},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			if want, have := test.want, modules[test.index]; !reflect.DeepEqual(want, have) {
				t.Errorf("want %v, have %v", want, have)
			}
		})
	}
}
