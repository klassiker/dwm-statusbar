package components

import (
	"log"
	"os"
)

var (
	BarHeight         = 21 - 2 * BarPadding
	BarPadding        = 1
)

type Basic = func(interval uint64) string

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func indexOf(arr []string, s string) int {
	for i, v := range arr {
		if v == s {
			return i
		}
	}

	return -1
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

	for unit = 0; *value > 1024.0 && unit < len(units); unit++ {
		*value /= 1024.0
	}

	return units[unit]
}