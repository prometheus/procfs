// Copyright 2018 The Prometheus Authors
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

// +build !windows

package sysfs

import (
	"errors"
	"path/filepath"
	"syscall"

	"github.com/prometheus/procfs/internal/util"
)

type ClassDrmCard struct {
	Name   string
	Enable uint64
	Driver string
}

func (fs FS) ClassDrmCard() ([]ClassDrmCard, error) {
	cards, err := filepath.Glob(fs.sys.Path("class/drm/card[0-9]"))
	if err != nil {
		return nil, err
	}

	stats := make([]ClassDrmCard, 0, len(cards))
	for _, card := range cards {
		cardStats, err := parseClassDrmCard(card)
		if err != nil {
			if errors.Is(err, syscall.ENODATA) {
				continue
			}
			return nil, err
		}
		cardStats.Name = filepath.Base(card)
		stats = append(stats, cardStats)
	}
	return stats, nil
}

func parseClassDrmCard(port string) (ClassDrmCard, error) {
	cardEnable, err := util.ReadIntFromFile(filepath.Join(port, "device/enable"))
	if err != nil {
		return ClassDrmCard{}, err
	}

	cardDriverPath, err := filepath.EvalSymlinks(filepath.Join(port, "device/driver"))
	if err != nil {
		return ClassDrmCard{}, err
	}

	return ClassDrmCard{
		Enable: uint64(cardEnable),
		Driver: filepath.Base(cardDriverPath),
	}, nil
}
