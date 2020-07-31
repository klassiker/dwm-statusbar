package components

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var (
	Batteries        = 2
	BatteryPath      = "/sys/class/power_supply"
	BatteryIconPlug  = "\uf1e6"
	BatteryIconFull  = "\uf240"
	BatteryIconHigh  = "\uf241"
	BatteryIconHalf  = "\uf242"
	BatteryIconLow   = "\uf243"
	BatteryIconEmpty = "\uf244"
	BatteryPerfect   = "^d^"
	BatteryGood      = "^c#00ff00^"
	BatteryOkay      = "^c#ffff00^"
	BatteryBad       = "^c#ff0000^"
)

func batteryStatus(path string) int {
	var capacity int

	if fileExists(path) {
		data, err := ioutil.ReadFile(path)
		check(err)

		capacity, err = strconv.Atoi(strings.TrimSpace(string(data)))
		check(err)
	}

	return capacity
}

func batteryStatusDrawing(capacity int) string {
	var icon string
	var status string

	switch {
	case capacity > 95:
		status = BatteryPerfect
		icon = BatteryIconFull
	case capacity > 75:
		status = BatteryGood
		icon = BatteryIconHigh
	case capacity > 50:
		status = BatteryOkay
		icon = BatteryIconHalf
	case capacity > 25:
		status = BatteryOkay
		icon = BatteryIconLow
	default:
		status = BatteryBad
		icon = BatteryIconEmpty
	}

	return fmt.Sprintf("%s%s %d%%", status, icon, capacity)
}

func batteryLoading() bool {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/AC/online", BatteryPath))
	check(err)

	online, err := strconv.Atoi(strings.TrimSpace(string(data)))
	check(err)

	return online == 1
}

func Battery(_ uint64) string {
	output := make([]string, Batteries)

	for i := 0; i < Batteries; i++ {
		capacity := batteryStatus(fmt.Sprintf("%s/BAT%d/capacity", BatteryPath, i))
		output[i] = batteryStatusDrawing(capacity)
	}

	acConnected := batteryLoading()

	if acConnected {
		output = append([]string{BatteryIconPlug}, output...)
	}

	return strings.Join(output, " ") + "^d^"
}
