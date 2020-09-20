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

// Because this shares so much logic with the net_udp.go module, it
// just calls the udp parsing functions and converts the types to the tcp
// equivalent. This could be an issue if the formats of /dev/net/tcp and /dev/net/udp
// diverge, but is an advantage if the format changes in sync

type (
	// NetTCP represents the contents of /proc/net/tcp{,6} file without the header.
	NetTCP []*netIPSocketLine

	// NetTCPSummary provides already computed values like the total queue lengths or
	// the total number of used sockets. In contrast to NetTCP it does not collect
	// the parsed lines into a slice.
	NetTCPSummary NetIPSocketSummary

	// netTCPLine represents the fields parsed from a single line
	// in /proc/net/tcp{,6}. Fields which are not used by TCP are skipped.
	// For the proc file format details, see https://linux.die.net/man/5/proc.
	netTCPLine netIPSocketLine
)

// NetTCP returns the IPv4 kernel/networking statistics for TCP datagrams
// read from /proc/net/tcp.
func (fs FS) NetTCP() (NetTCP, error) {
	return newNetTCP(fs.proc.Path("net/tcp"))
}

// NetTCP6 returns the IPv6 kernel/networking statistics for TCP datagrams
// read from /proc/net/tcp6.
func (fs FS) NetTCP6() (NetTCP, error) {
	return newNetTCP(fs.proc.Path("net/tcp6"))
}

// NetTCPSummary returns already computed statistics like the total queue lengths
// for TCP datagrams read from /proc/net/tcp.
func (fs FS) NetTCPSummary() (*NetTCPSummary, error) {
	return newNetTCPSummary(fs.proc.Path("net/tcp"))
}

// NetTCP6Summary returns already computed statistics like the total queue lengths
// for TCP datagrams read from /proc/net/tcp6.
func (fs FS) NetTCP6Summary() (*NetTCPSummary, error) {
	return newNetTCPSummary(fs.proc.Path("net/tcp6"))
}

func newNetTCP(file string) (NetTCP, error) {
	n, err := newNetIPSocket(file)
	n1 := NetTCP(n)
	return n1, err
}

func newNetTCPSummary(file string) (*NetTCPSummary, error) {
	n, err := newNetIPSocketSummary(file)
	if n == nil {
		return nil, err
	}
	n1 := NetTCPSummary(*n)
	return &n1, err
}
