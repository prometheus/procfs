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

package sysfs

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/prometheus/procfs/internal/parsers"
)

// SystemUncoreFrequency contains data from one uncore frequency domain under
// `/sys/devices/system/cpu/intel_uncore_frequency/`.
//
// Two directory layouts exist and are both parsed:
//
//   - Legacy TPMI:  package_<NN>_die_<MM>/
//   - Per-domain:   uncore<NN>/   (adds current_freq_khz plus package_id and
//     domain_id files)
//
// Frequency fields are pointers so that a missing or unreadable file (for
// example current_freq_khz, which is root-only on many platforms) yields a nil
// value rather than a parse error. PackageID and DomainID are likewise nil when
// unavailable; callers may fall back to the numeric suffix of Name.
//
// See the kernel documentation for the sysfs ABI:
// https://www.kernel.org/doc/html/latest/admin-guide/pm/intel_uncore_frequency_scaling.html
type SystemUncoreFrequency struct {
	Name              string // directory name, e.g. "package_00_die_00" or "uncore00"
	PackageID         *uint64
	DomainID          *uint64
	CurrentFreqKhz    *uint64
	MinFreqKhz        *uint64
	MaxFreqKhz        *uint64
	InitialMinFreqKhz *uint64
	InitialMaxFreqKhz *uint64
}

var (
	uncoreLegacyDirRE = regexp.MustCompile(`^package_\d+_die_\d+$`)
	uncoreDomainDirRE = regexp.MustCompile(`^uncore\d+$`)
)

// SystemUncoreFrequency returns the Intel uncore frequency domains found under
// `/sys/devices/system/cpu/intel_uncore_frequency/`.
//
// If the directory does not exist (a platform without the intel_uncore_frequency
// driver) it returns an empty slice and a nil error.
func (fs FS) SystemUncoreFrequency() ([]SystemUncoreFrequency, error) {
	root := fs.sys.Path("devices/system/cpu/intel_uncore_frequency")
	entries, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read %q: %w", root, err)
	}

	var out []SystemUncoreFrequency
	for _, e := range entries {
		name := e.Name()
		if !uncoreLegacyDirRE.MatchString(name) && !uncoreDomainDirRE.MatchString(name) {
			continue
		}
		out = append(out, parseUncoreFrequency(name, filepath.Join(root, name)))
	}
	return out, nil
}

// parseUncoreFrequency parses a single uncore frequency domain directory. Each
// attribute is optional: a missing or unreadable file leaves the corresponding
// field nil so that privileged (root-only) or layout-specific files do not fail
// the whole parse.
func parseUncoreFrequency(name, path string) SystemUncoreFrequency {
	return SystemUncoreFrequency{
		Name:              name,
		PackageID:         sysReadOptionalUint(filepath.Join(path, "package_id")),
		DomainID:          sysReadOptionalUint(filepath.Join(path, "domain_id")),
		CurrentFreqKhz:    sysReadOptionalUint(filepath.Join(path, "current_freq_khz")),
		MinFreqKhz:        sysReadOptionalUint(filepath.Join(path, "min_freq_khz")),
		MaxFreqKhz:        sysReadOptionalUint(filepath.Join(path, "max_freq_khz")),
		InitialMinFreqKhz: sysReadOptionalUint(filepath.Join(path, "initial_min_freq_khz")),
		InitialMaxFreqKhz: sysReadOptionalUint(filepath.Join(path, "initial_max_freq_khz")),
	}
}

// sysReadOptionalUint reads a sysfs file containing a single unsigned integer
// via parsers.SysReadFile. It returns nil if the file is missing, unreadable
// (e.g. permission denied), or not a valid integer.
func sysReadOptionalUint(path string) *uint64 {
	s, err := parsers.SysReadFile(path)
	if err != nil {
		return nil
	}
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return nil
	}
	return &v
}
