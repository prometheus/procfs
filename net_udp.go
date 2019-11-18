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
	"os"
	"strconv"
	"strings"
)

type (
	// NetUDPLine is a line parsed from /proc/net/udp
	// For the proc file format details, see https://linux.die.net/man/5/proc
	NetUDPLine struct {
		TxQueue uint64
		RxQueue uint64
	}

	NetUDP struct {
		TxQueueLength uint64
		RxQueueLength uint64
		UsedSockets   uint64
	}
)

// NetUDP returns kernel/networking statistics for udp datagrams read from /proc/net/udp.
func (fs FS) NetUDP() (*NetUDP, error) {
	return newNetUDP(fs.proc.Path("net/udp"))
}

// newNetUDP creates a new NetUDP from the contents of the given file.
func newNetUDP(file string) (*NetUDP, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	netUDP := &NetUDP{}
	s := bufio.NewScanner(f)
	s.Scan() // skip first line with headers
	for s.Scan() {
		fields := strings.Fields(s.Text())
		line, err := parseNetUDPLine(fields)
		if err != nil {
			return nil, err
		}
		netUDP.TxQueueLength += line.TxQueue
		netUDP.RxQueueLength += line.RxQueue
		netUDP.UsedSockets++
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return netUDP, nil
}

func parseNetUDPLine(fields []string) (*NetUDPLine, error) {
	line := &NetUDPLine{}
	if len(fields) < 5 {
		return nil, fmt.Errorf(
			"cannot parse net udp socket line as it has less then 5 columns: %s",
			strings.Join(fields, " "),
		)
	}
	q := strings.Split(fields[4], ":")
	if len(q) < 2 {
		return nil, fmt.Errorf(
			"cannot parse tx/rx queues in udp socket line as it has a missing colon: %s",
			fields[4],
		)
	}
	var err error // parse error
	if line.TxQueue, err = strconv.ParseUint(q[0], 16, 64); err != nil {
		return nil, fmt.Errorf("cannot parse tx_queue value in udp socket line: %s", err)
	}
	if line.RxQueue, err = strconv.ParseUint(q[1], 16, 64); err != nil {
		return nil, fmt.Errorf("cannot parse rx_queue value in udp socket line: %s", err)
	}
	return line, nil
}
