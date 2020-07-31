package components

import (
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"regexp"
	"strconv"
	"strings"
)

var (
	NetworkRegex       = regexp.MustCompile(`\s{2,}`)
	NetworkPath        = "/proc/net/dev"
	NetworkWired       = "eth0"
	NetworkWireless    = "wlan0"
	NetworkData        = make(map[string]NetworkDataStore)
	NetworkIconWIFI    = "\uf1eb"
	NetworkIconCable   = "\uf796"
	NetworkCache       = make(map[string]NetworkDataStore)
	NetworkUnits       = []string{"KB", "MB"}
	NetworkIconMap     = map[string]string{
		NetworkWired:    NetworkIconCable,
		NetworkWireless: NetworkIconWIFI,
	}
)

type NetworkDataStore struct {
	rxBytes, txBytes uint64
}

func networkReadData() {
	data, err := ioutil.ReadFile(NetworkPath)
	check(err)

	lines := strings.Split(NetworkRegex.ReplaceAllString(string(data), " "), "\n")

	for _, line := range lines {
		fields := strings.Split(line, " ")
		name := strings.TrimSuffix(fields[0], ":")

		if _, ok := NetworkIconMap[name]; !ok {
			continue
		}

		rx, err := strconv.ParseUint(fields[1], 10, 64)
		check(err)

		tx, err := strconv.ParseUint(fields[9], 10, 64)
		check(err)

		NetworkCache[name] = NetworkDataStore{
			rxBytes: rx,
			txBytes: tx,
		}
	}
}

func networkCalculateSpeed(iface string, interval uint64) string {
	cacheData, cached := NetworkCache[iface]

	if !cached {
		networkReadData()
		cacheData = NetworkCache[iface]
	}

	netData, ok := NetworkData[iface]

	if !ok {
		netData = cacheData
	}

	rxValue, txValue := float64((cacheData.rxBytes - netData.rxBytes) / interval * 1000) / 1024.0, float64((cacheData.txBytes - netData.txBytes) / interval * 1000) / 1024.0

	unitRx, unitTx := calculateUnit(&rxValue, NetworkUnits), calculateUnit(&txValue, NetworkUnits)

	rxString := strconv.FormatFloat(math.Round(rxValue * 100) / 100, 'f', 2, 64)
	txString := strconv.FormatFloat(math.Round(txValue * 100) / 100, 'f', 2, 64)

	NetworkData[iface] = cacheData
	delete(NetworkCache, iface)

	return fmt.Sprintf("%s%s / %s%s", txString, unitTx, rxString, unitRx)
}

func init() {
	conf := Config["network"]
	NetworkWired = conf["wired"]
	NetworkWireless = conf["wireless"]

	NetworkIconMap     = map[string]string{
		NetworkWired:    NetworkIconCable,
		NetworkWireless: NetworkIconWIFI,
	}
}

func Network(interval uint64) string {
	var output []string

	ifaces, err := net.Interfaces()
	check(err)

	for _, iface := range ifaces {
		isUp := net.FlagUp&iface.Flags == net.FlagUp

		if !isUp {
			continue
		}

		if icon, ok := NetworkIconMap[iface.Name]; ok {
			output = append(output, fmt.Sprintf("%s %s", icon, networkCalculateSpeed(iface.Name, interval)))
		}
	}

	return strings.Join(output, " ")
}
