package components

import (
	"testing"
)

func Test_memoryReadData(t *testing.T) {
	tests := []struct {
		name string
		want map[string]func(int) bool
	}{
		{
			name: "does something",
			want: map[string]func(int) bool{
				"Buffers":      positiveInteger,
				"Cached":       positiveInteger,
				"MemFree":      positiveInteger,
				"MemTotal":     positiveInteger,
				"SReclaimable": positiveInteger,
				"Shmem":        positiveInteger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := memoryReadData()

			for field, checker := range tt.want {
				v, has := data[field]
				if !has {
					t.Errorf("memoryReadData() = field %s is missing", field)
				}
				if !checker(v) {
					t.Errorf("memoryReadData() = field %s has invalid value %d", field, v)
				}
			}
		})
	}
}

func Test_memoryCalculateUsed(t *testing.T) {
	type args struct {
		data map[string]int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "does something",
			args: args{
				data: map[string]int{
					"MemTotal":     100000,
					"MemFree":      1000,
					"Buffers":      1000,
					"Cached":       1000,
					"SReclaimable": 1000,
					"Shmem":        10000,
				},
			},
			want: 106000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := memoryCalculateUsed(tt.args.data); got != tt.want {
				t.Errorf("memoryCalculateUsed() = %v, want %v", got, tt.want)
			}
		})
	}
}
