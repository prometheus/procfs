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
	Stats    map[string][]uint64
}

// Lnstat() retrieves stats from /proc/net/stat/
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

		lnstatFile := Lnstats{
			Filename: filepath.Base(filePath),
			Stats:    make(map[string][]uint64),
		}
		scanner := bufio.NewScanner(file)
		scanner.Scan()
		// First string is always a header for stats
		var headers []string
		for _, header := range strings.Fields(scanner.Text()) {
			headers = append(headers, header)
		}

		// Other strings represent per-CPU counters
		for scanner.Scan() {
			for num, counter := range strings.Fields(scanner.Text()) {
				value, err := strconv.ParseUint(counter, 16, 32)
				if err != nil {
					return nil, err
				}
				lnstatFile.Stats[headers[num]] = append(lnstatFile.Stats[headers[num]], value)
			}
		}
		lnstatsTotal = append(lnstatsTotal, lnstatFile)
	}
	return lnstatsTotal, nil
}
