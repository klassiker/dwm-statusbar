package components

import (
	"fmt"
	"github.com/godbus/dbus"
	"github.com/sqp/pulseaudio"
	"math"
	"strconv"
	"strings"
)

var (
	PulseaudioDevice             = "alsa_output.pci-XXXX_XX_XX.X.analog-stereo"
	PulseaudioHeadphonePort      = "analog-output-headphones"
	PulseaudioSink               = dbus.ObjectPath("/org/pulseaudio/coreX/sinkX")
	PulseaudioClient             = new(PulseaudioClientStruct)
	PulseaudioState              = new(PulseaudioStateStruct)
	PulseaudioVolumeFull         = 65536.0
	PulseaudioIconHeadphones     = "\uf025"
	PulseaudioIconVolumeNullIcon = "\uf026"
	PulseaudioIconVolumeLowIcon  = "\uf027"
	PulseaudioIconVolumeHighIcon = "\uf028"
	PulseaudioIconVolumeMuteIcon = "\uf6a9"
)

type PulseaudioStateStruct struct {
	muted         bool
	volume        int
	activePort    dbus.ObjectPath
	headphonePort dbus.ObjectPath
	volumeRead    []uint32
}

type PulseaudioClientStruct struct {
	*pulseaudio.Client
}

func (cl *PulseaudioClientStruct) DeviceActivePortUpdated(path dbus.ObjectPath, port dbus.ObjectPath) {
	if path != PulseaudioSink {
		return
	}

	PulseaudioState.activePort = port
}

func (cl *PulseaudioClientStruct) DeviceMuteUpdated(path dbus.ObjectPath, muted bool) {
	if path != PulseaudioSink {
		return
	}

	PulseaudioState.muted = muted
}

// We take the average of all channels
func (cl *PulseaudioClientStruct) DeviceVolumeUpdated(path dbus.ObjectPath, values []uint32) {
	if path != PulseaudioSink {
		return
	}

	var sum float64

	for _, value := range values {
		sum += math.Round(float64(value) / PulseaudioVolumeFull * 100)
	}

	PulseaudioState.volume = int(sum / float64(len(values)))
}

func init() {
	conf := Config["pulseaudio"]
	PulseaudioDevice = conf["device"]
	PulseaudioHeadphonePort = conf["headphonePort"]
	PulseaudioSink = dbus.ObjectPath(conf["sinkPath"])

	pulse, err := pulseaudio.New()
	check(err)

	PulseaudioClient.Client = pulse
	pulse.Register(PulseaudioClient)

	// Try to find the sink by the device name, if that fails use default
	var sinks []dbus.ObjectPath
	core := PulseaudioClient.Core()
	check(core.Get("Sinks", &sinks))

	for _, sink := range sinks {
		var name string
		conf := PulseaudioClient.Device(sink)
		check(conf.Get("Name", &name))

		if name == PulseaudioDevice {
			PulseaudioSink = sink
		}
	}

	device := PulseaudioClient.Device(PulseaudioSink)
	check(device.Get("Mute", &PulseaudioState.muted))
	check(device.Get("Volume", &PulseaudioState.volumeRead))
	check(device.Get("ActivePort", &PulseaudioState.activePort))
	check(device.Call("GetPortByName", 0, PulseaudioHeadphonePort).Store(&PulseaudioState.headphonePort))

	PulseaudioClient.DeviceVolumeUpdated(PulseaudioSink, PulseaudioState.volumeRead)

	go pulse.Listen()
}

func pulseaudioVolumeIcon() string {
	var volumeIcon string
	switch {
	case PulseaudioState.muted:
		volumeIcon = PulseaudioIconVolumeMuteIcon
	case PulseaudioState.volume >= 50:
		volumeIcon = PulseaudioIconVolumeHighIcon
	case PulseaudioState.volume >= 25:
		volumeIcon = PulseaudioIconVolumeLowIcon
	default:
		volumeIcon = PulseaudioIconVolumeNullIcon
	}
	return volumeIcon
}

func Pulseaudio(_ uint64) string {
	headphones := PulseaudioState.activePort == PulseaudioState.headphonePort
	volume := fmt.Sprintf("%s%%", strconv.Itoa(PulseaudioState.volume))
	volumeIcon := pulseaudioVolumeIcon()

	var output []string
	if headphones {
		output = []string{
			volumeIcon, PulseaudioIconHeadphones, volume,
		}
	} else {
		output = []string{
			volumeIcon, volume,
		}
	}

	return strings.Join(output, " ")
}
