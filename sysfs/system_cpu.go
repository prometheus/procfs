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

package sysfs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/prometheus/procfs/internal/util"
)

// SystemCpuCpufreq contains stats from devices/system/cpu/cpu[0-9]*/cpufreq/...
type SystemCpuCpufreq struct {
	CurrentFrequency   uint64
	MinimumFrequency   uint64
	MaximumFrequency   uint64
	TransitionLatency  uint64
	AvailableGovernors string
	Driver             string
	Govenor            string
	RelatedCpus        string
	SetSpeed           string
}

// SystemCpufreq is a collection of SystemCpuCpufreq for every CPU.
type SystemCpufreq map[string]SystemCpuCpufreq

// TODO: Add topology support.

// TODO: Add thermal_throttle support.

// NewSystemCpufreq returns CPU frequency metrics for all CPUs.
func NewSystemCpufreq() (SystemCpufreq, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return nil, err
	}

	return fs.NewSystemCpufreq()
}

// NewSystemCpufreq returns CPU frequency metrics for all CPUs.
func (fs FS) NewSystemCpufreq() (SystemCpufreq, error) {
	cpus, err := filepath.Glob(fs.Path("devices/system/cpu/cpu[0-9]*"))
	if err != nil {
		return SystemCpufreq{}, err
	}

	systemCpufreq := SystemCpufreq{}
	for _, cpu := range cpus {
		cpuName := filepath.Base(cpu)
		cpuNum := strings.TrimPrefix(cpuName, "cpu")

		cpufreq := SystemCpuCpufreq{}
		cpuCpufreqPath := filepath.Join(cpu, "cpufreq")
		if _, err := os.Stat(cpuCpufreqPath); os.IsNotExist(err) {
			continue
		} else {
			if _, err = os.Stat(filepath.Join(cpuCpufreqPath, "scaling_cur_freq")); err == nil {
				cpufreq, err = parseCpufreqCpuinfo("scaling", cpuCpufreqPath)
			} else if _, err = os.Stat(filepath.Join(cpuCpufreqPath, "cpuinfo_cur_freq")); err == nil {
				// Older kernels have metrics named `cpuinfo_...`.
				cpufreq, err = parseCpufreqCpuinfo("cpuinfo", cpuCpufreqPath)
			} else {
				return SystemCpufreq{}, fmt.Errorf("CPU %v is missing cpufreq", cpu)
			}
			if err != nil {
				return SystemCpufreq{}, err
			}
			systemCpufreq[cpuNum] = cpufreq
		}
	}

	return systemCpufreq, nil
}

func parseCpufreqCpuinfo(prefix string, cpuPath string) (SystemCpuCpufreq, error) {
	current, err := util.ReadUintFromFile(filepath.Join(cpuPath, prefix+"_cur_freq"))
	if err != nil {
		return SystemCpuCpufreq{}, err
	}
	maximum, err := util.ReadUintFromFile(filepath.Join(cpuPath, prefix+"_max_freq"))
	if err != nil {
		return SystemCpuCpufreq{}, err
	}
	minimum, err := util.ReadUintFromFile(filepath.Join(cpuPath, prefix+"_min_freq"))
	if err != nil {
		return SystemCpuCpufreq{}, err
	}
	transitionLatency, err := util.ReadUintFromFile(filepath.Join(cpuPath, "cpuinfo_transition_latency"))
	if err != nil {
		return SystemCpuCpufreq{}, err
	}
	fileContents, err := util.SysReadFile(filepath.Join(cpuPath, "scaling_available_governors"))
	if err != nil {
		return SystemCpuCpufreq{}, err
	}
	availableGovernors := strings.TrimSpace(string(fileContents))
	fileContents, err = util.SysReadFile(filepath.Join(cpuPath, "scaling_driver"))
	if err != nil {
		return SystemCpuCpufreq{}, err
	}
	driver := strings.TrimSpace(string(fileContents))
	fileContents, err = util.SysReadFile(filepath.Join(cpuPath, "scaling_governor"))
	if err != nil {
		return SystemCpuCpufreq{}, err
	}
	governor := strings.TrimSpace(string(fileContents))
	fileContents, err = util.SysReadFile(filepath.Join(cpuPath, "related_cpus"))
	if err != nil {
		return SystemCpuCpufreq{}, err
	}
	relatedCpus := strings.TrimSpace(string(fileContents))
	fileContents, err = util.SysReadFile(filepath.Join(cpuPath, "scaling_setspeed"))
	if err != nil {
		return SystemCpuCpufreq{}, err
	}
	setSpeed := strings.TrimSpace(string(fileContents))
	return SystemCpuCpufreq{
		CurrentFrequency:   current,
		MaximumFrequency:   maximum,
		MinimumFrequency:   minimum,
		TransitionLatency:  transitionLatency,
		AvailableGovernors: availableGovernors,
		Driver:             driver,
		Govenor:            governor,
		RelatedCpus:        relatedCpus,
		SetSpeed:           setSpeed,
	}, nil
}
