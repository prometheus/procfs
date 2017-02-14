package procfs

import (
	"os"
	"path/filepath"
)

// ProcFD models the content of /proc/<pid>/fd/<fd>.
type FDs map[string]string

// NewFD creates a new ProcFD instance from a given Proc instance.
func (p Proc) NewFD() (FDs, error) {
	// TODO: Use fileDescriptors() method.

	fds, err := filepath.Glob(p.path("fd/*"))
	if err != nil {
		return nil, err
	}

	if len(fds) == 0 {
		return nil, nil
	}

	pfds := make(FDs)

	for _, i := range fds {
		link, err := os.Readlink(i)
		if err != nil {
			continue
		}

		v := filepath.Base(i)

		pfds[v] = link
	}

	return pfds, nil
}
