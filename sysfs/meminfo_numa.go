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

//go:build linux
// +build linux

package sysfs

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/prometheus/procfs/internal/util"
)

// Meminfo represents memory statistics for NUMA node
type Meminfo struct {
	MemTotal        uint64
	MemFree         uint64
	MemUsed         uint64
	SwapCached      uint64
	Active          uint64
	Inactive        uint64
	ActiveAnon      uint64
	InactiveAnon    uint64
	ActiveFile      uint64
	InactiveFile    uint64
	Unevictable     uint64
	Mlocked         uint64
	Dirty           uint64
	Writeback       uint64
	FilePages       uint64
	Mapped          uint64
	AnonPages       uint64
	Shmem           uint64
	KernelStack     uint64
	PageTables      uint64
	NFS_Unstable    uint64
	Bounce          uint64
	WritebackTmp    uint64
	KReclaimable    uint64
	Slab            uint64
	SReclaimable    uint64
	SUnreclaim      uint64
	AnonHugePages   uint64
	ShmemHugePages  uint64
	ShmemPmdMapped  uint64
	FileHugePages   uint64
	FilePmdMapped   uint64
	HugePages_Total uint64
	HugePages_Free  uint64
	HugePages_Surp  uint64
}

func (fs FS) MeminfoNUMA() (map[int]Meminfo, error) {
	m := make(map[int]Meminfo)
	nodes, err := filepath.Glob(fs.sys.Path(nodePattern))
	if err != nil {
		return nil, err
	}

	for _, node := range nodes {
		nodeNumbers := nodeNumberRegexp.FindStringSubmatch(node)
		if len(nodeNumbers) != 2 {
			continue
		}
		nodeNumber, err := strconv.Atoi(nodeNumbers[1])
		if err != nil {
			return nil, err
		}
		b, err := util.ReadFileNoStat(filepath.Join(node, "meminfo"))
		if err != nil {
			return nil, err
		}
		meminfo, err := parseMeminfo(bytes.NewReader(b))
		if err != nil {
			return nil, err
		}
		m[nodeNumber] = meminfo
	}
	return m, nil
}

func parseMeminfo(r io.Reader) (Meminfo, error) {
	var m Meminfo
	s := bufio.NewScanner(r)
	for s.Scan() {
		// Each line has at least a name and value; we ignore the unit.
		// A line example: "Node 0 MemTotal:       395936028 kB"
		fields := strings.Fields(s.Text())
		if len(fields) < 4 {
			return Meminfo{}, fmt.Errorf("malformed meminfo line: %q", s.Text())
		}

		v, err := strconv.ParseUint(fields[3], 0, 64)
		if err != nil {
			return Meminfo{}, err
		}

		switch fields[2] {
		case "MemTotal:":
			m.MemTotal = v
		case "MemFree:":
			m.MemFree = v
		case "MemUsed:":
			m.MemUsed = v
		case "SwapCached:":
			m.SwapCached = v
		case "Active:":
			m.Active = v
		case "Inactive:":
			m.Inactive = v
		case "Active(anon):":
			m.ActiveAnon = v
		case "Inactive(anon):":
			m.InactiveAnon = v
		case "Active(file):":
			m.ActiveFile = v
		case "Inactive(file):":
			m.InactiveFile = v
		case "Unevictable:":
			m.Unevictable = v
		case "Mlocked:":
			m.Mlocked = v
		case "Dirty:":
			m.Dirty = v
		case "Writeback:":
			m.Writeback = v
		case "FilePages:":
			m.FilePages = v
		case "Mapped:":
			m.Mapped = v
		case "AnonPages:":
			m.AnonPages = v
		case "Shmem:":
			m.Shmem = v
		case "KernelStack:":
			m.KernelStack = v
		case "PageTables:":
			m.PageTables = v
		case "NFS_Unstable:":
			m.NFS_Unstable = v
		case "Bounce:":
			m.Bounce = v
		case "WritebackTmp:":
			m.WritebackTmp = v
		case "KReclaimable:":
			m.KReclaimable = v
		case "Slab:":
			m.Slab = v
		case "SReclaimable:":
			m.SReclaimable = v
		case "SUnreclaim:":
			m.SUnreclaim = v
		case "AnonHugePages:":
			m.AnonHugePages = v
		case "ShmemHugePages:":
			m.ShmemHugePages = v
		case "ShmemPmdMapped:":
			m.ShmemPmdMapped = v
		case "FileHugePages:":
			m.FileHugePages = v
		case "FilePmdMapped:":
			m.FilePmdMapped = v
		case "HugePages_Total:":
			m.HugePages_Total = v
		case "HugePages_Free:":
			m.HugePages_Free = v
		case "HugePages_Surp:":
			m.HugePages_Surp = v
		}
	}

	return m, nil
}
