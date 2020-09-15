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
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

var (
	deviceLineRE   = regexp.MustCompile(`(\w+\d+) : (\w+) ?(?:\(.+?\))? ?(\w+)? ((?:\w+\d*\[\d+\](?:\(\w\))? ?)+)`)
	statusLineRE   = regexp.MustCompile(`(\d+) blocks .*\[(\d+)/(\d+)\] \[[U_]+\]`)
	recoveryLineRE = regexp.MustCompile(`(?:(\d{1,3}\.\d)%) \((\d+)/\d+\).+?(\d+\.\d)min`)
	devicesStrRE   = regexp.MustCompile(`(\w+\d*)\[(\d+)\](?:\((\w)\))?`)
)

type MDAssignedDevice struct {
	Name  string
	Role  int64
	State string
}

// MDStat holds info parsed from /proc/mdstat.
type MDStat struct {
	// Name of the device.
	Name string
	// activity-state of the device.
	ActivityState string
	// Raid personality of the device.
	Personality string
	// Number of active disks.
	DisksActive int64
	// Total number of disks the device requires.
	DisksTotal int64
	// Number of failed disks.
	DisksFailed int64
	// Spare disks in the device.
	DisksSpare int64
	// Number of blocks the device holds.
	BlocksTotal int64
	// Number of blocks on the device that are in sync.
	BlocksSynced int64
	// Percentage of blocks on the device that are in sync.
	PercentSynced float64
	// Remaining minutes to complete sync
	RemainingSyncMinutes float64
	// List of assigned devices
	AssignedDevices []MDAssignedDevice
}

// MDStat parses an mdstat-file (/proc/mdstat) and returns a slice of
// structs containing the relevant info.  More information available here:
// https://raid.wiki.kernel.org/index.php/Mdstat
func (fs FS) MDStat() ([]MDStat, error) {
	data, err := ioutil.ReadFile(fs.proc.Path("mdstat"))
	if err != nil {
		return nil, err
	}
	mdstat, err := parseMDStat(data)
	if err != nil {
		return nil, fmt.Errorf("error parsing mdstat %s: %s", fs.proc.Path("mdstat"), err)
	}
	return mdstat, nil
}

// parseMDStat parses data from mdstat file (/proc/mdstat) and returns a slice of
// structs containing the relevant info.
func parseMDStat(mdStatData []byte) ([]MDStat, error) {
	mdStats := []MDStat{}
	lines := strings.Split(string(mdStatData), "\n")

	for i, line := range lines {
		if strings.TrimSpace(line) == "" || line[0] == ' ' ||
			strings.HasPrefix(line, "Personalities") ||
			strings.HasPrefix(line, "unused") {
			continue
		}

		matches := deviceLineRE.FindStringSubmatch(line)
		if len(matches) < 4 {
			return nil, fmt.Errorf("not enough fields in mdline (expected at least 4): %s", line)
		}
		mdName := matches[1] // mdx
		state := matches[2]  // active or inactive
		var personality string
		assignedDevicesStr := matches[3]

		if len(matches) == 5 {
			personality = matches[3]
			assignedDevicesStr = matches[4]
		}

		assignedDevices, err := evalAssignedDevices(assignedDevicesStr)

		if err != nil {
			return nil, err
		}

		if len(lines) <= i+3 {
			return nil, fmt.Errorf(
				"error parsing %s: too few lines for md device",
				mdName,
			)
		}

		// Failed disks have the suffix (F) & Spare disks have the suffix (S).
		fail := int64(strings.Count(line, "(F)"))
		spare := int64(strings.Count(line, "(S)"))
		active, total, size, err := evalStatusLine(lines[i], lines[i+1])

		if err != nil {
			return nil, fmt.Errorf("error parsing md device lines: %s", err)
		}

		syncLineIdx := i + 2
		if strings.Contains(lines[i+2], "bitmap") { // skip bitmap line
			syncLineIdx++
		}

		// If device is syncing at the moment, get the number of currently
		// synced bytes, otherwise that number equals the size of the device.
		syncedBlocks := size
		syncedPercent := float64(100)
		remainingSyncMinutes := float64(0)
		recovering := strings.Contains(lines[syncLineIdx], "recovery")
		resyncing := strings.Contains(lines[syncLineIdx], "resync")
		checking := strings.Contains(lines[syncLineIdx], "check")

		// Append recovery and resyncing state info.
		if recovering || resyncing || checking {
			if recovering {
				state = "recovering"
			} else if checking {
				state = "checking"
			} else {
				state = "resyncing"
			}

			// Handle case when resync=PENDING or resync=DELAYED.
			if strings.Contains(lines[syncLineIdx], "PENDING") ||
				strings.Contains(lines[syncLineIdx], "DELAYED") {
				syncedBlocks = 0
				syncedPercent = 0
			} else {
				syncedPercent, syncedBlocks, remainingSyncMinutes, err = evalRecoveryLine(lines[syncLineIdx])
				if err != nil {
					return nil, fmt.Errorf("error parsing sync line in md device %s: %s", mdName, err)
				}
			}
		}

		mdStats = append(mdStats, MDStat{
			Name:                 mdName,
			ActivityState:        state,
			Personality:          personality,
			DisksActive:          active,
			DisksFailed:          fail,
			DisksSpare:           spare,
			DisksTotal:           total,
			BlocksTotal:          size,
			BlocksSynced:         syncedBlocks,
			PercentSynced:        syncedPercent,
			RemainingSyncMinutes: remainingSyncMinutes,
			AssignedDevices:      assignedDevices,
		})
	}

	return mdStats, nil
}

