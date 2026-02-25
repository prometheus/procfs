// Copyright 2022 The Prometheus Authors
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

// Module represents a single module info in the kernel.
type Module struct {
	// Name is the name of the module.
	Name string
	// Size is the memory size of the module, in bytes
	Size uint64
	// Instances is the number of instances of the module are currently loaded.
	// A value of zero represents an unloaded module.
	Instances uint64
	// Dependencies is the list of modules that this module depends on.
	Dependencies []string
	// State is the state of the module is in: Live, Loading, or Unloading are the only possible values.
	State string
	// Offset is a memory offset for the loaded module
	Offset uint64
	// Taints is a list of taints that the module has.
	Taints []string
}

// Modules represents a list of Module structs.
type Modules []Module

// Modules parses the metrics from /proc/modules file and returns a slice of
// structs containing the relevant info.
func (fs FS) Modules() ([]Module, error) {
	data, err := util.ReadFileNoStat(fs.proc.Path("modules"))
	if err != nil {
		return nil, err
	}
	return parseModules(bytes.NewReader(data))
}

// parseModules parses the metrics from /proc/modules file
// and returns a []Module structure.
// - https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/4/html/reference_guide/s2-proc-modules
// - https://unix.stackexchange.com/a/152527/300614
func parseModules(r io.Reader) ([]Module, error) {
	var (
		scanner = bufio.NewScanner(r)
		modules = make([]Module, 0)
	)

	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) < 6 {
			return nil, fmt.Errorf("not enough fields in modules (expected at least 6 fields but got %d): %s", len(parts), parts)
		}

		module := Module{
			Name:         parts[0],
			Dependencies: []string{},
			State:        parts[4],
			Taints:       []string{},
		}

		if size, err := strconv.ParseUint(parts[1], 10, 64); err == nil {
			module.Size = size
		}

		if instances, err := strconv.ParseUint(parts[2], 10, 64); err == nil {
			module.Instances = instances
		}

		dependencies := parts[3]
		if dependencies != "-" {
			module.Dependencies = strings.Split(strings.TrimSuffix(dependencies, ","), ",")
		}

		if offset, err := strconv.ParseUint(parts[5], 10, 64); err == nil {
			module.Offset = offset
		}

		// Kernel Taint State is available if parts length is greater than 6.
		if len(parts) > 6 {
			taints := strings.TrimSuffix(strings.TrimPrefix(parts[6], "("), ")")
			module.Taints = strings.Split(taints, "")
		}

		modules = append(modules, module)
	}

	return modules, scanner.Err()
}
