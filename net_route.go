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
	"fmt"
	"io"
	"strings"

	"github.com/prometheus/procfs/internal/util"
)

const (
	BlackholeRepresentation string = "*"
	BlackholeIfaceName      string = "blackhole"
	RouteLineColumns        int    = 11
)

// A NetRouteLine represents one line from net/route
type NetRouteLine struct {
	Iface       string
	Destination uint32
	Gateway     uint32
	Flags       uint32
	RefCnt      uint32
	Use         uint32
	Metric      uint32
	Mask        uint32
	MTU         uint32
	Window      uint32
	IRTT        uint32
}

func (fs FS) NetRoute() ([]NetRouteLine, error) {
	return readNetRoute(fs.proc.Path("net", "route"))
}

func readNetRoute(path string) ([]NetRouteLine, error) {
	b, err := util.ReadFileNoStat(path)
	if err != nil {
		return nil, err
	}

	routelines, err := parseNetRoute(bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("failed to read net route from %s: %w", path, err)
	}
	return routelines, nil
}

func parseNetRoute(r io.Reader) ([]NetRouteLine, error) {
	var routelines []NetRouteLine

	scanner := bufio.NewScanner(r)
	scanner.Scan()
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		routeline, err := parseNetRouteLine(fields)
		if err != nil {
			return nil, err
		}
		routelines = append(routelines, *routeline)
	}
	return routelines, nil
}

func parseNetRouteLine(fields []string) (*NetRouteLine, error) {
	routeline := &NetRouteLine{
		Iface: fields[0],
	}
	if routeline.Iface == BlackholeRepresentation {
		routeline.Iface = BlackholeIfaceName
	}
	if len(fields) != RouteLineColumns {
		return nil, fmt.Errorf("invalid routeline, num of digits: %d", len(fields))
	}

	hexss := make([]string, 0, 3)
	hexss = append(hexss, fields[1], fields[2], fields[7])
	uss := make([]string, 0, 7)
	uss = append(uss, fields[3], fields[4], fields[5], fields[6], fields[8], fields[9], fields[10])
	hex, err := util.ParseHexUint64s(hexss)
	if err != nil {
		return nil, err
	}
	ss, err := util.ParseUint32s(uss)
	if err != nil {
		return nil, err
	}
	// hex
	routeline.Destination = uint32(*hex[0])
	routeline.Gateway = uint32(*hex[1])
	routeline.Mask = uint32(*hex[2])
	// uint32
	routeline.Flags = uint32(ss[0])
	routeline.RefCnt = uint32(ss[1])
	routeline.Use = uint32(ss[2])
	routeline.Metric = uint32(ss[3])
	routeline.MTU = uint32(ss[4])
	routeline.Window = uint32(ss[5])
	routeline.IRTT = uint32(ss[6])

	return routeline, nil
}
