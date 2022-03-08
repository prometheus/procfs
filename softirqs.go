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
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/prometheus/procfs/internal/util"
)

// Softirqs represents the softirq statistics.
type Softirqs struct {
	HI      []uint64
	TIMER   []uint64
	NETTX   []uint64
	NETRX   []uint64
	BLOCK   []uint64
	IRQPOLL []uint64
	TASKLET []uint64
	SCHED   []uint64
	HRTIMER []uint64
	RCU     []uint64
}

func (fs FS) Softirqs() (Softirqs, error) {
	fileName := fs.proc.Path("softirqs")
	data, err := util.ReadFileNoStat(fileName)
	if err != nil {
		return Softirqs{}, err
	}

	reader := bytes.NewReader(data)

	return parseSoftirqs(reader)
}

func parseSoftirqs(r io.Reader) (Softirqs, error) {
	var (
		softirqs = Softirqs{}
		scanner  = bufio.NewScanner(r)
	)

	if !scanner.Scan() {
		return Softirqs{}, fmt.Errorf("softirqs empty")
	}

	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		var err error

		// require at least one cpu
		if len(parts) < 2 {
			continue
		}
		switch {
		case parts[0] == "HI:":
			perCpu := parts[1:]
			softirqs.HI = make([]uint64, len(perCpu))
			for i, count := range perCpu {
				if softirqs.HI[i], err = strconv.ParseUint(count, 10, 64); err != nil {
					return Softirqs{}, fmt.Errorf("couldn't parse %q (HI%d): %w", count, i, err)
				}
			}
		case parts[0] == "TIMER:":
			perCpu := parts[1:]
			softirqs.TIMER = make([]uint64, len(perCpu))
			for i, count := range perCpu {
				if softirqs.TIMER[i], err = strconv.ParseUint(count, 10, 64); err != nil {
					return Softirqs{}, fmt.Errorf("couldn't parse %q (TIMER%d): %w", count, i, err)
				}
			}
		case parts[0] == "NET_TX:":
			perCpu := parts[1:]
			softirqs.NETTX = make([]uint64, len(perCpu))
			for i, count := range perCpu {
				if softirqs.NETTX[i], err = strconv.ParseUint(count, 10, 64); err != nil {
					return Softirqs{}, fmt.Errorf("couldn't parse %q (NET_TX%d): %w", count, i, err)
				}
			}
		case parts[0] == "NET_RX:":
			perCpu := parts[1:]
			softirqs.NETRX = make([]uint64, len(perCpu))
			for i, count := range perCpu {
				if softirqs.NETRX[i], err = strconv.ParseUint(count, 10, 64); err != nil {
					return Softirqs{}, fmt.Errorf("couldn't parse %q (NET_RX%d): %w", count, i, err)
				}
			}
		case parts[0] == "BLOCK:":
			perCpu := parts[1:]
			softirqs.BLOCK = make([]uint64, len(perCpu))
			for i, count := range perCpu {
				if softirqs.BLOCK[i], err = strconv.ParseUint(count, 10, 64); err != nil {
					return Softirqs{}, fmt.Errorf("couldn't parse %q (BLOCK%d): %w", count, i, err)
				}
			}
		case parts[0] == "IRQ_POLL:":
			perCpu := parts[1:]
			softirqs.IRQPOLL = make([]uint64, len(perCpu))
			for i, count := range perCpu {
				if softirqs.IRQPOLL[i], err = strconv.ParseUint(count, 10, 64); err != nil {
					return Softirqs{}, fmt.Errorf("couldn't parse %q (IRQ_POLL%d): %w", count, i, err)
				}
			}
		case parts[0] == "TASKLET:":
			perCpu := parts[1:]
			softirqs.TASKLET = make([]uint64, len(perCpu))
			for i, count := range perCpu {
				if softirqs.TASKLET[i], err = strconv.ParseUint(count, 10, 64); err != nil {
					return Softirqs{}, fmt.Errorf("couldn't parse %q (TASKLET%d): %w", count, i, err)
				}
			}
		case parts[0] == "SCHED:":
			perCpu := parts[1:]
			softirqs.SCHED = make([]uint64, len(perCpu))
			for i, count := range perCpu {
				if softirqs.SCHED[i], err = strconv.ParseUint(count, 10, 64); err != nil {
					return Softirqs{}, fmt.Errorf("couldn't parse %q (SCHED%d): %w", count, i, err)
				}
			}
		case parts[0] == "HRTIMER:":
			perCpu := parts[1:]
			softirqs.HRTIMER = make([]uint64, len(perCpu))
			for i, count := range perCpu {
				if softirqs.HRTIMER[i], err = strconv.ParseUint(count, 10, 64); err != nil {
					return Softirqs{}, fmt.Errorf("couldn't parse %q (HRTIMER%d): %w", count, i, err)
				}
			}
		case parts[0] == "RCU:":
			perCpu := parts[1:]
			softirqs.RCU = make([]uint64, len(perCpu))
			for i, count := range perCpu {
				if softirqs.RCU[i], err = strconv.ParseUint(count, 10, 64); err != nil {
					return Softirqs{}, fmt.Errorf("couldn't parse %q (RCU%d): %w", count, i, err)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return Softirqs{}, fmt.Errorf("couldn't parse softirqs: %w", err)
	}

	return softirqs, scanner.Err()
}
