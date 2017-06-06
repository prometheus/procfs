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

// Package bcache provides access to statistics exposed by the bcache (Linux
// block cache).
package bcache

// Stats contains bcache runtime statistics, parsed from /sys/fs/bcache/.
//
// The names and meanings of each statistic were taken from bcache.txt and
// files in drivers/md/bcache in the Linux kernel source. Counters are float64
// (in-kernel counters are mostly unsigned long).
type Stats struct {
	// The name of the bcache used to source these statistics.
	Name   string
	Bcache BcacheStats
	Bdevs  []BdevStats
	Caches []CacheStats
}

// BcacheStats contains statistics tied to a bcache ID.
type BcacheStats struct {
	AverageKeySize        float64
	BtreeCacheSize        float64
	CacheAvailablePercent float64
	Congested             float64
	RootUsagePercent      float64
	TreeDepth             float64
	Internal              InternalStats
	FiveMin               PeriodStats
	Total                 PeriodStats
}

// BdevStats contains statistics for one backing device.
type BdevStats struct {
	Name      string
	DirtyData float64
	FiveMin   PeriodStats
	Total     PeriodStats
}

// CacheStats contains statistics for one cache device.
type CacheStats struct {
	Name            string
	IOErrors        float64
	MetadataWritten float64
	Written         float64
	Priority        PriorityStats
}

// PriorityStats contains statistics from the priority_stats file.
type PriorityStats struct {
	UnusedPercent   float64
	MetadataPercent float64
}

// InternalStats contains internal bcache statistics.
type InternalStats struct {
	ActiveJournalEntries                float64
	BtreeNodes                          float64
	BtreeReadAverageDurationNanoSeconds float64
	CacheReadRaces                      float64
}

// PeriodStats contains statistics for a time period (5 min or total).
type PeriodStats struct {
	Bypassed            float64
	CacheBypassHits     float64
	CacheBypassMisses   float64
	CacheHits           float64
	CacheMissCollisions float64
	CacheMisses         float64
	CacheReadaheads     float64
}
