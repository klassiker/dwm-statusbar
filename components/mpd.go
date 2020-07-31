package components

import (
	"bytes"
	"fmt"
	"net"
	"strings"
)

var (
	MPDConn *net.UnixConn
	MPDCommand   = []byte("command_list_ok_begin\nstatus\ncurrentsong\ncommand_list_end\n")
	MPDSocket    = "~/.mpd/socket"
	MPDIcon      = "\uf001"
	MPDIconPlay  = "\uf04b"
	MPDIconPause = "\uf04c"
	MPDIconStop  = "\uf04d"
	MPDIconMap = map[string]string{
		"pause": MPDIconPause,
		"stop": MPDIconStop,
		"play": MPDIconPlay,
	}
)

func mpdRead(amount int) [][]byte {
	if MPDConn == nil {
		return [][]byte{}
	}

	read := make([]byte, amount)
	_, err := MPDConn.Read(read)
	check(err)

	read = bytes.Trim(read, "\x00")

	return bytes.Split(read, []byte{'\x0a'})
}

func mpdSend(msg []byte) {
	if MPDConn == nil {
		err := mpdConnect()

		if err != nil {
			return
		}
	}

	_, err := MPDConn.Write(msg)

	if err != nil {
		_ = mpdConnect()
	}
}

func mpdConnect() error {
	addr, err := net.ResolveUnixAddr("unix", MPDSocket)
	check(err)

	MPDConn, err = net.DialUnix("unix", nil, addr)

	if err == nil {
		_ = mpdRead(4096)
	}

	return err
}

func init() {
	MPDSocket = Config["mpd"]["socket"]
	_ = mpdConnect()
}

func MPD(_ uint64) string {
	state := map[string]string{
		"state": "",
		"file": "",
		"Artist": "",
		"Title": "",
	}

	mpdSend(MPDCommand)
	response := mpdRead(4096)

	for _, line := range response {
		parts := strings.Split(string(line), ":")
		key := parts[0]

		if _, ok := state[key]; ok {
			state[key] = strings.TrimSpace(parts[1])
		}
	}

	var title string

	if state["Artist"] != "" && state["Title"] != "" {
		title = fmt.Sprintf("%s - %s", state["Artist"], state["Title"])
	} else {
		title = state["file"]
	}

	output := []string{
		MPDIcon, MPDIconMap[state["state"]], title,
	}

	return strings.Join(output, " ")
}
