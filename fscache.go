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
	index_cookies_allocated uint64
	// data storage cookies allocated
	data_storage_cookies_allocated uint64
	// Number of special cookies allocated
	special_cookies_allocated uint64
	// Number of objects allocated
	objects_allocated uint64
	// Number of object allocation failures
	object_allocations_failure uint64
	// Number of objects that reached the available state
	objects_available uint64
	// Number of objects that reached the dead state
	objects_dead uint64
	// Number of objects that didn't have a coherency check
	objects_without_coherency_check uint64
	// Number of objects that passed a coherency check
	objects_with_coherency_check uint64
	// Number of objects that needed a coherency data update
	objects_need_coherency_check_update uint64
	// Number of objects that were declared obsolete
	objects_declared_obsolete uint64
	// Number of pages marked as being cached
	pages_markes_as_being_cached uint64
	// Number of uncache page requests seen
	uncache_pages_request_seen uint64
	// Number of acquire cookie requests seen
	acquire_cookies_requests_seen uint64
	// Number of acq reqs given a NULL parent
	acquire_requests_with_null_parent uint64
	// Number of acq reqs rejected due to no cache available
	acquire_requests_rejected_no_cache_available uint64
	// Number of acq reqs succeeded
	acquire_requests_succeeded uint64
	// Number of acq reqs rejected due to error
	acquire_requests_rejected_due_to_error uint64
	// Number of acq reqs failed on ENOMEM
	acquire_requests_failed_due_to_enomem uint64
	// Number of lookup calls made on cache backends
	lookups_number uint64
	// Number of negative lookups made
	lookups_negative uint64
	// Number of positive lookups made
	lookups_positive uint64
	// Number of objects created by lookup
	objects_created_by_lookup uint64
	// Number of lookups timed out and requeued
	lookups_timed_out_and_requeued uint64
	invalidations_number           uint64
	invals_running                 uint64
	// Number of update cookie requests seen
	update_cookie_request_seen uint64
	// Number of upd reqs given a NULL parent
	update_requests_with_null_parent uint64
	// Number of upd reqs granted CPU time
	update_requests_running uint64
	// Number of relinquish cookie requests seen
	relinquish_cookies_request_seen uint64
	// Number of rlq reqs given a NULL parent
	relinquish_cookies_with_null_parent uint64
	// Number of rlq reqs waited on completion of creation
	relinquish_requests_waiting_compete_creation uint64
	// Relinqs rtr
	relinqs_retries uint64
	// Number of attribute changed requests seen
	attribute_changed_requests_seen uint64
	// Number of attr changed requests queued
	attribute_changed_requests_queued uint64
	// Number of attr changed rejected -ENOBUFS
	attribute_changed_reject_due_to_enobufs uint64
	// Number of attr changed failed -ENOMEM
	attribute_changed_failed_due_to_enomem uint64
	// Number of attr changed ops given CPU time
	attribute_changed_ops uint64
	// Number of allocation requests seen
	allocations_requests_seen uint64
	// Number of successful alloc reqs
	allocations_ok_requests uint64
	// Number of alloc reqs that waited on lookup completion
	allocations_waiting_on_lookup uint64
	// Number of alloc reqs rejected -ENOBUFS
	allocations_rejected_due_to_enobufs uint64
	// Number of alloc reqs aborted -ERESTARTSYS
	allocations_aborted_due_to_erestartsys uint64
	// Number of alloc reqs submitted
	allocation_operations_submitted uint64
	// Number of alloc reqs waited for CPU time
	allocations_waited_for_cpu uint64
	// Number of alloc reqs aborted due to object death
	allocations_aborted_due_to_object_dead uint64
	// Number of retrieval (read) requests seen
	retrievals_read_requests uint64
	// Number of successful retr reqs
	retrievals_ok uint64
	// Number of retr reqs that waited on lookup completion
	retrievals_waiting_lookup_completion uint64
	// Number of retr reqs returned -ENODATA
	retrievals_returned_enodata uint64
	// Number of retr reqs rejected -ENOBUFS
	retrievals_rejected_due_to_enobufs uint64
	// Number of retr reqs aborted -ERESTARTSYS
	retrievals_aborted_due_to_erestartsys uint64
	// Number of retr reqs failed -ENOMEM
	retrievals_failed_due_to_enomem uint64
	// Number of retr reqs submitted
	retrievals_requests uint64
	// Number of retr reqs waited for CPU time
	retrievals_waiting_cpu uint64
	// Number of retr reqs aborted due to object death
	retrievals_aborted_due_to_object_dead uint64
	// Number of storage (write) requests seen
	store_write_requests uint64
	// Number of successful store reqs
	store_successful_requests uint64
	// Number of store reqs on a page already pending storage
	store_requests_on_pending_storage uint64
	// Number of store reqs rejected -ENOBUFS
	store_requests_rejected_due_to_enobufs uint64
	// Number of store reqs failed -ENOMEM
	store_requests_failed_due_to_enomem uint64
	// Number of store reqs submitted
	store_requests_submitted uint64
	// Number of store reqs granted CPU time
	store_requests_running uint64
	// Number of pages given store req processing time
	store_pages_with_requests_processing uint64
	// Number of store reqs deleted from tracking tree
	store_requests_deleted uint64
	// Number of store reqs over store limit
	store_requests_over_store_limit uint64
	// Number of release reqs against pages with no pending store
	release_requests_against_pages_with_no_pending_storage uint64
	// Number of release reqs against pages stored by time lock granted
	release_requests_against_pages_stored_by_time_lock_granted uint64
	// Number of release reqs ignored due to in-progress store
	release_requests_ignored_due_to_inprogress_store uint64
	// Number of page stores cancelled due to release req
	page_stores_cancelled_by_release_requests uint64
	vmscan_waiting                            uint64
	// Number of times async ops added to pending queues
	ops_pending uint64
	// Number of times async ops given CPU time
	ops_running uint64
	// Number of times async ops queued for processing
	ops_enqueued uint64
	// Number of async ops cancelled
	ops_cancelled uint64
	// Number of async ops rejected due to object lookup/create failure
	ops_rejected uint64
	// Number of async ops initialised
	ops_initialised uint64
	// Number of async ops queued for deferred release
	ops_deferred uint64
	// Number of async ops released (should equal ini=N when idle)
	ops_released uint64
	// Number of deferred-release async ops garbage collected
	ops_garbage_collected uint64
	// Number of in-progress alloc_object() cache ops
	cacheop_allocations_in_progress uint64
	// Number of in-progress lookup_object() cache ops
	cacheop_lookup_object_in_progress uint64
	// Number of in-progress lookup_complete() cache ops
	cacheop_lookup_complete_in_progress uint64
	// Number of in-progress grab_object() cache ops
	cacheop_grab_object_in_progress uint64
	cacheop_invalidations           uint64
	// Number of in-progress update_object() cache ops
	cacheop_update_object_in_progress uint64
	// Number of in-progress drop_object() cache ops
	cacheop_drop_object_in_progress uint64
	// Number of in-progress put_object() cache ops
	cacheop_put_object_in_progress uint64
	// Number of in-progress attr_changed() cache ops
	cacheop_attribute_change_in_progress uint64
	// Number of in-progress sync_cache() cache ops
	cacheop_sync_cache_in_progress uint64
	// Number of in-progress read_or_alloc_page() cache ops
	cacheop_read_or_alloc_page_in_progress uint64
	// Number of in-progress read_or_alloc_pages() cache ops
	cacheop_read_or_alloc_pages_in_progress uint64
	// Number of in-progress allocate_page() cache ops
	cacheop_allocate_page_in_progress uint64
	// Number of in-progress allocate_pages() cache ops
	cacheop_allocate_pages_in_progress uint64
	// Number of in-progress write_page() cache ops
	cacheop_write_pages_in_progress uint64
	// Number of in-progress uncache_page() cache ops
	cacheop_uncache_pages_in_progress uint64
	// Number of in-progress dissociate_pages() cache ops
	cacheop_dissociate_pages_in_progress uint64
	// Number of object lookups/creations rejected due to lack of space
	cacheev_lookups_and_creations_rejected_lack_space uint64
	// Number of stale objects deleted
	cacheev_stale_objects_deleted uint64
	// Number of objects retired when relinquished
	cacheev_retired_when_relinquished uint64
	// Number of objects culled
	cacheev_objects_culled uint64
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
			m.index_cookies_allocated, err = strconv.ParseUint(strings.Split(fields[1], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.data_storage_cookies_allocated, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.special_cookies_allocated, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
		case "Objects:":
			m.objects_allocated, err = strconv.ParseUint(strings.Split(fields[1], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.object_allocations_failure, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.objects_available, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.objects_dead, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
		case "ChkAux":
			m.objects_without_coherency_check, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.objects_with_coherency_check, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.objects_need_coherency_check_update, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.objects_declared_obsolete, err = strconv.ParseUint(strings.Split(fields[5], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
		case "Pages":
			m.pages_markes_as_being_cached, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.uncache_pages_request_seen, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
		case "Acquire:":
			m.acquire_cookies_requests_seen, err = strconv.ParseUint(strings.Split(fields[1], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.acquire_requests_with_null_parent, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.acquire_requests_rejected_no_cache_available, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.acquire_requests_succeeded, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.acquire_requests_rejected_due_to_error, err = strconv.ParseUint(strings.Split(fields[5], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.acquire_requests_failed_due_to_enomem, err = strconv.ParseUint(strings.Split(fields[6], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
		case "Lookups:":
			m.lookups_number, err = strconv.ParseUint(strings.Split(fields[1], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.lookups_negative, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.lookups_positive, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.objects_created_by_lookup, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.lookups_timed_out_and_requeued, err = strconv.ParseUint(strings.Split(fields[5], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
		case "Invals":
			m.invalidations_number, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.invals_running, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
		case "Updates:":
			m.update_cookie_request_seen, err = strconv.ParseUint(strings.Split(fields[1], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.update_requests_with_null_parent, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.update_requests_running, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
		case "Relinqs:":
			m.relinquish_cookies_request_seen, err = strconv.ParseUint(strings.Split(fields[1], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.relinquish_cookies_with_null_parent, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.relinquish_requests_waiting_compete_creation, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.relinqs_retries, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
		case "AttrChg:":
			m.attribute_changed_requests_seen, err = strconv.ParseUint(strings.Split(fields[1], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.attribute_changed_requests_queued, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.attribute_changed_reject_due_to_enobufs, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.attribute_changed_failed_due_to_enomem, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.attribute_changed_ops, err = strconv.ParseUint(strings.Split(fields[5], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
		case "Allocs":
			if strings.Split(fields[2], "=")[0] == "n" {
				m.allocations_requests_seen, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.allocations_ok_requests, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.allocations_waiting_on_lookup, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.allocations_rejected_due_to_enobufs, err = strconv.ParseUint(strings.Split(fields[5], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.allocations_aborted_due_to_erestartsys, err = strconv.ParseUint(strings.Split(fields[6], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
			} else {
				m.allocation_operations_submitted, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.allocations_waited_for_cpu, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.allocations_aborted_due_to_object_dead, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
			}
		case "Retrvls:":
			if strings.Split(fields[1], "=")[0] == "n" {
				m.retrievals_read_requests, err = strconv.ParseUint(strings.Split(fields[1], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.retrievals_ok, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.retrievals_waiting_lookup_completion, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.retrievals_returned_enodata, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.retrievals_rejected_due_to_enobufs, err = strconv.ParseUint(strings.Split(fields[5], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.retrievals_aborted_due_to_erestartsys, err = strconv.ParseUint(strings.Split(fields[6], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.retrievals_failed_due_to_enomem, err = strconv.ParseUint(strings.Split(fields[7], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
			} else {
				m.retrievals_requests, err = strconv.ParseUint(strings.Split(fields[1], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.retrievals_waiting_cpu, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.retrievals_aborted_due_to_object_dead, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
			}
		case "Stores":
			if strings.Split(fields[2], "=")[0] == "n" {
				m.store_write_requests, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.store_successful_requests, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.store_requests_on_pending_storage, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.store_requests_rejected_due_to_enobufs, err = strconv.ParseUint(strings.Split(fields[5], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.store_requests_failed_due_to_enomem, err = strconv.ParseUint(strings.Split(fields[6], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
			} else {
				m.store_requests_submitted, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.store_requests_running, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.store_pages_with_requests_processing, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.store_requests_deleted, err = strconv.ParseUint(strings.Split(fields[5], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.store_requests_over_store_limit, err = strconv.ParseUint(strings.Split(fields[6], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
			}
		case "VmScan":
			m.release_requests_against_pages_with_no_pending_storage, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.release_requests_against_pages_stored_by_time_lock_granted, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.release_requests_ignored_due_to_inprogress_store, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.page_stores_cancelled_by_release_requests, err = strconv.ParseUint(strings.Split(fields[5], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.vmscan_waiting, err = strconv.ParseUint(strings.Split(fields[6], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
		case "Ops":
			if strings.Split(fields[2], "=")[0] == "pend" {
				m.ops_pending, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.ops_running, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.ops_enqueued, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.ops_cancelled, err = strconv.ParseUint(strings.Split(fields[5], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.ops_rejected, err = strconv.ParseUint(strings.Split(fields[6], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
			} else {
				m.ops_initialised, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.ops_deferred, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.ops_released, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.ops_garbage_collected, err = strconv.ParseUint(strings.Split(fields[5], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
			}
		case "CacheOp:":
			if strings.Split(fields[1], "=")[0] == "alo" {
				m.cacheop_allocations_in_progress, err = strconv.ParseUint(strings.Split(fields[1], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_lookup_object_in_progress, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_lookup_complete_in_progress, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_grab_object_in_progress, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
			} else if strings.Split(fields[1], "=")[0] == "inv" {
				m.cacheop_invalidations, err = strconv.ParseUint(strings.Split(fields[1], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_update_object_in_progress, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_drop_object_in_progress, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_put_object_in_progress, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_attribute_change_in_progress, err = strconv.ParseUint(strings.Split(fields[5], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_sync_cache_in_progress, err = strconv.ParseUint(strings.Split(fields[6], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
			} else {
				m.cacheop_read_or_alloc_page_in_progress, err = strconv.ParseUint(strings.Split(fields[1], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_read_or_alloc_pages_in_progress, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_allocate_page_in_progress, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_allocate_pages_in_progress, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_write_pages_in_progress, err = strconv.ParseUint(strings.Split(fields[5], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_uncache_pages_in_progress, err = strconv.ParseUint(strings.Split(fields[6], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
				m.cacheop_dissociate_pages_in_progress, err = strconv.ParseUint(strings.Split(fields[7], "=")[1], 0, 64)
				if err != nil {
					return &m, err
				}
			}
		case "CacheEv:":
			m.cacheev_lookups_and_creations_rejected_lack_space, err = strconv.ParseUint(strings.Split(fields[1], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.cacheev_stale_objects_deleted, err = strconv.ParseUint(strings.Split(fields[2], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.cacheev_retired_when_relinquished, err = strconv.ParseUint(strings.Split(fields[3], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
			m.cacheev_objects_culled, err = strconv.ParseUint(strings.Split(fields[4], "=")[1], 0, 64)
			if err != nil {
				return &m, err
			}
		}
	}

	return &m, nil
}
