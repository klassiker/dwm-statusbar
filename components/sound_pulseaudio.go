package components

import (
	"github.com/godbus/dbus"
	"time"
)

type SoundPulseaudioStateStruct struct {
	streams, filtered map[dbus.ObjectPath]map[string]string
}

func NewSoundPulseaudioState() *SoundPulseaudioStateStruct {
	return &SoundPulseaudioStateStruct{
		streams:  make(map[dbus.ObjectPath]map[string]string),
		filtered: make(map[dbus.ObjectPath]map[string]string),
	}
}

func (ps *SoundPulseaudioStateStruct) AddStream(path dbus.ObjectPath) {
	time.Sleep(100 * time.Millisecond)

	var properties map[string]string
	check(PulseaudioClient.Stream(path).Get("PropertyList", &properties))
	ps.streams[path] = properties
	ps.filtered[path] = properties

	binary := properties["application.process.binary"]
	for _, listed := range SoundPulseaudioBlacklist {
		if binary == listed {
			delete(ps.filtered, path)
			return
		}
	}

	soundUpdate()
}

func (ps *SoundPulseaudioStateStruct) RemoveStream(path dbus.ObjectPath) {
	delete(ps.streams, path)
	delete(ps.filtered, path)
	soundUpdate()
}

func (ps *SoundPulseaudioStateStruct) Icons() []string {
	icons := make([]string, len(ps.streams))

	i := 0
	for _, props := range ps.streams {
		icons[i] = props["application.process.binary"]
		i++
	}

	return icons
}

func (ps *SoundPulseaudioStateStruct) Current() (string, string, string) {
	for _, props := range ps.filtered {
		return props["application.process.binary"], "unknown", props["media.name"]
	}

	return "", "unknown", ""
}

func (ps *SoundPulseaudioStateStruct) IsActive() bool {
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
		go SoundPulseaudioState.AddStream(path)
	}
}
