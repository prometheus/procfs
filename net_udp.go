// Copyright 2020 The Prometheus Authors
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
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	// readLimit is used by io.LimitReader while reading the content of the
	// /proc/net/udp{,6} files. The number of lines inside such a file is dynamic
	// as each line represents a single used socket.
	// In theory, the number of available sockets is 65535 (2^16 - 1) per IP.
	// With e.g. 150 Byte per line and the maximum number of 65535,
	// the reader needs to handle 150 Byte * 65535 =~ 10 MB for a single IP.
	readLimit = 4294967296 // Byte -> 4 GiB
)

type (
	// NetUDP represents the contents of /proc/net/udp{,6} file without the header.
	NetUDP []*netUDPLine

	// NetUDPSummary provides already computed values like the total queue lengths or
	// the total number of used sockets. In contrast to NetUDP it does not collect
	// the parsed lines into a slice.
	NetUDPSummary struct {
		// TxQueueLength shows the total queue length of all parsed tx_queue lengths.
		TxQueueLength uint64
		// RxQueueLength shows the total queue length of all parsed rx_queue lengths.
		RxQueueLength uint64
		// UsedSockets shows the total number of parsed lines representing the
		// number of used sockets.
		UsedSockets uint64
	}

	// netUDPLine represents the fields parsed from a single line
	// in /proc/net/udp{,6}. Fields which are not used by UDP are skipped.
	// For the proc file format details, see https://linux.die.net/man/5/proc.
	netUDPLine struct {
		Sl        uint64
		LocalAddr net.IP
		LocalPort uint64
		RemAddr   net.IP
		RemPort   uint64
		St        uint64
		TxQueue   uint64
		RxQueue   uint64
		UID       uint64
	}
)

// NetUDP returns the IPv4 kernel/networking statistics for UDP datagrams
// read from /proc/net/udp.
func (fs FS) NetUDP() (NetUDP, error) {
	return newNetUDP(fs.proc.Path("net/udp"))
}

// NetUDP6 returns the IPv6 kernel/networking statistics for UDP datagrams
// read from /proc/net/udp6.
func (fs FS) NetUDP6() (NetUDP, error) {
	return newNetUDP(fs.proc.Path("net/udp6"))
}

// NetUDPSummary returns already computed statistics like the total queue lengths
// for UDP datagrams read from /proc/net/udp.
func (fs FS) NetUDPSummary() (*NetUDPSummary, error) {
	return newNetUDPSummary(fs.proc.Path("net/udp"))
}

// NetUDP6Summary returns already computed statistics like the total queue lengths
// for UDP datagrams read from /proc/net/udp6.
func (fs FS) NetUDP6Summary() (*NetUDPSummary, error) {
	return newNetUDPSummary(fs.proc.Path("net/udp6"))
}

