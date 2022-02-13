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
	"errors"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/prometheus/procfs/internal/util"
)

// ProcSnmp6 models the content of /proc/<pid>/net/snmp6.
type ProcSnmp6 struct {
	// The process ID.
	PID int
	Ip6
	Icmp6
	Udp6
	Udp6Lite
}

type Ip6 struct {
	InReceives       float64
	InHdrErrors      float64
	InTooBigErrors   float64
	InNoRoutes       float64
	InAddrErrors     float64
	InUnknownProtos  float64
	InTruncatedPkts  float64
	InDiscards       float64
	InDelivers       float64
	OutForwDatagrams float64
	OutRequests      float64
	OutDiscards      float64
	OutNoRoutes      float64
	ReasmTimeout     float64
	ReasmReqds       float64
	ReasmOKs         float64
	ReasmFails       float64
	FragOKs          float64
	FragFails        float64
	FragCreates      float64
	InMcastPkts      float64
	OutMcastPkts     float64
	InOctets         float64
	OutOctets        float64
	InMcastOctets    float64
	OutMcastOctets   float64
	InBcastOctets    float64
	OutBcastOctets   float64
	InNoECTPkts      float64
	InECT1Pkts       float64
	InECT0Pkts       float64
	InCEPkts         float64
}

type Icmp6 struct {
	InMsgs                    float64
	InErrors                  float64
	OutMsgs                   float64
	OutErrors                 float64
	InCsumErrors              float64
	InDestUnreachs            float64
	InPktTooBigs              float64
	InTimeExcds               float64
	InParmProblems            float64
	InEchos                   float64
	InEchoReplies             float64
	InGroupMembQueries        float64
	InGroupMembResponses      float64
	InGroupMembReductions     float64
	InRouterSolicits          float64
	InRouterAdvertisements    float64
	InNeighborSolicits        float64
	InNeighborAdvertisements  float64
	InRedirects               float64
	InMLDv2Reports            float64
	OutDestUnreachs           float64
	OutPktTooBigs             float64
	OutTimeExcds              float64
	OutParmProblems           float64
	OutEchos                  float64
	OutEchoReplies            float64
	OutGroupMembQueries       float64
	OutGroupMembResponses     float64
	OutGroupMembReductions    float64
	OutRouterSolicits         float64
	OutRouterAdvertisements   float64
	OutNeighborSolicits       float64
	OutNeighborAdvertisements float64
	OutRedirects              float64
	OutMLDv2Reports           float64
	InType1                   float64
	InType134                 float64
	InType135                 float64
	InType136                 float64
	InType143                 float64
	OutType133                float64
	OutType135                float64
	OutType136                float64
	OutType143                float64
}

type Udp6 struct {
	InDatagrams  float64
	NoPorts      float64
	InErrors     float64
	OutDatagrams float64
	RcvbufErrors float64
	SndbufErrors float64
	InCsumErrors float64
	IgnoredMulti float64
}

type Udp6Lite struct {
	InDatagrams  float64
	NoPorts      float64
	InErrors     float64
	OutDatagrams float64
	RcvbufErrors float64
	SndbufErrors float64
	InCsumErrors float64
}

func (p Proc) Snmp6() (ProcSnmp6, error) {
	filename := p.path("net/snmp6")
	procSnmp6 := ProcSnmp6{PID: p.PID}

	data, err := util.ReadFileNoStat(filename)
	if err != nil {
		// On systems with IPv6 disabled, this file won't exist.
		// Do nothing.
		if errors.Is(err, os.ErrNotExist) {
			return procSnmp6, nil
		}

		return procSnmp6, err
	}

	netStats, err := parseSNMP6Stats(bytes.NewReader(data))
	if err != nil {
		return procSnmp6, err
	}

	mapStructureErr := mapstructure.Decode(netStats, &procSnmp6)
	if mapStructureErr != nil {
		return procSnmp6, mapStructureErr
	}

	return procSnmp6, nil

}

// parseSnmp6 parses the metrics from proc/<pid>/net/snmp6 file
// and returns a map contains those metrics.
func parseSNMP6Stats(r io.Reader) (map[string]map[string]float64, error) {
	var (
		netStats = map[string]map[string]float64{}
		scanner  = bufio.NewScanner(r)
	)

	for scanner.Scan() {
		stat := strings.Fields(scanner.Text())
		if len(stat) < 2 {
			continue
		}
		// Expect to have "6" in metric name, skip line otherwise
		if sixIndex := strings.Index(stat[0], "6"); sixIndex != -1 {
			protocol := stat[0][:sixIndex+1]
			name := stat[0][sixIndex+1:]
			if _, present := netStats[protocol]; !present {
				netStats[protocol] = map[string]float64{}
			}
			var err error
			netStats[protocol][name], err = strconv.ParseFloat(stat[1], 64)
			if err != nil {
				return nil, err
			}
		}
	}

	return netStats, scanner.Err()
}
