// Copyright 2022 The Prometheus Authors
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
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/prometheus/procfs/internal/util"
)

// ProcSnmp models the content of /proc/<pid>/net/snmp.
type ProcSnmp struct {
	// The process ID.
	PID int
	Ip
	Icmp
	IcmpMsg
	Tcp
	Udp
	UdpLite
}

type Ip struct {
	Forwarding      float64
	DefaultTTL      float64
	InReceives      float64
	InHdrErrors     float64
	InAddrErrors    float64
	ForwDatagrams   float64
	InUnknownProtos float64
	InDiscards      float64
	InDelivers      float64
	OutRequests     float64
	OutDiscards     float64
	OutNoRoutes     float64
	ReasmTimeout    float64
	ReasmReqds      float64
	ReasmOKs        float64
	ReasmFails      float64
	FragOKs         float64
	FragFails       float64
	FragCreates     float64
}

type Icmp struct {
	InMsgs           float64
	InErrors         float64
	InCsumErrors     float64
	InDestUnreachs   float64
	InTimeExcds      float64
	InParmProbs      float64
	InSrcQuenchs     float64
	InRedirects      float64
	InEchos          float64
	InEchoReps       float64
	InTimestamps     float64
	InTimestampReps  float64
	InAddrMasks      float64
	InAddrMaskReps   float64
	OutMsgs          float64
	OutErrors        float64
	OutDestUnreachs  float64
	OutTimeExcds     float64
	OutParmProbs     float64
	OutSrcQuenchs    float64
	OutRedirects     float64
	OutEchos         float64
	OutEchoReps      float64
	OutTimestamps    float64
	OutTimestampReps float64
	OutAddrMasks     float64
	OutAddrMaskReps  float64
}

type IcmpMsg struct {
	InType3  float64
	OutType3 float64
}

type Tcp struct {
	RtoAlgorithm float64
	RtoMin       float64
	RtoMax       float64
	MaxConn      float64
	ActiveOpens  float64
	PassiveOpens float64
	AttemptFails float64
	EstabResets  float64
	CurrEstab    float64
	InSegs       float64
	OutSegs      float64
	RetransSegs  float64
	InErrs       float64
	OutRsts      float64
	InCsumErrors float64
}

type Udp struct {
	InDatagrams  float64
	NoPorts      float64
	InErrors     float64
	OutDatagrams float64
	RcvbufErrors float64
	SndbufErrors float64
	InCsumErrors float64
	IgnoredMulti float64
}

type UdpLite struct {
	InDatagrams  float64
	NoPorts      float64
	InErrors     float64
	OutDatagrams float64
	RcvbufErrors float64
	SndbufErrors float64
	InCsumErrors float64
	IgnoredMulti float64
}

func (p Proc) Snmp() (ProcSnmp, error) {
	filename := p.path("net/snmp")
	procSnmp := ProcSnmp{PID: p.PID}

	data, err := util.ReadFileNoStat(filename)
	if err != nil {
		return procSnmp, err
	}

	netStats, err := parseSnmp(bytes.NewReader(data), filename)
	if err != nil {
		return procSnmp, err
	}

	mapStructureErr := mapstructure.Decode(netStats, &procSnmp)
	if mapStructureErr != nil {
		return procSnmp, mapStructureErr
	}

	return procSnmp, nil
}

// parseSnmp parses the metrics from proc/<pid>/net/snmp file
// and returns a map contains those metrics (e.g. {"Ip": {"Forwarding": 2}}).
func parseSnmp(r io.Reader, fileName string) (map[string]map[string]float64, error) {
	var (
		netStats = map[string]map[string]float64{}
		scanner  = bufio.NewScanner(r)
	)

	for scanner.Scan() {
		nameParts := strings.Split(scanner.Text(), " ")
		scanner.Scan()
		valueParts := strings.Split(scanner.Text(), " ")
		// Remove trailing :.
		protocol := nameParts[0][:len(nameParts[0])-1]
		netStats[protocol] = map[string]float64{}
		if len(nameParts) != len(valueParts) {
			return nil, fmt.Errorf("mismatch field count mismatch in %s: %s",
				fileName, protocol)
		}
		for i := 1; i < len(nameParts); i++ {
			var err error
			netStats[protocol][nameParts[i]], err = strconv.ParseFloat(valueParts[i], 64)
			if err != nil {
				return nil, err
			}
		}
	}

	return netStats, scanner.Err()
}
