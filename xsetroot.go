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

			xConn, err = xgb.NewConn()
			if err != nil {
				log.Fatal(err)
			}

			windowRoot = xproto.Setup(xConn).DefaultScreen(xConn).Root

			atomName, err = getAtom("WM_NAME")
			if err != nil {
				log.Fatal(err)
			}

			atomString, err = getAtom("UTF8_STRING")
			if err != nil {
				log.Fatal(err)
			}
		} else {
			break
		}
	}
}

func getAtom(name string) (xproto.Atom, error) {
	atom, err := xproto.InternAtom(xConn, true, uint16(len(name)), name).Reply()
	if err != nil {
		return 0, err
	}
	return atom.Atom, err
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

		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal(err)
		}
	}

	components.NoDraw = *nodraw
}

func connect() {
	var err error

	xConn, err = xgb.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	windowRoot = xproto.Setup(xConn).DefaultScreen(xConn).Root

	atomName, err = getAtom("WM_NAME")
	if err != nil {
		log.Fatal(err)
	}

	atomString, err = getAtom("UTF8_STRING")
	if err != nil {
		log.Fatal(err)
	}
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
		if err = pprof.WriteHeapProfile(f); err != nil {
			log.Fatal(err)
		}
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}

	os.Exit(1)
}
