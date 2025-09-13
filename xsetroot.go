package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"runtime/pprof"
	"sync"
	"syscall"
	"time"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
	"github.com/klassiker/dwm-statusbar/components"
)

var (
	xConn                *xgb.Conn
	windowRoot           xproto.Window
	atomName, atomString xproto.Atom
	shutdown             = make(chan bool)
	connection           = &sync.Mutex{}
)

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func xsetroot(name string) {
	if *debugFlag {
		log.Println(len(name), name)
		return
	}

	connection.Lock()
	defer connection.Unlock()

	for i := 0; i < 5; i++ {
		data := []byte(name)
		err := xproto.ChangePropertyChecked(xConn, xproto.PropModeReplace, windowRoot, atomName, atomString, 8, uint32(len(data)), data).Check()
		if err != nil {
			log.Println("xsetroot failed:", err, "trying reconnect", name)
			time.Sleep(5 * time.Millisecond)
			connect()
		} else {
			break
		}
	}
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
	connect()

	// this needs to be buffered
	c := make(chan os.Signal, 32)
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

func connect() {
	var err error

	xConn, err = xgb.NewConn()
	check(err)

	windowRoot = xproto.Setup(xConn).DefaultScreen(xConn).Root
	atomName = getAtom("WM_NAME")
	atomString = getAtom("UTF8_STRING")
}

func recovery() {
	if r := recover(); r != nil {
		err, ok := r.(error)

		if !ok {
			err = errors.New(r.(string))
		}

		cleanup(err)
	}
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
