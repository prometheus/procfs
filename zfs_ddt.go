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

package procfs

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/prometheus/procfs/internal/parsers"
)

// constants from https://github.com/openzfs/zfs/blob/master/lib/libspl/include/sys/kstat.h
// Kept as strings for comparison thus avoiding conversion to int.
const (
	kstatDataUint32 = "2"
	kstatDataUint64 = "4"
)

// ZFSDDTStats holds OpenZFS dedup table statistics for one pool and checksum,
// read from /proc/spl/kstat/zfs/<pool>/ddt_stats_<checksum> (OpenZFS 2.3+).
// See https://github.com/openzfs/zfs/blob/master/module/zfs/ddt.c (ddt_kstats_template).
type ZFSDDTStats struct {
	Pool     string
	Checksum string

	Lookup               uint64
	LookupNew            uint64
	LookupExisting       uint64
	LookupLiveHit        uint64
	LookupLiveWait       uint64
	LookupLiveMiss       uint64
	LookupLogHit         uint64
	LookupLogActiveHit   uint64
	LookupLogFlushingHit uint64
	LookupLogMiss        uint64
	LookupStoredHit      uint64
	LookupStoredMiss     uint64
	LogActiveEntries     uint64
	LogFlushingEntries   uint64
	LogIngestRate        uint64
	LogFlushRate         uint64
	LogFlushTimeRate     uint64
}

// ZFSDDTStats reads /proc/spl/kstat/zfs/*/ddt_stats_*. It returns an empty
// slice if no such files exist (OpenZFS older than 2.3, or no ZFS).
func (fs FS) ZFSDDTStats() ([]ZFSDDTStats, error) {
	paths, err := filepath.Glob(fs.proc.Path("spl/kstat/zfs/*/ddt_stats_*"))
	if err != nil {
		return nil, err
	}

	// Glob returns nil on no match, so older ZFS or no dedup gives an empty result for free.
	var stats []ZFSDDTStats
	for _, path := range paths {
		data, err := parsers.ReadFileNoStat(path)
		if err != nil {
			return nil, err
		}

		pool := filepath.Base(filepath.Dir(path))
		checksum := strings.TrimPrefix(filepath.Base(path), "ddt_stats_")

		stat, err := parseZFSDDTStats(pool, checksum, data)
		if err != nil {
			return nil, err
		}

		stats = append(stats, stat)
	}

	return stats, nil
}

func parseZFSDDTStats(pool, checksum string, data []byte) (ZFSDDTStats, error) {
	var stat ZFSDDTStats
	stat.Pool = pool
	stat.Checksum = checksum
	inTable := false

	// The kstat data is in the form of a table with 3 columns: name, type, value.
	lines := bytes.Split(data, []byte{'\n'})
	for _, line := range lines {
		fields := bytes.Fields(line)

		if !inTable && len(fields) == 3 && bytes.Equal(fields[0], []byte("name")) && bytes.Equal(fields[1], []byte("type")) && bytes.Equal(fields[2], []byte("data")) {
			inTable = true
			continue
		}

		if !inTable || len(fields) != 3 {
			continue
		}

		name := string(fields[0])
		var value uint64
		var err error

		switch string(fields[1]) {
		case kstatDataUint32, kstatDataUint64:
			value, err = strconv.ParseUint(string(fields[2]), 10, 64)
			if err != nil {
				return ZFSDDTStats{}, fmt.Errorf("%w: %s/%s: couldn't parse %s=%q: %w", ErrFileParse, pool, checksum, name, fields[2], err)
			}
		default:
			continue
		}

		switch name {
		case "lookup":
			stat.Lookup = value
		case "lookup_new":
			stat.LookupNew = value
		case "lookup_existing":
			stat.LookupExisting = value
		case "lookup_live_hit":
			stat.LookupLiveHit = value
		case "lookup_live_wait":
			stat.LookupLiveWait = value
		case "lookup_live_miss":
			stat.LookupLiveMiss = value
		case "lookup_log_hit":
			stat.LookupLogHit = value
		case "lookup_log_active_hit":
			stat.LookupLogActiveHit = value
		case "lookup_log_flushing_hit":
			stat.LookupLogFlushingHit = value
		case "lookup_log_miss":
			stat.LookupLogMiss = value
		case "lookup_stored_hit":
			stat.LookupStoredHit = value
		case "lookup_stored_miss":
			stat.LookupStoredMiss = value
		case "log_active_entries":
			stat.LogActiveEntries = value
		case "log_flushing_entries":
			stat.LogFlushingEntries = value
		case "log_ingest_rate":
			stat.LogIngestRate = value
		case "log_flush_rate":
			stat.LogFlushRate = value
		case "log_flush_time_rate":
			stat.LogFlushTimeRate = value
		default:
			// Ignore unknown fields, as they may be added in future versions of ZFS.
		}
	}

	if !inTable {
		return ZFSDDTStats{}, fmt.Errorf("%w: %s/%s: missing name/type/data header", ErrFileParse, pool, checksum)
	}

	return stat, nil
}
