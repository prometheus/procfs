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

func TestFscacheinfo(t *testing.T) {
	expected := Fscacheinfo{
		cookies_idx: 3,
		cookies_dat: 67877,
		cookies_spc: 0,
		objects_alc: 67473,
		objects_nal: 0,
		objects_avl: 67473,
		objects_ded: 388,
		chkaux_non:  12,
		chkaux_ok:   33,
		chkaux_upd:  44,
		chkaux_obs:  55,
		pages_mrk:   547164,
		pages_unc:   364577,
		acquire_n:   67880,
		acquire_nul: 98,
		acquire_noc: 25,
		acquire_ok:  67780,
		acquire_nbf: 39,
		acquire_oom: 26,
		lookups_n:   67473,
		lookups_neg: 67470,
		lookups_pos: 58,
		lookups_crt: 67473,
		lookups_tmo: 85,
		invals_n:    14,
		invals_run:  13,
		updates_n:   7,
		updates_nul: 3,
		updates_run: 8,
		relinqs_n:   394,
		relinqs_nul: 1,
		relinqs_wcr: 2,
		relinqs_rtr: 3,
		attrchg_n:   6,
		attrchg_ok:  5,
		attrchg_nbf: 4,
		attrchg_oom: 3,
		attrchg_run: 2,
		allocs_n:    20,
		allocs_ok:   19,
		allocs_wt:   18,
		allocs_nbf:  17,
		allocs_int:  16,
		allocs_ops:  15,
		allocs_owt:  14,
		allocs_abt:  13,
		retrvls_n:   151959,
		retrvls_ok:  82823,
		retrvls_wt:  23467,
		retrvls_nod: 69136,
		retrvls_nbf: 15,
		retrvls_int: 69,
		retrvls_oom: 43,
		retrvls_ops: 151959,
		retrvls_owt: 42747,
		retrvls_abt: 44,
		stores_n:    225565,
		stores_ok:   225565,
		stores_agn:  12,
		stores_nbf:  13,
		stores_oom:  14,
		stores_ops:  69156,
		stores_run:  294721,
		stores_pgs:  225565,
		stores_rxd:  225565,
		stores_olm:  43,
		vmscan_nos:  364512,
		vmscan_gon:  2,
		vmscan_bsy:  43,
		vmscan_can:  12,
		vmscan_wt:   66,
		ops_pend:    42753,
		ops_run:     221129,
		ops_enq:     628798,
		ops_can:     11,
		ops_rej:     88,
		ops_ini:     377538,
		ops_dfr:     27,
		ops_rel:     377538,
		ops_gc:      37,
		cacheop_alo: 1,
		cacheop_luo: 2,
		cacheop_luc: 3,
		cacheop_gro: 4,
		cacheop_inv: 5,
		cacheop_upo: 6,
		cacheop_dro: 7,
		cacheop_pto: 8,
		cacheop_atc: 9,
		cacheop_syn: 10,
		cacheop_rap: 11,
		cacheop_ras: 12,
		cacheop_alp: 13,
		cacheop_als: 14,
		cacheop_wrp: 15,
		cacheop_ucp: 16,
		cacheop_dsp: 17,
		cacheev_nsp: 18,
		cacheev_stl: 19,
		cacheev_rtr: 20,
		cacheev_cul: 21,
	}

	have, err := getProcFixtures(t).Fscacheinfo()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(have, expected) {
		t.Logf("have: %+v", have)
		t.Logf("expected: %+v", expected)
		t.Errorf("structs are not equal")
	}
}
