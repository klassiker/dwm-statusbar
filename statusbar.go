package main

import (
	"fmt"
	"github.com/klassiker/dwm-statusbar/components"
	"strings"
	"time"
)

const (
	second = 1000
	minute = 60 * second
)

var (
	baseInterval         = 900
	dwmSeparator         = ";"
	statusSeparatorStart = ""
	statusSeparatorMid   = " | "
	statusSeparatorEnd   = ""
	inputs               = [][]*Component{
		{
			{function: components.Network, interval: 1 * second},
			{register: components.Pulseaudio},
			{function: components.Memory, interval: 2 * second},
			{function: components.CPUPercentBar, interval: 1 * second},
			{function: components.CurrentTime, interval: 1 * second},
		},
		{
			{function: components.Uptime, interval: 60 * second},
			{function: components.Battery, interval: 10 * second},
			{function: components.Solaar, interval: 30 * minute},
			{function: components.Filesystem, interval: 60 * second},
			{function: components.Thermal, interval: 1 * second},
			{register: components.Sound},
		},
	}
)

type StatusUpdate struct {
	*Component
	status string
}

type Component struct {
	function   components.Basic
	register   components.Async
	interval   int64
	aggregator chan *StatusUpdate
	index, row int
	instant    bool
}

func (cp *Component) UpdateStatus() {
	cp.Update(cp.function(cp.interval))
}

func (cp *Component) Update(status string) {
	cp.aggregator <- &StatusUpdate{
		Component: cp,
		status:    strings.TrimSpace(status),
	}
}

func (cp *Component) Init(row, index int, agg chan *StatusUpdate) {
	cp.aggregator = agg
	cp.index = index
	cp.row = row
	cp.instant = cp.IsAsync()
}

func (cp *Component) IsAsync() bool {
	return cp.function == nil || cp.interval == 0
}

func (cp *Component) Run() {
	defer recovery()

	if cp.interval < 0 {
		cp.UpdateStatus()
	} else if cp.IsAsync() {
		cp.register(cp.Update)
	} else {
		cp.UpdateStatus()

		ticker := time.NewTicker(time.Duration(cp.interval) * time.Millisecond)

		for range ticker.C {
			cp.UpdateStatus()
		}
	}
}

func initComponents() chan *StatusUpdate {
	agg := make(chan *StatusUpdate)

	for row, cps := range inputs {
		for index, cp := range cps {
			cp.Init(row, index, agg)
			go cp.Run()
		}
	}

	return agg
}

func main() {
	status := make([][]string, len(inputs))
	output := make([]string, len(inputs)*4-1)

	for i := 0; i < len(inputs); i++ {
		status[i] = make([]string, len(inputs[i]))

		offset := i * 4

		output[offset] = statusSeparatorStart
		output[offset+2] = statusSeparatorEnd

		if i == 0 {
			continue
		}

		output[offset-1] = dwmSeparator
	}

	updateStatus := func() {
		text := strings.Join(output, "")
		if *debugFlag {
			fmt.Println("length:", len(text))
		}
		xsetroot(strings.Join(output, ""))
	}

	agg := initComponents()
	ticker := time.NewTicker(time.Duration(baseInterval) * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			updateStatus()
		case update := <-agg:
			status[update.row][update.index] = update.status
			output[update.row*4+1] = strings.Join(status[update.row], statusSeparatorMid)

			if update.instant {
				updateStatus()
			}
		case <-shutdown:
			cleanup(nil)
		}
	}

}
