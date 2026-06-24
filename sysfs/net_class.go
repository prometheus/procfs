// Copyright The Prometheus Authors
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

package sysfs

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/prometheus/procfs/internal/util"
)

const netclassPath = "class/net"

// NetClassBondAttrs contains info from files in /sys/class/net/<iface>/bonding
// for a bonding controller interface (iface).
type NetClassBondAttrs struct {
	ActiveDevice                           *NetClassIface    // /sys/class/net/<iface>/bonding/active_slave
	AdActorKey                             *uint64           // /sys/class/net/<iface>/bonding/ad_actor_key (Requires CAP_NET_ADMIN)
	AdActorSysPriority                     *uint64           // /sys/class/net/<iface>/bonding/ad_actor_sys_prio (Requires CAP_NET_ADMIN)
	AdActorSystem                          *net.HardwareAddr // /sys/class/net/<iface>/bonding/ad_actor_system (Requires CAP_NET_ADMIN)
	AdAggregator                           *uint64           // /sys/class/net/<iface>/bonding/ad_aggregator
	AdNumPorts                             *uint64           // /sys/class/net/<iface>/bonding/ad_num_ports
	AdPartnerKey                           *uint64           // /sys/class/net/<iface>/bonding/ad_partner_key (Requires CAP_NET_ADMIN)
	AdPartnerMac                           *net.HardwareAddr // /sys/class/net/<iface>/bonding/ad_partner_mac (Requires CAP_NET_ADMIN)
	AdSelect                               *string           // String representation from /sys/class/net/<iface>/bonding/ad_select
	AdSelectID                             *uint64           // ID from /sys/class/net/<iface>/bonding/ad_select
	AdUserPortKey                          *uint64           // /sys/class/net/<iface>/bonding/ad_user_port_key (Requires CAP_NET_ADMIN)
	AllDevicesActive                       *bool             // /sys/class/net/<iface>/bonding/all_slaves_active
	ARPAllTargets                          *string           // String representation from /sys/class/net/<iface>/bonding/arp_all_targets
	ARPAllTargetsID                        *uint64           // ID from /sys/class/net/<iface>/bonding/arp_all_targets
	ARPInterval                            *int64            // /sys/class/net/<iface>/bonding/arp_interval
	ARPIPTarget                            []net.IP          // /sys/class/net/<iface>/bonding/arp_ip_target
	ARPValidate                            *string           // String representation from /sys/class/net/<iface>/bonding/arp_validate
	ARPValidateID                          *uint64           // ID from /sys/class/net/<iface>/bonding/arp_validate
	DownDelay                              *int64            // /sys/class/net/<iface>/bonding/downdelay
	FailoverMac                            *string           // String representation from /sys/class/net/<iface>/bonding/failover_mac
	FailoverMacID                          *uint64           // ID from /sys/class/net/<iface>/bonding/failover_mac
	LACPRate                               *string           // String representation from /sys/class/net/<iface>/bonding/lacp_rate
	LACPRateID                             *uint64           // ID from /sys/class/net/<iface>/bonding/lacp_rate
	LPInterval                             *int64            // /sys/class/net/<iface>/bonding/lp_interval
	MIIMon                                 *int64            // /sys/class/net/<iface>/bonding/miimon
	MIIStatus                              *bool             // /sys/class/net/<iface>/bonding/mii_status
	MinLinks                               *uint64           // /sys/class/net/<iface>/bonding/min_links
	Mode                                   *string           // String representation from /sys/class/net/<iface>/bonding/mode
	ModeID                                 *uint64           // ID from /sys/class/net/<iface>/bonding/mode
	NumberGratuitousArp                    *uint64           // /sys/class/net/<iface>/bonding/num_grat_arp
	NumberUnsolicitedNeighborAdvertisement *uint64           // /sys/class/net/<iface>/bonding/num_unsol_na
	PacketsPerDevice                       *int64            // /sys/class/net/<iface>/bonding/packets_per_slave
	PrimaryDevice                          *string           // /sys/class/net/<iface>/bonding/primary
	PrimaryReselect                        *string           // String representation from /sys/class/net/<iface>/bonding/primary_reselect
	PrimaryReselectID                      *uint64           // ID from /sys/class/net/<iface>/bonding/primary_reselect
	DeviceQueueIDs                         map[string]uint64 // /sys/class/net/<iface>/bonding/queue_id
	ResendIgmp                             *int64            // /sys/class/net/<iface>/bonding/resend_igmp
	Devices                                []*NetClassIface  // /sys/class/net/<iface>/bonding/slaves
	TLBDynamicLB                           *int64            // /sys/class/net/<iface>/bonding/tlb_dynamic_lb
	UpDelay                                *int64            // /sys/class/net/<iface>/bonding/updelay
	UseCarrier                             *bool             // /sys/class/net/<iface>/bonding/use_carrier
	TransmitHashPolicy                     *string           // String representation from /sys/class/net/<iface>/bonding/xmit_hash_policy
	TransmitHashPolicyID                   *uint64           // ID from /sys/class/net/<iface>/bonding/xmit_hash_policy
}

