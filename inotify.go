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

package procfs

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

// InotifyStat represents a file descriptor's inotify watch count.
type InotifyStat struct {
	// File descriptor
	FD string
	// List of inotify lines in the fdinfo file
	Lines []string
}

// InotifyStat constructor. Only available on kernel 3.8+.  On older kernels,
// an empty slice of *InotifyStat.Lines will be returned.
func (p Proc) newInotifyStat(fd string) (*InotifyStat, error) {
	fdinfo, err := p.fileDescriptorInfo(fd)
	if err != nil {
		return nil, err
	}

	lines := []string{}

	scanner := bufio.NewScanner(strings.NewReader(string(fdinfo)))
	for scanner.Scan() {
		r := regexp.MustCompile("^inotify")
		if r.MatchString(scanner.Text()) {
			lines = append(lines, scanner.Text())
		}
	}

	i := &InotifyStat{
		FD:    fd,
		Lines: lines,
	}

	return i, nil
}

// InotifyStats returns inotify info of a process.
func (p Proc) InotifyStats() ([]InotifyStat, error) {
	names, err := p.fileDescriptors()
	if err != nil {
		return nil, err
	}

	i := []InotifyStat{}

	for _, n := range names {
		target, err := os.Readlink(p.path("fd", n))
		if err == nil {
			if strings.Contains(target, "inotify") {
				newstat, err := p.newInotifyStat(n)
				if err != nil {
					return nil, err
				}
				i = append(i, *newstat)
			}
		}
	}

	return i, nil
}

// InotifyWatchLen returns the total number of inotify watches used
// by a process.
func (p Proc) InotifyWatchLen() (int, error) {
	stats, err := p.InotifyStats()
	if err != nil {
		return 0, err
	}

	t := 0

	for _, s := range stats {
		t += len(s.Lines)
	}

	return t, nil
}