// newNetUDP creates a new NetUDP{,6} from the contents of the given file.
func newNetUDP(file string) (NetUDP, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	netUDP := NetUDP{}

	lr := io.LimitReader(f, readLimit)
	s := bufio.NewScanner(lr)
	s.Scan() // skip first line with headers
	for s.Scan() {
		fields := strings.Fields(s.Text())
		line, err := parseNetUDPLine(fields)
		if err != nil {
			return nil, err
		}
		netUDP = append(netUDP, line)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return netUDP, nil
}

// newNetUDPSummary creates a new NetUDP{,6} from the contents of the given file.
func newNetUDPSummary(file string) (*NetUDPSummary, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	netUDPSummary := &NetUDPSummary{}

	lr := io.LimitReader(f, readLimit)
	s := bufio.NewScanner(lr)
	s.Scan() // skip first line with headers
	for s.Scan() {
		fields := strings.Fields(s.Text())
		line, err := parseNetUDPLine(fields)
		if err != nil {
			return nil, err
		}
		netUDPSummary.TxQueueLength += line.TxQueue
		netUDPSummary.RxQueueLength += line.RxQueue
		netUDPSummary.UsedSockets++
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return netUDPSummary, nil
}

// NOTE: this function assumes the architecture is little endian, like x86 and amd64
// /proc/net/{t,u}dp{,6} files are network byte order for ipv4
// for ipv6 the address is four words consisting of four bytes each. In each of those four words the four bytes are written in reverse order.

func parseIP(hexIP string) (net.IP, error) {
	var byteIP []byte
	byteIP, err := hex.DecodeString(hexIP)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot parse address field in socket line: %s", hexIP)
	}
	switch len(byteIP) {
	case 4:
		return net.IP{byteIP[3], byteIP[2], byteIP[1], byteIP[0]}, nil
	case 16:
		i := net.IP{
			byteIP[3], byteIP[2], byteIP[1], byteIP[0],
			byteIP[7], byteIP[6], byteIP[5], byteIP[4],
			byteIP[11], byteIP[10], byteIP[9], byteIP[8],
			byteIP[15], byteIP[14], byteIP[13], byteIP[12],
		}
		return i, nil
	default:
		return nil, fmt.Errorf("Unable to parse IP %s", hexIP)
	}
}

// parseNetUDPLine parses a single line, represented by a list of fields.
func parseNetUDPLine(fields []string) (*netUDPLine, error) {
	line := &netUDPLine{}
	if len(fields) < 8 {
		return nil, fmt.Errorf(
			"cannot parse net udp socket line as it has less then 8 columns: %s",
			strings.Join(fields, " "),
		)
	}
	var err error // parse error

	// sl
	s := strings.Split(fields[0], ":")
	if len(s) != 2 {
		return nil, fmt.Errorf(
			"cannot parse sl field in udp socket line: %s", fields[0])
	}

	if line.Sl, err = strconv.ParseUint(s[0], 0, 64); err != nil {
		return nil, fmt.Errorf("cannot parse sl value in udp socket line: %s", err)
	}
	// local_address
	l := strings.Split(fields[1], ":")
	if len(l) != 2 {
		return nil, fmt.Errorf(
			"cannot parse local_address field in udp socket line: %s", fields[1])
	}
	if line.LocalAddr, err = parseIP(l[0]); err != nil {
		return nil, err
	}
	if line.LocalPort, err = strconv.ParseUint(l[1], 16, 64); err != nil {
		return nil, fmt.Errorf(
			"cannot parse local_address port value in udp socket line: %s", err)
	}

	// remote_address
	r := strings.Split(fields[2], ":")
	if len(r) != 2 {
		return nil, fmt.Errorf(
			"cannot parse rem_address field in udp socket line: %s", fields[1])
	}
	if line.RemAddr, err = parseIP(r[0]); err != nil {
		return nil, err
	}
	if line.RemPort, err = strconv.ParseUint(r[1], 16, 64); err != nil {
		return nil, fmt.Errorf(
			"cannot parse rem_address port value in udp socket line: %s", err)
	}

	// st
	if line.St, err = strconv.ParseUint(fields[3], 16, 64); err != nil {
		return nil, fmt.Errorf(
			"cannot parse st value in udp socket line: %s", err)
	}

	// tx_queue and rx_queue
	q := strings.Split(fields[4], ":")
	if len(q) != 2 {
		return nil, fmt.Errorf(
			"cannot parse tx/rx queues in udp socket line as it has a missing colon: %s",
			fields[4],
		)
	}
	if line.TxQueue, err = strconv.ParseUint(q[0], 16, 64); err != nil {
		return nil, fmt.Errorf("cannot parse tx_queue value in udp socket line: %s", err)
	}
	if line.RxQueue, err = strconv.ParseUint(q[1], 16, 64); err != nil {
		return nil, fmt.Errorf("cannot parse rx_queue value in udp socket line: %s", err)
	}

	// uid
	if line.UID, err = strconv.ParseUint(fields[7], 0, 64); err != nil {
		return nil, fmt.Errorf(
			"cannot parse uid value in udp socket line: %s", err)
	}

	return line, nil
}
