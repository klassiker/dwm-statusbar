package components

import (
	"strings"
)

var (
	SoundUpdate          func(string)
	SoundPulseaudioState = NewSoundPulseaudioState()
	SoundMPDState        = NewSoundMPDState()
	SoundMPRISState      = NewSoundMPRISState()
	SoundStates          = SoundStatesStruct{
		states: []SoundState{
			SoundMPDState,
			SoundMPRISState,
			SoundPulseaudioState,
		},
		defaultState: SoundMPDState,
	}
)

type SoundState interface {
	Current() (string, string, string)
	IsActive() bool
}

type SoundStatesStruct struct {
	states       []SoundState
	defaultState SoundState
}

func (sss *SoundStatesStruct) Current() (string, string, string) {
	for _, state := range sss.states {
		if state.IsActive() {
			return state.Current()
		}
	}

	return sss.defaultState.Current()
}

func (sss *SoundStatesStruct) Player(player string) string {
	return mapValueOrDefault(SoundPlayerIconMap, player, IconSoundPlayerUnknown)
}

func (sss *SoundStatesStruct) State(state string) string {
	return mapValueOrDefault(SoundStateIconMap, state, IconSoundStateUnknown)
}

func (sss *SoundStatesStruct) Icons(current string) []string {
	icons := SoundPulseaudioState.Icons()

	var playerIcons []string
	for _, icon := range icons {
		if icon != current {
			playerIcons = append(playerIcons, sss.Player(icon))
		}
	}

	return playerIcons
}

func (sss *SoundStatesStruct) Title(title, state string) string {
	if state == "stop" {
		return ""
	} else if len(title) > SoundTitleMaxLength {
		return strings.TrimSpace(title[:SoundTitleMaxLength-3]) + "..."
	} else {
		return title
	}
}

func (sss *SoundStatesStruct) Output() string {
	player, state, title := sss.Current()

	title = sss.Title(title, state)
	stateIcon := sss.State(state)
	playerIcons := sss.Icons(player)
	playerIcon := sss.Player(player)

	var output []string
	if len(playerIcons) > 0 {
		output = append(playerIcons, "/")
	}
	output = append(output, playerIcon, stateIcon, title)

	return strings.Join(output, " ")
}

func soundUpdate() {
	SoundUpdate(SoundStates.Output())
}

func Sound(update func(string)) {
	SoundUpdate = update
	go soundPulseaudio()
	go soundMPDListen()
	go soundMPRISListen()
	soundUpdate()
}
