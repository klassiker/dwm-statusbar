package components

import (
	"fmt"
	"math"
	"syscall"
	"time"
)

func Uptime(_ uint64) string {
	var info syscall.Sysinfo_t
	check(syscall.Sysinfo(&info))

	duration := time.Duration(info.Uptime) * time.Second

	days := math.Floor(duration.Hours() / 24)
	hours := math.Floor(math.Mod(duration.Hours(), 24))
	minutes := math.Floor(math.Mod(duration.Minutes(), 60))
	seconds := math.Floor(math.Mod(duration.Seconds(), 60))

	var durationString string

	switch {
	case days > 0:
		durationString = fmt.Sprintf("%0.0fd %0.0fh", days, hours)
	case hours > 0:
		durationString = fmt.Sprintf("%0.0fh %0.0fm", hours, minutes)
	case minutes > 0:
		durationString = fmt.Sprintf("%0.0fm %0.0fs", minutes, seconds)
	default:
		durationString = fmt.Sprintf("%0.0fs", seconds)
	}

	return IconUptime + " " + durationString
}
