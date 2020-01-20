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
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/procfs/internal/util"
)

// CPUInfo contains general information about a system CPU found in /proc/cpuinfo
type CPUInfo struct {
	Processor       uint
	VendorID        string
	CPUFamily       string
	Model           string
	ModelName       string
	Stepping        string
	Microcode       string
	CPUMHz          float64
	CacheSize       string
	PhysicalID      string
	Siblings        uint
	CoreID          string
	CPUCores        uint
	APICID          string
	InitialAPICID   string
	FPU             string
	FPUException    string
	CPUIDLevel      uint
	WP              string
	Flags           []string
	Bugs            []string
	BogoMips        float64
	CLFlushSize     uint
	CacheAlignment  uint
	AddressSizes    string
	PowerManagement string
}

const (
	platformX86     = "x86"
	platformARM     = "arm"
	platformS390X   = "s390x"
	platformPPC     = "ppc"
	platformUnknown = "unknown"
)

var (
	cpuinfoX86Regexp   = regexp.MustCompile(`(?m)^\s*processor\s*:\s*\d+\s*vendor`)
	cpuinfoARMRegexp   = regexp.MustCompile(`^\s*Processor\s*:\s*ARM`)
	cpuinfoS390XRegexp = regexp.MustCompile(`^\s*vendor_id\s*:\s*IBM/S390`)
	cpuinfoPPCRegexp   = regexp.MustCompile(`(?m)^\s*processor\s*:\s*\d+\s+cpu\s+:\s+POWER`)

	cpuinfoClockRegexp          = regexp.MustCompile(`([\d.]+)`)
	cpuinfoS390XProcessorRegexp = regexp.MustCompile(`^processor\s+(\d+):.*`)
)

// cpuinfoDetectFormat attempts to determine the format used by the cpuinfo.
// This format corresponds to the platform generating the /proc/cpuinfo file.
// Returns "unknown"
func cpuinfoDetectFormat(info []byte) string {
	switch {
	case cpuinfoX86Regexp.Match(info):
		return platformX86
	case cpuinfoARMRegexp.Match(info):
		return platformARM
	case cpuinfoPPCRegexp.Match(info):
		return platformPPC
	case cpuinfoS390XRegexp.Match(info):
		return platformS390X
	}
	return platformUnknown
}

// CPUInfo returns information about current system CPUs.
// See https://www.kernel.org/doc/Documentation/filesystems/proc.txt
func (fs FS) CPUInfo() ([]CPUInfo, error) {
	data, err := util.ReadFileNoStat(fs.proc.Path("cpuinfo"))
	if err != nil {
		return nil, err
	}
	return parseCPUInfo(data)
}

// parseCPUInfo parses data from /proc/cpuinfo
func parseCPUInfo(info []byte) ([]CPUInfo, error) {
	platform := cpuinfoDetectFormat(info)
	switch platform {
	case platformX86:
		return parseCPUInfoX86(info)
	case platformARM:
		return parseCPUInfoARM(info)
	case platformS390X:
		return parseCPUInfoS390X(info)
	case platformPPC:
		return parseCPUInfoPPC(info)
	}
	return nil, errors.New("unable to determine format of 'cpuinfo'")
}

