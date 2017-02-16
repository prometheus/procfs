package procfs

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// NetSockstat stats on /proc/net/sockstat
type NetSockstat struct {
	Sockets struct {
		Used int64
	}

	TCP struct {
		InUse  int64
		Orphan int64
		Tw     int64
		Alloc  int64
		Mem    int64
	}

	UDP struct {
		InUse int64
		Mem   int64
	}

	UDPLite struct {
		InUse int64
	}

	RAW struct {
		InUse int64
	}

	FRAG struct {
		InUse  int64
		Memory int64
	}
}

// NewNetSockstat returns kernel/system statistics read from /proc/net/sockstat.
func NewNetSockstat() (NetSockstat, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return NetSockstat{}, err
	}

	return fs.NewNetSockstat()
}

// NewNetSockstat returns an information about current kernel/system statistics.
func (fs FS) NewNetSockstat() (m NetSockstat, err error) {
	f, err := os.Open(fs.Path("net/sockstat"))
	if err != nil {
		return NetSockstat{}, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()

		switch {
		case strings.HasPrefix(line, "sockets:"):
			parsedLines := strings.Fields(line)
			for i := range parsedLines {
				switch parsedLines[i] {
				case "used":
					m.Sockets.Used = parseInt64(parsedLines[i+1])
				}
			}
		case strings.HasPrefix(line, "TCP:"):
			parsedLines := strings.Fields(line)
			for i := range parsedLines {
				switch parsedLines[i] {
				case "inuse":
					m.TCP.InUse = parseInt64(parsedLines[i+1])
				case "orphan":
					m.TCP.Orphan = parseInt64(parsedLines[i+1])
				case "tw":
					m.TCP.Tw = parseInt64(parsedLines[i+1])
				case "alloc":
					m.TCP.Alloc = parseInt64(parsedLines[i+1])
				case "mem":
					m.TCP.Mem = parseInt64(parsedLines[i+1])
				}
			}
		case strings.HasPrefix(line, "UDP:"):
			parsedLines := strings.Fields(line)
			for i := range parsedLines {
				switch parsedLines[i] {
				case "inuse":
					m.UDP.InUse = parseInt64(parsedLines[i+1])
				case "mem":
					m.UDP.Mem = parseInt64(parsedLines[i+1])
				}
			}
		case strings.HasPrefix(line, "UDPLITE:"):
			parsedLines := strings.Fields(line)
			for i := range parsedLines {
				switch parsedLines[i] {
				case "inuse":
					m.UDPLite.InUse = parseInt64(parsedLines[i+1])
				}
			}
		case strings.HasPrefix(line, "RAW:"):
			parsedLines := strings.Fields(line)
			for i := range parsedLines {
				switch parsedLines[i] {
				case "inuse":
					m.RAW.InUse = parseInt64(parsedLines[i+1])
				}
			}
		case strings.HasPrefix(line, "FRAG:"):
			parsedLines := strings.Fields(line)
			for i := range parsedLines {
				switch parsedLines[i] {
				case "inuse":
					m.FRAG.InUse = parseInt64(parsedLines[i+1])
				case "memory":
					m.FRAG.Memory = parseInt64(parsedLines[i+1])
				}
			}
		}
	}

	return m, nil
}

func parseInt64(val string) int64 {
	v, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0
	}

	return v
}
