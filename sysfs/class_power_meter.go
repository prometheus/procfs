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

	"github.com/prometheus/procfs/internal/parsers"
)

// PowerMeter represents a single ACPI 4.0 power meter device, as exposed by
// the kernel driver drivers/hwmon/acpi_power_meter.c.
//
// All optional fields use pointer types and are nil when the firmware does
// not expose the corresponding sysfs attribute. Fields documented as RW in
// the kernel (AverageMin, AverageMax, AverageInterval, Cap) can be read but
// not written through this API; callers needing to set values must write to
// sysfs directly.
//
// Typical sysfs path: /sys/bus/acpi/drivers/power_meter/ACPI000D:XX/.
//
// See: https://docs.kernel.org/hwmon/acpi_power_meter.html
type PowerMeter struct {
	// Standard hwmon attributes (always exposed by the driver).
	Average            *int64 // power1_average (microWatt)
	AverageMin         *int64 // power1_average_min (microWatt, RW)
	AverageMax         *int64 // power1_average_max (microWatt, RW)
	AverageInterval    *int64 // power1_average_interval (millisecond, RW)
	AverageIntervalMin *int64 // power1_average_interval_min (millisecond)
	AverageIntervalMax *int64 // power1_average_interval_max (millisecond)
	Alarm              *int64 // power1_alarm (0 or 1)

	// Optional capping attributes (present only when the platform supports
	// power capping; on non-IBM systems the kernel module must be loaded
	// with force_cap_on=1 on kernel >= 4.14).
	Cap     *int64 // power1_cap (microWatt, RW)
	CapMin  *int64 // power1_cap_min (microWatt)
	CapMax  *int64 // power1_cap_max (microWatt)
	CapHyst *int64 // power1_cap_hyst (microWatt)

	// ACPI extension attributes (firmware-dependent; may not be present on
	// all platforms).
	//
	// Note: the kernel emits power1_accuracy as a string with a percent
	// suffix (e.g. "1.50%"), so it is preserved as-is rather than parsed
	// to a number.
	Accuracy     string // power1_accuracy (e.g. "1.50%")
	IsBattery    *int64 // power1_is_battery (0 or 1)
	ModelNumber  string // power1_model_number
	SerialNumber string // power1_serial_number
	OEMInfo      string // power1_oem_info

	// Metadata.
	Name     string   // device directory name (e.g. "ACPI000D:00")
	Path     string   // full sysfs path
	Measures []string // device names from the measures/ subdirectory
}

// PowerMeterClass is the collection of all ACPI power meters enumerated from
// /sys/bus/acpi/drivers/power_meter/.
type PowerMeterClass []PowerMeter

// GetPowerMeters returns a slice of PowerMeter, one for each ACPI 4.0 power
// meter discovered via /sys/bus/acpi/drivers/power_meter/. Returns nil, nil
// when the ACPI power meter bus path does not exist (i.e. the host has no
// ACPI power meter device).
func GetPowerMeters(fs FS) (PowerMeterClass, error) {
	pattern := fs.sys.Path("bus/acpi/drivers/power_meter/ACPI000D:*")

	dirs, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to glob power meter devices: %w", err)
	}
	if len(dirs) == 0 {
		return PowerMeterClass{}, os.ErrNotExist
	}

	meters := make(PowerMeterClass, 0, len(dirs))
	for _, d := range dirs {
		pm, err := parsePowerMeter(d)
		if err != nil {
			return nil, fmt.Errorf("failed to parse power meter %q: %w", d, err)
		}
		pm.Name = filepath.Base(d)
		pm.Path = d
		meters = append(meters, *pm)
	}
	return meters, nil
}

// parsePowerMeter reads every attribute file inside a single power meter
// sysfs directory and returns a populated PowerMeter.
func parsePowerMeter(path string) (*PowerMeter, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var pm PowerMeter
	for _, f := range files {
		// Skip subdirectories (measures/ is handled separately) and
		// non-regular files.
		if !f.Type().IsRegular() {
			continue
		}

		name := filepath.Join(path, f.Name())
		value, err := parsers.SysReadFile(name)
		if err != nil {
			// Tolerate: device not ready / attribute unsupported /
			// permission denied. Matches the strategy used by
			// parsePowerSupply in class_power_supply.go.
			if os.IsNotExist(err) ||
				err.Error() == "operation not supported" ||
				err.Error() == "no such device" ||
				errors.Is(err, os.ErrInvalid) {
				continue
			}
			return nil, fmt.Errorf("failed to read file %q: %w", name, err)
		}

		vp := parsers.NewValueParser(value)

		switch f.Name() {
		case "power1_average":
			pm.Average = vp.PInt64()
		case "power1_average_min":
			pm.AverageMin = vp.PInt64()
		case "power1_average_max":
			pm.AverageMax = vp.PInt64()
		case "power1_average_interval":
			pm.AverageInterval = vp.PInt64()
		case "power1_average_interval_min":
			pm.AverageIntervalMin = vp.PInt64()
		case "power1_average_interval_max":
			pm.AverageIntervalMax = vp.PInt64()
		case "power1_alarm":
			pm.Alarm = vp.PInt64()
		case "power1_cap":
			pm.Cap = vp.PInt64()
		case "power1_cap_min":
			pm.CapMin = vp.PInt64()
		case "power1_cap_max":
			pm.CapMax = vp.PInt64()
		case "power1_cap_hyst":
			pm.CapHyst = vp.PInt64()
		case "power1_accuracy":
			pm.Accuracy = value
		case "power1_is_battery":
			pm.IsBattery = vp.PInt64()
		case "power1_model_number":
			pm.ModelNumber = value
		case "power1_serial_number":
			pm.SerialNumber = value
		case "power1_oem_info":
			pm.OEMInfo = value
		}

		if err := vp.Err(); err != nil {
			return nil, fmt.Errorf("failed to parse %q: %w", f.Name(), err)
		}
	}

	// measures/ failure is non-fatal; the directory may be empty or absent.
	pm.Measures, _ = parsePowerMeterMeasures(path)
	return &pm, nil
}

// parsePowerMeterMeasures reads the measures/ subdirectory of a power meter
// and returns the basenames of every symlink target (i.e. the device names
// this meter measures).
func parsePowerMeterMeasures(meterPath string) ([]string, error) {
	measureDir := filepath.Join(meterPath, "measures")
	entries, err := os.ReadDir(measureDir)
	if err != nil {
		return nil, err
	}

	var devices []string
	for _, e := range entries {
		target, err := os.Readlink(filepath.Join(measureDir, e.Name()))
		if err != nil {
			continue
		}
		devices = append(devices, filepath.Base(target))
	}
	return devices, nil
}