func parseCPUInfoX86(info []byte) ([]CPUInfo, error) {
	scanner := bufio.NewScanner(bytes.NewReader(info))

	// find the first "processor" line
	firstLine := firstNonEmptyLine(scanner)
	if !strings.HasPrefix(firstLine, "processor") || !strings.Contains(firstLine, ":") {
		return nil, errors.New("invalid cpuinfo file: " + firstLine)
	}
	field := strings.SplitN(firstLine, ": ", 2)
	v, err := strconv.ParseUint(field[1], 0, 32)
	if err != nil {
		return nil, err
	}
	firstcpu := CPUInfo{Processor: uint(v)}
	cpuinfo := []CPUInfo{firstcpu}
	i := 0

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, ":") {
			continue
		}
		field := strings.SplitN(line, ": ", 2)
		switch strings.TrimSpace(field[0]) {
		case "processor":
			cpuinfo = append(cpuinfo, CPUInfo{}) // start of the next processor
			i++
			v, err := strconv.ParseUint(field[1], 0, 32)
			if err != nil {
				return nil, err
			}
			cpuinfo[i].Processor = uint(v)
		case "vendor", "vendor_id":
			cpuinfo[i].VendorID = field[1]
		case "cpu family":
			cpuinfo[i].CPUFamily = field[1]
		case "model":
			cpuinfo[i].Model = field[1]
		case "model name":
			cpuinfo[i].ModelName = field[1]
		case "stepping":
			cpuinfo[i].Stepping = field[1]
		case "microcode":
			cpuinfo[i].Microcode = field[1]
		case "cpu MHz":
			v, err := strconv.ParseFloat(field[1], 64)
			if err != nil {
				return nil, err
			}
			cpuinfo[i].CPUMHz = v
		case "cache size":
			cpuinfo[i].CacheSize = field[1]
		case "physical id":
			cpuinfo[i].PhysicalID = field[1]
		case "siblings":
			v, err := strconv.ParseUint(field[1], 0, 32)
			if err != nil {
				return nil, err
			}
			cpuinfo[i].Siblings = uint(v)
		case "core id":
			cpuinfo[i].CoreID = field[1]
		case "cpu cores":
			v, err := strconv.ParseUint(field[1], 0, 32)
			if err != nil {
				return nil, err
			}
			cpuinfo[i].CPUCores = uint(v)
		case "apicid":
			cpuinfo[i].APICID = field[1]
		case "initial apicid":
			cpuinfo[i].InitialAPICID = field[1]
		case "fpu":
			cpuinfo[i].FPU = field[1]
		case "fpu_exception":
			cpuinfo[i].FPUException = field[1]
		case "cpuid level":
			v, err := strconv.ParseUint(field[1], 0, 32)
			if err != nil {
				return nil, err
			}
			cpuinfo[i].CPUIDLevel = uint(v)
		case "wp":
			cpuinfo[i].WP = field[1]
		case "flags":
			cpuinfo[i].Flags = strings.Fields(field[1])
		case "bugs":
			cpuinfo[i].Bugs = strings.Fields(field[1])
		case "bogomips":
			v, err := strconv.ParseFloat(field[1], 64)
			if err != nil {
				return nil, err
			}
			cpuinfo[i].BogoMips = v
		case "clflush size":
			v, err := strconv.ParseUint(field[1], 0, 32)
			if err != nil {
				return nil, err
			}
			cpuinfo[i].CLFlushSize = uint(v)
		case "cache_alignment":
			v, err := strconv.ParseUint(field[1], 0, 32)
			if err != nil {
				return nil, err
			}
			cpuinfo[i].CacheAlignment = uint(v)
		case "address sizes":
			cpuinfo[i].AddressSizes = field[1]
		case "power management":
			cpuinfo[i].PowerManagement = field[1]
		}
	}
	return cpuinfo, nil
}

func parseCPUInfoARM(info []byte) ([]CPUInfo, error) {
	scanner := bufio.NewScanner(bytes.NewReader(info))

	firstLine := firstNonEmptyLine(scanner)
	if !strings.HasPrefix(firstLine, "Processor") || !strings.Contains(firstLine, ":") {
		return nil, errors.New("invalid cpuinfo file: " + firstLine)
	}
	field := strings.SplitN(firstLine, ": ", 2)
	commonCPUInfo := CPUInfo{VendorID: field[1]}

	cpuinfo := []CPUInfo{}
	i := -1
	featuresLine := ""

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, ":") {
			continue
		}
		field := strings.SplitN(line, ": ", 2)
		switch strings.TrimSpace(field[0]) {
		case "processor":
			cpuinfo = append(cpuinfo, commonCPUInfo) // start of the next processor
			i++
			v, err := strconv.ParseUint(field[1], 0, 32)
			if err != nil {
				return nil, err
			}
			cpuinfo[i].Processor = uint(v)
		case "BogoMIPS":
			v, err := strconv.ParseFloat(field[1], 64)
			if err != nil {
				return nil, err
			}
			cpuinfo[i].BogoMips = v
		case "Features":
			featuresLine = line
		}
	}
	fields := strings.SplitN(featuresLine, ": ", 2)
	for i := range cpuinfo {
		cpuinfo[i].Flags = strings.Fields(fields[1])
	}
	return cpuinfo, nil
}

