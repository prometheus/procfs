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
	"github.com/prometheus/procfs/internal/util"
	"strconv"
	"strings"
)

// Cgroup models one line from /proc/[pid]/cgroup. Each Cgroup struct describes the the placement of a PID inside a
// specific control hierarchy. The kernel has two cgroup APIs, v1 and v2. v1 has one hierarchy per available resource
// controller, while v2 has one unified hierarchy shared by all controllers. Regardless of v1 or v2, all hierarchies
// contain all running processes, so the question answerable with a Cgroup struct is 'where is this process in
// this hierarchy' (where==what path). By prefixing this path with the mount point of this hierarchy, you can locate
// the relevant pseudo-files needed to read/set the data for this PID in this hierarchy
//
// Also see http://man7.org/linux/man-pages/man7/cgroups.7.html
type Cgroup struct {
	// HierarchyId for cgroups V2 is always 0. For cgroups v1 this is a unique
	// ID number that can be matched to a hierarchy ID found in /proc/cgroups
	HierarchyId int
	// Controllers using this hierarchy of processes. Controllers are also known as subsystems.
	// TODO is it a problem that len(Controllers)==1 even when there are no controllers?
	Controllers []string
	// Path of this control group, relative to the mount point of the various controllers
	Path string
}

// parseCgroupString parses each line of the /proc/[pid]/cgroup file
func parseCgroupString(cgroupStr string) (*Cgroup, error) {
	var err error

	fields := strings.Split(cgroupStr, ":")
	if len(fields) != 3 {
		return nil, fmt.Errorf("incorrect number of fields (%d) in cgroup string: %s", len(fields), cgroupStr)
	}

	cgroup := &Cgroup{
		Path: fields[2],
	}
	cgroup.HierarchyId, err = strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse hierarchy ID")
	}
	// fields[1] may be ""
	ssNames := strings.Split(fields[1], ",")
	cgroup.Controllers = append(cgroup.Controllers, ssNames...)

	return cgroup, nil
}

// parseCgroups reads each line of the /proc/[pid]/cgroup file
func parseCgroups(data []byte) ([]*Cgroup, error) {
	cgroups := []*Cgroup{}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		mountString := scanner.Text()
		parsedMounts, err := parseCgroupString(mountString)
		if err != nil {
			return nil, err
		}
		cgroups = append(cgroups, parsedMounts)
	}

	err := scanner.Err()
	return cgroups, err
}

// GetCgroups returns a Cgroup struct for all process control hierarchies running on this system
func GetCgroups(pid int) ([]*Cgroup, error) {
	data, err := util.ReadFileNoStat(fmt.Sprintf("/proc/%d/cgroup", pid))
	if err != nil {
		return nil, err
	}
	return parseCgroups(data)
}
