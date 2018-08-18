// Copyright 2018 The Prometheus Authors
// Copyright 2018 Sam Kottler <shkottler@gmail.com>

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
	"net"
	"os"
	"strings"
)

type ARPEntry struct {
	// IP address
	IPAddr net.IP
	// MAC address
	HWAddr net.HardwareAddr
	// Name of the device
	Device string
}

// GatherARPEntries retrieves all the ARP entries, parse the relevant columns,
// and then return a slice of ARPEntry's.
func GatherARPEntries() ([]ARPEntry, error) {
	fs, err := NewFS(DefaultMountPoint)

	if err != nil {
		return nil, err
	}

	file, err := os.Open(fs.Path("net/arp"))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	entries := make([]ARPEntry, 0)

	for scanner.Scan() {
		columns := strings.Fields(scanner.Text())
		entry := parseARPEntry(columns)
		entries = append(entries, entry)
	}

	return entries, nil
}

func parseARPEntry(columns []string) ARPEntry {
	ip := net.ParseIP(columns[0])
	mac := net.HardwareAddr(columns[3])

	entry := ARPEntry{
		IPAddr: ip,
		HWAddr: mac,
		Device: columns[5],
	}

	return entry
}