// NetClassBondDeviceAttrs contains info from files in /sys/class/net/<iface>/bonding_slave
// for a bonding device interface (iface).
type NetClassBondDeviceAttrs struct {
	Controller                    *NetClassIface    // /sys/class/net/<iface>/master
	AdActorOperationalPortState   *uint64           // /sys/class/net/<iface>/bonding_slave/ad_actor_oper_port_state
	AdAggregatorID                *uint64           // /sys/class/net/<iface>/bonding_slave/ad_aggregator_id
	AdPartnerOperationalPortState *uint64           // /sys/class/net/<iface>/bonding_slave/ad_partner_oper_port_state
	LinkFailureCount              *uint64           // /sys/class/net/<iface>/bonding_slave/link_failure_count
	MiiStatus                     *bool             // /sys/class/net/<iface>/bonding_slave/mii_status
	PermamentHWAddress            *net.HardwareAddr // /sys/class/net/<iface>/bonding_slave/perm_hwaddr
	QueueID                       *uint64           // /sys/class/net/<iface>/bonding_slave/queue_id
	State                         *uint64           // /sys/class/net/<iface>/bonding_slave/state
}

// NetClassIface contains info from files in /sys/class/net/<iface>
// for single interface (iface).
type NetClassIface struct {
	Name             string                   // Interface name
	AddrAssignType   *int64                   // /sys/class/net/<iface>/addr_assign_type
	AddrLen          *int64                   // /sys/class/net/<iface>/addr_len
	Address          string                   // /sys/class/net/<iface>/address
	Broadcast        string                   // /sys/class/net/<iface>/broadcast
	BondAttrs        *NetClassBondAttrs       // /sys/class/net/<iface>/bonding
	BondDeviceAttrs  *NetClassBondDeviceAttrs // /sys/class/net/<iface>/bonding_slave
	Carrier          *int64                   // /sys/class/net/<iface>/carrier
	CarrierChanges   *int64                   // /sys/class/net/<iface>/carrier_changes
	CarrierUpCount   *int64                   // /sys/class/net/<iface>/carrier_up_count
	CarrierDownCount *int64                   // /sys/class/net/<iface>/carrier_down_count
	DevID            *int64                   // /sys/class/net/<iface>/dev_id
	Dormant          *int64                   // /sys/class/net/<iface>/dormant
	Duplex           string                   // /sys/class/net/<iface>/duplex
	Flags            *int64                   // /sys/class/net/<iface>/flags
	IfAlias          string                   // /sys/class/net/<iface>/ifalias
	IfIndex          *int64                   // /sys/class/net/<iface>/ifindex
	IfLink           *int64                   // /sys/class/net/<iface>/iflink
	LinkMode         *int64                   // /sys/class/net/<iface>/link_mode
	MTU              *int64                   // /sys/class/net/<iface>/mtu
	NameAssignType   *int64                   // /sys/class/net/<iface>/name_assign_type
	NetDevGroup      *int64                   // /sys/class/net/<iface>/netdev_group
	OperState        string                   // /sys/class/net/<iface>/operstate
	PhysPortID       string                   // /sys/class/net/<iface>/phys_port_id
	PhysPortName     string                   // /sys/class/net/<iface>/phys_port_name
	PhysSwitchID     string                   // /sys/class/net/<iface>/phys_switch_id
	Speed            *int64                   // /sys/class/net/<iface>/speed
	TxQueueLen       *int64                   // /sys/class/net/<iface>/tx_queue_len
	Type             *int64                   // /sys/class/net/<iface>/type
}

// NetClass is collection of info for every interface (iface) in /sys/class/net. The map keys
// are interface (iface) names.
type NetClass map[string]NetClassIface

// NetClassDevices scans /sys/class/net for devices and returns them as a list of names.
func (fs FS) NetClassDevices() ([]string, error) {
	var res []string
	path := fs.sys.Path(netclassPath)

	devices, err := os.ReadDir(path)
	if err != nil {
		return res, fmt.Errorf("cannot access dir %q: %w", path, err)
	}

	for _, deviceDir := range devices {
		if deviceDir.Type().IsRegular() {
			continue
		}
		res = append(res, deviceDir.Name())
	}

	return res, nil
}

