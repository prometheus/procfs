package sysfs

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/prometheus/procfs/internal/util"
)

type PowerSupply struct {
	Name                     string // Power Supply Name
	Type                     string `fileName:"type"`                        // /sys/class/power_supply/<Name>/type
	Status                   string `fileName:"status"`                      // /sys/class/power_supply/<Name>/status
	VoltageNow               *int64 `fileName:"voltage_now"`                 // /sys/class/power_supply/<Name>/voltage_now
	EnergyNow                *int64 `fileName:"energy_now"`                  // /sys/class/power_supply/<Name>/energy_now
	ChargeType               string `fileName:"charge_type"`                 // /sys/class/power_supply/<Name>/charge_type
	Authentic                *int64 `fileName:"authentic"`                   // /sys/class/power_suppy/<Name>/authentic
	Health                   string `fileName:"health"`                      // /sys/class/power_suppy/<Name>/health
	VoltageOCV               *int64 `fileName:"voltage_ocv"`                 // /sys/class/power_suppy/<Name>/voltage_ocv
	VoltageMinDesign         *int64 `fileName:"voltage_min_design"`          // /sys/class/power_suppy/<Name>/voltage_min_design
	VoltageMaxDesign         *int64 `fileName:"voltage_max_design"`          // /sys/class/power_suppy/<Name>/voltage_max_design
	VoltageMin               *int64 `fileName:"voltage_min"`                 // /sys/class/power_suppy/<Name>/voltage_min
	VoltageMax               *int64 `fileName:"voltage_max"`                 // /sys/class/power_suppy/<Name>/voltage_max
	VoltageBoot              *int64 `fileName:"voltage_boot"`                // /sys/class/power_suppy/<Name>/voltage_boot
	CurrentBoot              *int64 `fileName:"current_boot"`                // /sys/class/power_suppy/<Name>/current_boot
	ChargeEmptyDesign        *int64 `fileName:"charge_empty_design"`         // /sys/class/power_suppy/<Name>/charge_empty_design
	ChargeFullDesign         *int64 `fileName:"charge_full_design"`          // /sys/class/power_suppy/<Name>/charge_full_design
	EnergyEmptyDesign        *int64 `fileName:"energy_empty_design"`         // /sys/class/power_suppy/<Name>/energy_empty_design
	EnergyFullDesign         *int64 `fileName:"energy_full_design"`          // /sys/class/power_suppy/<Name>/energy_full_design
	ChargeCounter            *int64 `fileName:"charge_counter"`              // /sys/class/power_suppy/<Name>/charge_counter
	PrechargeCurrent         *int64 `fileName:"precharge_current"`           // /sys/class/power_suppy/<Name>/precharge_current
	ChargeTermCurrent        *int64 `fileName:"charge_term_current"`         // /sys/class/power_suppy/<Name>/charge_term_current
	ConstantChargeCurrent    *int64 `fileName:"constant_charge_current"`     // /sys/class/power_suppy/<Name>/constant_charge_current
	ConstantChargeCurrentMax *int64 `fileName:"constant_charge_current_max"` // /sys/class/power_suppy/<Name>/constant_charge_current_max
	ConstantChargeVoltage    *int64 `fileName:"constant_charge_voltage"`     // /sys/class/power_suppy/<Name>/constant_charge_voltage
	ConstantChargeVoltageMax *int64 `fileName:"constant_charge_voltage_max"` // /sys/class/power_suppy/<Name>/constant_charge_voltage_max
	InputCurrentLimit        *int64 `fileName:"input_current_limit"`         // /sys/class/power_suppy/<Name>/input_current_limit
	ChargeControlLimit       *int64 `fileName:"charge_control_limit"`        // /sys/class/power_suppy/<Name>/charge_control_limit
	ChargeControlLimitMax    *int64 `fileName:"charge_control_limit_max"`    // /sys/class/power_suppy/<Name>/charge_control_limit_max
	Capacity                 *int64 `fileName:"capacity"`                    // /sys/class/power_suppy/<Name>/capacity
	CapacityAlertMin         *int64 `fileName:"capacity_alert_min"`          // /sys/class/power_suppy/<Name>/capacity_alert_min
	CapacityAlertMax         *int64 `fileName:"capacity_alert_max"`          // /sys/class/power_suppy/<Name>/capacity_alert_max
	CapacityLevel            string `fileName:"capacity_level"`              // /sys/class/power_suppy/<Name>/capacity_level
	Temp                     *int64 `fileName:"temp"`                        // /sys/class/power_suppy/<Name>/temp
	TempAlertMin             *int64 `fileName:"temp_alert_min"`              // /sys/class/power_suppy/<Name>/temp_alert_min
	TempAlertMax             *int64 `fileName:"temp_alert_max"`              // /sys/class/power_suppy/<Name>/temp_alert_max
	TempAmbient              *int64 `fileName:"temp_ambient"`                // /sys/class/power_suppy/<Name>/temp_ambient
	TempAmbientMin           *int64 `fileName:"temp_ambient_min"`            // /sys/class/power_suppy/<Name>/temp_ambient_min
	TempAmbientMax           *int64 `fileName:"temp_ambient_max"`            // /sys/class/power_suppy/<Name>/temp_ambient_max
	TempMin                  *int64 `fileName:"temp_min"`                    // /sys/class/power_suppy/<Name>/temp_min
	TempMax                  *int64 `fileName:"temp_max"`                    // /sys/class/power_suppy/<Name>/temp_max
	TimeToEmpty              *int64 `fileName:"time_to_empty"`               // /sys/class/power_suppy/<Name>/time_to_empty
	TimeToFull               *int64 `fileName:"time_to_full"`                // /sys/class/power_suppy/<Name>/time_to_full
}

