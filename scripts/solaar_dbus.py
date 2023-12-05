#!/usr/bin/env python
# -*- python-mode -*-

from logitech_receiver import base as _base
from solaar import NAME, __version__
# noinspection PyProtectedMember
from solaar.cli import _receivers
# noinspection PyProtectedMember
from solaar.cli.show import _battery_text

from dbus.mainloop.glib import DBusGMainLoop
import dbus
import dbus.service

from gi.repository import GLib


class DBusService(dbus.service.Object):
    def __init__(self, receivers, conn=None, object_path=None, bus_name=None):
        dbus.service.Object.__init__(self, conn, object_path, bus_name)
        self._receivers = receivers

    # noinspection PyPep8Naming
    @dbus.service.method(
        dbus_interface='io.github.pwr_solaar.solaar.Status',
        in_signature='s',
        out_signature='a{sv}',
        sender_keyword='sender',
        connection_keyword='conn'
    )
    def Battery(self, serial, sender=None, conn=None):
        device = self.device_for_serial(serial)

        result = {
            'level': 'N/A',
            'next': 'N/A',
            'status': 'N/A',
            'voltage': 'N/A',
        }

        try:
            device.ping()
        except _base.NoSuchDevice:
            print('Device not found: %s' % serial)
            result['status'] = 'unknown'
            return dbus.Dictionary(result, signature='sv')

        if not device.online:
            print('Device not online: %s' % serial)
            result['status'] = 'offline'
            return dbus.Dictionary(result, signature='sv')

        level, nextLevel, status, voltage = device.battery()
        result['level'] = _battery_text(level)
        result['next'] = _battery_text(nextLevel)
        result['status'] = str(status)
        result['voltage'] = voltage if voltage is not None else 'N/A'
        return dbus.Dictionary(result, signature='sv')

    def device_for_serial(self, serial):
        for receiver in self._receivers:
            count = receiver.count()
            if not count:
                continue

            for device in receiver:
                if device.serial == serial:
                    return device
        return None

    # noinspection PyMethodMayBeStatic
    def run(self):
        loop = GLib.MainLoop()
        loop.run()


def run(receivers, args=None, find_receiver=None, find_device=None):
    assert receivers

    print('%s version %s' % (NAME, __version__))

    DBusGMainLoop(set_as_default=True)
    bus = dbus.SessionBus()
    assert bus

    busname = dbus.service.BusName('io.github.pwr_solaar.solaar.Status', bus)
    service = DBusService(receivers, bus, '/io/github/pwr_solaar/solaar/Status', busname)
    service.run()


if __name__ == "__main__":
    run(list(_receivers()))