// NetClassByIface returns info for a single net interfaces (iface).
func (fs FS) NetClassByIface(devicePath string) (*NetClassIface, error) {
	devices, err := fs.NetClass()
	if err != nil {
		return nil, err
	}
	if device, found := devices[devicePath]; found {
		return &device, nil
	}
	return nil, fmt.Errorf("device %s not found", devicePath)
}

// NetClass returns info for all net interfaces (iface) read from /sys/class/net/<iface>.
func (fs FS) NetClass() (NetClass, error) {
	devices, err := fs.NetClassDevices()
	if err != nil {
		return nil, err
	}

	path := fs.sys.Path(netclassPath)
	netClass := NetClass{}
	for _, devicePath := range devices {
		interfaceClass, err := parseNetClassIface(filepath.Join(path, devicePath))
		if err != nil {
			return nil, err
		}
		interfaceClass.Name = devicePath
		netClass[devicePath] = *interfaceClass
	}
	if err := fs.resolveBondingRelationships(&netClass); err != nil {
		return nil, err
	}

	return netClass, nil
}

// canIgnoreError returns true if the error is non-fatal and can be ignored.
// Some kernels and some devices don't expose specific attributes or return
// errors when reading those attributes; we can ignore these errors and the
// attribute that caused them.
func canIgnoreError(err error) bool {
	var errno syscall.Errno

	switch {
	case os.IsNotExist(err):
		return true
	case os.IsPermission(err):
		return true
	case err.Error() == "operation not supported":
		return true
	case errors.Is(err, os.ErrInvalid):
		return true
	case errors.As(err, &errno) && (errno == syscall.EINVAL):
		return true
	}
	// all other errors are fatal
	return false
}

// ParseNetClassAttribute parses a given file in /sys/class/net/<iface>
// and sets the value in a given NetClassIface object if the value was readable.
// It returns an error if the file cannot be read and the error is fatal.
func ParseNetClassAttribute(devicePath, attrName string, interfaceClass *NetClassIface) error {
	attrPath := filepath.Join(devicePath, attrName)
	value, err := util.SysReadFile(attrPath)
	if err != nil {
		if canIgnoreError(err) {
			return nil
		}
		return fmt.Errorf("failed to read file %q: %w", attrPath, err)
	}

	vp := util.NewValueParser(value)
	switch attrName {
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
	default:
		return nil
	}

	return nil
}

// parseNetClassIface scans predefined files in /sys/class/net/<iface>
// directory and gets their contents.
func parseNetClassIface(devicePath string) (*NetClassIface, error) {
	interfaceClass := NetClassIface{}

	files, err := os.ReadDir(devicePath)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if !f.Type().IsRegular() {
			continue
		}
		if err := ParseNetClassAttribute(devicePath, f.Name(), &interfaceClass); err != nil {
			return nil, err
		}
	}
	bondingPath := filepath.Join(devicePath, "bonding")
	if _, err := os.Stat(bondingPath); !os.IsNotExist(err) {
		interfaceClass.BondAttrs, err = parseNetClassBondAttrs(bondingPath)
		if err != nil {
			return nil, err
		}
	}
	bondingDevicesPath := filepath.Join(devicePath, "bonding_slave")
	if _, err := os.Stat(bondingDevicesPath); !os.IsNotExist(err) {
		interfaceClass.BondDeviceAttrs, err = parseNetClassBondDeviceAttrs(bondingDevicesPath)
		if err != nil {
			return nil, err
		}
	}

	return &interfaceClass, nil
}

