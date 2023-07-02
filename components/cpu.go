package components

// Taken from: https://rosettacode.org/wiki/Linux_CPU_utilization#Go

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var (
	CPUPath              = "/proc/stat"
	CPUCores             = 4
	CPUData              = make([]CPUDataStore, CPUCores+1)
	CPUBarWidth          = 5
	CPUBarForeground     = "#ff0000"
	CPUBarBackground     = "#000000"
	CPUBarDrawBackground = fmt.Sprintf("^c%s^^r0,%d,%d,%d^", CPUBarBackground, BarPadding, CPUBarWidth, BarHeight)
	CPUBarDrawForeground = fmt.Sprintf("^c%s^^r0,%s,%d,%s^^f%d^", CPUBarForeground, "%d", CPUBarWidth, "%d", CPUBarWidth+2*BarPadding)
	CPUBARDrawLast       = fmt.Sprintf("^f%d^^d^", -2*BarPadding)
	CPUBarDraw           = CPUBarDrawBackground + CPUBarDrawForeground
)

type CPUDataStore struct {
	prevIdleTime  uint64
	prevTotalTime uint64
	filled        bool
}

func cpuReadData() []float64 {
	data, err := os.ReadFile(CPUPath)
	check(err)

	info := filter(strings.Split(string(data), "\n"), func(s string) bool {
		return strings.HasPrefix(s, "cpu")
	})

	cpuPercent := make([]float64, CPUCores+1)
	for i, line := range info {
		split := strings.Split(line[5:], " ")

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

	var draw string
	for _, v := range cpuPercent {
		height := int(math.Round(float64(BarHeight) * v))
		offset := BarHeight - height
		draw += fmt.Sprintf(CPUBarDraw, offset, height)
	}

	output := []string{
		cpuText,
		draw,
		CPUBARDrawLast,
	}

	return strings.Join(output, "")
}
