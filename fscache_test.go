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
		index_cookies_allocated:                                    3,
		data_storage_cookies_allocated:                             67877,
		special_cookies_allocated:                                  0,
		objects_allocated:                                          67473,
		object_allocations_failure:                                 0,
		objects_available:                                          67473,
		objects_dead:                                               388,
		objects_without_coherency_check:                            12,
		objects_with_coherency_check:                               33,
		objects_need_coherency_check_update:                        44,
		objects_declared_obsolete:                                  55,
		pages_markes_as_being_cached:                               547164,
		uncache_pages_request_seen:                                 364577,
		acquire_cookies_requests_seen:                              67880,
		acquire_requests_with_null_parent:                          98,
		acquire_requests_rejected_no_cache_available:               25,
		acquire_requests_succeeded:                                 67780,
		acquire_requests_rejected_due_to_error:                     39,
		acquire_requests_failed_due_to_enomem:                      26,
		lookups_number:                                             67473,
		lookups_negative:                                           67470,
		lookups_positive:                                           58,
		objects_created_by_lookup:                                  67473,
		lookups_timed_out_and_requeued:                             85,
		invalidations_number:                                       14,
		invals_running:                                             13,
		update_cookie_request_seen:                                 7,
		update_requests_with_null_parent:                           3,
		update_requests_running:                                    8,
		relinquish_cookies_request_seen:                            394,
		relinquish_cookies_with_null_parent:                        1,
		relinquish_requests_waiting_compete_creation:               2,
		relinqs_retries:                                            3,
		attribute_changed_requests_seen:                            6,
		attribute_changed_requests_queued:                          5,
		attribute_changed_reject_due_to_enobufs:                    4,
		attribute_changed_failed_due_to_enomem:                     3,
		attribute_changed_ops:                                      2,
		allocations_requests_seen:                                  20,
		allocations_ok_requests:                                    19,
		allocations_waiting_on_lookup:                              18,
		allocations_rejected_due_to_enobufs:                        17,
		allocations_aborted_due_to_erestartsys:                     16,
		allocation_operations_submitted:                            15,
		allocations_waited_for_cpu:                                 14,
		allocations_aborted_due_to_object_dead:                     13,
		retrievals_read_requests:                                   151959,
		retrievals_ok:                                              82823,
		retrievals_waiting_lookup_completion:                       23467,
		retrievals_returned_enodata:                                69136,
		retrievals_rejected_due_to_enobufs:                         15,
		retrievals_aborted_due_to_erestartsys:                      69,
		retrievals_failed_due_to_enomem:                            43,
		retrievals_requests:                                        151959,
		retrievals_waiting_cpu:                                     42747,
		retrievals_aborted_due_to_object_dead:                      44,
		store_write_requests:                                       225565,
		store_successful_requests:                                  225565,
		store_requests_on_pending_storage:                          12,
		store_requests_rejected_due_to_enobufs:                     13,
		store_requests_failed_due_to_enomem:                        14,
		store_requests_submitted:                                   69156,
		store_requests_running:                                     294721,
		store_pages_with_requests_processing:                       225565,
		store_requests_deleted:                                     225565,
		store_requests_over_store_limit:                            43,
		release_requests_against_pages_with_no_pending_storage:     364512,
		release_requests_against_pages_stored_by_time_lock_granted: 2,
		release_requests_ignored_due_to_inprogress_store:           43,
		page_stores_cancelled_by_release_requests:                  12,
		vmscan_waiting:                                             66,
		ops_pending:                                                42753,
		ops_running:                                                221129,
		ops_enqueued:                                               628798,
		ops_cancelled:                                              11,
		ops_rejected:                                               88,
		ops_initialised:                                            377538,
		ops_deferred:                                               27,
		ops_released:                                               377538,
		ops_garbage_collected:                                      37,
		cacheop_allocations_in_progress:                            1,
		cacheop_lookup_object_in_progress:                          2,
		cacheop_lookup_complete_in_progress:                        3,
		cacheop_grab_object_in_progress:                            4,
		cacheop_invalidations:                                      5,
		cacheop_update_object_in_progress:                          6,
		cacheop_drop_object_in_progress:                            7,
		cacheop_put_object_in_progress:                             8,
		cacheop_attribute_change_in_progress:                       9,
		cacheop_sync_cache_in_progress:                             10,
		cacheop_read_or_alloc_page_in_progress:                     11,
		cacheop_read_or_alloc_pages_in_progress:                    12,
		cacheop_allocate_page_in_progress:                          13,
		cacheop_allocate_pages_in_progress:                         14,
		cacheop_write_pages_in_progress:                            15,
		cacheop_uncache_pages_in_progress:                          16,
		cacheop_dissociate_pages_in_progress:                       17,
		cacheev_lookups_and_creations_rejected_lack_space:          18,
		cacheev_stale_objects_deleted:                              19,
		cacheev_retired_when_relinquished:                          20,
		cacheev_objects_culled:                                     21,
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
