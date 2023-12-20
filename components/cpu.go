package components

// Taken from: https://rosettacode.org/wiki/Linux_CPU_utilization#Go

import (
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

var (
	CPUPath              = "/proc/stat"
	CPUFile              *os.File
	CPUData              = make([]CPUDataStore, CPUCores+1)
	CPUBarWidth          = 5
	CPUBarBackground     = "^c#000000^"
	CPUBarForeground     = "^c#ff0000^"
	CPUBarDrawBackground = fmt.Sprintf("^r0,1,%d,%d^^f%d^", CPUBarWidth, BarHeight, CPUBarWidth+2*BarPadding)
	CPUBarDrawForeground = fmt.Sprintf("^r0,%s,%d,%s^^f%d^", "%d", CPUBarWidth, "%d", CPUBarWidth+2*BarPadding)
	CPUBARDrawLast       = fmt.Sprintf("^f%d^^d^", -2*BarPadding)
)

type CPUDataStore struct {
	prevIdleTime  uint64
	prevTotalTime uint64
	filled        bool
}

func init() {
	var err error
	CPUFile, err = os.Open(CPUPath)
	check(err)
}

func cpuReadData() []float64 {
	_, err := CPUFile.Seek(0, 0)
	check(err)

	data, err := io.ReadAll(CPUFile)
	check(err)

	// TODO dynamic cpu core count
	info := filter(strings.Split(string(data), "\n"), func(s string) bool {
		return strings.HasPrefix(s, "cpu")
	})[:CPUCores]

	cpuPercent := make([]float64, CPUCores+1)
	for i, line := range info {
		split := strings.Fields(line)[1:]

		idleTime, err := strconv.ParseUint(split[3], 10, 64)
		check(err)

		var totalTime uint64
		for _, s := range split {
			u, err := strconv.ParseUint(s, 10, 64)
			check(err)

			totalTime += u
		}

		cpuData := &CPUData[i]

		if cpuData.filled {
			deltaIdleTime := idleTime - cpuData.prevIdleTime
			deltaTotalTime := totalTime - cpuData.prevTotalTime
			cpuPercent[i] = 1.0 - float64(deltaIdleTime)/float64(deltaTotalTime)
		} else {
			cpuData.filled = true
		}

		cpuData.prevIdleTime = idleTime
		cpuData.prevTotalTime = totalTime
	}

	return cpuPercent
}

func CPUPercentBar(_ int64) string {
	cpuPercent := cpuReadData()
	cpuText := fmt.Sprintf("%s %0.0f%% ", IconCPU, math.Round(cpuPercent[0]*100.0))

	if NoDraw {
		return cpuText
	}

	cpuPercent = cpuPercent[1:]

	// TODO cleanup this mess
	var draw string
	draw += CPUBarBackground
	for range cpuPercent {
		draw += CPUBarDrawBackground
	}
	draw += fmt.Sprintf("^f-%d^", CPUCores * (CPUBarWidth+2*BarPadding))
	draw += CPUBarForeground
	for _, v := range cpuPercent {
		height := int(math.Round(float64(BarHeight) * v))
		offset := BarHeight - height + BarPadding
		draw += fmt.Sprintf(CPUBarDrawForeground, offset, height)
	}

	output := []string{
		cpuText,
		draw,
		CPUBARDrawLast,
	}

	return strings.Join(output, "")
}
