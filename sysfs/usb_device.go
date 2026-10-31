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
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/prometheus/procfs/internal/parsers"
)

const usbDevicesPath = "bus/usb/devices"

// USBRemovable represents whether a USB device can be removed, as reported
// by /sys/bus/usb/devices/<dev>/removable.
type USBRemovable string

const (
	USBRemovableUnknown   USBRemovable = "unknown"
	USBRemovableFixed     USBRemovable = "fixed"
	USBRemovableRemovable USBRemovable = "removable"
)

// String returns the string representation of the removable state.
func (r USBRemovable) String() string {
	return string(r)
}

// USBInterface contains info from files in /sys/bus/usb/devices for a
// single USB interface (e.g. 1-2:1.0).
type USBInterface struct {
	Configuration string // configuration value from the entry name
	Interface     string // interface number from the entry name
	Class         string // /sys/bus/usb/devices/<iface>/bInterfaceClass
	SubClass      string // /sys/bus/usb/devices/<iface>/bInterfaceSubClass
	Protocol      string // /sys/bus/usb/devices/<iface>/bInterfaceProtocol
	Driver        string // /sys/bus/usb/devices/<iface>/driver (empty if unbound)
}

// USBController contains info from files in /sys/bus/usb/devices for a
// root hub / USB controller (e.g. usb1).
type USBController struct {
	BusNum   string // /sys/bus/usb/devices/usb<N>/busnum
	BCDUSB   string // /sys/bus/usb/devices/usb<N>/version
	Speed    uint64 // /sys/bus/usb/devices/usb<N>/speed (Mbps)
	MaxChild uint64 // /sys/bus/usb/devices/usb<N>/maxchild
}

// USBDevice contains info from files in /sys/bus/usb/devices for a
// single USB device. Root hubs (usb1, usb2, ...) and interface entries
// (e.g. 1-2:1.0) are not represented here.
type USBDevice struct {
	Name              string // entry name below /sys/bus/usb/devices (e.g. "1-2")
	BusNum            string // /sys/bus/usb/devices/<dev>/busnum
	DevNum            string // /sys/bus/usb/devices/<dev>/devnum
	IDVendor          string // /sys/bus/usb/devices/<dev>/idVendor
	IDProduct         string // /sys/bus/usb/devices/<dev>/idProduct
	Manufacturer      string // /sys/bus/usb/devices/<dev>/manufacturer
	Product           string // /sys/bus/usb/devices/<dev>/product
	Serial            string // /sys/bus/usb/devices/<dev>/serial
	DeviceClass       string // /sys/bus/usb/devices/<dev>/bDeviceClass
	DeviceSubClass    string // /sys/bus/usb/devices/<dev>/bDeviceSubClass
	DeviceProtocol    string // /sys/bus/usb/devices/<dev>/bDeviceProtocol
	BCDDevice         string // /sys/bus/usb/devices/<dev>/bcdDevice
	BCDUSB            string // /sys/bus/usb/devices/<dev>/version
	Speed             uint64 // /sys/bus/usb/devices/<dev>/speed (Mbps)
	NumConfigurations uint64 // /sys/bus/usb/devices/<dev>/bNumConfigurations
	NumInterfaces     uint64 // /sys/bus/usb/devices/<dev>/bNumInterfaces
	BMaxPower         uint64 // /sys/bus/usb/devices/<dev>/bMaxPower (mA)
	// Authorized is nil when /sys/bus/usb/devices/<dev>/authorized is unreadable.
	Authorized *bool
	// Removable is nil when /sys/bus/usb/devices/<dev>/removable is unreadable.
	Removable              *USBRemovable
	URBNum                 uint64         // /sys/bus/usb/devices/<dev>/urbnum
	PowerControl           string         // /sys/bus/usb/devices/<dev>/power/control
	RuntimeStatus          string         // /sys/bus/usb/devices/<dev>/power/runtime_status
	AutosuspendDelayMillis *int64         // /sys/bus/usb/devices/<dev>/power/autosuspend_delay_ms
	Interfaces             []USBInterface // /sys/bus/usb/devices/<dev>/<dev>:<cfg>.<iface>
}

