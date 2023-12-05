//go:build ignore
// +build ignore

package components

import (
	"github.com/godbus/dbus"
)

type HidType int

type ConfigSolaar struct {
	serial  string
	hidType HidType
}

type ConfigFilesystem struct {
	path, icon string
}

type ConfigNetwork struct {
	name   dbus.ObjectPath
	status string
}

var (
	NoDraw = false

	BarHeight  = 21 - 2*BarPadding
	BarPadding = 1

	Batteries = 2

	SolaarDevices = []ConfigSolaar{
		{"", HidKeyboard},
		{"", HidMouse},
	}

	FilesystemMounts = []ConfigFilesystem{
		{"/", IconFilesystemRoot},
		{"/home", IconFilesystemHome},
	}

	NetworkInterfaces = map[string]string{
		"eth0":  IconNetworkCable,
		"wlan0": IconNetworkWifi,
		"wwp0":  IconNetworkMobile,
	}
	NetworkVPNServices = []ConfigNetwork{
		{"openvpn-client@profile.service", "Initialization Sequence Completed"},
	}

	PulseaudioDevice        = "alsa_output.pci-0000_00_00.0."
	PulseaudioHeadphonePort = "analog-output-headphones"
	PulseaudioSinkPrefix    = "/org/pulseaudio/core1/sink"

	SoundMPDSocket           = "/home/user/.mpd/socket"
	SoundTitleMaxLength      = 64
	SoundPulseaudioBlacklist = []string{"mpd", "chrome", "chromium", "firefox"}
	SoundPlayerIconMap       = map[string]string{
		"mpd":      IconSoundPlayerMPD,
		"mpv":      IconSoundPlayerMPV,
		"chrome":   IconSoundPlayerChrome,
		"chromium": IconSoundPlayerChrome,
		"firefox":  IconSoundPlayerFirefox,
		"unknown":  IconSoundStateUnknown,
	}
	SoundStateIconMap = map[string]string{
		"pause":   IconSoundStatePause,
		"stop":    IconSoundStateStop,
		"play":    IconSoundStatePlay,
		"unknown": IconSoundStateUnknown,
	}

	ThermalInputs = []string{
		"coretemp",
	}
)

const (
	HidMouse HidType = iota
	HidKeyboard
)

const (
	IconBatteryPlug          = "\uf1e6" // plug
	IconBatteryFull          = "\uf240" // battery-full
	IconBatterHigh           = "\uf241" // battery-three-quarters
	IconBatteryHalf          = "\uf242" // battery-half
	IconBatteryLow           = "\uf243" // battery-quarter
	IconBatteryEmpty         = "\uf244" // battery-empty
	IconSolaarCharging       = "\uf0e7" // bolt
	IconSolaarOffline        = "\uf059" // circle-question
	IconSolaarUnknown        = "\uf06a" // circle-exclamation
	IconSolaarWarning        = "\uf071" // triangle-exclamation
	IconSolaarError          = "\uf00d" // xmark
	IconSolaarMouse          = "\uf8cc" // computer-mouse
	IconSolaarKeyboard       = "\uf11c" // keyboard
	IconCPU                  = "\uf2db" // microchip
	IconCurrentTimeCalendar  = "\uf073" // calendar-alt
	IconCurrentTimeClock     = "\uf017" // clock
	IconFilesystemRoot       = "\uf0a0" // hdd
	IconFilesystemHome       = "\uf015" // home
	IconMemory               = "\uf538" // memory
	IconNetworkWifi          = "\uf1eb" // wifi
	IconNetworkCable         = "\uf796" // ethernet
	IconNetworkMobile        = "\uf519" // broadcast-tower
	IconNetworkVPN           = "\uf084" // key
	IconPulseaudioHeadphones = "\uf025" // headphones
	IconPulseaudioVolumeNull = "\uf026" // volume-off
	IconPulseaudioVolumeLow  = "\uf027" // volume-down
	IconPulseaudioVolumeHigh = "\uf028" // volume-up
	IconPulseaudioVolumeMute = "\uf6a9" // volume-mute
	IconSoundStatePlay       = "\uf04b" // play
	IconSoundStatePause      = "\uf04c" // pause
	IconSoundStateStop       = "\uf04d" // stop
	IconSoundStateUnknown    = "\uf128" // question
	IconSoundPlayerMPD       = "\uf001" // music
	IconSoundPlayerChrome    = "\uf268" // chrome
	IconSoundPlayerFirefox   = "\ue007" // firefox-browser
	IconSoundPlayerMPV       = "\uf87c" // photo-video
	IconSoundPlayerMail      = "\uf0e0" // envelope
	IconSoundPlayerUnknown   = "\uf059" // question-circle
	IconThermalCold          = "\uf2cb" // thermometer-empty
	IconThermalLow           = "\uf2ca" // thermometer-quarter
	IconThermalOkay          = "\uf2c9" // thermometer-half
	IconThermalHigh          = "\uf2c8" // thermometer-three-quarters
	IconThermalBurn          = "\uf2c7" // thermometer-full
	IconUptime               = "\uf062" // arrow-up
)
