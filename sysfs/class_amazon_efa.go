// Copyright 2023 Amazon Web Services
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

//go:build linux
// +build linux

package sysfs

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/prometheus/procfs/internal/util"
)

// Amazon Elastic Fabric Adapter counters are exposed similarly
// to InfiniBand counters in SysFS. The same structure is used
// in this class for consistency.

const AmazonEfaPath = "class/infiniband"

// AmazonEfaCounters contains counter values from files in
// /sys/class/infiniband/<Name>/ports/<Port>/hw_counters
// for a single port of one Amazon Elastic Fabric Adapter device.
type AmazonEfaCounters struct {
	AllocPdErr        *uint64 // hw_counters/alloc_pd_err
	AllocUcontextErr  *uint64 // hw_counters/alloc_ucontext_err
	CmdsErr           *uint64 // hw_counters/cmds_err
	CompletedCmds     *uint64 // hw_counters/completed_cmds
	CreateAhErr       *uint64 // hw_counters/create_ah_err
	CreateCqErr       *uint64 // hw_counters/create_cq_err
	CreateQpErr       *uint64 // hw_counters/create_qp_err
	KeepAliveRcvd     *uint64 // hw_counters/keep_alive_rcvd
	Lifespan          *uint64 // hw_counters/lifespan
	MmapErr           *uint64 // hw_counters/mmap_err
	NoCompletionCmds  *uint64 // hw_counters/no_completion_cmds
	RdmaReadBytes     *uint64 // hw_counters/rdma_read_bytes
	RdmaReadRespBytes *uint64 // hw_counters/rdma_read_resp_bytes
	RdmaReadWrErr     *uint64 // hw_counters/rdma_read_wr_err
	RdmaReadWrs       *uint64 // hw_counters/rdma_read_wrs
	RecvBytes         *uint64 // hw_counters/recv_bytes
	RecvWrs           *uint64 // hw_counters/recv_wrs
	RegMrErr          *uint64 // hw_counters/reg_mr_err
	RxBytes           *uint64 // hw_counters/rx_bytes
	RxDrops           *uint64 // hw_counters/rx_drops
	RxPkts            *uint64 // hw_counters/rx_pkts
	SendBytes         *uint64 // hw_counters/send_bytes
	SendWrs           *uint64 // hw_counters/send_wrs
	SubmittedCmds     *uint64 // hw_counters/submitted_cmds
	TxBytes           *uint64 // hw_counters/tx_bytes
	TxPkts            *uint64 // hw_counters/tx_pkts
}

// AmazonEfaPort contains info from files in
// /sys/class/infiniband/<Name>/ports/<Port>
// for a single port of one Amazon Elastic Fabric Adapter device.
type AmazonEfaPort struct {
	Name        string
	Port        uint
	State       string // String representation from /sys/class/infiniband/<Name>/ports/<Port>/state
	StateID     uint   // ID from /sys/class/infiniband/<Name>/ports/<Port>/state
	PhysState   string // String representation from /sys/class/infiniband/<Name>/ports/<Port>/phys_state
	PhysStateID uint   // String representation from /sys/class/infiniband/<Name>/ports/<Port>/phys_state
	Rate        uint64 // in bytes/second from /sys/class/infiniband/<Name>/ports/<Port>/rate
	Counters    AmazonEfaCounters
}

// AmazonEfaDevice contains info from files in /sys/class/infiniband for a
// single Amazon Elastic Fabric Adapter (EFA) device.
type AmazonEfaDevice struct {
	Name  string
	Ports map[uint]AmazonEfaPort
}

// AmazonEfaClass is a collection of every Amazon Elastic Fabric Adapter (EFA) device in
// /sys/class/infiniband.
//
// The map keys are the names of the Amazon Elastic Fabric Adapter (EFA) devices.
type AmazonEfaClass map[string]AmazonEfaDevice

// AmazonEfaClass returns info for all Amazon Elastic Fabric Adapter (EFA) devices read from
// /sys/class/infiniband.
func (fs FS) AmazonEfaClass() (AmazonEfaClass, error) {
	path := fs.sys.Path(AmazonEfaPath)

	dirs, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	aec := make(AmazonEfaClass, len(dirs))
	for _, d := range dirs {
		device, err := fs.parseAmazonEfaDevice(d.Name())
		if err != nil {
			return nil, err
		}

		aec[device.Name] = *device
	}

	return aec, nil
}

