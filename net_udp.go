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
	NetUDP []*netIPSocketLine

	// NetUDPSummary provides already computed values like the total queue lengths or
	// the total number of used sockets. In contrast to NetUDP it does not collect
	// the parsed lines into a slice.
	NetUDPSummary NetIPSocketSummary

	// netUDPLine represents the fields parsed from a single line
	// in /proc/net/udp{,6}. Fields which are not used by UDP are skipped.
	// For the proc file format details, see https://linux.die.net/man/5/proc.
	netUDPLine netIPSocketLine
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
	n, err := newNetUDPSummary(fs.proc.Path("net/udp"))
	n1 := NetUDPSummary(*n)
	return &n1, err
}

// NetUDP6Summary returns already computed statistics like the total queue lengths
// for UDP datagrams read from /proc/net/udp6.
func (fs FS) NetUDP6Summary() (*NetUDPSummary, error) {
	n, err := newNetUDPSummary(fs.proc.Path("net/udp6"))
	n1 := NetUDPSummary(*n)
	return &n1, err
}

// newNetUDP creates a new NetUDP{,6} from the contents of the given file.
func newNetUDP(file string) (NetUDP, error) {
	n, err := newNetIPSocket(file)
	n1 := NetUDP(n)
	return n1, err
}

func newNetUDPSummary(file string) (*NetUDPSummary, error) {
	n, err := newNetIPSocketSummary(file)
	if n == nil {
		return nil, err
	}
	n1 := NetUDPSummary(*n)
	return &n1, err
}

