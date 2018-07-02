// Copyright 2015 The Prometheus Authors
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

package sysfs

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

type IscsiStats struct {
	SessionName   string
	TargetName    string
	InitiatorName string
	State         string
	IfaceName     string
}

func NewIscsi() ([]IscsiStats, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return nil, err
	}

	return fs.NewIscsi()
}

func (fs FS) NewIscsi() ([]IscsiStats, error) {
	sessionDirs, err := filepath.Glob(fs.Path("class/iscsi_session/session*"))
	if err != nil {
		return nil, err
	}

	var result []IscsiStats
	for _, sessionDir := range sessionDirs {
		var err error
		var stats IscsiStats
		stats.SessionName = filepath.Base(sessionDir)

		stats.TargetName, err = readIscsiInfo(sessionDir, "targetname")
		if err != nil {
			return nil, err
		}

		stats.InitiatorName, err = readIscsiInfo(sessionDir, "initiatorname")
		if err != nil {
			return nil, err
		}

		stats.State, err = readIscsiInfo(sessionDir, "state")
		if err != nil {
			return nil, err
		}

		stats.IfaceName, err = readIscsiInfo(sessionDir, "ifacename")
		if err != nil {
			return nil, err
		}

		result = append(result, stats)
	}

	return result, nil
}

func readIscsiInfo(dir, file string) (string, error) {
	filename := filepath.Join(dir, file)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}
