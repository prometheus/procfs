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

// SystemCPUCpufreq contains stats from devices/system/cpu/cpu[0-9]*/cpufreq/...
type SystemCPUCpufreq struct {
	Name               string
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

// SystemCpufreq is a collection of SystemCPUCpufreq for every CPU.
type SystemCpufreq []SystemCPUCpufreq

// TODO: Add topology support.

// TODO: Add thermal_throttle support.

// NewSystemCpufreq returns CPU frequency metrics for all CPUs.
func NewSystemCpufreq() (SystemCpufreq, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return SystemCpufreq{}, err
	}

	return fs.NewSystemCpufreq()
}

// NewSystemCpufreq returns CPU frequency metrics for all CPUs.
func (fs FS) NewSystemCpufreq() (SystemCpufreq, error) {
	var cpufreq = &SystemCPUCpufreq{}

	cpus, err := filepath.Glob(fs.Path("devices/system/cpu/cpu[0-9]*"))
	if err != nil {
		return SystemCpufreq{}, err
	}

	systemCpufreq := SystemCpufreq{}
	for _, cpu := range cpus {
		cpuName := filepath.Base(cpu)
		cpuNum := strings.TrimPrefix(cpuName, "cpu")

		cpuCpufreqPath := filepath.Join(cpu, "cpufreq")
		if _, err := os.Stat(cpuCpufreqPath); os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return SystemCpufreq{}, err
		}

		if _, err = os.Stat(filepath.Join(cpuCpufreqPath, "scaling_cur_freq")); err == nil {
			cpufreq, err = systemCpufreq.parseCpufreqCpuinfo("scaling", cpuCpufreqPath)
		} else if _, err = os.Stat(filepath.Join(cpuCpufreqPath, "cpuinfo_cur_freq")); err == nil {
			// Older kernels have metrics named `cpuinfo_...`.
			cpufreq, err = systemCpufreq.parseCpufreqCpuinfo("cpuinfo", cpuCpufreqPath)
		} else {
			return SystemCpufreq{}, fmt.Errorf("CPU %v is missing cpufreq", cpu)
		}
		if err != nil {
			return SystemCpufreq{}, err
		}
		cpufreq.Name = cpuNum
		systemCpufreq = append(systemCpufreq, *cpufreq)
	}

	return systemCpufreq, nil
}

func (s SystemCpufreq) parseCpufreqCpuinfo(prefix string, cpuPath string) (*SystemCPUCpufreq, error) {
	systemCPUCpufreqClass := SystemCPUCpufreq{}

	uintFiles := [4]string{
		prefix + "_cur_freq",
		prefix + "_max_freq",
		prefix + "_min_freq",
		"cpuinfo_transition_latency",
	}
	var uintOut [4]uint64

	for i, f := range uintFiles {
		v, err := util.ReadUintFromFile(filepath.Join(cpuPath, f))
		if err != nil {
			return &systemCPUCpufreqClass, err
		}

		uintOut[i] = v
	}

	systemCPUCpufreqClass.CurrentFrequency = uintOut[0]
	systemCPUCpufreqClass.MaximumFrequency = uintOut[1]
	systemCPUCpufreqClass.MinimumFrequency = uintOut[2]
	systemCPUCpufreqClass.TransitionLatency = uintOut[3]

	stringFiles := [5]string{
		"scaling_available_governors",
		"scaling_driver",
		"scaling_governor",
		"related_cpus",
		"scaling_setspeed",
	}
	var stringOut [5]string

	for i, f := range stringFiles {
		fileContents, err := util.SysReadFile(filepath.Join(cpuPath, f))
		if err != nil {
			return &systemCPUCpufreqClass, err
		}

		stringOut[i] = strings.TrimSpace(string(fileContents))
	}

	systemCPUCpufreqClass.AvailableGovernors = stringOut[0]
	systemCPUCpufreqClass.Driver = stringOut[1]
	systemCPUCpufreqClass.Govenor = stringOut[2]
	systemCPUCpufreqClass.RelatedCpus = stringOut[3]
	systemCPUCpufreqClass.SetSpeed = stringOut[4]

	return &systemCPUCpufreqClass, nil
}
