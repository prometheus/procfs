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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/prometheus/procfs/internal/util"
)

const netclassPath = "class/net"

// NetClassIface contains info from files in /sys/class/net/<iface>
// for single interface (iface).
type NetClassIface struct {
	Name             string // Interface name
	AddrAssignType   *int64 // /sys/class/net/<iface>/addr_assign_type
	AddrLen          *int64 // /sys/class/net/<iface>/addr_len
	Address          string // /sys/class/net/<iface>/address
	Broadcast        string // /sys/class/net/<iface>/broadcast
	Carrier          *int64 // /sys/class/net/<iface>/carrier
	CarrierChanges   *int64 // /sys/class/net/<iface>/carrier_changes
	CarrierUpCount   *int64 // /sys/class/net/<iface>/carrier_up_count
	CarrierDownCount *int64 // /sys/class/net/<iface>/carrier_down_count
	DevID            *int64 // /sys/class/net/<iface>/dev_id
	Dormant          *int64 // /sys/class/net/<iface>/dormant
	Duplex           string // /sys/class/net/<iface>/duplex
	Flags            *int64 // /sys/class/net/<iface>/flags
	IfAlias          string // /sys/class/net/<iface>/ifalias
	IfIndex          *int64 // /sys/class/net/<iface>/ifindex
	IfLink           *int64 // /sys/class/net/<iface>/iflink
	LinkMode         *int64 // /sys/class/net/<iface>/link_mode
	MTU              *int64 // /sys/class/net/<iface>/mtu
	NameAssignType   *int64 // /sys/class/net/<iface>/name_assign_type
	NetDevGroup      *int64 // /sys/class/net/<iface>/netdev_group
	OperState        string // /sys/class/net/<iface>/operstate
	PhysPortID       string // /sys/class/net/<iface>/phys_port_id
	PhysPortName     string // /sys/class/net/<iface>/phys_port_name
	PhysSwitchID     string // /sys/class/net/<iface>/phys_switch_id
	Speed            *int64 // /sys/class/net/<iface>/speed
	TxQueueLen       *int64 // /sys/class/net/<iface>/tx_queue_len
	Type             *int64 // /sys/class/net/<iface>/type
	Collisions       *int64 // /sys/class/net/<iface>/statistics/collisions
	Multicast        *int64 // /sys/class/net/<iface>/statistics/multicast
	RxBytes          *int64 // /sys/class/net/<iface>/statistics/rx_bytes
	RxCompressed     *int64 // /sys/class/net/<iface>/statistics/rx_compressed
	RxCrcErrors      *int64 // /sys/class/net/<iface>/statistics/rx_crc_errors
	RxDropped        *int64 // /sys/class/net/<iface>/statistics/rx_dropped
	RxErrors         *int64 // /sys/class/net/<iface>/statistics/rx_errors
	RxFifoErrors     *int64 // /sys/class/net/<iface>/statistics/rx_fifo_errors
	RxFrameErrors    *int64 // /sys/class/net/<iface>/statistics/rx_frame_errors
	RxLengthErrors   *int64 // /sys/class/net/<iface>/statistics/rx_length_errors
	RxMissedErrors   *int64 // /sys/class/net/<iface>/statistics/rx_missed_errors
	RxNoHandler      *int64 // /sys/class/net/<iface>/statistics/rx_nohandler
	RxOverErrors     *int64 // /sys/class/net/<iface>/statistics/rx_over_errors
	RxPackets        *int64 // /sys/class/net/<iface>/statistics/rx_over_packets
	TxAbortedErrors  *int64 // /sys/class/net/<iface>/statistics/tx_aborted_errors
	TxBytes          *int64 // /sys/class/net/<iface>/statistics/tx_bytes
	TxCarrierErrors  *int64 // /sys/class/net/<iface>/statistics/tx_carrier_errors
	TxCompressed     *int64 // /sys/class/net/<iface>/statistics/tx_compressed
	TxDropped        *int64 // /sys/class/net/<iface>/statistics/tx_dropped
	TxErrors         *int64 // /sys/class/net/<iface>/statistics/tx_errors
	TxFifoErrors     *int64 // /sys/class/net/<iface>/statistics/tx_fifo_errors
	TxHeartbtErrors  *int64 // /sys/class/net/<iface>/statistics/tx_heartbeat_errors
	TxPackets        *int64 // /sys/class/net/<iface>/statistics/tx_packets
	TxWindowErrors   *int64 // /sys/class/net/<iface>/statistics/tx_window_errors
}

// NetClass is collection of info for every interface (iface) in /sys/class/net. The map keys
// are interface (iface) names.
type NetClass map[string]NetClassIface

// NetClassDevices scans /sys/class/net for devices and returns them as a list of names.
func (fs FS) NetClassDevices() ([]string, error) {
	var res []string
	path := fs.sys.Path(netclassPath)

	devices, err := ioutil.ReadDir(path)
	if err != nil {
		return res, fmt.Errorf("cannot access %s dir %s", path, err)
	}

	for _, deviceDir := range devices {
		if deviceDir.Mode().IsRegular() {
			continue
		}
		res = append(res, deviceDir.Name())
	}

	return res, nil
}

// NetClass returns info for all net interfaces (iface) read from /sys/class/net/<iface>.
func (fs FS) NetClass() (NetClass, error) {
	devices, err := fs.NetClassDevices()
	if err != nil {
		return nil, err
	}

	path := fs.sys.Path(netclassPath)
	netClass := NetClass{}
	for _, deviceDir := range devices {
		interfaceClass, err := netClass.parseNetClassIface(filepath.Join(path, deviceDir))
		if err != nil {
			return nil, err
		}
		interfaceClass.Name = deviceDir
		netClass[deviceDir] = *interfaceClass
	}
	return netClass, nil
}

