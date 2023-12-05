package components

import (
	"errors"
	"fmt"
	"github.com/godbus/dbus"
	"strconv"
	"strings"
)

var (
	solaarDbus *dbus.Conn
)

// lib/logitech_receiver/hidpp20.py
// BATTERY_STATUS = _NamedInts(
//
//	discharging=0x00,
//	recharging=0x01,
//	almost_full=0x02,
//	full=0x03,
//	slow_recharge=0x04,
//	invalid_battery=0x05,
//	thermal_error=0x06
//
// )
const (
	solaarStatusDischarging    = "discharging"
	solaarStatusRecharging     = "recharging"
	solaarStatusAlmostFull     = "almost full"
	solaarStatusFull           = "full"
	solaarStatusSlowRecharge   = "slow recharge"
	solaarStatusInvalidBattery = "invalid battery"
	solaarStatusThermalError   = "thermal error"

	solaarStatusOffline = "offline"
	solaarStatusUnknown = "unknown"
)

var (
	solaarError   = drawColor("#ffff00")
	solaarWarning = drawColor("#ff0000")

	solaarStatusStatic = map[string]string{
		solaarStatusOffline:        IconSolaarOffline,
		solaarStatusUnknown:        IconSolaarUnknown,
		solaarStatusInvalidBattery: IconSolaarWarning,
		solaarStatusThermalError:   IconSolaarWarning,
	}
	solaarStatusCharging = map[string]string{
		solaarStatusDischarging:  "",
		solaarStatusAlmostFull:   "",
		solaarStatusFull:         "",
		solaarStatusRecharging:   IconSolaarCharging,
		solaarStatusSlowRecharge: IconSolaarCharging,
	}
	solaarIconDevice = map[HidType]string{
		HidMouse:    IconSolaarMouse,
		HidKeyboard: IconSolaarKeyboard,
	}
)

func init() {
	var err error
	solaarDbus, err = dbus.SessionBusPrivate()
	check(err)

	check(dbusPrivate(solaarDbus))
}

func solaarBattery(serial string) (map[string]dbus.Variant, error) {
	object := solaarDbus.Object("io.github.pwr_solaar.solaar.Status", "/io/github/pwr_solaar/solaar/Status")

	call := object.Call("Battery", dbus.FlagNoAutoStart, serial)

	if call.Err != nil {
		return nil, call.Err
	}

	if len(call.Body) != 1 {
		return nil, errors.New("unknown body length")
	}

	return call.Body[0].(map[string]dbus.Variant), nil
}

func solaarStatusBattery(rawLevel string) string {
	cutLevel, ok := strings.CutSuffix(rawLevel, "%")
	if !ok {
		return IconSolaarUnknown
	}

	level, err := strconv.ParseInt(cutLevel, 10, 0)
	if err != nil {
		fmt.Println("solaar: parseInt:", err)
		return IconSolaarError
	}

	return batteryDraw(int(level))
}

func solaarStatus(device ConfigSolaar) string {
	typeIcon := solaarIconDevice[device.hidType] + " "
	errorIcon := solaarError + IconSolaarError + DrawReset
	if NoDraw {
		errorIcon = IconSolaarError
	}

	res, err := solaarBattery(device.serial)

	if err != nil {
		fmt.Println("solaar: solaarBattery:", err)
		return typeIcon + errorIcon
	}

	status, ok := res["status"].Value().(string)
	if !ok {
		fmt.Println("solaar: status not string")
		return typeIcon + errorIcon
	}

	level, ok := res["level"].Value().(string)
	if !ok {
		fmt.Println("solaar: level not string")
		return typeIcon + errorIcon
	}

	if static, ok := solaarStatusStatic[status]; ok {
		if NoDraw {
			return typeIcon + static
		} else {
			return typeIcon + solaarWarning + static + DrawReset
		}
	}

	chargingIcon, ok := solaarStatusCharging[status]
	if !ok {
		fmt.Println("solaar: unknown status:", status)
		return typeIcon + errorIcon
	}

	if chargingIcon != "" {
		chargingIcon = " " + chargingIcon
	}

	return typeIcon + solaarStatusBattery(level) + chargingIcon
}

func Solaar(_ int64) string {
	output := make([]string, len(SolaarDevices))

	for i, device := range SolaarDevices {
		output[i] = solaarStatus(device)
	}

	return strings.Join(output, " ")
}
