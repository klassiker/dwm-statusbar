# DWM Statusbar

## Description

A simple statusbar for dwm (6.2) written in go without using external commands.

With support for multiple batteries, thermal monitors, pulseaudio streams, mpris metadata, systemd service monitoring.

Dependencies:
- ttf-font-awesome
- github.com/godbus/dbus (for dbus communication - pulseaudio, service monitoring)
- github.com/sqp/pulseaudio (to communicate with the pulseaudio daemon)
- github.com/BurntSushi/xgb (to avoid using xsetroot as external command)

Example output (there are more colors in use depending on the sensor levels):

![](statusbar.png)

## Install

Since we use colors, draw stuff and use two bars the patch `dwm-status2d-extrabar-6.2.diff` is required: https://dwm.suckless.org/patches/status2d/

## Config

Outputs are configured in `statusbar.go` with refresh interval in miliseconds. The interval is passed to the component function. There is an additional `baseInterval` which limits the update rate of the bar.

Async components have a register function and use a channel to instantly update the statusbar.

Copy the `config.def.go` in the components folder to `config.go` and remove the first line.

Where to find configuration values:

- The hwmon names found in `/sys/class/hwmon/hwmon*/name`. All available inputs are used.
- Pulseaudio device `pacmd list-sinks | grep 'name:'`
- Pulseaudio headphone port `pacmd list-sinks | grep ports -A 20`. Most likely `analog-output-headphones`
- Pulseaudio sink path `pacmd list-sinks | grep 'index:'` /org/pulseaudio/core1/sink`X`
- Network VPN services (names from systemctl)
- Network VPN services property/required value: `d-feet` or `gdbus introspect --session --dest=org.freedesktop.systemd1 --object-path /org/freedesktop/systemd1/unit/unit_2eservice`
