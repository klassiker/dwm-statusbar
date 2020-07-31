//+build ignore

package components

var Config = map[string]map[string]string{
	"thermal": {
		"hwmons": "acpitz,coretemp",
	},
	"pulseaudio": {
		"device": "alsa_output.pci-XXXX_XX_XX.X.analog-stereo",
		"headphonePort": "analog-output-headphones",
		"sinkPath": "/org/pulseaudio/coreX/sinkX",
	},
	"network": {
		"wired": "eth0",
		"wireless": "wlan0",
	},
	"mpd": {
		"socket": "/home/X/.mpd/socket",
	},
	"filesystem": {
		"root": "/",
		"home": "/home",
	},
}
