// Copyright 2018 The Prometheus Authors
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

// +build !windows

package procfs

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type ProcSMaps []*ProcSMap

type ProcSMap struct {
	*ProcMap
	// Size of the mapping
	Size uint64
	// Amount of the mapping that is currently resident in RAM
	Rss uint64
	// Process's proportional share of this mapping
	Pss uint64
	// Size in bytes of clean shared pages
	SharedClean uint64
	// Size in bytes of dirty shared pages
	SharedDirty uint64
	// Size in bytes of clean private pages
	PrivateClean uint64
	// Size in bytes of dirty private pages
	PrivateDirty uint64
	// Amount of memory currently marked as referenced or accessed
	Referenced uint64
	// Amount of memory that does not belong to any file
	Anonymous uint64
	// Amount would-be-anonymous memory currently on swap
	Swap uint64
	// Process's proportional memory on swap
	SwapPss uint64
	// Page size used by the kernel to back the virtual memory area
	KernelPageSize uint64
	// Page size used by the processor MMU to back the virtual memory area
	MMUPageSize uint64
	// Kernel flags associated with the virtual memory area
	VMFlags string
}

// SizeSum returns the sum of Pss from all mappings
func (s ProcSMaps) SizeSum() (sum uint64) {
	for _, x := range s {
		sum += x.Size
	}

	return
}

// RssSum returns the sum of Pss from all mappings
func (s ProcSMaps) RssSum() (sum uint64) {
	for _, x := range s {
		sum += x.Rss
	}

	return
}

// PssSum returns the sum of Pss from all mappings
func (s ProcSMaps) PssSum() (sum uint64) {
	for _, x := range s {
		sum += x.Pss
	}

	return
}

// SharedCleanSum returns the sum of Pss from all mappings
func (s ProcSMaps) SharedCleanSum() (sum uint64) {
	for _, x := range s {
		sum += x.SharedClean
	}

	return
}

// SharedDirtySum returns the sum of Pss from all mappings
func (s ProcSMaps) SharedDirtySum() (sum uint64) {
	for _, x := range s {
		sum += x.SharedDirty
	}

	return
}

// PrivateCleanSum returns the sum of Pss from all mappings
func (s ProcSMaps) PrivateCleanSum() (sum uint64) {
	for _, x := range s {
		sum += x.PrivateClean
	}

	return
}

// PrivateDirtySum returns the sum of Pss from all mappings
func (s ProcSMaps) PrivateDirtySum() (sum uint64) {
	for _, x := range s {
		sum += x.PrivateDirty
	}

	return
}

// ReferencedSum returns the sum of Pss from all mappings
func (s ProcSMaps) ReferencedSum() (sum uint64) {
	for _, x := range s {
		sum += x.Referenced
	}

	return
}

// AnonymousSum returns the sum of Pss from all mappings
func (s ProcSMaps) AnonymousSum() (sum uint64) {
	for _, x := range s {
		sum += x.Anonymous
	}

	return
}

// SwapSum returns the sum of Pss from all mappings
func (s ProcSMaps) SwapSum() (sum uint64) {
	for _, x := range s {
		sum += x.Swap
	}

	return
}

// SwapPssSum returns the sum of Pss from all mappings
func (s ProcSMaps) SwapPssSum() (sum uint64) {
	for _, x := range s {
		sum += x.SwapPss
	}

	return
}

// ProcSMaps reads from /proc/[pid]/maps to get the memory-mappings of the
// process.
func (p Proc) ProcSMaps() (ProcSMaps, error) {
	file, err := os.Open(p.path("smaps"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var currentMap *ProcSMap

	smaps := []*ProcSMap{}
	scan := bufio.NewScanner(file)

	for scan.Scan() {
		line := scan.Text()
		// first line of a mapping start with an hexadecimal address in lower-case
		// All other line start with a capitalized words.
		if (line[0] >= '0' && line[0] <= '9') || (line[0] >= 'a' && line[0] <= 'f') {
			if currentMap != nil {
				smaps = append(smaps, currentMap)
			}

			m, err := parseProcMap(line)
			if err != nil {
				return nil, err
			}

			currentMap = &ProcSMap{
				ProcMap: m,
			}

			continue
		}

		kv := strings.SplitN(line, ":", 2)
		if len(kv) != 2 {
			continue
		}

		k := kv[0]
		v := strings.TrimSpace(kv[1])
		v = strings.TrimRight(v, " kB")

		vKBytes, _ := strconv.ParseUint(v, 10, 64)
		vBytes := vKBytes * 1024

		currentMap.fillValue(k, v, vKBytes, vBytes)
	}

	if currentMap != nil {
		smaps = append(smaps, currentMap)
	}

	return smaps, nil
}

func (s *ProcSMap) fillValue(k string, vString string, vUint uint64, vUintBytes uint64) {
	switch k {
	case "Size":
		s.Size = vUintBytes
	case "Rss":
		s.Rss = vUintBytes
	case "Pss":
		s.Pss = vUintBytes
	case "Shared_Clean":
		s.SharedClean = vUintBytes
	case "Shared_Dirty":
		s.SharedDirty = vUintBytes
	case "Private_Clean":
		s.PrivateClean = vUintBytes
	case "Private_Dirty":
		s.PrivateDirty = vUintBytes
	case "Referenced":
		s.Referenced = vUintBytes
	case "Anonymous":
		s.Anonymous = vUintBytes
	case "Swap":
		s.Swap = vUintBytes
	case "SwapPss":
		s.SwapPss = vUintBytes
	case "KernelPageSize":
		s.KernelPageSize = vUintBytes
	case "MMUPageSize":
		s.MMUPageSize = vUintBytes
	case "VmFlags":
		s.VMFlags = vString
	}
}
