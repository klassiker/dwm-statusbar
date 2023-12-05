package components

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	BatteryPath    = "/sys/class/power_supply"
	BatteryPerfect = drawColor("")
	BatteryGood    = drawColor("#00ff00")
	BatteryOkay    = drawColor("#ffff00")
	BatteryBad     = drawColor("#ff8800")
	BatteryDead    = drawColor("#ff0000")
)

func batteryCapacity(index int) int {
	path := fmt.Sprintf("%s/BAT%d/capacity", BatteryPath, index)
	var capacity int

	if fileExists(path) {
		data, err := os.ReadFile(path)
		check(err)

		capacity, err = strconv.Atoi(strings.TrimSpace(string(data)))
		check(err)
	}

	return capacity
}

func batteryAcOnline() bool {
	data, err := os.ReadFile(fmt.Sprintf("%s/AC/online", BatteryPath))
	check(err)

	online, err := strconv.Atoi(strings.TrimSpace(string(data)))
	check(err)

	return online == 1
}

func Battery(_ int64) string {
	output := make([]string, Batteries)

	for i := 0; i < Batteries; i++ {
		output[i] = batteryDraw(batteryCapacity(i))
	}

	if batteryAcOnline() {
		output = append([]string{IconBatteryPlug}, output...)
	}

	return strings.Join(output, " ")
}
