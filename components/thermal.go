package components

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	ThermalInputPattern = regexp.MustCompile(`^temp[0-9]+_input$`)
	ThermalPath         = "/sys/class/hwmon/"
	ThermalPerfect      = "^d^"
	ThermalGood         = "^c#00ff00^"
	ThermalOkay         = "^c#ffff00^"
	ThermalBad          = "^c#ff0000^"
)

func thermalTemperature(input string) float64 {
	data, err := ioutil.ReadFile(input)
	check(err)

	tempRaw, err := strconv.Atoi(strings.TrimSpace(string(data)))
	check(err)

	return float64(tempRaw) / 1000.0
}

func thermalTemperatureDrawing(temperature float64) string {
	var icon string
	var status string

	switch {
	case temperature < 40.0:
		status = ThermalPerfect
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

	return fmt.Sprintf("%s%s %0.1f°C", status, icon, temperature)
}

func thermalInputsByNames(names []string) []string {
	var out []string

	tmp := make([][]string, len(names))

	for i := 0; i < len(tmp); i++ {
		tmp[i] = []string{}
	}

	dirs, err := ioutil.ReadDir(ThermalPath)
	check(err)

	for _, dir := range dirs {
		path, err := filepath.EvalSymlinks(ThermalPath + dir.Name())
		check(err)

		resolved, err := os.Stat(path)
		check(err)

		if !resolved.IsDir() {
			continue
		}

		files, err := ioutil.ReadDir(path)
		check(err)

		var index int

		for _, file := range files {
			if file.Name() != "name" {
				continue
			}

			nameRaw, err := ioutil.ReadFile(path + "/name")
			check(err)

			name := strings.TrimSpace(string(nameRaw))

			index = indexOf(names, name)
		}

		if index == -1 {
			continue
		}

		for _, file := range files {
			if ThermalInputPattern.Match([]byte(file.Name())) {
				tmp[index] = append(tmp[index], path+"/"+file.Name())
			}
		}
	}

	for _, arr := range tmp {
		out = append(out, arr...)
	}

	return out
}

func init() {
	start := time.Now()

	ThermalInputs = thermalInputsByNames(ThermalInputs)

	if len(ThermalInputs) == 0 {
		panic(errors.New("thermal: no thermal input found"))
	}

	profilingLog(start)
}

func Thermal(_ uint64) string {
	output := make([]string, len(ThermalInputs))

	for i, input := range ThermalInputs {
		temperature := thermalTemperature(input)
		output[i] = thermalTemperatureDrawing(temperature)
	}

	return strings.Join(output, " ") + "^d^"
}
