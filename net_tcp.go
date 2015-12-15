package procfs

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

type NetTcpLine struct {
	Sl            string `net_tcp:"sl"`
	LocalAddress  string `net_tcp:"local_address"`
	RemoteAddress string `net_tcp:"rem_address"`
	St            string `net_tcp:"st"`
	TxQueue       string `net_tcp:"tx_queue"`
	RxQueue       string `net_tcp:"rx_queue"`
	Tr            string `net_tcp:"tr"`
	TmWhen        string `net_tcp:"tm->when"`
	Retrnsmt      string `net_tcp:"retrnsmt"`
	Uid           string `net_tcp:"uid"`
	Timeout       string `net_tcp:"timeout"`
	Inode         string `net_tcp:"inode"`
	RefCount      string `net_tcp:""`
	MemoryAddress string `net_tcp:""`

	// There are optional attributes which I am not capturing
	// including: Retransmit Timeout, Predicted Tick, Ack.quick,
	// Sending Congestion Window, Slow Start Size Threshold.
}

// NetTcp stats on /proc/net/sockstat
type NetTcp []NetTcpLine

// NewNetTcp returns kernel/system statistics read from /proc/net/sockstat.
func NewNetTcp() (NetTcp, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return NetTcp{}, err
	}

	return fs.NewNetTcp()
}

// NewNetTcp returns an information about current kernel/system statistics.
func (fs FS) NewNetTcp() (m NetTcp, err error) {
	f, err := os.Open(fs.Path("net/sockstat"))
	if err != nil {
		return NetTcp{}, err
	}
	defer f.Close()

	re := regexp.MustCompile(m.regex())
	s := bufio.NewScanner(f)

	for s.Scan() {
		line := s.Text()

		l, err := m.netTcpLine(re, line)
		if err == nil {
			m = append(m, l)
		}
	}

	return m, nil
}

func (m NetTcp) netTcpLine(re *regexp.Regexp, line string) (l NetTcpLine, err error) {
	matches := re.FindAllStringSubmatch(line, 21)
	if matches == nil {
		return NetTcpLine{}, fmt.Errorf("Invalid NetTcp Line")
	}
	r := matches[0]

	l = NetTcpLine{
		Sl:            r[1],
		LocalAddress:  r[2] + ":" + r[3],
		RemoteAddress: r[4] + ":" + r[5],
		St:            r[6],
		TxQueue:       r[7],
		RxQueue:       r[8],
		Tr:            r[9],
		TmWhen:        r[10],
		Retrnsmt:      r[11],
		Uid:           r[12],
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
func (n NetTcp) regex() string {
	return `^\s*(\d+):\s([\dA-F]{8}(?:[\dA-F]{24})?):([\dA-F]{4})\s([\dA-F]{8}(?:[\dA-F]{24})?):([\dA-F]{4})\s([\dA-F]{2})\s([\dA-F]{8}):([\dA-F]{8})\s(\d\d):([\dA-F]{8}|F{9,}|1AD7F[\dA-F]{6\})\s([\dA-F]{8})\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+((?:[\dA-F]{8}){1,2})(?:\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(-?\d+))?\s*(.*)$`
}

/*
  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 00000000:0016 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 16071 1 0000000000000000 100 0 0 10 0
   1: 00000000:0056 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 16481 1 0000000000000000 100 0 0 10 0
   2: 00000000:0057 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 16482 1 0000000000000000 100 0 0 10 0
   3: 00000000:0BB8 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 516865617 1 0000000000000000 100 0 0 10 0
   4: 00000000:0058 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 16476 1 0000000000000000 100 0 0 10 0
   5: 00000000:0059 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 16480 1 0000000000000000 100 0 0 10 0
   6: 00000000:005A 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 16479 1 0000000000000000 100 0 0 10 0
   7: 00000000:01BB 00000000:0000 0A 00000000:00000000 02:00000002 00000000     0        0 3994549147 2 0000000000000000 100 0 1 10 0
   8: 260FC90A:01BB AE6857C8:C747 03 00000000:00000000 01:000000DC 00000005     0        0 0 2 0000000000000000
   9: 260FC90A:01BB 0E7F51DE:7A79 03 00000000:00000000 01:0000084A 00000005     0        0 0 2 0000000000000000
  10: 260FC90A:01BB 89E2E3C9:749B 03 00000000:00000000 01:00000A16 00000005     0        0 0 2 0000000000000000
  11: 260FC90A:01BB AE6857C8:C752 03 00000000:00000000 01:0000008E 00000002     0        0 0 2 0000000000000000
  12: 260FC90A:01BB D81B38AC:64F9 03 00000000:00000000 01:00000063 00000000     0        0 0 2 0000000000000000
  13: 260FC90A:01BB 14288875:155E 03 00000000:00000000 01:00000A16 00000005     0        0 0 2 0000000000000000
  14: 260FC90A:01BB 53202464:C9D6 03 00000000:00000000 01:00000051 00000000     0        0 0 2 0000000000000000
  15: 260FC90A:01BB 5C49E8DD:F022 03 00000000:00000000 01:0000021F 00000004     0        0 0 2 0000000000000000
  16: 260FC90A:01BB 9A99935B:AFF3 03 00000000:00000000 01:00000246 00000004     0        0 0 2 0000000000000000
  17: 260FC90A:01BB 231D4B4B:EF00 03 00000000:00000000 01:000001BA 00000004     0        0 0 2 0000000000000000
  18: 260FC90A:01BB 162491A5:C524 03 00000000:00000000 01:00000272 00000004     0        0 0 2 0000000000000000
*/
