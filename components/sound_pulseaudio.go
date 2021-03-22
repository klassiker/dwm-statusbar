package components

import (
	"github.com/godbus/dbus"
	"log"
	"sync"
	"time"
)

type SoundPulseaudioStateStruct struct {
	*sync.Mutex
	streams, filtered map[dbus.ObjectPath]map[string]string
}

func NewSoundPulseaudioState() *SoundPulseaudioStateStruct {
	return &SoundPulseaudioStateStruct{
		Mutex:    &sync.Mutex{},
		streams:  make(map[dbus.ObjectPath]map[string]string),
		filtered: make(map[dbus.ObjectPath]map[string]string),
	}
}

func (ps *SoundPulseaudioStateStruct) AddStream(path dbus.ObjectPath) {
	// we need to wait some time for the information to be available
	time.Sleep(100 * time.Millisecond)

	var properties map[string]string
	err := PulseaudioClient.Stream(path).Get("PropertyList", &properties)

	if err != nil {
		log.Println(err.Error(), path)
		return
	}

	ps.Lock()
	ps.streams[path] = properties

	binary := properties["application.process.binary"]
	for _, listed := range SoundPulseaudioBlacklist {
		if binary == listed {
			ps.Unlock()
			return
		}
	}

	ps.filtered[path] = properties
	ps.Unlock()

	soundUpdate()
}

func (ps *SoundPulseaudioStateStruct) RemoveStream(path dbus.ObjectPath) {
	ps.Lock()
	delete(ps.streams, path)
	delete(ps.filtered, path)
	ps.Unlock()

	soundUpdate()
}

func (ps *SoundPulseaudioStateStruct) Icons() []string {
	ps.Lock()
	defer ps.Unlock()

	icons := make([]string, len(ps.streams))

	i := 0
	for _, props := range ps.streams {
		icons[i] = props["application.process.binary"]
		i++
	}

	return icons
}

func (ps *SoundPulseaudioStateStruct) Current() (string, string, string) {
	ps.Lock()
	defer ps.Unlock()

	for _, props := range ps.filtered {
		return props["application.process.binary"], "unknown", props["media.name"]
	}

	return "", "unknown", ""
}

func (ps *SoundPulseaudioStateStruct) IsActive() bool {
	ps.Lock()
	defer ps.Unlock()

	return len(ps.filtered) > 0
}

func (cl *PulseaudioClientStruct) NewPlaybackStream(path dbus.ObjectPath) {
	//log.Println("NewPlaybackStream", path)
	SoundPulseaudioState.AddStream(path)
}

func (cl *PulseaudioClientStruct) PlaybackStreamRemoved(path dbus.ObjectPath) {
	//log.Println("PlaybackStreamRemoved", path)
	SoundPulseaudioState.RemoveStream(path)
}

func soundPulseaudio() {
	var list []dbus.ObjectPath
	check(PulseaudioCore.Get("PlaybackStreams", &list))

	for _, path := range list {
		SoundPulseaudioState.AddStream(path)
	}
}
