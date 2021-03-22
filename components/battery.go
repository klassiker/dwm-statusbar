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

type BatteryStruct struct {
	index int
}

func (b *BatteryStruct) Capacity() int {
	path := fmt.Sprintf("%s/BAT%d/capacity", BatteryPath, b.index)
	var capacity int

	if fileExists(path) {
		data, err := os.ReadFile(path)
		check(err)

		capacity, err = strconv.Atoi(strings.TrimSpace(string(data)))
		check(err)
	}

	return capacity
}

func (b *BatteryStruct) Drawing() string {
	var icon string
	var status string

	capacity := b.Capacity()
	reset := DrawReset

	switch {
	case capacity > 95:
		status = BatteryPerfect
		icon = IconBatteryFull
	case capacity > 75:
		status = BatteryGood
		icon = IconBatterHigh
	case capacity > 50:
		status = BatteryOkay
		icon = IconBatteryHalf
	case capacity > 25:
		status = BatteryBad
		icon = IconBatteryLow
	default:
		status = BatteryDead
		icon = IconBatteryEmpty
	}

	if NoDraw {
		status = ""
	}

	if status == "" {
		reset = ""
	}

	return fmt.Sprintf("%s%s %d%%%s", status, icon, capacity, reset)
}

func acOnline() bool {
	data, err := os.ReadFile(fmt.Sprintf("%s/AC/online", BatteryPath))
	check(err)

	online, err := strconv.Atoi(strings.TrimSpace(string(data)))
	check(err)

	return online == 1
}

func Battery(_ uint64) string {
	output := make([]string, Batteries)

	for i := 0; i < Batteries; i++ {
		battery := &BatteryStruct{i}
		output[i] = battery.Drawing()
	}

	if acOnline() {
		output = append([]string{IconBatteryPlug}, output...)
	}

	return strings.Join(output, " ")
}
