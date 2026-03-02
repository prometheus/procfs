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

// MEIClass returns Management Engine Interface (DMI) information read from /sys/class/mei/.
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
