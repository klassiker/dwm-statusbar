package components

import (
	"fmt"
	"github.com/godbus/dbus"
	"github.com/sqp/pulseaudio"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

var (
	PulseaudioCore   *pulseaudio.Object
	PulseaudioClient *PulseaudioClientStruct
	PulseaudioState  = new(PulseaudioStateStruct)
)

type PulseaudioStateStruct struct {
	muted         bool
	volume        int
	activePort    dbus.ObjectPath
	headphonePort dbus.ObjectPath
	volumeRead    []uint32
	channel       chan string
}

func (ps *PulseaudioStateStruct) Headphones() bool {
	return PulseaudioState.activePort == PulseaudioState.headphonePort
}

func (ps *PulseaudioStateStruct) Icon() string {
	switch {
	case ps.muted:
		return IconPulseaudioVolumeMute
	case ps.volume >= 50:
		return IconPulseaudioVolumeHigh
	case ps.volume >= 25:
		return IconPulseaudioVolumeLow
	default:
		return IconPulseaudioVolumeNull
	}
}

func (ps *PulseaudioStateStruct) Port(port dbus.ObjectPath) {
	ps.activePort = port
	ps.Update()
}

func (ps *PulseaudioStateStruct) Muted(muted bool) {
	ps.muted = muted
	ps.Update()
}

func (ps *PulseaudioStateStruct) Volume(values []uint32) {
	var sum float64

	for _, value := range values {
		sum += math.Round(float64(value) / 65536.0 * 100)
	}

	ps.volume = int(sum / float64(len(values)))
	ps.Update()
}

func (ps *PulseaudioStateStruct) Reset() {
	var sinks []dbus.ObjectPath
	check(PulseaudioCore.Get("Sinks", &sinks))

	for _, sink := range sinks {
		var name string
		conf := PulseaudioClient.Device(sink)
		check(conf.Get("Name", &name))

		if strings.HasPrefix(name, PulseaudioDevice) {
			obj := PulseaudioClient.Device(sink)

			check(obj.Get("Mute", &ps.muted))
			check(obj.Get("Volume", &ps.volumeRead))
			check(obj.Get("ActivePort", &ps.activePort))

			_ = obj.Call("GetPortByName", 0, PulseaudioHeadphonePort).Store(&ps.headphonePort)

			PulseaudioClient.DeviceVolumeUpdated(sink, ps.volumeRead)

			return
		}
	}

	log.Println("Pulseaudio: no active sink with device prefix found!")
}

func (ps *PulseaudioStateStruct) Channel(channel chan string) {
	ps.channel = channel
	ps.Update()
}

func (ps *PulseaudioStateStruct) Update() {
	volume := fmt.Sprintf("%s%%", strconv.Itoa(ps.volume))
	output := []string{ps.Icon()}

	if ps.Headphones() {
		output = append(output, IconPulseaudioHeadphones, volume)
	} else {
		output = append(output, volume)
	}

	ps.channel <- strings.Join(output, " ")
}

type PulseaudioClientStruct struct {
	*pulseaudio.Client
}

func (cl *PulseaudioClientStruct) NewSink(_ dbus.ObjectPath) {
	PulseaudioState.Reset()
}

func (cl *PulseaudioClientStruct) SinkRemoved(_ dbus.ObjectPath) {
	PulseaudioState.Reset()
}

func (cl *PulseaudioClientStruct) DeviceActivePortUpdated(path dbus.ObjectPath, port dbus.ObjectPath) {
	if strings.HasPrefix(string(path), PulseaudioSinkPrefix) {
		PulseaudioState.Port(port)
	}
}

func (cl *PulseaudioClientStruct) DeviceMuteUpdated(path dbus.ObjectPath, muted bool) {
	if strings.HasPrefix(string(path), PulseaudioSinkPrefix) {
		PulseaudioState.Muted(muted)
	}
}

func (cl *PulseaudioClientStruct) DeviceVolumeUpdated(path dbus.ObjectPath, values []uint32) {
	if strings.HasPrefix(string(path), PulseaudioSinkPrefix) {
		PulseaudioState.Volume(values)
	}
}

func init() {
	start := time.Now()

	pulse, err := pulseaudio.New()
	check(err)

	PulseaudioClient = &PulseaudioClientStruct{pulse}
	PulseaudioCore = PulseaudioClient.Core()

	pulse.Register(PulseaudioClient)

	go PulseaudioState.Reset()
	go pulse.Listen()

	profilingLog(start)
}

func Pulseaudio(channel chan string) {
	PulseaudioState.Channel(channel)
}
