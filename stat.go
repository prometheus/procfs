package procfs

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Stat represents kernel/system statistics.
type Stat struct {
	// Boot time in seconds since the Epoch.
	BootTime int64
}

// NewStat returns kernel/system statistics read from /proc/stat.
func NewStat() (Stat, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return Stat{}, err
	}

	return fs.NewStat()
}

// NewStat returns an information about current kernel/system statistics.
func (fs FS) NewStat() (Stat, error) {
	f, err := os.Open(fs.Path("stat"))
	if err != nil {
		return Stat{}, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {

		/*
			# collect cpu and other misc load
			if (open(FILE, "/proc/stat")) {
			    while (my $line = <FILE>) {
			        $line =~ s/^\s+|\s+$//g;
			        $load{boot_time} = $1 if ($line =~ /^btime\s+(\d+)/);
			        $load{d_cpu} += $1 + $2 + $3 if ($line =~ /^cpu\s+(\d+)\s+(\d+)\s+(\d+)/);
			        $load{d_intrs} = $1 if ($line =~ /^intr\s+(\d+)/);
			        $load{d_proc_switches} = $1 if ($line =~ /^ctxt\s+(\d+)/);
			        $load{d_proc_forks} = $1 if ($line =~ /^processes\s+(\d+)/);
			        $load{uptime} = $load{time} - $load{boot_time};
			        $load{cpus}++ if ($line =~ /^cpu\d+/);
			    }
			    close FILE;
			}
		*/

		line := s.Text()
		if !strings.HasPrefix(line, "btime") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return Stat{}, fmt.Errorf("couldn't parse %s line %s", f.Name(), line)
		}
		i, err := strconv.ParseInt(fields[1], 10, 32)
		if err != nil {
			return Stat{}, fmt.Errorf("couldn't parse %s: %s", fields[1], err)
		}
		return Stat{BootTime: i}, nil
	}
	if err := s.Err(); err != nil {
		return Stat{}, fmt.Errorf("couldn't parse %s: %s", f.Name(), err)
	}

	return Stat{}, fmt.Errorf("couldn't parse %s, missing btime", f.Name())
}
