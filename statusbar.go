package main

import (
	"errors"
	"github.com/klassiker/dwm-statusbar/components"
	"strings"
	"time"
)

var (
	baseInterval         = 900
	dwmSeparator         = ";"
	statusSeparatorStart = "["
	statusSeparatorMid   = "] ["
	statusSeparatorEnd   = "]"
	inputs               = [][]*Component{
		{
			{function: components.Network, interval: 2 * 1000},
			{register: components.Pulseaudio},
			{function: components.Memory, interval: 2 * 1000},
			{function: components.CPUPercentBar, interval: 2 * 1000},
			{function: components.Battery, interval: 10 * 1000},
			{function: components.CurrentTime, interval: 1 * 1000},
		},
		{
			{function: components.Uptime, interval: 60 * 1000},
			{function: components.Filesystem, interval: 60 * 1000},
			{function: components.Thermal, interval: 1 * 1000},
			{register: components.Sound},
		},
	}
)

type StatusUpdate struct {
	index, row int
	instant    bool
	status     string
}

type Component struct {
	function   components.Basic
	register   components.Async
	interval   uint64
	aggregator chan *StatusUpdate
	index, row int
	instant    bool
}

func (cp *Component) UpdateStatus() {
	cp.Update(cp.Status())
}

func (cp *Component) Update(status string) {
	cp.aggregator <- &StatusUpdate{
		index:   cp.index,
		row:     cp.row,
		instant: cp.instant,
		status:  strings.TrimSpace(status),
	}
}

func (cp *Component) Status() string {
	return cp.function(cp.interval)
}

func (cp *Component) IsAsync() bool {
	return cp.function == nil || cp.interval == 0
}

func (cp *Component) Run() {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)

			if !ok {
				err = errors.New(r.(string))
			}

			cleanup(err)
		}
	}()

	if cp.interval < 0 {
		cp.UpdateStatus()
	} else if cp.IsAsync() {
		channel := make(chan string)
		go cp.register(channel)

		for msg := range channel {
			cp.Update(msg)
		}
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
			cp.aggregator = agg
			cp.index = index
			cp.row = row
			cp.instant = cp.IsAsync()
			go cp.Run()
		}
	}

	return agg
}

func main() {
	status := make([][]string, len(inputs))
	output := []string{statusSeparatorStart, "", statusSeparatorEnd}

	for i := 0; i < len(inputs); i++ {
		status[i] = make([]string, len(inputs[i]))

		if i == 0 {
			continue
		}

		output = append(output, dwmSeparator, statusSeparatorStart, "", statusSeparatorEnd)
	}

	updateStatus := func() {
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
