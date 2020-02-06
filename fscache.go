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
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/prometheus/procfs/internal/util"
)

// Fscache represents memory statistics.
type Fscacheinfo struct {
	// Number of index cookies allocated
	cookies_idx uint64
	// data storage cookies allocated
	cookies_dat uint64
	// Number of special cookies allocated
	cookies_spc uint64
	// Number of objects allocated
	objects_alc uint64
	// Number of object allocation failures
	objects_nal uint64
	// Number of objects that reached the available state
	objects_avl uint64
	// Number of objects that reached the dead state
	objects_ded uint64
	// Number of objects that didn't have a coherency check
	chkaux_non uint64
	// Number of objects that passed a coherency check
	chkaux_ok uint64
	// Number of objects that needed a coherency data update
	chkaux_upd uint64
	// Number of objects that were declared obsolete
	chkaux_obs uint64
	// Number of pages marked as being cached
	pages_mrk uint64
	// Number of uncache page requests seen
	pages_unc uint64
	// Number of acquire cookie requests seen
	acquire_n uint64
	// Number of acq reqs given a NULL parent
	acquire_nul uint64
	// Number of acq reqs rejected due to no cache available
	acquire_noc uint64
	// Number of acq reqs succeeded
	acquire_ok uint64
	// Number of acq reqs rejected due to error
	acquire_nbf uint64
	// Number of acq reqs failed on ENOMEM
	acquire_oom uint64
	// Number of lookup calls made on cache backends
	lookups_n uint64
	// Number of negative lookups made
	lookups_neg uint64
	// Number of positive lookups made
	lookups_pos uint64
	// Number of objects created by lookup
	lookups_crt uint64
	// Number of lookups timed out and requeued
	lookups_tmo uint64
	invals_n    uint64
	invals_run  uint64
	// Number of update cookie requests seen
	updates_n uint64
	// Number of upd reqs given a NULL parent
	updates_nul uint64
	// Number of upd reqs granted CPU time
	updates_run uint64
	// Number of relinquish cookie requests seen
	relinqs_n uint64
	// Number of rlq reqs given a NULL parent
	relinqs_nul uint64
	// Number of rlq reqs waited on completion of creation
	relinqs_wcr uint64
	// Relinqs rtr
	relinqs_rtr uint64
	// Number of attribute changed requests seen
	attrchg_n uint64
	// Number of attr changed requests queued
	attrchg_ok uint64
	// Number of attr changed rejected -ENOBUFS
	attrchg_nbf uint64
	// Number of attr changed failed -ENOMEM
	attrchg_oom uint64
	// Number of attr changed ops given CPU time
	attrchg_run uint64
	// Number of allocation requests seen
	allocs_n uint64
	// Number of successful alloc reqs
	allocs_ok uint64
	// Number of alloc reqs that waited on lookup completion
	allocs_wt uint64
	// Number of alloc reqs rejected -ENOBUFS
	allocs_nbf uint64
	// Number of alloc reqs aborted -ERESTARTSYS
	allocs_int uint64
	// Number of alloc reqs submitted
	allocs_ops uint64
	// Number of alloc reqs waited for CPU time
	allocs_owt uint64
	// Number of alloc reqs aborted due to object death
	allocs_abt uint64
	// Number of retrieval (read) requests seen
	retrvls_n uint64
	// Number of successful retr reqs
	retrvls_ok uint64
	// Number of retr reqs that waited on lookup completion
	retrvls_wt uint64
	// Number of retr reqs returned -ENODATA
	retrvls_nod uint64
	// Number of retr reqs rejected -ENOBUFS
	retrvls_nbf uint64
	// Number of retr reqs aborted -ERESTARTSYS
	retrvls_int uint64
	// Number of retr reqs failed -ENOMEM
	retrvls_oom uint64
	// Number of retr reqs submitted
	retrvls_ops uint64
	// Number of retr reqs waited for CPU time
	retrvls_owt uint64
	// Number of retr reqs aborted due to object death
	retrvls_abt uint64
	// Number of storage (write) requests seen
	stores_n uint64
	// Number of successful store reqs
	stores_ok uint64
	// Number of store reqs on a page already pending storage
	stores_agn uint64
	// Number of store reqs rejected -ENOBUFS
	stores_nbf uint64
	// Number of store reqs failed -ENOMEM
	stores_oom uint64
	// Number of store reqs submitted
	stores_ops uint64
	// Number of store reqs granted CPU time
	stores_run uint64
	// Number of pages given store req processing time
	stores_pgs uint64
	// Number of store reqs deleted from tracking tree
	stores_rxd uint64
	// Number of store reqs over store limit
	stores_olm uint64
	// Number of release reqs against pages with no pending store
	vmscan_nos uint64
	// Number of release reqs against pages stored by time lock granted
	vmscan_gon uint64
	// Number of release reqs ignored due to in-progress store
	vmscan_bsy uint64
	// Number of page stores cancelled due to release req
	vmscan_can uint64
	vmscan_wt  uint64
	// Number of times async ops added to pending queues
	ops_pend uint64
	// Number of times async ops given CPU time
	ops_run uint64
	// Number of times async ops queued for processing
	ops_enq uint64
	// Number of async ops cancelled
	ops_can uint64
	// Number of async ops rejected due to object lookup/create failure
	ops_rej uint64
	// Number of async ops initialised
	ops_ini uint64
	// Number of async ops queued for deferred release
	ops_dfr uint64
	// Number of async ops released (should equal ini=N when idle)
	ops_rel uint64
	// Number of deferred-release async ops garbage collected
	ops_gc uint64
	// Number of in-progress alloc_object() cache ops
	cacheop_alo uint64
	// Number of in-progress lookup_object() cache ops
	cacheop_luo uint64
	// Number of in-progress lookup_complete() cache ops
	cacheop_luc uint64
	// Number of in-progress grab_object() cache ops
	cacheop_gro uint64
	cacheop_inv uint64
	// Number of in-progress update_object() cache ops
	cacheop_upo uint64
	// Number of in-progress drop_object() cache ops
	cacheop_dro uint64
	// Number of in-progress put_object() cache ops
	cacheop_pto uint64
	// Number of in-progress sync_cache() cache ops
	cacheop_syn uint64
	// Number of in-progress attr_changed() cache ops
	cacheop_atc uint64
	// Number of in-progress read_or_alloc_page() cache ops
	cacheop_rap uint64
	// Number of in-progress read_or_alloc_pages() cache ops
	cacheop_ras uint64
	// Number of in-progress allocate_page() cache ops
	cacheop_alp uint64
	// Number of in-progress allocate_pages() cache ops
	cacheop_als uint64
	// Number of in-progress write_page() cache ops
	cacheop_wrp uint64
	// Number of in-progress uncache_page() cache ops
	cacheop_ucp uint64
	// Number of in-progress dissociate_pages() cache ops
	cacheop_dsp uint64
	// Number of object lookups/creations rejected due to lack of space
	cacheev_nsp uint64
	// Number of stale objects deleted
	cacheev_stl uint64
	// Number of objects retired when relinquished
	cacheev_rtr uint64
	// Number of objects culled
	cacheev_cul uint64
}

