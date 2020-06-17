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
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/prometheus/procfs/internal/fs"
)

// Lnstats contains statistics for one counter for all cpus
type Lnstats struct {
	Filename string
	Name     string
	Value    map[uint64]uint64
}

func Lnstat() ([]Lnstats, error) {
	fs, err := NewFS(fs.DefaultProcMountPoint)
	if err != nil {
		return nil, err
	}
	statFiles, err := filepath.Glob(fs.proc.Path("net/stat/*"))
	if err != nil {
		return nil, err
	}

	var lnstatsTotal []Lnstats

	for _, filePath := range statFiles {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}

		var lnstatsOnce []Lnstats
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		scanner.Scan()
		// First string is always a header for stats
		for _, header := range strings.Fields(scanner.Text()) {
			lnstat := Lnstats{
				Filename: filepath.Base(filePath),
				Name:     header,
			}
			lnstat.Value = make(map[uint64]uint64)
			lnstatsOnce = append(lnstatsOnce, lnstat)
		}

		// Other strings represent per-CPU counters
		var cpu uint64 = 0
		for scanner.Scan() {
			for num, counter := range strings.Fields(scanner.Text()) {
				lnstatsOnce[num].Value[cpu], err = strconv.ParseUint(counter, 16, 32)
				if err != nil {
					return nil, err
				}
			}
			cpu++
		}
		lnstatsTotal = append(lnstatsTotal, lnstatsOnce...)
	}
	return lnstatsTotal, nil
}
