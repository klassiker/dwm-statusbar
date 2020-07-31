package components

import (
	"fmt"
	"github.com/godbus/dbus"
	"log"
	"time"
)

var (
	dbusConnection *dbus.Conn
)

func dbusIsMemberChar(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || c == '_'
}

func dbusStringProperty(unit dbus.ObjectPath, property, requiredValue string) bool {
	path := "/org/freedesktop/systemd1/unit/" + dbusEscape(unit)
	value := dbusProperty("org.freedesktop.systemd1", path, property).Value()
	return value == requiredValue
}

func dbusProperty(destination string, path dbus.ObjectPath, property string) dbus.Variant {
	object := dbusConnection.Object(destination, path)

	if !object.Path().IsValid() {
		log.Fatal("Invalid dbus path: ", object.Path())
	}

	variant, err := object.GetProperty(property)
	check(err)

	return variant
}

func dbusEscape(path dbus.ObjectPath) dbus.ObjectPath {
	if path.IsValid() {
		return path
	}

	var output []rune

	for _, char := range path {
		if dbusIsMemberChar(char) {
			output = append(output, char)
		} else {
			output = append(output, []rune(fmt.Sprintf("_%x", char))...)
		}
	}

	return dbus.ObjectPath(output)
}

func init() {
	start := time.Now()

	var err error
	dbusConnection, err = dbus.SystemBus()
	check(err)

	profilingLog(start)
}