// Meminfo returns an information about current kernel/system memory statistics.
// See https://www.kernel.org/doc/Documentation/filesystems/proc.txt
func (fs FS) Fscacheinfo() (Fscacheinfo, error) {
	b, err := util.ReadFileNoStat(fs.proc.Path("fs/fscache/stats"))
	if err != nil {
		return Fscacheinfo{}, err
	}

	m, err := parseFscacheinfo(bytes.NewReader(b))
	if err != nil {
		return Fscacheinfo{}, fmt.Errorf("failed to parse Fscacheinfo: %v", err)
	}

	return *m, nil
}

func parseFscacheinfo(r io.Reader) (*Fscacheinfo, error) {
	var m Fscacheinfo
	var err error
	s := bufio.NewScanner(r)
	for s.Scan() {
		// Each line has at least a name and value; we ignore the unit.
		fields := strings.Fields(s.Text())
		if len(fields) < 2 {
			return nil, fmt.Errorf("malformed Fscacheinfo line: %q", s.Text())
		}

		switch fields[0] {
		case "Cookies:":
			m.cookies_idx, err = strconv.ParseUint(strings.Split(fields[1], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.cookies_dat, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.cookies_spc, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
		case "Objects:":
			m.objects_alc, err = strconv.ParseUint(strings.Split(fields[1], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.objects_nal, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.objects_avl, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.objects_ded, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
		case "ChkAux":
			m.chkaux_non, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.chkaux_ok, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.chkaux_upd, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.chkaux_obs, err = strconv.ParseUint(strings.Split(fields[5], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
		case "Pages":
			m.pages_mrk, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.pages_unc, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
		case "Acquire:":
			m.acquire_n, err = strconv.ParseUint(strings.Split(fields[1], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.acquire_nul, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.acquire_noc, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.acquire_ok, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.acquire_nbf, err = strconv.ParseUint(strings.Split(fields[5], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.acquire_oom, err = strconv.ParseUint(strings.Split(fields[6], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
		case "Lookups:":
			m.lookups_n, err = strconv.ParseUint(strings.Split(fields[1], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.lookups_neg, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.lookups_pos, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.lookups_crt, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.lookups_tmo, err = strconv.ParseUint(strings.Split(fields[5], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
		case "Invals":
			m.invals_n, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.invals_run, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
		case "Updates:":
			m.updates_n, err = strconv.ParseUint(strings.Split(fields[1], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.updates_nul, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.updates_run, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
		case "Relinqs:":
			m.relinqs_n, err = strconv.ParseUint(strings.Split(fields[1], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.relinqs_nul, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.relinqs_wcr, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.relinqs_rtr, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
		case "AttrChg:":
			m.attrchg_n, err = strconv.ParseUint(strings.Split(fields[1], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.attrchg_ok, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.attrchg_nbf, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.attrchg_oom, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.attrchg_run, err = strconv.ParseUint(strings.Split(fields[5], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
		case "Allocs":
			if strings.Split(fields[2], "=")[0] == "n" {
				m.allocs_n, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.allocs_ok, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.allocs_wt, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.allocs_nbf, err = strconv.ParseUint(strings.Split(fields[5], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.allocs_int, err = strconv.ParseUint(strings.Split(fields[6], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
			} else {
				m.allocs_ops, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.allocs_owt, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.allocs_abt, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
			}
		case "Retrvls:":
			if strings.Split(fields[1], "=")[0] == "n" {
				m.retrvls_n, err = strconv.ParseUint(strings.Split(fields[1], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.retrvls_ok, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.retrvls_wt, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.retrvls_nod, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.retrvls_nbf, err = strconv.ParseUint(strings.Split(fields[5], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.retrvls_int, err = strconv.ParseUint(strings.Split(fields[6], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.retrvls_oom, err = strconv.ParseUint(strings.Split(fields[7], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
			} else {
				m.retrvls_ops, err = strconv.ParseUint(strings.Split(fields[1], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.retrvls_owt, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.retrvls_abt, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
			}
		case "Stores":
			if strings.Split(fields[2], "=")[0] == "n" {
				m.stores_n, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.stores_ok, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.stores_agn, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.stores_nbf, err = strconv.ParseUint(strings.Split(fields[5], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.stores_oom, err = strconv.ParseUint(strings.Split(fields[6], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
			} else {
				m.stores_ops, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.stores_run, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.stores_pgs, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.stores_rxd, err = strconv.ParseUint(strings.Split(fields[5], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.stores_olm, err = strconv.ParseUint(strings.Split(fields[6], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
			}
		case "VmScan":
			m.vmscan_nos, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.vmscan_gon, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.vmscan_bsy, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.vmscan_can, err = strconv.ParseUint(strings.Split(fields[5], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.vmscan_wt, err = strconv.ParseUint(strings.Split(fields[6], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
		case "Ops":
			if strings.Split(fields[2], "=")[0] == "pend" {
				m.ops_pend, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.ops_run, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.ops_enq, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.ops_can, err = strconv.ParseUint(strings.Split(fields[5], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.ops_rej, err = strconv.ParseUint(strings.Split(fields[6], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
			} else {
				m.ops_ini, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.ops_dfr, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.ops_rel, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.ops_gc, err = strconv.ParseUint(strings.Split(fields[5], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
			}
		case "CacheOp:":
			if strings.Split(fields[1], "=")[0] == "alo" {
				m.cacheop_alo, err = strconv.ParseUint(strings.Split(fields[1], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_luo, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_luc, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_gro, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
			} else if strings.Split(fields[1], "=")[0] == "inv" {
				m.cacheop_inv, err = strconv.ParseUint(strings.Split(fields[1], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_upo, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_dro, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_pto, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_atc, err = strconv.ParseUint(strings.Split(fields[5], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_syn, err = strconv.ParseUint(strings.Split(fields[6], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
			} else {
				m.cacheop_rap, err = strconv.ParseUint(strings.Split(fields[1], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_ras, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_alp, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_als, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_wrp, err = strconv.ParseUint(strings.Split(fields[5], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_ucp, err = strconv.ParseUint(strings.Split(fields[6], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_dsp, err = strconv.ParseUint(strings.Split(fields[7], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
			}
		case "CacheEv:":
			m.cacheev_nsp, err = strconv.ParseUint(strings.Split(fields[1], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.cacheev_stl, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.cacheev_rtr, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.cacheev_cul, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
		}
	}

	return &m, nil
}