func parseCPUInfoS390X(info []byte) ([]CPUInfo, error) {
	scanner := bufio.NewScanner(bytes.NewReader(info))

	firstLine := firstNonEmptyLine(scanner)
	if !strings.HasPrefix(firstLine, "vendor_id") || !strings.Contains(firstLine, ":") {
		return nil, errors.New("invalid cpuinfo file: " + firstLine)
	}
	field := strings.SplitN(firstLine, ": ", 2)
	cpuinfo := []CPUInfo{}
	commonCPUInfo := CPUInfo{VendorID: field[1]}

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, ":") {
			continue
		}
		field := strings.SplitN(line, ": ", 2)
		switch strings.TrimSpace(field[0]) {
		case "bogomips per cpu":
			v, err := strconv.ParseFloat(field[1], 64)
			if err != nil {
				return nil, err
			}
			commonCPUInfo.BogoMips = v
		case "features":
			commonCPUInfo.Flags = strings.Fields(field[1])
		}
		if strings.HasPrefix(line, "processor") {
			match := cpuinfoS390XProcessorRegexp.FindStringSubmatch(line)
			if len(match) < 2 {
				return nil, errors.New("Invalid line found in cpuinfo: " + line)
			}
			cpu := commonCPUInfo
			v, err := strconv.ParseUint(match[1], 0, 32)
			if err != nil {
				return nil, err
			}
			cpu.Processor = uint(v)
			cpuinfo = append(cpuinfo, cpu)
		}
		if strings.HasPrefix(line, "cpu number") {
			break
		}
	}

	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, ":") {
			continue
		}
		field := strings.SplitN(line, ": ", 2)
		switch strings.TrimSpace(field[0]) {
		case "cpu number":
			i++
		case "cpu MHz dynamic":
			clock := cpuinfoClockRegexp.FindString(strings.TrimSpace(field[1]))
			v, err := strconv.ParseFloat(clock, 64)
			if err != nil {
				return nil, err
			}
			cpuinfo[i].CPUMHz = v
		}
	}

	return cpuinfo, nil
}

func parseCPUInfoPPC(info []byte) ([]CPUInfo, error) {
	scanner := bufio.NewScanner(bytes.NewReader(info))

	firstLine := firstNonEmptyLine(scanner)
	if !strings.HasPrefix(firstLine, "processor") || !strings.Contains(firstLine, ":") {
		return nil, errors.New("invalid cpuinfo file: " + firstLine)
	}
	field := strings.SplitN(firstLine, ": ", 2)
	v, err := strconv.ParseUint(field[1], 0, 32)
	if err != nil {
		return nil, err
	}
	firstcpu := CPUInfo{Processor: uint(v)}
	cpuinfo := []CPUInfo{firstcpu}
	i := 0

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, ":") {
			continue
		}
		field := strings.SplitN(line, ": ", 2)
		switch strings.TrimSpace(field[0]) {
		case "processor":
			cpuinfo = append(cpuinfo, CPUInfo{}) // start of the next processor
			i++
			v, err := strconv.ParseUint(field[1], 0, 32)
			if err != nil {
				return nil, err
			}
			cpuinfo[i].Processor = uint(v)
		case "cpu":
			cpuinfo[i].VendorID = field[1]
		case "clock":
			clock := cpuinfoClockRegexp.FindString(strings.TrimSpace(field[1]))
			v, err := strconv.ParseFloat(clock, 64)
			if err != nil {
				return nil, err
			}
			cpuinfo[i].CPUMHz = v
		}
	}
	return cpuinfo, nil
}

// firstNonEmptyLine advances the scanner to the first non-empty line
// and returns the contents of that line
func firstNonEmptyLine(scanner *bufio.Scanner) string {
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" {
			return line
		}
	}
	return ""
}