func parseNetClassBondAttrs(devicePath string) (*NetClassBondAttrs, error) {
	attrs := NetClassBondAttrs{}

	files, err := os.ReadDir(devicePath)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if !f.Type().IsRegular() {
			continue
		}
		fileName := filepath.Join(devicePath, f.Name())
		value, err := util.SysReadFile(fileName)
		if err != nil {
			if os.IsNotExist(err) || os.IsPermission(err) || err.Error() == "operation not supported" || err.Error() == "invalid argument" {
				continue
			}
			return nil, fmt.Errorf("failed to read file %q: %w", fileName, err)
		}
		vp := util.NewValueParser(value)
		switch f.Name() {
		case "ad_actor_key":
			attrs.AdActorKey = vp.PUInt64()
		case "ad_actor_sys_prio":
			attrs.AdActorSysPriority = vp.PUInt64()
		case "ad_actor_system":
			mac, err := net.ParseMAC(value)
			if err != nil {
				return nil, err
			}
			attrs.AdActorSystem = &mac
		case "ad_aggregator":
			attrs.AdAggregator = vp.PUInt64()
		case "ad_num_ports":
			attrs.AdNumPorts = vp.PUInt64()
		case "ad_partner_key":
			attrs.AdPartnerKey = vp.PUInt64()
		case "ad_partner_mac":
			mac, err := net.ParseMAC(value)
			if err != nil {
				return nil, err
			}
			attrs.AdPartnerMac = &mac
		case "ad_select":
			name, id, err := parseBondOption(value)
			if err != nil {
				return nil, fmt.Errorf("could not parse ad_select file in %q: %w", fileName, err)
			}
			attrs.AdSelect = &name
			attrs.AdSelectID = &id
		case "ad_user_port_key":
			attrs.AdUserPortKey = vp.PUInt64()
		case "all_slaves_active":
			attrs.AllDevicesActive = util.ParseBool(value)
		case "arp_all_targets":
			name, id, err := parseBondOption(value)
			if err != nil {
				return nil, fmt.Errorf("could not parse arp_all_targets file in %q: %w", fileName, err)
			}
			attrs.ARPAllTargets = &name
			attrs.ARPAllTargetsID = &id
		case "arp_interval":
			attrs.ARPInterval = vp.PInt64()
		case "arp_ip_target":
			ips, err := parseArpTargets(value)
			if err != nil {
				return nil, err
			}
			attrs.ARPIPTarget = ips
		case "arp_validate":
			name, id, err := parseBondOption(value)
			if err != nil {
				return nil, fmt.Errorf("could not parse arp_validate file in %q: %w", fileName, err)
			}
			attrs.ARPValidate = &name
			attrs.ARPValidateID = &id
		case "downdelay":
			attrs.DownDelay = vp.PInt64()
		case "fail_over_mac":
			name, id, err := parseBondOption(value)
			if err != nil {
				return nil, fmt.Errorf("could not parse fail_over_mac file in %q: %w", fileName, err)
			}
			attrs.FailoverMac = &name
			attrs.FailoverMacID = &id
		case "lacp_rate":
			name, id, err := parseBondOption(value)
			if err != nil {
				return nil, fmt.Errorf("could not parse lacp_rate file in %q: %w", fileName, err)
			}
			attrs.LACPRate = &name
			attrs.LACPRateID = &id
		case "lp_interval":
			attrs.LPInterval = vp.PInt64()
		case "miimon":
			attrs.MIIMon = vp.PInt64()
		case "mii_status":
			attrs.MIIStatus = util.ParseBool(value)
		case "min_links":
			attrs.MinLinks = vp.PUInt64()
		case "mode":
			name, id, err := parseBondOption(value)
			if err != nil {
				return nil, fmt.Errorf("could not parse mode file in %q: %w", fileName, err)
			}
			attrs.Mode = &name
			attrs.ModeID = &id
		case "num_grat_arp":
			attrs.NumberGratuitousArp = vp.PUInt64()
		case "num_unsol_na":
			attrs.NumberUnsolicitedNeighborAdvertisement = vp.PUInt64()
		case "queue_id":
			ids, err := parseDeviceQueueIDs(value)
			if err != nil {
				return nil, err
			}
			attrs.DeviceQueueIDs = ids
		case "packets_per_slave":
			attrs.PacketsPerDevice = vp.PInt64()
		case "primary_reselect":
			name, id, err := parseBondOption(value)
			if err != nil {
				return nil, fmt.Errorf("could not parse primary_reselect file in %q: %w", fileName, err)
			}
			attrs.PrimaryReselect = &name
			attrs.PrimaryReselectID = &id
		case "resend_igmp":
			attrs.ResendIgmp = vp.PInt64()
		case "tlb_dynamic_lb":
			attrs.TLBDynamicLB = vp.PInt64()
		case "updelay":
			attrs.UpDelay = vp.PInt64()
		case "use_carrier":
			attrs.UseCarrier = util.ParseBool(value)
		case "xmit_hash_policy":
			name, id, err := parseBondOption(value)
			if err != nil {
				return nil, fmt.Errorf("could not parse xmit_hash_policy file in %q: %w", fileName, err)
			}
			attrs.TransmitHashPolicy = &name
			attrs.TransmitHashPolicyID = &id
		}
	}
	return &attrs, nil
}