// parseNetClassIface scans predefined files in /sys/class/net/<iface>
// directory and gets their contents.
func (nc NetClass) parseNetClassIface(devicePath string) (*NetClassIface, error) {
	interfaceClass := NetClassIface{}

	files, err := ioutil.ReadDir(devicePath)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if !f.Mode().IsRegular() && f.Name() != "statistics" {
			continue
		}

		if f.Name() == "statistics" {
			statisticsFiles, err := ioutil.ReadDir(devicePath + "/statistics")
			if err != nil {
				return nil, err
			}
			for _, statisticsFile := range statisticsFiles {
				statisticName := filepath.Join(devicePath+"/statistics", statisticsFile.Name())
				value, err := util.SysReadFile(statisticName)
				if err != nil {
					if os.IsNotExist(err) || os.IsPermission(err) || err.Error() == "operation not supported" || err.Error() == "invalid argument" {
						continue
					}
					return nil, fmt.Errorf("failed to read file %q: %v", statisticName, err)
				}
				vp := util.NewValueParser(value)
				switch statisticsFile.Name() {
				case "collisions":
					interfaceClass.Collisions = vp.PInt64()
				case "multicast":
					interfaceClass.Multicast = vp.PInt64()
				case "rx_bytes":
					interfaceClass.RxBytes = vp.PInt64()
				case "rx_compressed":
					interfaceClass.RxCompressed = vp.PInt64()
				case "rx_crc_errors":
					interfaceClass.RxCrcErrors = vp.PInt64()
				case "rx_dropped":
					interfaceClass.RxDropped = vp.PInt64()
				case "rx_errors":
					interfaceClass.RxErrors = vp.PInt64()
				case "rx_fifo_errors":
					interfaceClass.RxFifoErrors = vp.PInt64()
				case "rx_frame_errors":
					interfaceClass.RxFrameErrors = vp.PInt64()
				case "rx_length_errors":
					interfaceClass.RxLengthErrors = vp.PInt64()
				case "rx_missed_errors":
					interfaceClass.RxMissedErrors = vp.PInt64()
				case "rx_nohandler":
					interfaceClass.RxNoHandler = vp.PInt64()
				case "rx_over_errors":
					interfaceClass.RxOverErrors = vp.PInt64()
				case "rx_packets":
					interfaceClass.RxPackets = vp.PInt64()
				case "tx_aborted_errors":
					interfaceClass.TxAbortedErrors = vp.PInt64()
				case "tx_bytes":
					interfaceClass.TxBytes = vp.PInt64()
				case "tx_carrier_errors":
					interfaceClass.TxCarrierErrors = vp.PInt64()
				case "tx_compressed":
					interfaceClass.TxCompressed = vp.PInt64()
				case "tx_dropped":
					interfaceClass.TxDropped = vp.PInt64()
				case "tx_errors":
					interfaceClass.TxErrors = vp.PInt64()
				case "tx_fifo_errors":
					interfaceClass.TxFifoErrors = vp.PInt64()
				case "tx_heartbeat_errors":
					interfaceClass.TxHeartbtErrors = vp.PInt64()
				case "tx_packets":
					interfaceClass.TxPackets = vp.PInt64()
				case "tx_window_errors":
					interfaceClass.TxWindowErrors = vp.PInt64()
				}
			}
		} else {
			name := filepath.Join(devicePath, f.Name())
			value, err := util.SysReadFile(name)
			if err != nil {
				if os.IsNotExist(err) || os.IsPermission(err) || err.Error() == "operation not supported" || err.Error() == "invalid argument" {
					continue
				}
				return nil, fmt.Errorf("failed to read file %q: %v", name, err)
			}
			vp := util.NewValueParser(value)
			switch f.Name() {
			case "addr_assign_type":
				interfaceClass.AddrAssignType = vp.PInt64()
			case "addr_len":
				interfaceClass.AddrLen = vp.PInt64()
			case "address":
				interfaceClass.Address = value
			case "broadcast":
				interfaceClass.Broadcast = value
			case "carrier":
				interfaceClass.Carrier = vp.PInt64()
			case "carrier_changes":
				interfaceClass.CarrierChanges = vp.PInt64()
			case "carrier_up_count":
				interfaceClass.CarrierUpCount = vp.PInt64()
			case "carrier_down_count":
				interfaceClass.CarrierDownCount = vp.PInt64()
			case "dev_id":
				interfaceClass.DevID = vp.PInt64()
			case "dormant":
				interfaceClass.Dormant = vp.PInt64()
			case "duplex":
				interfaceClass.Duplex = value
			case "flags":
				interfaceClass.Flags = vp.PInt64()
			case "ifalias":
				interfaceClass.IfAlias = value
			case "ifindex":
				interfaceClass.IfIndex = vp.PInt64()
			case "iflink":
				interfaceClass.IfLink = vp.PInt64()
			case "link_mode":
				interfaceClass.LinkMode = vp.PInt64()
			case "mtu":
				interfaceClass.MTU = vp.PInt64()
			case "name_assign_type":
				interfaceClass.NameAssignType = vp.PInt64()
			case "netdev_group":
				interfaceClass.NetDevGroup = vp.PInt64()
			case "operstate":
				interfaceClass.OperState = value
			case "phys_port_id":
				interfaceClass.PhysPortID = value
			case "phys_port_name":
				interfaceClass.PhysPortName = value
			case "phys_switch_id":
				interfaceClass.PhysSwitchID = value
			case "speed":
				interfaceClass.Speed = vp.PInt64()
			case "tx_queue_len":
				interfaceClass.TxQueueLen = vp.PInt64()
			case "type":
				interfaceClass.Type = vp.PInt64()
			}
		}
	}
	return &interfaceClass, nil

}
