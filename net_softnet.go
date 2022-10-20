// Copyright 2019 The Prometheus Authors
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
	"strconv"
	"strings"

	"github.com/prometheus/procfs/internal/util"
)

// For the proc file format details,
// See:
// * Linux 2.6.23 https://elixir.bootlin.com/linux/v2.6.23/source/net/core/dev.c#L2343
// * Linux 4.17 https://elixir.bootlin.com/linux/v4.17/source/net/core/net-procfs.c#L162
// * Linux 5.10 https://elixir.bootlin.com/linux/v5.10/source/net/core/net-procfs.c#L172
// * https://elixir.bootlin.com/linux/v4.17/source/include/linux/netdevice.h#L2810
// * https://elixir.bootlin.com/linux/v5.10/source/include/linux/netdevice.h#L3181

// SoftnetStat contains a single row of data from /proc/net/softnet_stat.
type SoftnetStat struct {
	// Number of processed packets.
	Processed uint32
	// Number of dropped packets.
	Dropped uint32
	// Number of times processing packets ran out of quota.
	TimeSqueezed uint32
	// Number of collision occur while obtaining device lock while transmitting.
	CPUCollision uint32
	// Number of times cpu woken up received_rps.
	ReceivedRps uint32
	// number of times flow limit has been reached.
	FlowLimitCount uint32
	// Softnet Backlog status.
	SoftnetBacklogLen uint32
	// CPU id owning this softnet_data.
	Index uint32
	// softnet_data's Width
	Width int
}

var softNetProcFile = "net/softnet_stat"

// NetSoftnetStat reads data from /proc/net/softnet_stat.
func (fs FS) NetSoftnetStat() ([]SoftnetStat, error) {
	b, err := util.ReadFileNoStat(fs.proc.Path(softNetProcFile))
	if err != nil {
		return nil, err
	}

	entries, err := parseSoftnet(bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("failed to parse /proc/net/softnet_stat: %w", err)
	}

	return entries, nil
}

func parseSoftnet(r io.Reader) ([]SoftnetStat, error) {
	const minColumns = 9

	s := bufio.NewScanner(r)

	var stats []SoftnetStat
	for s.Scan() {
		columns := strings.Fields(s.Text())
		width := len(columns)

		if width < minColumns {
			return nil, fmt.Errorf("%d columns were detected, but at least %d were expected", width, minColumns)
		}

		// * Linux 2.6.23 https://elixir.bootlin.com/linux/v2.6.23/source/net/core/dev.c#L2343
		if width == minColumns {
			us, err := parseHexUint32s(columns[0:9])
			if err != nil {
				return nil, err
			}

			stats = append(stats, SoftnetStat{
				Processed:    us[0],
				Dropped:      us[1],
				TimeSqueezed: us[2],
				CPUCollision: us[8],
				Width:        width,
			})
		}

		// * Linux 4.17 https://elixir.bootlin.com/linux/v4.17/source/net/core/net-procfs.c#L162
		if width == 11 {
			us, err := parseHexUint32s(columns[0:11])
			if err != nil {
				return nil, err
			}

			stats = append(stats, SoftnetStat{
				Processed:      us[0],
				Dropped:        us[1],
				TimeSqueezed:   us[2],
				CPUCollision:   us[8],
				ReceivedRps:    us[9],
				FlowLimitCount: us[10],
				Width:          width,
			})
		}

		// * Linux 5.10 https://elixir.bootlin.com/linux/v5.10/source/net/core/net-procfs.c#L172
		if width >= 13 {
			us, err := parseHexUint32s(columns[0:13])
			if err != nil {
				return nil, err
			}

			stats = append(stats, SoftnetStat{
				Processed:         us[0],
				Dropped:           us[1],
				TimeSqueezed:      us[2],
				CPUCollision:      us[8],
				ReceivedRps:       us[9],
				FlowLimitCount:    us[10],
				SoftnetBacklogLen: us[11],
				Index:             us[12],
				Width:             width,
			})
		}
	}

	return stats, nil
}

func parseHexUint32s(ss []string) ([]uint32, error) {
	us := make([]uint32, 0, len(ss))
	for _, s := range ss {
		u, err := strconv.ParseUint(s, 16, 32)
		if err != nil {
			return nil, err
		}

		us = append(us, uint32(u))
	}

	return us, nil
}
