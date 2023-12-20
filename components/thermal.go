package components

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	ThermalInputPattern = regexp.MustCompile(`^temp[0-9]+_input$`)
	ThermalPath         = "/sys/class/hwmon/"
	ThermalPerfect      = drawColor("")
	ThermalCold         = drawColor("#0000ff")
	ThermalGood         = drawColor("#00ff00")
	ThermalOkay         = drawColor("#ffff00")
	ThermalBad          = drawColor("#ff0000")
)

func thermalTemperature(hwmon ConfigThermal) string {
	out := make([]string, len(hwmon.inputs))

	for i, input := range hwmon.inputs {
		_, err := input.Seek(0, 0)
		check(err)

		data, err := io.ReadAll(input)
		check(err)

		tempRaw, err := strconv.Atoi(strings.TrimSpace(string(data)))
		check(err)

		out[i] = thermalTemperatureDrawing(float64(tempRaw) / 1000.0)
	}

	return strings.Join(out, " ")
}

func thermalTemperatureDrawing(temperature float64) string {
	var icon string
	var status string

	switch {
	case temperature < 40.0:
		status = ThermalCold
		icon = IconThermalCold
	case temperature < 55.0:
		status = ThermalPerfect
		icon = IconThermalLow
	case temperature < 70.0:
		status = ThermalGood
		icon = IconThermalOkay
	case temperature < 80.0:
		status = ThermalOkay
		icon = IconThermalHigh
	default:
		status = ThermalBad
		icon = IconThermalBurn
	}

	if NoDraw {
		status = ""
	}

	reset := DrawReset
	if status == "" {
		reset = ""
	}

	return fmt.Sprintf("%s%s %0.1f°C%s", status, icon, temperature, reset)
}

func thermalInputsByNames(thermals []ConfigThermal) {
	hwmons, err := os.ReadDir(ThermalPath)
	check(err)

	for _, hwmon := range hwmons {
		path, err := filepath.EvalSymlinks(ThermalPath + hwmon.Name())
		check(err)

		resolved, err := os.Stat(path)
		check(err)

		if !resolved.IsDir() {
			continue
		}

		// TODO allow find by device or multiple devices with the same name
		nameRaw, err := os.ReadFile(filepath.Join(path, "name"))
		if err != nil {
			fmt.Println("thermal: read name", path, "error:", err)
			continue
		}

		name := strings.TrimSpace(string(nameRaw))

		files, err := os.ReadDir(path)
		check(err)

		var inputs []*os.File
		for _, file := range files {
			if ThermalInputPattern.MatchString(file.Name()) {
				file, err := os.Open(filepath.Join(path, file.Name()))
				check(err)
				inputs = append(inputs, file)
			}
		}

		for i, thermal := range thermals {
			if thermal.name == name {
				thermals[i].inputs = inputs
			}
		}
	}
}

func init() {
	thermalInputsByNames(ThermalHwmons)

	for _, hwmon:= range ThermalHwmons {
		if len(hwmon.inputs) == 0 {
			panic(fmt.Errorf("thermal: no thermal input found for %s", hwmon.name))
		}
	}
}

func Thermal(_ int64) string {
	output := make([]string, len(ThermalHwmons))

	for i, hwmon := range ThermalHwmons {
		output[i] = thermalTemperature(hwmon)
	}

	return strings.Join(output, " ")
}
