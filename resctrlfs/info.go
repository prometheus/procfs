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
	"strings"

	"github.com/prometheus/procfs/internal/util"
)

const l3MonInfoPath = "info/L3_MON"

// L3MonInfo describes the L3 monitoring capabilities from info/L3_MON.
type L3MonInfo struct {
	// NumRMIDs is the number of resource monitoring ids the hardware offers.
	// One RMID is used per monitoring group, so this is the upper limit of
	// groups that can be monitored at the same time.
	NumRMIDs uint64
	// MonFeatures lists the counters the hardware supports, in the order the
	// kernel reports them, e.g. "llc_occupancy", "mbm_total_bytes",
	// "mbm_local_bytes".
	MonFeatures []string
}

// L3MonInfo returns the L3 monitoring capabilities of the CPU.
//
// It errors if info/L3_MON is missing, which means the CPU or the kernel does
// not support L3 monitoring. The error wraps the underlying os error, so it can
// be tested with os.IsNotExist.
func (fs FS) L3MonInfo() (L3MonInfo, error) {
	var info L3MonInfo

	numRMIDs, err := util.ReadUintFromFile(fs.resctrl.Path(l3MonInfoPath, "num_rmids"))
	if err != nil {
		return info, err
	}
	info.NumRMIDs = numRMIDs

	data, err := util.ReadFileNoStat(fs.resctrl.Path(l3MonInfoPath, "mon_features"))
	if err != nil {
		return info, err
	}

	for _, line := range strings.Split(string(data), "\n") {
		feature := strings.TrimSpace(line)
		if feature == "" {
			continue
		}
		info.MonFeatures = append(info.MonFeatures, feature)
	}

	return info, nil
}
