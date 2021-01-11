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

package procfs

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
)

// ProcLimits represents the soft limits for each of the process's resource
// limits. For more information see getrlimit(2):
// http://man7.org/linux/man-pages/man2/getrlimit.2.html.
type ProcLimits struct {
	// CPU time limit in seconds.
	CPUTime uint64
	// Maximum size of files that the process may create.
	FileSize uint64
	// Maximum size of the process's data segment (initialized data,
	// uninitialized data, and heap).
	DataSize uint64
	// Maximum size of the process stack in bytes.
	StackSize uint64
	// Maximum size of a core file.
	CoreFileSize uint64
	// Limit of the process's resident set in pages.
	ResidentSet uint64
	// Maximum number of processes that can be created for the real user ID of
	// the calling process.
	Processes uint64
	// Value one greater than the maximum file descriptor number that can be
	// opened by this process.
	OpenFiles uint64
	// Maximum number of bytes of memory that may be locked into RAM.
	LockedMemory uint64
	// Maximum size of the process's virtual memory address space in bytes.
	AddressSpace uint64
	// Limit on the combined number of flock(2) locks and fcntl(2) leases that
	// this process may establish.
	FileLocks uint64
	// Limit of signals that may be queued for the real user ID of the calling
	// process.
	PendingSignals uint64
	// Limit on the number of bytes that can be allocated for POSIX message
	// queues for the real user ID of the calling process.
	MsqqueueSize uint64
	// Limit of the nice priority set using setpriority(2) or nice(2).
	NicePriority uint64
	// Limit of the real-time priority set using sched_setscheduler(2) or
	// sched_setparam(2).
	RealtimePriority uint64
	// Limit (in microseconds) on the amount of CPU time that a process
	// scheduled under a real-time scheduling policy may consume without making
	// a blocking system call.
	RealtimeTimeout uint64
}

const (
	limitsFields    = 3
	limitsUnlimited = "unlimited"
)

var (
	limitsDelimiter = regexp.MustCompile("  +")
)

// NewLimits returns the current soft limits of the process.
//
// Deprecated: use p.Limits() instead
func (p Proc) NewLimits() (ProcLimits, error) {
	return p.Limits()
}

// Limits returns the current soft limits of the process.
func (p Proc) Limits() (ProcLimits, error) {
	f, err := os.Open(p.path("limits"))
	if err != nil {
		return ProcLimits{}, err
	}
	defer f.Close()

	var (
		l = ProcLimits{}
		s = bufio.NewScanner(f)
	)
	for s.Scan() {
		fields := limitsDelimiter.Split(s.Text(), limitsFields)
		if len(fields) != limitsFields {
			return ProcLimits{}, fmt.Errorf(
				"couldn't parse %s line %s", f.Name(), s.Text())
		}

		switch fields[0] {
		case "Max cpu time":
			l.CPUTime, err = parseLimitsUInt(fields[1])
		case "Max file size":
			l.FileSize, err = parseLimitsUInt(fields[1])
		case "Max data size":
			l.DataSize, err = parseLimitsUInt(fields[1])
		case "Max stack size":
			l.StackSize, err = parseLimitsUInt(fields[1])
		case "Max core file size":
			l.CoreFileSize, err = parseLimitsUInt(fields[1])
		case "Max resident set":
			l.ResidentSet, err = parseLimitsUInt(fields[1])
		case "Max processes":
			l.Processes, err = parseLimitsUInt(fields[1])
		case "Max open files":
			l.OpenFiles, err = parseLimitsUInt(fields[1])
		case "Max locked memory":
			l.LockedMemory, err = parseLimitsUInt(fields[1])
		case "Max address space":
			l.AddressSpace, err = parseLimitsUInt(fields[1])
		case "Max file locks":
			l.FileLocks, err = parseLimitsUInt(fields[1])
		case "Max pending signals":
			l.PendingSignals, err = parseLimitsUInt(fields[1])
		case "Max msgqueue size":
			l.MsqqueueSize, err = parseLimitsUInt(fields[1])
		case "Max nice priority":
			l.NicePriority, err = parseLimitsUInt(fields[1])
		case "Max realtime priority":
			l.RealtimePriority, err = parseLimitsUInt(fields[1])
		case "Max realtime timeout":
			l.RealtimeTimeout, err = parseLimitsUInt(fields[1])
		}
		if err != nil {
			return ProcLimits{}, err
		}
	}

	return l, s.Err()
}

func parseLimitsUInt(s string) (uint64, error) {
	if s == limitsUnlimited {
		// From Linux include/uapi/linux/resource.h:
		//   #define RLIM64_INFINITY (~0ULL)
		return math.MaxUint64, nil
	}
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("couldn't parse value %s: %s", s, err)
	}
	return i, nil
}