func parseNetClassBondDeviceAttrs(devicePath string) (*NetClassBondDeviceAttrs, error) {
	attrs := NetClassBondDeviceAttrs{}
	files, err := os.ReadDir(devicePath)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if !f.Type().IsRegular() {
			continue
		}
		name := filepath.Join(devicePath, f.Name())
		value, err := util.SysReadFile(name)
		if err != nil {
			if os.IsNotExist(err) || os.IsPermission(err) || err.Error() == "operation not supported" || err.Error() == "invalid argument" {
				continue
			}
			return nil, fmt.Errorf("failed to read file %q: %w", name, err)
		}
		vp := util.NewValueParser(value)
		switch f.Name() {
		case "ad_actor_oper_port_state":
			attrs.AdActorOperationalPortState = vp.PUInt64()
		case "ad_aggregator_id":
			attrs.AdAggregatorID = vp.PUInt64()
		case "ad_partner_oper_port_state":
			attrs.AdPartnerOperationalPortState = vp.PUInt64()
		case "link_failure_count":
			attrs.LinkFailureCount = vp.PUInt64()
		case "mii_status":
			attrs.MiiStatus = util.ParseBool(value)
		case "perm_hwaddr":
			mac, err := net.ParseMAC(value)
			if err != nil {
				return nil, err
			}
			attrs.PermamentHWAddress = &mac
		case "queue_id":
			attrs.QueueID = vp.PUInt64()
		case "state":
			attrs.State = vp.PUInt64()
		}
	}
	return &attrs, nil
}

func (fs FS) resolveBondingRelationships(netClass *NetClass) error {
	for _, netClassIface := range *netClass {
		if netClassIface.BondAttrs != nil {
			path := filepath.Join(fs.sys.Path(netclassPath), netClassIface.Name, "bonding")
			if _, err := os.Stat(filepath.Join(path, "active_slave")); !os.IsNotExist(err) {
				activeSlave, err := util.SysReadFile(filepath.Join(path, "active_slave"))
				if err != nil {
					return fmt.Errorf("unable to read %s", filepath.Join(path, "active_slave"))
				}
				if activeSlave != "" {
					if intf, exists := (*netClass)[activeSlave]; exists {
						netClassIface.BondAttrs.ActiveDevice = &intf
					} else {
						return fmt.Errorf("unable to find device %s", activeSlave)
					}
				}
			}
			if _, err := os.Stat(filepath.Join(path, "slaves")); !os.IsNotExist(err) {
				devices, err := util.SysReadFile(filepath.Join(path, "slaves"))
				if err != nil {
					return fmt.Errorf("unable to read %s", filepath.Join(path, "slaves"))
				}
				if devices != "" {
					for _, device := range strings.Split(devices, " ") {
						if intf, exists := (*netClass)[device]; exists {
							netClassIface.BondAttrs.Devices = append(netClassIface.BondAttrs.Devices, &intf)
						} else {
							return fmt.Errorf("unable to find device %s", device)
						}
					}
				}
			}
		}
		if netClassIface.BondDeviceAttrs != nil {
			path := filepath.Join(fs.sys.Path(netclassPath), netClassIface.Name, "master")
			controller, err := filepath.EvalSymlinks(path)
			if err != nil {
				return fmt.Errorf("unable to read %s", path)
			}
			name := filepath.Base(controller)
			if intf, exists := (*netClass)[name]; exists {
				netClassIface.BondDeviceAttrs.Controller = &intf
			} else {
				return fmt.Errorf("unable to find device %s", controller)
			}
		}
	}
	return nil
}

// parseBondOption parses a "<string representation> <id>" string into its component parts.
func parseBondOption(s string) (string, uint64, error) {
	parts := strings.Split(s, " ")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("failed to split %s into 'NAME ID'", s)
	}
	value, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return "", 0, fmt.Errorf("failed to convert %s into uint64", parts[1])
	}
	return parts[0], value, nil
}

func parseDeviceQueueIDs(data string) (queueIDs map[string]uint64, err error) {
	queueIDs = make(map[string]uint64)
	for _, line := range strings.Split(data, " ") {
		sep := strings.LastIndex(line, ":")
		if queueIDs[line[:sep]], err = strconv.ParseUint(line[sep+1:], 10, 64); err != nil {
			return nil, err
		}
	}
	return
}

func parseArpTargets(data string) (arpTargets []net.IP, err error) {
	if data == "" {
		return
	}
	for _, ipString := range strings.Split(data, " ") {
		ip := net.ParseIP(ipString)
		if ip == nil {
			return nil, fmt.Errorf("could not parse ip %s", ip)
		}
		arpTargets = append(arpTargets, ip)
	}
	return
}