type PowerSupplyClass map[string]PowerSupply

func NewPowerSupplyClass() (PowerSupplyClass, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return nil, err
	}

	return fs.NewPowerSupplyClass()
}

func (fs FS) NewPowerSupplyClass() (PowerSupplyClass, error) {
	path := fs.Path("class/power_supply")

	powerSupplyDirs, err := ioutil.ReadDir(path)
	if err != nil {
		return PowerSupplyClass{}, fmt.Errorf("cannot access %s dir %s", path, err)
	}

	powerSupplyClass := PowerSupplyClass{}
	for _, powerSupplyDir := range powerSupplyDirs {
		powerSupply, err := powerSupplyClass.parsePowerSupply(path + "/" + powerSupplyDir.Name())
		if err != nil {
			return nil, err
		}
		powerSupply.Name = powerSupplyDir.Name()
		powerSupplyClass[powerSupplyDir.Name()] = *powerSupply
	}
	return powerSupplyClass, nil
}

func (psc PowerSupplyClass) parsePowerSupply(powerSupplyPath string) (*PowerSupply, error) {
	powerSupply := PowerSupply{}
	powerSupplyElem := reflect.ValueOf(&powerSupply).Elem()
	powerSupplyType := reflect.TypeOf(powerSupply)

	//start from 1 - skip the Name field
	for i := 1; i < powerSupplyElem.NumField(); i++ {
		fieldType := powerSupplyType.Field(i)
		fieldValue := powerSupplyElem.Field(i)

		if fieldType.Tag.Get("fileName") == "" {
			panic(fmt.Errorf("field %s does not have a filename tag", fieldType.Name))
		}

		value, err := util.SysReadFile(powerSupplyPath + "/" + fieldType.Tag.Get("fileName"))

		if err != nil {
			if os.IsNotExist(err) || err.Error() == "operation not supported" || err.Error() == "invalid argument" {
				continue
			}
			return nil, fmt.Errorf("could not access file %s: %s", fieldType.Tag.Get("fileName"), err)
		}

		switch fieldValue.Kind() {
		case reflect.String:
			fieldValue.SetString(value)
		case reflect.Ptr:
			var int64ptr *int64
			switch fieldValue.Type() {
			case reflect.TypeOf(int64ptr):
				var intValue int64
				if strings.HasPrefix(value, "0x") {
					intValue, err = strconv.ParseInt(value[2:], 16, 64)
					if err != nil {
						return nil, fmt.Errorf("expected hex value for %s, got: %s", fieldType.Name, value)
					}
				} else {
					intValue, err = strconv.ParseInt(value, 10, 64)
					if err != nil {
						return nil, fmt.Errorf("expected Uint64 value for %s, got: %s", fieldType.Name, value)
					}
				}
				fieldValue.Set(reflect.ValueOf(&intValue))
			default:
				return nil, fmt.Errorf("unhandled pointer type %q", fieldValue.Type())
			}
		default:
			return nil, fmt.Errorf("unhandled type %q", fieldValue.Kind())
		}
	}

	return &powerSupply, nil
}
