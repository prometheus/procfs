// Copyright 2017 Prometheus Team
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

// Package procfs provides functions to retrieve system, kernel and process
// metrics from the pseudo-filesystem proc.
package procfs

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// XfrmStat models the contents of /proc/net/xfrm_stat.
type XfrmStat struct {
	// All errors which are not matched by other
	XfrmInError int
	// No buffer is left
	XfrmInBufferError int
	// Header Error
	XfrmInHdrError int
	// No state found
	// i.e. either inbound SPI, address, or IPSEC protocol at SA is wrong
	XfrmInNoStates int
	// Transformation protocol specific error
	// e.g. SA Key is wrong
	XfrmInStateProtoError int
	// Transformation mode specific error
	XfrmInStateModeError int
	// Sequence error
	// e.g. sequence number is out of window
	XfrmInStateSeqError int
	// State is expired
	XfrmInStateExpired int
	// State has mismatch option
	// e.g. UDP encapsulation type is mismatched
	XfrmInStateMismatch int
	// State is invalid
	XfrmInStateInvalid int
	// No matching template for states
	// e.g. Inbound SAs are correct but SP rule is wrong
	XfrmInTmplMismatch int
	// No policy is found for states
	// e.g. Inbound SAs are correct but no SP is found
	XfrmInNoPols int
	// Policy discards
	XfrmInPolBlock int
	// Policy error
	XfrmInPolError int
	// All errors which are not matched by others
	XfrmOutError int
	// Bundle generation error
	XfrmOutBundleGenError int
	// Bundle check error
	XfrmOutBundleCheckError int
	// No state was found
	XfrmOutNoStates int
	// Transformation protocol specific error
	XfrmOutStateProtoError int
	// Transportation mode specific error
	XfrmOutStateModeError int
	// Sequence error
	// i.e sequence number overflow
	XfrmOutStateSeqError int
	// State is expired
	XfrmOutStateExpired int
	// Policy discads
	XfrmOutPolBlock int
	// Policy is dead
	XfrmOutPolDead int
	// Policy Error
	XfrmOutPolError     int
	XfrmFwdHdrError     int
	XfrmOutStateInvalid int
	XfrmAcquireError    int
}

// NewXfrmStat reads the xfrm_stat statistics.
func NewXfrmStat() (XfrmStat, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return XfrmStat{}, err
	}

	return fs.NewXfrmStat()
}

// ParseXfrmStat reads the xfrm_stat statistics from the 'proc' filesystem.
func (fs FS) NewXfrmStat() (XfrmStat, error) {

	file, err := os.Open(fs.Path("net/xfrm_stat"))

	if err != nil {
		return XfrmStat{}, err
	}
	defer file.Close()

	var (
		x = XfrmStat{}
		s = bufio.NewScanner(file)
	)

	for s.Scan() {
		fields := strings.Fields(s.Text())
		name := fields[0]
		value := fields[1]

		switch name {
		case "XfrmInError":
			x.XfrmInError, err = strconv.Atoi(value)
		case "XfrmInBufferError":
			x.XfrmInBufferError, err = strconv.Atoi(value)
		case "XfrmInHdrError":
			x.XfrmInHdrError, err = strconv.Atoi(value)
		case "XfrmInNoStates":
			x.XfrmInNoStates, err = strconv.Atoi(value)
		case "XfrmInStateProtoError":
			x.XfrmInStateProtoError, err = strconv.Atoi(value)
		case "XfrmInStateModeError":
			x.XfrmInStateModeError, err = strconv.Atoi(value)
		case "XfrmInStateSeqError":
			x.XfrmInStateSeqError, err = strconv.Atoi(value)
		case "XfrmInStateExpired":
			x.XfrmInStateExpired, err = strconv.Atoi(value)
		case "XfrmInStateInvalid":
			x.XfrmInStateInvalid, err = strconv.Atoi(value)
		case "XfrmInTmplMismatch":
			x.XfrmInTmplMismatch, err = strconv.Atoi(value)
		case "XfrmInNoPols":
			x.XfrmInNoPols, err = strconv.Atoi(value)
		case "XfrmInPolBlock":
			x.XfrmInPolBlock, err = strconv.Atoi(value)
		case "XfrmInPolError":
			x.XfrmInPolError, err = strconv.Atoi(value)
		case "XfrmOutError":
			x.XfrmOutError, err = strconv.Atoi(value)
		case "XfrmInStateMismatch":
			x.XfrmInStateMismatch, err = strconv.Atoi(value)
		case "XfrmOutBundleGenError":
			x.XfrmOutBundleGenError, err = strconv.Atoi(value)
		case "XfrmOutBundleCheckError":
			x.XfrmOutBundleCheckError, err = strconv.Atoi(value)
		case "XfrmOutNoStates":
			x.XfrmOutNoStates, err = strconv.Atoi(value)
		case "XfrmOutStateProtoError":
			x.XfrmOutStateProtoError, err = strconv.Atoi(value)
		case "XfrmOutStateModeError":
			x.XfrmOutStateModeError, err = strconv.Atoi(value)
		case "XfrmOutStateSeqError":
			x.XfrmOutStateSeqError, err = strconv.Atoi(value)
		case "XfrmOutStateExpired":
			x.XfrmOutStateExpired, err = strconv.Atoi(value)
		case "XfrmOutPolBlock":
			x.XfrmOutPolBlock, err = strconv.Atoi(value)
		case "XfrmOutPolDead":
			x.XfrmOutPolDead, err = strconv.Atoi(value)
		case "XfrmOutPolError":
			x.XfrmOutPolError, err = strconv.Atoi(value)
		case "XfrmFwdHdrError":
			x.XfrmFwdHdrError, err = strconv.Atoi(value)
		case "XfrmOutStateInvalid":
			x.XfrmOutStateInvalid, err = strconv.Atoi(value)
		case "XfrmAcquireError":
			x.XfrmAcquireError, err = strconv.Atoi(value)
		}

		if err != nil {
			return XfrmStat{}, err
		}
	}

	return x, s.Err()
}