// USBDevices returns info for all real USB devices read from
// /sys/bus/usb/devices. Root hubs and interface entries are skipped.
func (fs FS) USBDevices() ([]USBDevice, error) {
	path := fs.sys.Path(usbDevicesPath)

	dirs, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var usbDevs []USBDevice
	for _, d := range dirs {
		name := d.Name()
		// Root hubs are named "usb1", "usb2", ...; interface entries
		// contain a colon (e.g. "1-2:1.0"). Neither represents a real
		// device, so skip both.
		if strings.HasPrefix(name, "usb") || strings.Contains(name, ":") {
			continue
		}

		device, err := fs.parseUSBDevice(filepath.Join(path, name))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				// Not a real device (no idVendor attribute).
				continue
			}
			return nil, err
		}

		usbDevs = append(usbDevs, *device)
	}

	return usbDevs, nil
}

// USBControllers returns info for all root hubs / USB controllers read from
// /sys/bus/usb/devices. Interfaces of the root hubs (e.g. usb1/1-0:1.0) are
// not represented here.
func (fs FS) USBControllers() ([]USBController, error) {
	path := fs.sys.Path(usbDevicesPath)

	dirs, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var controllers []USBController
	for _, d := range dirs {
		name := d.Name()
		if !strings.HasPrefix(name, "usb") || strings.Contains(name, ":") {
			continue
		}

		controllers = append(controllers, fs.parseUSBController(filepath.Join(path, name)))
	}

	return controllers, nil
}

// parseUSBController reads the attributes of a root hub on a best-effort
// basis; attributes that are not available are left at their zero value.
func (fs FS) parseUSBController(controllerPath string) USBController {
	controller := USBController{}

	if value, err := parsers.SysReadFile(filepath.Join(controllerPath, "busnum")); err == nil {
		controller.BusNum = strings.TrimSpace(value)
	}
	if value, err := parsers.SysReadFile(filepath.Join(controllerPath, "version")); err == nil {
		controller.BCDUSB = strings.TrimSpace(value)
	}
	if value, err := parsers.SysReadUintFromFile(filepath.Join(controllerPath, "speed")); err == nil {
		controller.Speed = value
	}
	if value, err := parsers.SysReadUintFromFile(filepath.Join(controllerPath, "maxchild")); err == nil {
		controller.MaxChild = value
	}

	return controller
}

