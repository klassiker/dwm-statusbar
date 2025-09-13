package components

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/godbus/dbus"
)

const (
	DrawReset = "^d^"
)

func drawColor(color string) string {
	if color == "" {
		return ""
	} else {
		return fmt.Sprintf("^c%s^", color)
	}
}

type Basic = func(interval int64) string

type Async = func(channel func(string))

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)

	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}

	return vsf
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func calculateUnit(value *float64, units []string) string {
	var unit int

	for unit = 0; *value > 1024.0 && unit < len(units)-1; unit++ {
		*value /= 1024.0
	}

	return units[unit]
}

func mapValueOrDefault(valueMap map[string]string, key, defaultValue string) string {
	if value, ok := valueMap[key]; ok {
		return value
	} else if key != "" {
		pc, _, _, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			log.Printf("unknown key: %s - %s", details.Name(), key)
		}
	}

	return defaultValue
}

func batteryDraw(level int) string {
	var status string
	var icon string
	switch {
	case level > 95:
		status = BatteryPerfect
		icon = IconBatteryFull
	case level > 75:
		status = BatteryGood
		icon = IconBatterHigh
	case level > 50:
		status = BatteryOkay
		icon = IconBatteryHalf
	case level > 25:
		status = BatteryBad
		icon = IconBatteryLow
	default:
		status = BatteryDead
		icon = IconBatteryEmpty
	}

	if NoDraw {
		status = ""
	}

	reset := DrawReset
	if status == "" {
		reset = ""
	}

	return fmt.Sprintf("%s%s %d%%%s", status, icon, level, reset)
}

func dbusPrivate(conn *dbus.Conn) error {
	if err := conn.Auth(nil); err != nil {
		// we already have an error, so we just try to close it here
		_ = conn.Close()
		return err
	}

	if err := conn.Hello(); err != nil {
		// we already have an error, so we just try to close it here
		_ = conn.Close()
		return err
	}

	return nil
}

func dbusIsMemberChar(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || c == '_'
}

func dbusEscape(path dbus.ObjectPath) dbus.ObjectPath {
	if path.IsValid() {
		return path
	}

	var output []rune

	for _, char := range path {
		if dbusIsMemberChar(char) {
			output = append(output, char)
		} else {
			output = append(output, []rune(fmt.Sprintf("_%x", char))...)
		}
	}

	return dbus.ObjectPath(output)
}
