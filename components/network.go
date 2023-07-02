package components

import (
	"fmt"
	"math"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	NetworkRegex = regexp.MustCompile(`\s{2,}`)
	NetworkPath  = "/proc/net/dev"
	NetworkData  = make(map[string]NetworkDataStore)
	NetworkCache = make(map[string]NetworkDataStore)
	NetworkUnits = []string{"KB", "MB"}
)

type NetworkDataStore struct {
	rxBytes, txBytes uint64
}

func networkReadData() {
	data, err := os.ReadFile(NetworkPath)
	check(err)

	lines := strings.Split(NetworkRegex.ReplaceAllString(string(data), " "), "\n")

	for _, line := range lines {
		fields := strings.Split(line, " ")
		name := strings.TrimSuffix(fields[0], ":")

		if _, ok := NetworkInterfaces[name]; !ok || len(fields) <= 1 {
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

	rxValue, txValue := float64((cacheData.rxBytes-netData.rxBytes)/interval*1000)/1024.0, float64((cacheData.txBytes-netData.txBytes)/interval*1000)/1024.0

	unitRx, unitTx := calculateUnit(&rxValue, NetworkUnits), calculateUnit(&txValue, NetworkUnits)

	rxString := strconv.FormatFloat(math.Round(rxValue*100)/100, 'f', 2, 64)
	txString := strconv.FormatFloat(math.Round(txValue*100)/100, 'f', 2, 64)

	NetworkData[iface] = cacheData
	delete(NetworkCache, iface)

	return fmt.Sprintf("%s%s / %s%s", txString, unitTx, rxString, unitRx)
}

func Network(interval int64) string {
	var output []string

	// TODO use a passive dbus listener to reduce traffic, reduces execution time by 15ms
	for _, unit := range NetworkVPNServices {
		if dbusStringProperty(unit.unit, unit.property, unit.required) {
			output = append(output, IconNetworkVPN)
			break
		}
	}

	ifaces, err := net.Interfaces()
	check(err)

	for _, iface := range ifaces {
		if !(net.FlagUp&iface.Flags == net.FlagUp) {
			continue
		}

		if icon, ok := NetworkInterfaces[iface.Name]; ok {
			speed := networkCalculateSpeed(iface.Name, uint64(interval))
			output = append(output, fmt.Sprintf("%s %s", icon, speed))
		}
	}

	return strings.Join(output, " ")
}
