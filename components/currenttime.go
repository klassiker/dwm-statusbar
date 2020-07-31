package components

import (
	"strings"
	"time"
)

var (
	CalendarIcon = "\uf073"
	ClockIcon = "\uf017"
)

func CurrentTime(_ uint64) string {
	timeNow := time.Now()
	output := []string{
		CalendarIcon, timeNow.Format("2006-01-02"),
		ClockIcon, timeNow.Format("15:04:05"),
	}
	return strings.Join(output, " ")
}
