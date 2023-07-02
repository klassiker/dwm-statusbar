package components

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var (
	MemoryPath              = "/proc/meminfo"
	MemoryBarWidth          = 40
	MemoryBarForeground     = "#0000ff"
	MemoryBarBackground     = "#000000"
	MemoryBarDrawBackground = fmt.Sprintf("^c%s^^r0,%d,%d,%d^", MemoryBarBackground, BarPadding, MemoryBarWidth, BarHeight)
	MemoryBarDrawForeground = fmt.Sprintf("^c%s^^r0,%d,%s,%d^^f%d^^d^", MemoryBarForeground, BarPadding, "%d", BarHeight, MemoryBarWidth)
	MemoryBarDraw           = MemoryBarDrawBackground + MemoryBarDrawForeground
	MemoryUnits             = []string{"KB", "MB", "GB"}
	MemoryData              = map[string]int{
		"MemTotal":     0,
		"MemFree":      0,
		"Shmem":        0,
		"Buffers":      0,
		"Cached":       0,
		"SReclaimable": 0,
	}
)

func memoryCalculateBar(percent float64) string {
	width := int(math.Round(percent * float64(MemoryBarWidth)))
	return fmt.Sprintf(MemoryBarDraw, width)
}

func memoryReadData() {
	data, err := os.ReadFile(MemoryPath)
	check(err)

	for _, line := range strings.Split(string(data), "\n") {
		parts := strings.SplitAfter(line, ":")

		if len(parts) == 2 {
			name, valueRaw := parts[0], parts[1]
			name = strings.Split(name, ":")[0]

			if _, ok := MemoryData[name]; ok {
				var value int
				_, err = fmt.Sscanf(strings.TrimSpace(valueRaw), "%d kB", &value)
				check(err)

				MemoryData[name] = value
			}
		}
	}
}

func Memory(_ int64) string {
	memoryReadData()

	memUsed := float64(MemoryData["MemTotal"] - MemoryData["MemFree"] - MemoryData["Buffers"] - MemoryData["Cached"] - MemoryData["SReclaimable"] + MemoryData["Shmem"])

	memBar := memoryCalculateBar(memUsed / float64(MemoryData["MemTotal"]))
	memUnit := calculateUnit(&memUsed, MemoryUnits)
	memUsedString := strconv.FormatFloat(math.Round(memUsed*100)/100, 'f', 2, 64)

	output := []string{
		IconMemory, memUsedString + memUnit,
	}

	if !NoDraw {
		output = append(output, memBar)
	}

	return strings.Join(output, " ")
}
