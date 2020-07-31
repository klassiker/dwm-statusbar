package components

import (
	"strings"
	"time"
)

func CurrentTime(_ uint64) string {
	timeNow := time.Now()
	output := []string{
		IconCurrentTimeCalendar, timeNow.Format("2006-01-02"),
		IconCurrentTimeClock, timeNow.Format("15:04:05"),
	}
	return strings.Join(output, " ")
}