func (fs FS) parseUSBDevice(devicePath string) (*USBDevice, error) {
	// idVendor is the one attribute that identifies a real USB device;
	// all other attributes are optional and read on a best-effort basis.
	vendor, err := parsers.SysReadFile(filepath.Join(devicePath, "idVendor"))
	if err != nil {
		return nil, fmt.Errorf("failed to read idVendor for %q: %w", devicePath, err)
	}

	device := &USBDevice{
		Name:     filepath.Base(devicePath),
		IDVendor: strings.TrimSpace(vendor),
	}

	if value, err := parsers.SysReadFile(filepath.Join(devicePath, "busnum")); err == nil {
		device.BusNum = strings.TrimSpace(value)
	}

	if value, err := parsers.SysReadFile(filepath.Join(devicePath, "devnum")); err == nil {
		device.DevNum = strings.TrimSpace(value)
	}

	if value, err := parsers.SysReadFile(filepath.Join(devicePath, "idProduct")); err == nil {
		device.IDProduct = strings.TrimSpace(value)
	}

	if value, err := parsers.SysReadFile(filepath.Join(devicePath, "manufacturer")); err == nil {
		device.Manufacturer = strings.TrimSpace(value)
	}

	if value, err := parsers.SysReadFile(filepath.Join(devicePath, "product")); err == nil {
		device.Product = strings.TrimSpace(value)
	}

	if value, err := parsers.SysReadFile(filepath.Join(devicePath, "serial")); err == nil {
		device.Serial = strings.TrimSpace(value)
	}

	if value, err := parsers.SysReadFile(filepath.Join(devicePath, "bDeviceClass")); err == nil {
		device.DeviceClass = strings.TrimSpace(value)
	}

	if value, err := parsers.SysReadFile(filepath.Join(devicePath, "bDeviceSubClass")); err == nil {
		device.DeviceSubClass = strings.TrimSpace(value)
	}

	if value, err := parsers.SysReadFile(filepath.Join(devicePath, "bDeviceProtocol")); err == nil {
		device.DeviceProtocol = strings.TrimSpace(value)
	}

	if value, err := parsers.SysReadFile(filepath.Join(devicePath, "bcdDevice")); err == nil {
		device.BCDDevice = strings.TrimSpace(value)
	}

	if value, err := parsers.SysReadFile(filepath.Join(devicePath, "version")); err == nil {
		device.BCDUSB = strings.TrimSpace(value)
	}

	if value, err := parsers.SysReadUintFromFile(filepath.Join(devicePath, "speed")); err == nil {
		device.Speed = value
	}

	if value, err := parsers.SysReadUintFromFile(filepath.Join(devicePath, "bNumConfigurations")); err == nil {
		device.NumConfigurations = value
	}

	if value, err := parsers.SysReadUintFromFile(filepath.Join(devicePath, "bNumInterfaces")); err == nil {
		device.NumInterfaces = value
	}

	if value, err := parsers.SysReadFile(filepath.Join(devicePath, "bMaxPower")); err == nil {
		bmax := strings.TrimSuffix(strings.TrimSpace(value), "mA")
		if v, err := strconv.ParseUint(strings.TrimSpace(bmax), 10, 64); err == nil {
			device.BMaxPower = v
		}
	}

	if value, err := parsers.SysReadFile(filepath.Join(devicePath, "authorized")); err == nil {
		authorized := strings.TrimSpace(value) == "1"
		device.Authorized = &authorized
	}

	if value, err := parsers.SysReadFile(filepath.Join(devicePath, "removable")); err == nil {
		removable := USBRemovable(strings.TrimSpace(value))
		device.Removable = &removable
	}

	if value, err := parsers.SysReadUintFromFile(filepath.Join(devicePath, "urbnum")); err == nil {
		device.URBNum = value
	}

	if value, err := parsers.SysReadFile(filepath.Join(devicePath, "power", "control")); err == nil {
		device.PowerControl = strings.TrimSpace(value)
	}
	if value, err := parsers.SysReadFile(filepath.Join(devicePath, "power", "runtime_status")); err == nil {
		device.RuntimeStatus = strings.TrimSpace(value)
	}
	if value, err := parsers.SysReadFile(filepath.Join(devicePath, "power", "autosuspend_delay_ms")); err == nil {
		s := strings.TrimSpace(value)
		if v, err := strconv.ParseInt(s, 10, 64); err == nil {
			device.AutosuspendDelayMillis = &v
		}
	}

	// Interfaces are exposed as separate entries inside the device directory,
	// e.g. /sys/bus/usb/devices/1-2/1-2:1.0. They are read on a best-effort
	// basis: a device whose interface entries cannot be listed is still
	// reported, just without interfaces.
	ifaces, err := fs.parseUSBInterfaces(devicePath)
	if err == nil {
		device.Interfaces = ifaces
	}

	return device, nil
}

func (fs FS) parseUSBInterfaces(basePath string) ([]USBInterface, error) {
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	var ifaces []USBInterface
	for _, e := range entries {
		name := e.Name()
		if !strings.Contains(name, ":") {
			continue
		}
		parts := strings.SplitN(name, ":", 2)
		if len(parts) != 2 {
			continue
		}
		cfgIface := parts[1]
		cfgParts := strings.SplitN(cfgIface, ".", 2)
		if len(cfgParts) != 2 {
			continue
		}
		ifacePath := filepath.Join(basePath, name)
		iface := USBInterface{
			Configuration: cfgParts[0],
			Interface:     cfgParts[1],
		}
		if value, err := parsers.SysReadFile(filepath.Join(ifacePath, "bInterfaceClass")); err == nil {
			iface.Class = strings.TrimSpace(value)
		}
		if value, err := parsers.SysReadFile(filepath.Join(ifacePath, "bInterfaceSubClass")); err == nil {
			iface.SubClass = strings.TrimSpace(value)
		}
		if value, err := parsers.SysReadFile(filepath.Join(ifacePath, "bInterfaceProtocol")); err == nil {
			iface.Protocol = strings.TrimSpace(value)
		}
		// Driver is empty when no driver is bound to the interface: in that
		// case the driver symlink does not exist, which is indistinguishable
		// from an unreadable symlink and both leave Driver empty.
		if link, err := os.Readlink(filepath.Join(ifacePath, "driver")); err == nil {
			iface.Driver = filepath.Base(link)
		}
		ifaces = append(ifaces, iface)
	}

	return ifaces, nil
}
