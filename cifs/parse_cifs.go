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

package cifs

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// parseHeader parses our SMB header
func parseHeader(line string, header map[string]uint64) error {
	for _, regexpHeader := range regexpHeaders {
		match := regexpHeader.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		for index, name := range regexpHeader.SubexpNames() {
			if index == 0 || name == "" {
				continue
			}
			value, err := strconv.ParseUint(match[index], 10, 64)
			if nil != err {
				return fmt.Errorf("invalid value in header")
			}
			header[name] = value
		}
	}
	return nil
}

// parseSMBStats parses a SMB block
func parseSMBStats(line string, stats map[string]uint64, sessionIDs *SessionIDs) error {
	for _, regexpSMB := range regexpSMBs {
		match := regexpSMB.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		for index, name := range regexpSMB.SubexpNames() {
			if index == 0 || name == "" {
				continue
			}
			switch name {
			case "sessionID":
				value, err := strconv.ParseUint(match[index], 10, 64)
				if nil != err {
					return fmt.Errorf("type mismatch for sessionID")
				}
				sessionIDs.SessionID = value
			case "server":
				if match[index] != "" {
					sessionIDs.Server = match[index]
				}
			case "share":
				if match[index] != "" {
					sessionIDs.Share = match[index]
				}
			default:
				value, err := strconv.ParseUint(match[index], 10, 64)
				if nil != err {
					return fmt.Errorf("invalid value in SMB Statistics")
				}
				stats[name] = value
			}
		}
		return nil
	}
	return nil
}

// ParseClientStats returns stats read from /proc/fs/cifs/Stats
func ParseClientStats(r io.Reader) (*ClientStats, error) {
	stats := &ClientStats{}
	stats.Header = make(map[string]uint64)
	scanner := bufio.NewScanner(r)
	var currentSMBBlock *SMBStats
	var currentSMBMetrics map[string]uint64
	var currentSMBSessionIDs *SessionIDs
	for scanner.Scan() {
		line := scanner.Text()
		// if line is empty we can go back to start
		if line == "" {
			continue
		}
		parseHeader(line, stats.Header)
		// If we see a new SMB block we are initializing all necessary structs and hashmaps
		if strings.Contains(line, ") \\") {
			currentSMBMetrics = make(map[string]uint64)
			currentSMBSessionIDs = &SessionIDs{}
			currentSMBBlock = &SMBStats{
				SessionIDs: *currentSMBSessionIDs,
				Stats:      currentSMBMetrics,
			}
			stats.SMBStatsList = append(stats.SMBStatsList, currentSMBBlock)
		}
		// Only parseSMBStats if we have a SMB block
		if currentSMBSessionIDs != nil {
			parseSMBStats(line, currentSMBMetrics, &currentSMBBlock.SessionIDs)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning SMB file: %s", err)
	}

	if len(stats.Header) == 0 {
		// We should never have an empty Header. Otherwise the file is invalid
		return nil, fmt.Errorf("error scanning SMB file: header is empty")
	}
	return stats, nil
}
