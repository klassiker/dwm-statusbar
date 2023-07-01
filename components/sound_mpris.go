package components

import (
	"github.com/godbus/dbus"
	"strings"
)

var (
	SoundMPRISMatchPlayer = "type='signal',path='/org/mpris/MediaPlayer2',interface='org.freedesktop.DBus.Properties',member='PropertiesChanged',arg0='org.mpris.MediaPlayer2.Player'"
	SoundMPRISMatchStream = "type='signal',path='/org/mpris/MediaPlayer2',interface='org.freedesktop.DBus.Properties',member='PropertiesChanged',arg0='com.github.altdesktop.playerctld'"
)

type SoundMPRISStateStruct struct {
	artist, title, status, stream string
}

func NewSoundMPRISState() *SoundMPRISStateStruct {
	return &SoundMPRISStateStruct{
		artist: "",
		title:  "",
		status: "",
	}
}

func (ms *SoundMPRISStateStruct) State(state dbus.Variant) {
	ms.status, _ = state.Value().(string)
}

func (ms *SoundMPRISStateStruct) Metadata(data dbus.Variant) {
	metadata := data.Value().(map[string]dbus.Variant)
	artist, ok := metadata["xesam:artist"].Value().([]string)
	if ok {
		ms.artist = artist[0]
	}
	ms.title, ok = metadata["xesam:title"].Value().(string)
}

func (ms *SoundMPRISStateStruct) Stream(data dbus.Variant) {
	full, ok := data.Value().([]string)
	if ok && len(full) > 0 {
		ms.stream = strings.Split(full[0], ".")[3]
	}
}

func (ms *SoundMPRISStateStruct) Current() (string, string, string) {
	status := ms.status == "Playing"
	statusText := ms.status

	if status {
		statusText = "play"
	}

	return ms.stream, statusText, ms.title
}

func (ms *SoundMPRISStateStruct) IsActive() bool {
	return ms.status == "Playing"
}

func soundMPRISListen() {
	session, err := dbus.SessionBus()
	check(err)

	object := session.Object("org.mpris.MediaPlayer2.playerctld", "/org/mpris/MediaPlayer2")

	metadata, err := object.GetProperty("org.mpris.MediaPlayer2.Player.Metadata")
	if err == nil {
		SoundMPRISState.Metadata(metadata)
	}

	playbackStatus, err := object.GetProperty("org.mpris.MediaPlayer2.Player.PlaybackStatus")
	if err == nil {
		SoundMPRISState.State(playbackStatus)
	}

	stream, err := object.GetProperty("com.github.altdesktop.playerctld.PlayerNames")
	if err == nil {
		SoundMPRISState.Stream(stream)
	}

	_ = session.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, SoundMPRISMatchPlayer)
	_ = session.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, SoundMPRISMatchStream)

	ch := make(chan *dbus.Message, 1000)
	session.Eavesdrop(ch)
	for m := range ch {
		var ok bool
		var value dbus.Variant
		arg, ok := m.Body[1].(map[string]dbus.Variant)
		if !ok {
			continue
		}

		if value, ok = arg["PlaybackStatus"]; ok {
			SoundMPRISState.State(value)
		} else if value, ok = arg["Metadata"]; ok {
			SoundMPRISState.Metadata(value)
		} else if value, ok = arg["PlayerNames"]; ok {
			SoundMPRISState.Stream(value)
		}

		if ok {
			soundUpdate()
		}
	}
}
