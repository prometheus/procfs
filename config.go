package procfs

import (
	"fmt"
	"github.com/prometheus/procfs/internal/util"
	"strings"
)

func (fs FS) LoadSysConfig(name string) (string, error) {
	paths := strings.Split(name, ".")

	content, err := util.ReadFileNoStat(fs.proc.Path(paths...))
	if err !=nil {
		return "", fmt.Errorf("read %s config err: %w", name, err)
	}

	return string(content), nil
}