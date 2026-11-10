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
	"path/filepath"
	"strconv"

	"github.com/prometheus/procfs/internal/parsers"
)

// NUMANodeCPUList returns the online CPUs belonging to each NUMA node, keyed
// by node number, as read from `/sys/devices/system/node/nodeN/cpulist`.
// The kernel only lists online CPUs in this file. Memory-only nodes have an
// empty CPU list.
// See https://www.kernel.org/doc/Documentation/ABI/stable/sysfs-devices-node
func (fs FS) NUMANodeCPUList() (map[int][]uint16, error) {
	nodes, err := filepath.Glob(fs.sys.Path(nodePattern))
	if err != nil {
		return nil, err
	}

	m := make(map[int][]uint16, len(nodes))
	for _, node := range nodes {
		nodeNumbers := nodeNumberRegexp.FindStringSubmatch(node)
		if len(nodeNumbers) != 2 {
			continue
		}
		nodeNumber, err := strconv.Atoi(nodeNumbers[1])
		if err != nil {
			return nil, err
		}
		cpulist, err := parsers.ReadFileNoStat(filepath.Join(node, "cpulist"))
		if err != nil {
			return nil, err
		}
		cpus, err := parseCPURange(cpulist)
		if err != nil {
			return nil, fmt.Errorf("failed to parse cpulist of NUMA node %d: %w", nodeNumber, err)
		}
		m[nodeNumber] = cpus
	}
	return m, nil
}
