// Copyright The Prometheus Authors
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

//go:build linux

package sysfs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/prometheus/procfs/internal/util"
)

const meiClassPath = "class/mei/mei0"

type MEIClass struct {
	Dev           *string
	DevState      *string
	FWStatus      *string
	FWVersion     *string
	HBMVersion    *string
	HBMVersionDrv *string
	Kind          *string
	Trc           *string
	TxQueueLimit  *string
}

// MEIClass returns Management Engine Interface (MEI) information read from /sys/class/mei/.
func (fs FS) MEIClass() (*MEIClass, error) {
	path := fs.sys.Path(meiClassPath)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %q: %w", path, err)
	}

	var mei MEIClass
	for _, f := range files {
		if !f.Type().IsRegular() {
			continue
		}

		name := f.Name()
		if name == "uevent" {
			continue
		}

		filename := filepath.Join(path, name)
		value, err := util.SysReadFile(filename)
		if err != nil {
			// no check for perms since all files (well apart from tx_queue_limit) are 0444
			return nil, fmt.Errorf("failed to read file %q: %w", filename, err)
		}

		switch name {
		case "dev":
			mei.Dev = &value
		case "dev_state":
			mei.DevState = &value
		case "fw_status":
			mei.FWStatus = &value
		case "fw_ver":
			mei.FWVersion = &value
		case "hbm_ver":
			mei.HBMVersion = &value
		case "hbm_ver_drv":
			mei.HBMVersionDrv = &value
		case "kind":
			mei.Kind = &value
		case "trc":
			mei.Trc = &value
		case "tx_queue_limit":
			mei.TxQueueLimit = &value
		}
	}

	return &mei, nil
}
