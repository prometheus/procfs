// Copyright 2020 Aleksei Zakharov
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

package lnstat

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"strings"
  "path/filepath"

	"github.com/prometheus/procfs/internal/util"
)

type Lnstats struct {
  Filename string
  Name string
  Value uint64
  CPU uint64
}

func Lnstat() (Lnstats[], error) {
  statFiles, err := filepath.Glob(fs.proc.Path("net/stat/*"))
  if err != nil {
    return err
  }

  var lnstats Lnstats[]

  for filePath range statsFiles {
    file, err := os.Open(filePath)
    if err != nil {
      return nil, err
    }

    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanLines)
    scanner.Scan()
    for _, header := range strings.Fields(scanner.Text()) {
      lnstat := Lnstats {
        Filename: filepath.Base(filePath),
        Name: header,
      }
      lnstats = append(lnstats. lnstat)
    }

    var cpu uint32 = 0
    for scanner.Scan() {
      for num, counter := range strings.Fields(scanner.Text()) {
        lnstats[num].CPU = cpu
        lnstats[unm].Value, err = strconv.ParseUint(counter, 16, 32)
        if err != nil {
          return nil, err
        }
      }
      cpu++
    }
  }
  return lnstats, nil
}
