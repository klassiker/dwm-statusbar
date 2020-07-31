//+build ignore

package components

import (
	"github.com/godbus/dbus"
)

type ConfigFilesystem struct {
	path, icon string
}

type ConfigDBusService struct {
	unit               dbus.ObjectPath
	property, required string
}

var (
	BarHeight  = 21 - 2*BarPadding
	BarPadding = 1

	Batteries   = 2
	BatteryPath = "/sys/class/power_supply"

	FilesystemMounts = []ConfigFilesystem{
		{"/", IconFilesystemRoot},
		{"/home", IconFilesystemHome},
	}

	NetworkInterfaces = map[string]string{
		"eth0":  IconNetworkCable,
		"wlan0": IconNetworkWifi,
		"wwp0":  IconNetworkMobile,
	}
	NetworkVPNServices = []ConfigDBusService{
		{"openvpn-client@profile.service", "org.freedesktop.systemd1.Service.StatusText", "Initialization Sequence Completed"},
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
	IconBatteryPlug          = "\uf1e6"
	IconBatteryFull          = "\uf240"
	IconBatterHigh           = "\uf241"
	IconBatteryHalf          = "\uf242"
	IconBatteryLow           = "\uf243"
	IconBatteryEmpty         = "\uf244"
	IconCPU                  = "\uf2db"
	IconCurrentTimeCalendar  = "\uf073"
	IconCurrentTimeClock     = "\uf017"
	IconFilesystemRoot       = "\uf0a0"
	IconFilesystemHome       = "\uf015"
	IconMemory               = "\uf538"
	IconNetworkWifi          = "\uf1eb"
	IconNetworkCable         = "\uf796"
	IconNetworkMobile        = "\uf519"
	IconNetworkVPN           = "\uf084"
	IconPulseaudioHeadphones = "\uf025"
	IconPulseaudioVolumeNull = "\uf026"
	IconPulseaudioVolumeLow  = "\uf027"
	IconPulseaudioVolumeHigh = "\uf028"
	IconPulseaudioVolumeMute = "\uf6a9"
	IconSoundStatePlay       = "\uf04b"
	IconSoundStatePause      = "\uf04c"
	IconSoundStateStop       = "\uf04d"
	IconSoundStateUnknown    = "\uf128"
	IconSoundPlayerMPD       = "\uf001"
	IconSoundPlayerChrome    = "\uf268"
	IconSoundPlayerFirefox   = "\ue007"
	IconSoundPlayerMPV       = "\uf87c"
	IconSoundPlayerUnknown   = "\uf059"
	IconThermalCold          = "\uf2cb"
	IconThermalLow           = "\uf2ca"
	IconThermalOkay          = "\uf2c9"
	IconThermalHigh          = "\uf2c8"
	IconThermalBurn          = "\uf2c7"
)
