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
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/prometheus/procfs/internal/fs"
	"github.com/prometheus/procfs/internal/util"
)

// CPUStat shows how much time the cpu spend in various stages.
type CPUStat struct {
	User      float64
	Nice      float64
	System    float64
	Idle      float64
	Iowait    float64
	IRQ       float64
	SoftIRQ   float64
	Steal     float64
	Guest     float64
	GuestNice float64
}

// SoftIRQStat represent the softirq statistics as exported in the procfs stat file.
// A nice introduction can be found at https://0xax.gitbooks.io/linux-insides/content/interrupts/interrupts-9.html
// It is possible to get per-cpu stats by reading `/proc/softirqs`.
type SoftIRQStat struct {
	Hi          uint64
	Timer       uint64
	NetTx       uint64
	NetRx       uint64
	Block       uint64
	BlockIoPoll uint64
	Tasklet     uint64
	Sched       uint64
	Hrtimer     uint64
	Rcu         uint64
}

// Stat represents kernel/system statistics.
type Stat struct {
	// Boot time in seconds since the Epoch.
	BootTime *uint64
	// Summed up cpu statistics.
	CPUTotal *CPUStat
	// Per-CPU statistics.
	CPU []*CPUStat
	// Number of times interrupts were handled, which contains numbered and unnumbered IRQs.
	IRQTotal *uint64
	// Number of times a numbered IRQ was triggered.
	IRQ []uint64
	// Number of times a context switch happened.
	ContextSwitches *uint64
	// Number of times a process was created.
	ProcessCreated *uint64
	// Number of processes currently running.
	ProcessesRunning *uint64
	// Number of processes currently blocked (waiting for IO).
	ProcessesBlocked *uint64
	// Number of times a softirq was scheduled.
	SoftIRQTotal *uint64
	// Detailed softirq statistics.
	SoftIRQ *SoftIRQStat
}

// Parse a cpu statistics line and returns the CPUStat struct plus the cpu id (or -1 for the overall sum).
func parseCPUStat(line string) (CPUStat, int64, error) {
	cpuStat := CPUStat{}
	var cpu string

	count, err := fmt.Sscanf(line, "%s %f %f %f %f %f %f %f %f %f %f",
		&cpu,
		&cpuStat.User, &cpuStat.Nice, &cpuStat.System, &cpuStat.Idle,
		&cpuStat.Iowait, &cpuStat.IRQ, &cpuStat.SoftIRQ, &cpuStat.Steal,
		&cpuStat.Guest, &cpuStat.GuestNice)

	if err != nil && err != io.EOF {
		return CPUStat{}, -1, fmt.Errorf("couldn't parse %q (cpu): %w", line, err)
	}
	if count == 0 {
		return CPUStat{}, -1, fmt.Errorf("couldn't parse %q (cpu): 0 elements parsed", line)
	}

	cpuStat.User /= userHZ
	cpuStat.Nice /= userHZ
	cpuStat.System /= userHZ
	cpuStat.Idle /= userHZ
	cpuStat.Iowait /= userHZ
	cpuStat.IRQ /= userHZ
	cpuStat.SoftIRQ /= userHZ
	cpuStat.Steal /= userHZ
	cpuStat.Guest /= userHZ
	cpuStat.GuestNice /= userHZ

	if cpu == "cpu" {
		return cpuStat, -1, nil
	}

	cpuID, err := strconv.ParseInt(cpu[3:], 10, 64)
	if err != nil {
		return CPUStat{}, -1, fmt.Errorf("couldn't parse %q (cpu/cpuid): %w", line, err)
	}

	return cpuStat, cpuID, nil
}

// Parse a softirq line.
func parseSoftIRQStat(line string) (SoftIRQStat, uint64, error) {
	softIRQStat := SoftIRQStat{}
	var total uint64
	var prefix string

	_, err := fmt.Sscanf(line, "%s %d %d %d %d %d %d %d %d %d %d %d",
		&prefix, &total,
		&softIRQStat.Hi, &softIRQStat.Timer, &softIRQStat.NetTx, &softIRQStat.NetRx,
		&softIRQStat.Block, &softIRQStat.BlockIoPoll,
		&softIRQStat.Tasklet, &softIRQStat.Sched,
		&softIRQStat.Hrtimer, &softIRQStat.Rcu)

	if err != nil {
		return SoftIRQStat{}, 0, fmt.Errorf("couldn't parse %q (softirq): %w", line, err)
	}

	return softIRQStat, total, nil
}

