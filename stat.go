package procfs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// CPUStat shows how much time the cpu spend in various stages.
type CPUStat struct {
	User      float64
	Nice      float64
	System    float64
	Idle      float64
	Iowait    float64
	Irq       float64
	Softirq   float64
	Steal     float64
	Guest     float64
	GuestNice float64
}

// SoftirqStat represent the softirq statistics as exported in the procfs stat file.
// A nice introduction can be found at https://0xax.gitbooks.io/linux-insides/content/interrupts/interrupts-9.html
// It is possible to get per-cpu stats by reading /proc/softirqs
type SoftirqStat struct {
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
	BootTime uint64
	// Summed up cpu statistics.
	CPUTotal CPUStat
	// Per-CPU statistics.
	CPU []CPUStat
	// Number of times interrupts were handled, which contains numbered and unnumbered IRQs.
	IRQTotal uint64
	// Number of times a numbered IRQ was triggered.
	IRQ []uint64
	// Number of times a context switch happened.
	ContextSwitches uint64
	// Number of times a process was created.
	ProcessCreated uint64
	// Number of processes currently running.
	ProcessesRunning uint64
	// Number of processes currently blocked (waiting for IO).
	ProcessesBlocked uint64
	// Number of times a softirq was scheduled.
	SoftirqTotal uint64
	// Detailed softirq statistics.
	Softirq SoftirqStat
}

// NewStat returns kernel/system statistics read from /proc/stat.
func NewStat() (Stat, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return Stat{}, err
	}

	return fs.NewStat()
}

// Parse a cpu statistics line and returns the CPUStat struct plus the cpu id (or -1 for the overall sum).
func parseCPUStat(line string) (CPUStat, int64, error) {
	cpuStat := CPUStat{}
	var cpu string

	count, err := fmt.Sscanf(line, "%s %f %f %f %f %f %f %f %f %f %f",
		&cpu,
		&cpuStat.User, &cpuStat.Nice, &cpuStat.System, &cpuStat.Idle,
		&cpuStat.Iowait, &cpuStat.Irq, &cpuStat.Softirq, &cpuStat.Steal,
		&cpuStat.Guest, &cpuStat.GuestNice)

	if err != nil && err != io.EOF {
		return CPUStat{}, -1, fmt.Errorf("couldn't parse %s (cpu): %s", line, err)
	}
	if count == 0 {
		return CPUStat{}, -1, fmt.Errorf("couldn't parse %s (cpu): 0 elements parsed", line)
	}

	cpuStat.User /= userHZ
	cpuStat.Nice /= userHZ
	cpuStat.System /= userHZ
	cpuStat.Idle /= userHZ
	cpuStat.Iowait /= userHZ
	cpuStat.Irq /= userHZ
	cpuStat.Softirq /= userHZ
	cpuStat.Steal /= userHZ
	cpuStat.Guest /= userHZ
	cpuStat.GuestNice /= userHZ

	if cpu == "cpu" {
		return cpuStat, -1, nil
	}

	cpuID, err := strconv.ParseInt(cpu[3:], 10, 64)
	if err != nil {
		return CPUStat{}, -1, fmt.Errorf("couldn't parse %s (cpu/cpuid): %s", line, err)
	}

	return cpuStat, cpuID, nil
}

// Parse a softirq line.
func parseSoftirqStat(line string) (SoftirqStat, uint64, error) {
	softirqStat := SoftirqStat{}
	var total uint64
	var prefix string

	_, err := fmt.Sscanf(line, "%s %d %d %d %d %d %d %d %d %d %d %d",
		&prefix, &total,
		&softirqStat.Hi, &softirqStat.Timer, &softirqStat.NetTx, &softirqStat.NetRx,
		&softirqStat.Block, &softirqStat.BlockIoPoll,
		&softirqStat.Tasklet, &softirqStat.Sched,
		&softirqStat.Hrtimer, &softirqStat.Rcu)

	if err != nil {
		return SoftirqStat{}, 0, fmt.Errorf("couldn't parse %s (softirq): %s", line, err)
	}

	return softirqStat, total, nil
}

// NewStat returns an information about current kernel/system statistics.
func (fs FS) NewStat() (Stat, error) {
	// See https://www.kernel.org/doc/Documentation/filesystems/proc.txt

	f, err := os.Open(fs.Path("stat"))
	if err != nil {
		return Stat{}, err
	}
	defer f.Close()

	stat := Stat{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(scanner.Text())
		// require at least <key> <value>
		if len(parts) < 2 {
			continue
		}
		switch {
		case parts[0] == "btime":
			if stat.BootTime, err = strconv.ParseUint(parts[1], 10, 64); err != nil {
				return Stat{}, fmt.Errorf("couldn't parse %s (btime): %s", parts[1], err)
			}
		case parts[0] == "intr":
			if stat.IRQTotal, err = strconv.ParseUint(parts[1], 10, 64); err != nil {
				return Stat{}, fmt.Errorf("couldn't parse %s (intr): %s", parts[1], err)
			}
			numberedIrqs := parts[2:]
			stat.IRQ = make([]uint64, len(numberedIrqs))
			for i, count := range numberedIrqs {
				if stat.IRQ[i], err = strconv.ParseUint(count, 10, 64); err != nil {
					return Stat{}, fmt.Errorf("couldn't parse %s (intr%d): %s", count, i, err)
				}
			}
		case parts[0] == "ctxt":
			if stat.ContextSwitches, err = strconv.ParseUint(parts[1], 10, 64); err != nil {
				return Stat{}, fmt.Errorf("couldn't parse %s (ctxt): %s", parts[1], err)
			}
		case parts[0] == "processes":
			if stat.ProcessCreated, err = strconv.ParseUint(parts[1], 10, 64); err != nil {
				return Stat{}, fmt.Errorf("couldn't parse %s (processes): %s", parts[1], err)
			}
		case parts[0] == "procs_running":
			if stat.ProcessesRunning, err = strconv.ParseUint(parts[1], 10, 64); err != nil {
				return Stat{}, fmt.Errorf("couldn't parse %s (procs_running): %s", parts[1], err)
			}
		case parts[0] == "procs_blocked":
			if stat.ProcessesBlocked, err = strconv.ParseUint(parts[1], 10, 64); err != nil {
				return Stat{}, fmt.Errorf("couldn't parse %s (procs_blocked): %s", parts[1], err)
			}
		case parts[0] == "softirq":
			softIrqStats, total, err := parseSoftirqStat(line)
			if err != nil {
				return Stat{}, fmt.Errorf("couldn't parse %s (softirq): %s", line, err)
			}
			stat.SoftirqTotal = total
			stat.Softirq = softIrqStats
		case strings.HasPrefix(parts[0], "cpu"):
			cpuStat, cpuID, err := parseCPUStat(line)
			if err != nil {
				return Stat{}, err
			}
			if cpuID == -1 {
				stat.CPUTotal = cpuStat
			} else {
				for int64(len(stat.CPU)) <= cpuID {
					stat.CPU = append(stat.CPU, CPUStat{})
				}
				stat.CPU[cpuID] = cpuStat
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return Stat{}, fmt.Errorf("couldn't parse %s: %s", f.Name(), err)
	}

	return stat, nil
}
