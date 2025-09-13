package components

import (
	"slices"
	"testing"
	"time"
)

// TODO mock CPUFile to provide precise data
func Test_cpuReadData(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "does something",
		},
	}

	cores := 8
	CPUCores = cores

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := cpuReadData()

			// we get the overall data as well
			if len(data) != cores+1 {
				t.Errorf("got %d, want %d", len(data), cores+1)
			}

			if !slices.Equal(data, make([]float64, cores+1)) {
				t.Errorf("got %v, want empty", data)
			}

			time.Sleep(time.Second)

			data = cpuReadData()
			if slices.Equal(data, make([]float64, cores+1)) {
				t.Error("got empty, want at least something")
			}
		})
	}
}
