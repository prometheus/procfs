// Copyright 2017 The Prometheus Authors
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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/prometheus/procfs/bcache"
	"github.com/prometheus/procfs/cpu"
	"github.com/prometheus/procfs/xfs"
	"github.com/prometheus/procfs/internal/fs"
)

// FS represents the pseudo-filesystem sys, which provides an interface to
// kernel data structures.
type FS struct {
	sys fs.FS
}

// DefaultMountPoint is the common mount point of the sys filesystem.
const DefaultMountPoint = fs.DefaultSysMountPoint

// NewDefaultFS returns a new FS mounted under the default mountPoint. It will error
// if the mount point can't be read.
func NewDefaultFS() (FS, error) {
	return NewFS(DefaultMountPoint)
}

// NewFS returns a new FS mounted under the given mountPoint. It will error
// if the mount point can't be read.
func NewFS(mountPoint string) (FS, error) {
	fs, err := fs.NewFS(mountPoint)
	if err != nil {
		return FS{}, err
	}
	return FS{fs}, nil
}

// BcacheStats retrieves a map of vulnerability names to their mitigations.
func (fs FS) CPUVulnerabilities() ([]cpu.Vulnerability, error) {
	matches, err := filepath.Glob(fs.Path("devices/system/cpu/vulnerabilities/*"))
	if err != nil {
		return nil, err
	}

	vulnerabilities := make([]cpu.Vulnerability, 0, len(matches))
	for _, match := range matches {
		name := filepath.Base(match)

		value, err := ioutil.ReadFile(match)
		if err != nil {
			return nil, err
		}

		v, err := cpu.ParseVulnerability(name, strings.TrimSpace(string(value)))
		if err != nil {
			return nil, err
		}

		vulnerabilities = append(vulnerabilities, v)
	}

	return vulnerabilities, nil
}
