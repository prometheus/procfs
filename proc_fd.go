package procfs

import (
	"os"
	"path/filepath"
)

// ProcFD models the content of /proc/<pid>/fd/<fd>.
type ProcFD map[string]string

// NewFD creates a new ProcFD instance from a given Proc instance.
func (p Proc) NewFD() (ProcFD, error) {
	fds, err := filepath.Glob(p.path("fd/*"))
	if err != nil {
		return nil, err
	}

	if len(fds) == 0 {
		return nil, nil
	}

	pfds := make(ProcFD)

	for _, i := range fds {
		link, err := os.Readlink(i)
		if err != nil {
			continue
		}

		pfds[filepath.Base(i)] = link
	}

	return pfds, nil
}
