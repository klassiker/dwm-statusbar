package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"statusbar/components"
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
	inputs               = [2][]Component{
		{
			{components.Network, 2 * 1000, nil},
			{components.Pulseaudio, 2 * 1000, nil},
			{components.Memory, 2 * 1000, nil},
			{components.CPUPercentBar, 2 * 1000, nil},
			{components.Battery, 10 * 1000, nil},
			{components.CurrentTime, 1 * 1000, nil},
		},
		{
			{components.Uptime, 60 * 1000, nil},
			{components.Filesystem, 60 * 1000, nil},
			{components.Thermal, 1 * 1000, nil},
			{components.MPD, 2 * 1000, nil},
		},
	}
)

type StatusUpdate struct {
	index  int
	status string
}

type Component struct {
	function components.Basic
	interval uint64
	channel  chan *StatusUpdate
}

func (c Component) Run(position int) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)

			if !ok {
				err = errors.New(r.(string))
			}

			cleanup(err)
		}
	}()

	status := &StatusUpdate{index: position}

	if c.interval < 0 {
		res := c.function(c.interval)

		status.status = res

		c.channel <- status
		return
	}

	for {
		res := c.function(c.interval)

		status.status = strings.TrimSpace(res)

		c.channel <- status

		time.Sleep(time.Duration(c.interval) * time.Millisecond)
	}
}

func fillChannels(inputs []Component, channels []chan *StatusUpdate, offset int) {
	for i, component := range inputs {
		component.channel = make(chan *StatusUpdate)
		channels[i+offset] = component.channel
		go component.Run(i + offset)
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

	statusTop := make([]string, inputsTop)
	statusBot := make([]string, inputsBot)

	go func() {
		outTop := []string{statusSeparatorStart, "", statusSeparatorEnd}
		outBot := []string{statusSeparatorStart, "", statusSeparatorEnd}

		for {
			time.Sleep(time.Duration(baseInterval) * time.Millisecond)

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
