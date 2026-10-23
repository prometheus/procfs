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

//go:build linux

package resctrlfs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/prometheus/procfs/internal/parsers"
)

const monDataPath = "mon_data"

// MonData holds the monitoring counters of one domain of one resource,
// read from mon_data/mon_<resource>_<id>.
type MonData struct {
	Resource string // e.g. "L3"
	ID       string // domain id with leading zeros stripped, e.g. "0"
	// Counters are nil when the file is absent (feature not supported)
	// or reads "Unavailable" (no sample available right now). For example
	// Arm MPAM never offers MBMLocalBytes, because the hardware can not
	// tell local from remote traffic.
	LLCOccupancy  *uint64
	MBMTotalBytes *uint64
	MBMLocalBytes *uint64
}

// MonData returns the monitoring counters of the root control group, one entry
// per monitoring domain. A domain is a socket, or a CCX on AMD.
//
// It errors only if the mon_data directory itself can't be read. A domain
// counter that is missing, unreadable or unavailable stays nil, because the
// hardware may support only some of the features and samples are not always
// available. A mount without any monitoring domain returns an empty slice.
func (fs FS) MonData() ([]MonData, error) {
	path := fs.resctrl.Path(monDataPath)

	dirs, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list monitoring domains at %q: %w", path, err)
	}

	// os.ReadDir sorts by filename, so the result is ordered by domain.
	mons := make([]MonData, 0, len(dirs))
	for _, d := range dirs {
		resource, id, ok := parseMonDirName(d.Name())
		if !ok {
			continue
		}

		domain := fs.resctrl.Path(monDataPath, d.Name())
		mons = append(mons, MonData{
			Resource:     resource,
			ID:           id,
			LLCOccupancy: readCounter(domain, "llc_occupancy"),
			// The mbm counters are hardware registers of limited width. They
			// are free running and wrap around; this library reports the raw
			// value and leaves the wrap handling to the caller.
			MBMTotalBytes: readCounter(domain, "mbm_total_bytes"),
			MBMLocalBytes: readCounter(domain, "mbm_local_bytes"),
		})
	}

	return mons, nil
}

// readCounter reads a single counter file of a monitoring domain. It returns
// nil if the file does not exist, can not be read, or does not hold a number.
// The kernel writes the literal "Unavailable" when the hardware can not
// deliver a sample, which is a normal and transient state.
func readCounter(domain, name string) *uint64 {
	value, err := parsers.ReadUintFromFile(filepath.Join(domain, name))
	if err != nil {
		return nil
	}
	return &value
}

// parseMonDirName splits a mon_data directory name into resource and domain id.
// The name is "mon_<resource>_<id>", for example "mon_L3_00". The resource name
// may itself contain underscores, so the split is at the last one. The id is
// zero padded by the kernel; the padding is stripped, so "00" becomes "0".
// The second return value is false for a name that is not a monitoring domain.
func parseMonDirName(name string) (string, string, bool) {
	rest, ok := strings.CutPrefix(name, "mon_")
	if !ok {
		return "", "", false
	}

	i := strings.LastIndex(rest, "_")
	if i <= 0 {
		return "", "", false
	}

	resource, id := rest[:i], rest[i+1:]
	if !isDigits(id) {
		return "", "", false
	}

	id = strings.TrimLeft(id, "0")
	if id == "" {
		id = "0"
	}

	return resource, id, true
}

// isDigits reports whether s is a non-empty run of decimal digits.
func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