// NewStat returns information about current cpu/process statistics.
// See https://www.kernel.org/doc/Documentation/filesystems/proc.txt
//
// Deprecated: Use fs.Stat() instead.
func NewStat() (*Stat, error) {
	fs, err := NewFS(fs.DefaultProcMountPoint)
	if err != nil {
		return nil, err
	}
	return fs.Stat()
}

// NewStat returns information about current cpu/process statistics.
// See: https://www.kernel.org/doc/Documentation/filesystems/proc.txt
//
// Deprecated: Use fs.Stat() instead.
func (fs FS) NewStat() (*Stat, error) {
	return fs.Stat()
}

// Stat returns information about current cpu/process statistics.
// See: https://www.kernel.org/doc/Documentation/filesystems/proc.txt
func (fs FS) Stat() (*Stat, error) {
	fileName := fs.proc.Path("stat")
	data, err := util.ReadFileNoStat(fileName)
	if err != nil {
		return nil, err
	}

	stat := Stat{}
	errMsgs := []string{}
	erroredWhenParseCPUStat := false
	scanner := bufio.NewScanner(bytes.NewReader(data))
	// in order to avoid issue https://github.com/prometheus/node_exporter/issues/1882
	// we try best to parse data in file /proc/stat. if failed to parse some metric ,we
	// record the error,then keep going, do not return error directly
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(scanner.Text())
		// require at least <key> <value>
		if len(parts) < 2 {
			continue
		}
		switch {
		case parts[0] == "btime":
			btime, err := strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				errMsgs = append(errMsgs, fmt.Sprintf("couldn't parse %q (btime): %s", parts[1], err))
				continue
			}
			stat.BootTime = &btime
		case parts[0] == "intr":
			irqTotal, err := strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				errMsgs = append(errMsgs, fmt.Sprintf("couldn't parse %q (intr): %s", parts[1], err))
				continue
			}
			numberedIRQs := parts[2:]
			irqs := make([]uint64, len(numberedIRQs))
			errored := false
			for i, count := range numberedIRQs {
				irq, err := strconv.ParseUint(count, 10, 64)
				if err != nil {
					errMsgs = append(errMsgs, fmt.Sprintf("couldn't parse %q (intr%d): %s", count, i, err))
					errored = true
					break
				}
				irqs[i] = irq
			}
			if !errored {
				stat.IRQTotal = &irqTotal
				stat.IRQ = irqs
			}
		case parts[0] == "ctxt":
			ctxt, err := strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				errMsgs = append(errMsgs, fmt.Sprintf("couldn't parse %q (ctxt): %s", parts[1], err))
				continue
			}
			stat.ContextSwitches = &ctxt
		case parts[0] == "processes":
			processes, err := strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				errMsgs = append(errMsgs, fmt.Sprintf("couldn't parse %q (processes): %s", parts[1], err))
				continue
			}
			stat.ProcessCreated = &processes
		case parts[0] == "procs_running":
			procsRunning, err := strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				errMsgs = append(errMsgs, fmt.Sprintf("couldn't parse %q (procs_running): %s", parts[1], err))
				continue
			}
			stat.ProcessesRunning = &procsRunning
		case parts[0] == "procs_blocked":
			procsBlocked, err := strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				errMsgs = append(errMsgs, fmt.Sprintf("couldn't parse %q (procs_blocked): %s", parts[1], err))
				continue
			}
			stat.ProcessesBlocked = &procsBlocked
		case parts[0] == "softirq":
			softIRQStats, total, err := parseSoftIRQStat(line)
			if err != nil {
				errMsgs = append(errMsgs, err.Error())
				continue
			}
			stat.SoftIRQTotal = &total
			stat.SoftIRQ = &softIRQStats
		case strings.HasPrefix(parts[0], "cpu"):
			cpuStat, cpuID, err := parseCPUStat(line)
			if err != nil {
				erroredWhenParseCPUStat = true
				errMsgs = append(errMsgs, err.Error())
				continue
			}
			if cpuID == -1 {
				stat.CPUTotal = &cpuStat
			} else {
				for int64(len(stat.CPU)) <= cpuID {
					var ptr *CPUStat
					stat.CPU = append(stat.CPU, ptr)
				}
				stat.CPU[cpuID] = &cpuStat
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("couldn't parse %q: %w", fileName, err)
	}
	// we can't make sure the order of the cpu stat data if there is an error
	// when parse stat data of cpu, so we reset the cpu stat data to nil
	if erroredWhenParseCPUStat {
		stat.CPUTotal = nil
		stat.CPU = []*CPUStat{}
	}
	if len(errMsgs) > 0 {
		return &stat, fmt.Errorf("%s", strings.Join(errMsgs, ","))
	}
	return &stat, nil
}