// Parse one AmazonEfa device.
func (fs FS) parseAmazonEfaDevice(name string) (*AmazonEfaDevice, error) {
	path := fs.sys.Path(AmazonEfaPath, name)
	device := AmazonEfaDevice{Name: name}

	portsPath := filepath.Join(path, "ports")
	ports, err := os.ReadDir(portsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to list AmazonEfa ports at %q: %w", portsPath, err)
	}

	device.Ports = make(map[uint]AmazonEfaPort, len(ports))
	for _, d := range ports {
		port, err := fs.parseAmazonEfaPort(name, d.Name())
		if err != nil {
			return nil, err
		}

		device.Ports[port.Port] = *port
	}

	return &device, nil
}

// Scans predefined files in /sys/class/infiniband/<device>/ports/<port>
// directory and gets their contents.
func (fs FS) parseAmazonEfaPort(name string, port string) (*AmazonEfaPort, error) {
	portNumber, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to convert %s into uint", port)
	}
	aep := AmazonEfaPort{Name: name, Port: uint(portNumber)}

	portPath := fs.sys.Path(AmazonEfaPath, name, "ports", port)
	content, err := os.ReadFile(filepath.Join(portPath, "state"))
	if err != nil {
		return nil, err
	}
	id, name, err := parseState(string(content))
	if err != nil {
		return nil, fmt.Errorf("could not parse state file in %q: %w", portPath, err)
	}
	aep.State = name
	aep.StateID = id

	content, err = os.ReadFile(filepath.Join(portPath, "phys_state"))
	if err != nil {
		return nil, err
	}
	id, name, err = parseState(string(content))
	if err != nil {
		return nil, fmt.Errorf("could not parse phys_state file in %q: %w", portPath, err)
	}
	aep.PhysState = name
	aep.PhysStateID = id

	content, err = os.ReadFile(filepath.Join(portPath, "rate"))
	if err != nil {
		return nil, err
	}
	aep.Rate, err = parseRate(string(content))
	if err != nil {
		return nil, fmt.Errorf("could not parse rate file in %q: %w", portPath, err)
	}

	counters, err := parseAmazonEfaCounters(portPath)
	if err != nil {
		return nil, err
	}
	aep.Counters = *counters

	return &aep, nil
}

func parseAmazonEfaCounters(portPath string) (*AmazonEfaCounters, error) {
	var counters AmazonEfaCounters

	path := filepath.Join(portPath, "hw_counters")
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if !f.Type().IsRegular() {
			continue
		}

		name := filepath.Join(path, f.Name())
		value, err := util.SysReadFile(name)
		if err != nil {
			if os.IsNotExist(err) || os.IsPermission(err) || err.Error() == "operation not supported" || err.Error() == "invalid argument" {
				continue
			}
			return nil, fmt.Errorf("failed to read file %q: %w", name, err)
		}

		vp := util.NewValueParser(value)

		switch f.Name() {

		case "lifespan":
			counters.Lifespan = vp.PUInt64()
		case "rdma_read_bytes":
			counters.RdmaReadBytes = vp.PUInt64()
		case "rdma_read_resp_bytes":
			counters.RdmaReadRespBytes = vp.PUInt64()
		case "rdma_read_wr_err":
			counters.RdmaReadWrErr = vp.PUInt64()
		case "rdma_read_wrs":
			counters.RdmaReadWrs = vp.PUInt64()
		case "recv_bytes":
			counters.RecvBytes = vp.PUInt64()
		case "recv_wrs":
			counters.RecvWrs = vp.PUInt64()
		case "rx_bytes":
			counters.RxBytes = vp.PUInt64()
		case "rx_drops":
			counters.RxDrops = vp.PUInt64()
		case "rx_pkts":
			counters.RxPkts = vp.PUInt64()
		case "send_bytes":
			counters.SendBytes = vp.PUInt64()
		case "send_wrs":
			counters.SendWrs = vp.PUInt64()
		case "tx_bytes":
			counters.TxBytes = vp.PUInt64()
		case "tx_pkts":
			counters.TxPkts = vp.PUInt64()

			if err != nil {
				// Ugly workaround for handling https://github.com/prometheus/node_exporter/issues/966
				// when counters are `N/A (not available)`.
				// This was already patched and submitted, see
				// https://www.spinics.net/lists/linux-rdma/msg68596.html
				// Remove this as soon as the fix lands in the enterprise distros.
				if strings.Contains(value, "N/A (no PMA)") {
					continue
				}
				return nil, err
			}
		}
	}

	return &counters, nil
}
