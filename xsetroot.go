package main

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"log"
)

var (
	XConn                *xgb.Conn
	windowRoot           xproto.Window
	atomName, atomString xproto.Atom
)

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func xsetroot(name string) {
	data := []byte(name)
	check(xproto.ChangePropertyChecked(XConn, xproto.PropModeReplace, windowRoot, atomName, atomString, 8, uint32(len(data)), data).Check())
}

func getAtom(name string) xproto.Atom {
	atom, err := xproto.InternAtom(XConn, true, uint16(len(name)), name).Reply()
	check(err)
	return atom.Atom
}

func init() {
	var err error

	XConn, err = xgb.NewConn()
	check(err)

	windowRoot = xproto.Setup(XConn).DefaultScreen(XConn).Root
	atomName = getAtom("WM_NAME")
	atomString = getAtom("UTF8_STRING")
}
