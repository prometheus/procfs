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
	"bytes"
	"strconv"
	"strings"

	"github.com/prometheus/procfs/internal/util"
)

// ProtocolStats stores the contents from /proc/net/protocols
type ProtocolStats map[string]ProtocolStatLine

// ProtocolStatLine contains a single line parsed from /proc/net/protocols. We
// only care about the first six columns as the rest are not likely to change
// and only serve to provide a set of capabilities for each protocol.
type ProtocolStatLine struct {
	Name      string // The name of the protocol
	Size      uint64 // The size, in bytes, of a given protocol structure. e.g. sizeof(struct tcp_sock) or sizeof(struct unix_sock)
	Sockets   int64  // Number of sockets in use by this protocol
	Memory    int64  // Number of 4KB pages allocated by all sockets of this protocol
	Pressure  bool   // This is either yes, no, or NI (not implemented). For the sake of simplicity we treat NI as not experiencing memory pressure.
	MaxHeader uint64 // Protocol specific max header size
}

// Protocols reads stats from /proc/net/protocols and returns a map of
// PortocolStatLine entries
func (fs FS) Protocols() (ProtocolStats, error) {
	data, err := util.ReadFileNoStat(fs.proc.Path("net/protocols"))
	if err != nil {
		return ProtocolStats{}, err
	}
	return parseProtocols(bufio.NewScanner(bytes.NewReader(data)))
}

func parseProtocols(s *bufio.Scanner) (ProtocolStats, error) {
	protocolStats := ProtocolStats{}

	// Skip the header line
	s.Scan()

	for s.Scan() {
		line, err := protocolStats.parseLine(s.Text())
		if err != nil {
			return ProtocolStats{}, err
		}

		protocolStats[line.Name] = *line
	}
	return protocolStats, nil
}

func (ps ProtocolStats) parseLine(rawLine string) (*ProtocolStatLine, error) {
	line := &ProtocolStatLine{}
	var err error

	// We only care about the first 6 columns
	// protocol  size sockets  memory press maxhdr
	fields := strings.Fields(rawLine)[0:6]

	line.Name = fields[0]
	line.Size, err = strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return nil, err
	}
	line.Sockets, err = strconv.ParseInt(fields[2], 10, 64)
	if err != nil {
		return nil, err
	}
	line.Memory, err = strconv.ParseInt(fields[3], 10, 64)
	if err != nil {
		return nil, err
	}
	line.Pressure = false
	if fields[4] == "yes" {
		line.Pressure = true
	}
	line.MaxHeader, err = strconv.ParseUint(fields[5], 10, 64)
	if err != nil {
		return nil, err
	}
	return line, nil
}
