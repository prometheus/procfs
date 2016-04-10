package procfs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Sys map[string]string

func NewSys() (Sys, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return nil, err
	}

	return fs.NewSys()
}

func (fs FS) NewSys() (m Sys, err error) {
	m = make(Sys)

	filepath.Walk(filepath.Join(string(fs), "sys"), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			body, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			m[path] = strings.TrimSpace(string(body))
		}

		return nil
	})

	return m, nil
}
