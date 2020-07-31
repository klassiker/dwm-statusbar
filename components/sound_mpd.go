package components

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"time"
)

var (
	SoundMPDCommand = []byte("command_list_ok_begin\nstatus\ncurrentsong\nidle player\ncommand_list_end\n")
)

type SoundMPDStateStruct struct {
	values     map[string]string
	connection *net.UnixConn
}

func NewSoundMPDState() *SoundMPDStateStruct {
	return &SoundMPDStateStruct{
		values: make(map[string]string),
	}
}

func (ms *SoundMPDStateStruct) Set(key, value string) {
	ms.values[strings.ToLower(key)] = value
}

func (ms *SoundMPDStateStruct) Has(key string) (ok bool) {
	_, ok = ms.values[strings.ToLower(key)]
	return
}

func (ms *SoundMPDStateStruct) Reset() {
	ms.values = map[string]string{"state": "stop", "file": "", "artist": "", "title": "", "changed": ""}
}

func (ms *SoundMPDStateStruct) Connect() {
	ms.Reset()

	addr, err := net.ResolveUnixAddr("unix", SoundMPDSocket)
	check(err)

	ms.connection, err = net.DialUnix("unix", nil, addr)
	if err != nil {
		if strings.HasSuffix(err.Error(), "connection refused") {
			return
		} else {
			check(err)
		}
	}

	_, err = ms.connection.Read(make([]byte, 1024))
	check(err)
}

func (ms *SoundMPDStateStruct) Connected() bool {
	return !ms.Disconnected()
}

func (ms *SoundMPDStateStruct) Disconnected() bool {
	return ms.connection == nil
}

func (ms *SoundMPDStateStruct) Parse(data []byte) {
	response := bytes.Split(bytes.Trim(data, "\x00"), []byte{'\x0a'})

	for _, line := range response {
		//log.Printf("%s", line)
		parts := strings.Split(string(line), ":")

		if ms.Has(parts[0]) {
			ms.Set(parts[0], strings.TrimSpace(parts[1]))
		}
	}

	if ms.values["changed"] == "player" {
		ms.Reset()
		ms.Update()
	}
}

func (ms *SoundMPDStateStruct) Update() {
	_, err := ms.connection.Write(SoundMPDCommand)
	check(err)

	responseRead := make([]byte, 4096)
	_, err = ms.connection.Read(responseRead)
	check(err)

	ms.Parse(responseRead)
	if SoundChannel != nil {
		soundUpdate()
	}
}

func (ms *SoundMPDStateStruct) Current() (string, string, string) {
	var title string

	if ms.values["artist"] != "" && ms.values["title"] != "" {
		title = fmt.Sprintf("%s - %s", ms.values["artist"], ms.values["title"])
	} else {
		title = ms.values["file"]
	}

	return "mpd", ms.values["state"], title
}

func (ms *SoundMPDStateStruct) IsActive() bool {
	return ms.values["state"] == "play"
}

func soundMPDListen() {
	for {
		if SoundMPDState.Disconnected() {
			SoundMPDState.Connect()

			if SoundMPDState.Disconnected() {
				time.Sleep(10 * time.Second)
				continue
			}
		}

		SoundMPDState.Update()

		_, err := SoundMPDState.connection.Read(make([]byte, 1024))
		if err != nil {
			if err.Error() == "EOF" {
				SoundMPDState.Connect()
			} else {
				check(err)
			}
		}
	}
}
