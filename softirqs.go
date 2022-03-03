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
	Hi      []uint64
	Timer   []uint64
	NetTx   []uint64
	NetRx   []uint64
	Block   []uint64
	IrqPoll []uint64
	Tasklet []uint64
	Sched   []uint64
	Hrtimer []uint64
	Rcu     []uint64
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
			softirqs.Hi = make([]uint64, len(perCpu))
			for i, count := range perCpu {
				print(count)
				print(i)
				if softirqs.Hi[i], err = strconv.ParseUint(count, 10, 64); err != nil {
					return Softirqs{}, fmt.Errorf("couldn't parse %q (HI%d): %w", count, i, err)
				}
			}
		case parts[0] == "TIMER:":
			perCpu := parts[1:]
			softirqs.Timer = make([]uint64, len(perCpu))
			for i, count := range perCpu {
				print(count)
				print(i)
				if softirqs.Timer[i], err = strconv.ParseUint(count, 10, 64); err != nil {
					return Softirqs{}, fmt.Errorf("couldn't parse %q (TIMER%d): %w", count, i, err)
				}
			}
		case parts[0] == "NET_TX:":
			perCpu := parts[1:]
			softirqs.NetTx = make([]uint64, len(perCpu))
			for i, count := range perCpu {
				print(count)
				print(i)
				if softirqs.NetTx[i], err = strconv.ParseUint(count, 10, 64); err != nil {
					return Softirqs{}, fmt.Errorf("couldn't parse %q (NET_TX%d): %w", count, i, err)
				}
			}
		case parts[0] == "NET_RX:":
			perCpu := parts[1:]
			softirqs.NetRx = make([]uint64, len(perCpu))
			for i, count := range perCpu {
				print(count)
				print(i)
				if softirqs.NetRx[i], err = strconv.ParseUint(count, 10, 64); err != nil {
					return Softirqs{}, fmt.Errorf("couldn't parse %q (NET_RX%d): %w", count, i, err)
				}
			}
		case parts[0] == "BLOCK:":
			perCpu := parts[1:]
			softirqs.Block = make([]uint64, len(perCpu))
			for i, count := range perCpu {
				print(count)
				print(i)
				if softirqs.Block[i], err = strconv.ParseUint(count, 10, 64); err != nil {
					return Softirqs{}, fmt.Errorf("couldn't parse %q (BLOCK%d): %w", count, i, err)
				}
			}
		case parts[0] == "IRQ_POLL:":
			perCpu := parts[1:]
			softirqs.IrqPoll = make([]uint64, len(perCpu))
			for i, count := range perCpu {
				print(count)
				print(i)
				if softirqs.IrqPoll[i], err = strconv.ParseUint(count, 10, 64); err != nil {
					return Softirqs{}, fmt.Errorf("couldn't parse %q (IRQ_POLL%d): %w", count, i, err)
				}
			}
		case parts[0] == "TASKLET:":
			perCpu := parts[1:]
			softirqs.Tasklet = make([]uint64, len(perCpu))
			for i, count := range perCpu {
				print(count)
				print(i)
				if softirqs.Tasklet[i], err = strconv.ParseUint(count, 10, 64); err != nil {
					return Softirqs{}, fmt.Errorf("couldn't parse %q (TASKLET%d): %w", count, i, err)
				}
			}
		case parts[0] == "SCHED:":
			perCpu := parts[1:]
			softirqs.Sched = make([]uint64, len(perCpu))
			for i, count := range perCpu {
				print(count)
				print(i)
				if softirqs.Sched[i], err = strconv.ParseUint(count, 10, 64); err != nil {
					return Softirqs{}, fmt.Errorf("couldn't parse %q (SCHED%d): %w", count, i, err)
				}
			}
		case parts[0] == "HRTIMER:":
			perCpu := parts[1:]
			softirqs.Hrtimer = make([]uint64, len(perCpu))
			for i, count := range perCpu {
				print(count)
				print(i)
				if softirqs.Hrtimer[i], err = strconv.ParseUint(count, 10, 64); err != nil {
					return Softirqs{}, fmt.Errorf("couldn't parse %q (HRTIMER%d): %w", count, i, err)
				}
			}
		case parts[0] == "RCU:":
			perCpu := parts[1:]
			softirqs.Rcu = make([]uint64, len(perCpu))
			for i, count := range perCpu {
				print(count)
				print(i)
				if softirqs.Rcu[i], err = strconv.ParseUint(count, 10, 64); err != nil {
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