package components

import (
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/godbus/dbus"
)

var (
	NetworkPath  = "/proc/net/dev"
	NetworkData  = make(map[string]*NetworkDataStore)
	NetworkUnits = []string{"KB", "MB"}

	networkDbus *dbus.Conn
)

func init() {
	var err error
	networkDbus, err = dbus.SystemBusPrivate()
	check(err)

	check(dbusPrivate(networkDbus))

	for name := range NetworkInterfaces {
		NetworkData[name] = &NetworkDataStore{
			last: &NetworkDataStore{},
		}
	}
}

type NetworkDataStore struct {
	rx, tx uint64
	last   *NetworkDataStore
}

func networkUnitActive(unit ConfigNetwork) bool {
	path := "/org/freedesktop/systemd1/unit/" + dbusEscape(unit.name)
	object := networkDbus.Object("org.freedesktop.systemd1", path)

	if !object.Path().IsValid() {
		log.Fatal("Invalid dbus path: ", object.Path())
	}

	variant, err := object.GetProperty("org.freedesktop.systemd1.Service.StatusText")
	check(err)

	return variant.Value() == unit.status
}

func networkReadData() {
	networkFile, err := os.Open(NetworkPath)
	check(err)

	data, err := io.ReadAll(networkFile)
	check(err)

	fields := strings.FieldsFunc(string(data), func(r rune) bool { return r == ' ' })
	lines := strings.Split(strings.Join(fields, " "), "\n")
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

		last := NetworkData[name].last
		NetworkData[name].rx = rx - last.rx
		NetworkData[name].tx = tx - last.tx
		last.rx = rx
		last.tx = tx
	}
}

func networkCalculateSpeed(iface string, interval float64) string {
	data := NetworkData[iface]
	// interval is in ms, we need seconds
	rx, tx := float64(data.rx)/interval*1000/1024.0, float64(data.tx)/interval*1000/1024.0
	rxUnit, txUnit := calculateUnit(&rx, NetworkUnits), calculateUnit(&tx, NetworkUnits)

	rxString := strconv.FormatFloat(math.Round(rx*100)/100, 'f', 2, 64)
	txString := strconv.FormatFloat(math.Round(tx*100)/100, 'f', 2, 64)

	return fmt.Sprintf("%s%s / %s%s", txString, txUnit, rxString, rxUnit)
}

func Network(interval int64) string {
	var output []string

	// TODO use a passive dbus listener to reduce traffic, reduces execution time by 15ms
	for _, unit := range NetworkVPNServices {
		if networkUnitActive(unit) {
			output = append(output, IconNetworkVPN)
			break
		}
	}

	ifaces, err := net.Interfaces()
	check(err)

	networkReadData()

	for _, iface := range ifaces {
		if net.FlagUp&iface.Flags != net.FlagUp {
			continue
		}

		if icon, ok := NetworkInterfaces[iface.Name]; ok {
			speed := networkCalculateSpeed(iface.Name, float64(interval))
			output = append(output, fmt.Sprintf("%s %s", icon, speed))
		}
	}

	return strings.Join(output, " ")
}