func evalStatusLine(deviceLine, statusLine string) (active, total, size int64, err error) {

	sizeStr := strings.Fields(statusLine)[0]
	size, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("unexpected statusLine %s: %s", statusLine, err)
	}

	if strings.Contains(deviceLine, "raid0") || strings.Contains(deviceLine, "linear") {
		// In the device deviceLine, only disks have a number associated with them in [].
		total = int64(strings.Count(deviceLine, "["))
		return total, total, size, nil
	}

	if strings.Contains(deviceLine, "inactive") {
		return 0, 0, size, nil
	}

	matches := statusLineRE.FindStringSubmatch(statusLine)
	if len(matches) != 4 {
		return 0, 0, 0, fmt.Errorf("couldn't find all the substring matches: %s", statusLine)
	}

	total, err = strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("unexpected statusLine %s: %s", statusLine, err)
	}

	active, err = strconv.ParseInt(matches[3], 10, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("unexpected statusLine %s: %s", statusLine, err)
	}

	return active, total, size, nil
}

func evalRecoveryLine(recoveryLine string) (syncedPercent float64, syncedBlocks int64, remainMinutes float64, err error) {
	matches := recoveryLineRE.FindStringSubmatch(recoveryLine)
	if len(matches) != 4 {
		return 0.0, 0, 0.0, fmt.Errorf("unexpected recoveryLine: %s", recoveryLine)
	}

	syncedPercent, err = strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0.0, 0, 0.0, fmt.Errorf("%s in recoveryLine: %s", err, recoveryLine)
	}

	syncedBlocks, err = strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return 0.0, 0, 0.0, fmt.Errorf("%s in recoveryLine: %s", err, recoveryLine)
	}

	remainMinutes, err = strconv.ParseFloat(matches[3], 64)
	if err != nil {
		return 0.0, 0, 0.0, fmt.Errorf("%s in recoveryLine: %s", err, recoveryLine)
	}

	return syncedPercent, syncedBlocks, remainMinutes, nil
}

func evalAssignedDevices(assignedDevicesStr string) ([]MDAssignedDevice, error) {
	fields := strings.Fields(assignedDevicesStr)
	assignedDevices := make([]MDAssignedDevice, len(fields))

	for i, d := range fields {
		matches := devicesStrRE.FindStringSubmatch(d)
		if len(matches) < 3 {
			return nil, fmt.Errorf("couldn't find all the substring matches: %s", d)
		}

		name := matches[1]
		state := "active"
		role, err := strconv.ParseInt(matches[2], 10, 64)
		if err != nil {
			return nil, err
		}

		if len(matches) == 4 {
			switch matches[3] {
			case "S":
				state = "spare"
			case "F":
				state = "failed"
			}
		}

		assignedDevices[i] = MDAssignedDevice{
			Name:  name,
			Role:  role,
			State: state,
		}
	}

	return assignedDevices, nil
}
