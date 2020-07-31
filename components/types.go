package components

import (
	"log"
	"os"
	"runtime"
	"time"
)

type Basic = func(interval uint64) string

type Async = func(channel chan string)

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

func mapValueOrDefault(valueMap map[string]string, key, defaultValue string) string {
	if value, ok := valueMap[key]; ok {
		return value
	} else {
		pc, _, _, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			log.Printf("unknown key: %s - %s", details.Name(), key)
		}
		return defaultValue
	}
}

func profilingLog(start time.Time) {
	//_, file, _, _ := runtime.Caller(1)
	//log.Printf("init-%s: %s", filepath.Base(file), time.Since(start))
}
