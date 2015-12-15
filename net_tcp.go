package procfs

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

// NetTCPLine is a parsed line of /proc/net/tcp
type NetTCPLine struct {
	Sl            string `net_tcp:"sl"`
	LocalAddress  string `net_tcp:"local_address"`
	RemoteAddress string `net_tcp:"rem_address"`
	St            string `net_tcp:"st"`
	TxQueue       string `net_tcp:"tx_queue"`
	RxQueue       string `net_tcp:"rx_queue"`
	Tr            string `net_tcp:"tr"`
	TmWhen        string `net_tcp:"tm->when"`
	Retrnsmt      string `net_tcp:"retrnsmt"`
	UID           string `net_tcp:"uid"`
	Timeout       string `net_tcp:"timeout"`
	Inode         string `net_tcp:"inode"`
	RefCount      string `net_tcp:""`
	MemoryAddress string `net_tcp:""`

	// There are optional attributes which I am not capturing
	// including: Retransmit Timeout, Predicted Tick, Ack.quick,
	// Sending Congestion Window, Slow Start Size Threshold.
}

// NetTCP stats on /proc/net/tcp
type NetTCP []NetTCPLine

// NewNetTCP returns kernel/system statistics read from /proc/net/tcp.
func NewNetTCP() (NetTCP, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return NetTCP{}, err
	}

	return fs.NewNetTCP()
}

// NewNetTCP returns an information about current kernel/system statistics.
func (fs FS) NewNetTCP() (m NetTCP, err error) {
	f, err := os.Open(fs.Path("net/tcp"))
	if err != nil {
		return NetTCP{}, err
	}
	defer f.Close()

	re := regexp.MustCompile(m.regex())
	s := bufio.NewScanner(f)

	for s.Scan() {
		line := s.Text()

		l, err := m.netTCPLine(re, line)
		if err == nil {
			m = append(m, l)
		}
	}

	return m, nil
}

func (n NetTCP) netTCPLine(re *regexp.Regexp, line string) (l NetTCPLine, err error) {
	matches := re.FindAllStringSubmatch(line, 21)
	if matches == nil {
		return NetTCPLine{}, fmt.Errorf("Invalid NetTCP Line")
	}
	r := matches[0]

	l = NetTCPLine{
		Sl:            r[1],
		LocalAddress:  r[2] + ":" + r[3],
		RemoteAddress: r[4] + ":" + r[5],
		St:            r[6],
		TxQueue:       r[7],
		RxQueue:       r[8],
		Tr:            r[9],
		TmWhen:        r[10],
		Retrnsmt:      r[11],
		UID:           r[12],
		Timeout:       r[13],
		Inode:         r[14],
		RefCount:      r[15],
		MemoryAddress: r[16],
	}

	return
}

/*
https://metacpan.org/source/SALVA/Linux-Proc-Net-TCP-0.07/README

qr/^\s*
    (\d+):\s                                     # sl                        -  0
    ([\dA-F]{8}(?:[\dA-F]{24})?):([\dA-F]{4})\s  # local address and port    -  1 &  2
    ([\dA-F]{8}(?:[\dA-F]{24})?):([\dA-F]{4})\s  # remote address and port   -  3 &  4
    ([\dA-F]{2})\s                               # st                        -  5
    ([\dA-F]{8}):([\dA-F]{8})\s                  # tx_queue and rx_queue     -  6 &  7
    (\d\d):([\dA-F]{8}|F{9,}|1AD7F[\dA-F]{6})\s  # tr and tm->when           -  8 &  9
    ([\dA-F]{8})\s+                              # retrnsmt                  - 10
    (\d+)\s+                                     # uid                       - 11
    (\d+)\s+                                     # timeout                   - 12
    (\d+)\s+                                     # inode                     - 13
    (\d+)\s+                                     # ref count                 - 14
    ((?:[\dA-F]{8}){1,2})                        # memory address            - 15
    (?:
	\s+
	(\d+)\s+                                 # retransmit timeout        - 16
	(\d+)\s+                                 # predicted tick            - 17
	(\d+)\s+                                 # ack.quick                 - 18
	(\d+)\s+                                 # sending congestion window - 19
	(-?\d+)                                  # slow start size threshold - 20
    )?
    \s*
    (.*)                                         # more                      - 21
    $
/xi,
*/
func (n NetTCP) regex() string {
	return `^\s*(\d+):\s([\dA-F]{8}(?:[\dA-F]{24})?):([\dA-F]{4})\s([\dA-F]{8}(?:[\dA-F]{24})?):([\dA-F]{4})\s([\dA-F]{2})\s([\dA-F]{8}):([\dA-F]{8})\s(\d\d):([\dA-F]{8}|F{9,}|1AD7F[\dA-F]{6\})\s([\dA-F]{8})\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+((?:[\dA-F]{8}){1,2})(?:\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(-?\d+))?\s*(.*)$`
}
