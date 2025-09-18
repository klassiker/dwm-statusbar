# DWM Statusbar

## TODO

- passive listener for solaar

## Description

A simple statusbar for dwm (6.2) written in go without using external commands.

With support for multiple batteries, thermal monitors, pulseaudio streams, mpris metadata, systemd service monitoring.

Dependencies:
- ttf-font-awesome
- github.com/godbus/dbus (for dbus communication - pulseaudio, service monitoring)
- github.com/sqp/pulseaudio (to communicate with the pulseaudio daemon)
- github.com/jezek/xgb (to avoid using xsetroot as external command)

Example output (there are more colors in use depending on the sensor levels):

![](statusbar.png)

## Install

`git clone`, configure, `go install`. Done.

To use the Solaar component for battery status, you need to enable a service running the script `scripts/solaar_dbus.py`, which uses the solaar libraries to provide a DBus service.

Since this bar uses colors, draws stuff and uses two bars this patch `dwm-status2d-extrabar-6.2.diff` is highly recommended: https://dwm.suckless.org/patches/status2d/

I use a slightly modified version of the patch which allows you to use a configurable amount of additional bars that are automatically hidden when there is nothing to display (statusbar can not handle this yet, I only use 1 bar at the moment) or manually reduced if you press a shortcut. It also fixes a calculation bug in `status2d` which prevented the `ClkStatusText` button press from working and crashes on incomplete formatting input. Diffs for 6.2 and current master are in `diffs/`, if you find them useful feel free to submit them to the official `patches` pages.

If you don't want to use two separate bars use only one array of components in the inputs in `statusbar.go`.

If you don't want drawings and colors use the `-nodraw` flag.

## Config

Outputs are configured in `statusbar.go` with refresh interval in miliseconds. The interval is passed to the component function. There is an additional `baseInterval` which limits the update rate of the bar.

Async components have a register function and use a channel to instantly update the statusbar.

Modify `config.go` in the `components` folder.

Where to find configuration values:

- The hwmon names found in `/sys/class/hwmon/hwmon*/name`. All available inputs are used.
- Pulseaudio device `pacmd list-sinks | grep 'name:'`
- Pulseaudio headphone port `pacmd list-sinks | grep ports -A 20`. Most likely `analog-output-headphones`
- Pulseaudio sink path `pacmd list-sinks | grep 'index:'` /org/pulseaudio/core1/sink`X`
- Network VPN services (names from systemctl)
- Network VPN services property/required value: `d-feet` or `gdbus introspect --session --dest=org.freedesktop.systemd1 --object-path /org/freedesktop/systemd1/unit/unit_2eservice`
