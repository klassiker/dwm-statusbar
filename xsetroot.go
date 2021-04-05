package main

import (
	"flag"
	"fmt"
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
	"github.com/klassiker/dwm-statusbar/components"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"runtime/pprof"
	"syscall"
)

var (
	xConn                *xgb.Conn
	windowRoot           xproto.Window
	atomName, atomString xproto.Atom
	shutdown             = make(chan bool)
)

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func xsetroot(name string) {
	if *debugFlag {
		log.Println(name)
		return
	}
	data := []byte(name)
	check(xproto.ChangePropertyChecked(xConn, xproto.PropModeReplace, windowRoot, atomName, atomString, 8, uint32(len(data)), data).Check())
}

func getAtom(name string) xproto.Atom {
	atom, err := xproto.InternAtom(xConn, true, uint16(len(name)), name).Reply()
	check(err)
	return atom.Atom
}

var memprofile = flag.String("memprofile", "", "write memory profile to file")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var debugFlag = flag.Bool("debug", false, "don't call xsetroot")
var nodraw = flag.Bool("nodraw", false, "don't use drawings and colors")

func init() {
	var err error

	xConn, err = xgb.NewConn()
	check(err)

	windowRoot = xproto.Setup(xConn).DefaultScreen(xConn).Root
	atomName = getAtom("WM_NAME")
	atomString = getAtom("UTF8_STRING")

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("- Ctrl+C pressed in Terminal")
		shutdown <- true
	}()

	flag.Parse()

	if *cpuprofile != "" {
		log.Println("CPU profiling enable for file:", *cpuprofile)
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		check(pprof.StartCPUProfile(f))
	}

	components.NoDraw = *nodraw
}

func cleanup(err error) {
	if xConn != nil {
		if err != nil {
			xsetroot(err.Error())
		}

		xConn.Close()
	}

	if err != nil {
		debug.PrintStack()
	}

	pprof.StopCPUProfile()
	if *memprofile != "" {
		log.Println("Memory profiling enabled for file:", *memprofile)
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		check(pprof.WriteHeapProfile(f))
		check(f.Close())
	}

	os.Exit(1)
}
