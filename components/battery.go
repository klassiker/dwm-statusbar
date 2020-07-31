package components

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type BatteryStruct struct {
	index int
}

func (b *BatteryStruct) Capacity() int {
	path := fmt.Sprintf("%s/BAT%d/capacity", BatteryPath, b.index)
	var capacity int

	if fileExists(path) {
		data, err := ioutil.ReadFile(path)
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

	switch {
	case capacity > 95:
		status = "^d^"
		icon = IconBatteryFull
	case capacity > 75:
		status = "^c#00ff00^"
		icon = IconBatterHigh
	case capacity > 50:
		status = "^c#ffff00^"
		icon = IconBatteryHalf
	case capacity > 25:
		status = "^c#ff8800^"
		icon = IconBatteryLow
	default:
		status = "^c#ff0000^"
		icon = IconBatteryEmpty
	}

	return fmt.Sprintf("%s%s %d%%", status, icon, capacity)
}

func Battery(_ uint64) string {
	output := make([]string, Batteries)

	for i := 0; i < Batteries; i++ {
		battery := &BatteryStruct{i}
		output[i] = battery.Drawing()
	}

	data, err := ioutil.ReadFile(fmt.Sprintf("%s/AC/online", BatteryPath))
	check(err)

	online, err := strconv.Atoi(strings.TrimSpace(string(data)))
	check(err)

	if online == 1 {
		output = append([]string{IconBatteryPlug}, output...)
	}

	return strings.Join(output, " ") + "^d^"
}
