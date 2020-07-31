package main

import (
	"errors"
	"fmt"
	"github.com/klassiker/dwm-statusbar/components"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"
	"time"
)

var (
	baseInterval         = 900
	dwmSeparator         = ";"
	statusSeparatorStart = "["
	statusSeparatorMid   = "] ["
	statusSeparatorEnd   = "]"
	inputs               = [2][]*Component{
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
	index   int
	status  string
	instant bool
}

type Component struct {
	function components.Basic
	register components.Async
	interval uint64
	channel  chan *StatusUpdate
	update   *StatusUpdate
}

func (cp *Component) UpdateStatus() {
	cp.Update(cp.Status())
}

func (cp *Component) Update(status string) {
	cp.update.status = strings.TrimSpace(status)
	cp.channel <- cp.update
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

		for {
			select {
			case msg := <-channel:
				cp.Update(msg)
			}
		}
	} else {
		cp.UpdateStatus()

		for range time.Tick(time.Duration(cp.interval) * time.Millisecond) {
			go cp.UpdateStatus()
		}
	}
}

func fillChannels(inputs []*Component, channels []chan *StatusUpdate, offset int) {
	for i, component := range inputs {
		position := i + offset

		component.channel = make(chan *StatusUpdate)
		component.update = &StatusUpdate{
			index:   position,
			instant: component.IsAsync(),
		}

		channels[position] = component.channel

		go component.Run()
	}
}

func init() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("- Ctrl+C pressed in Terminal")
		cleanup(nil)
	}()
}

func main() {
	inputsTop := len(inputs[0])
	inputsBot := len(inputs[1])
	inputsInt := inputsTop + inputsBot

	channels := make([]chan *StatusUpdate, inputsInt)
	fillChannels(inputs[0], channels, 0)
	fillChannels(inputs[1], channels, inputsTop)

	agg := make(chan *StatusUpdate)

	for _, ch := range channels {
		go func(c chan *StatusUpdate) {
			for msg := range c {
				agg <- msg
			}
		}(ch)
	}

	trigger := make(chan bool)
	statusTop := make([]string, inputsTop)
	statusBot := make([]string, inputsBot)

	go func() {
		outTop := []string{statusSeparatorStart, "", statusSeparatorEnd}
		outBot := []string{statusSeparatorStart, "", statusSeparatorEnd}

		for {
			// wait for ticker or trigger
			select {
			case <-time.Tick(time.Duration(baseInterval) * time.Millisecond):
			case <-trigger:
			}

			outTop[1] = strings.Join(statusTop, statusSeparatorMid)
			outBot[1] = strings.Join(statusBot, statusSeparatorMid)

			xsetroot(strings.Join(outTop, "") + dwmSeparator + strings.Join(outBot, ""))
		}
	}()

	var input *StatusUpdate
	for {
		input = <-agg

		if input.index < inputsTop {
			statusTop[input.index] = input.status
		} else {
			statusBot[input.index-inputsTop] = input.status
		}

		if input.instant {
			trigger <- true
		}
	}

}

func cleanup(err error) {
	if XConn != nil {
		if err != nil {
			xsetroot(err.Error())
		}

		XConn.Close()
	}

	if err != nil {
		debug.PrintStack()
	}

	os.Exit(1)
}
