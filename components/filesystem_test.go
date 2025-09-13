package components

import (
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_filesystemStatus(t *testing.T) {
	type args struct {
		path   string
		mounts map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "drive is there",
			args: args{
				path:   "/home",
				mounts: map[string]string{"/home": "OK"},
			},
			want: "OK",
		},
		{
			name: "drive is missing",
			args: args{
				path:   "/home",
				mounts: map[string]string{},
			},
			want: "REMOVED",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filesystemStatus(tt.args.path, tt.args.mounts); got != tt.want {
				t.Errorf("filesystemStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filesystemData(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want *regexp.Regexp
	}{
		{
			name: "disk usage stats",
			args: args{
				path: "/",
			},
			want: regexp.MustCompile(`^[\d.]+[A-Z]?B/[\d.]+[A-Z]?B$`),
		},
		{
			name: "tmpfs usage stats",
			args: args{
				path: "/tmp",
			},
			want: regexp.MustCompile(`^[\d.]+[A-Z]?B/[\d.]+[A-Z]?B$`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filesystemData(tt.args.path); !tt.want.MatchString(got) {
				t.Errorf("filesystemData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filesystemMounts(t *testing.T) {
	tests := []struct {
		name string
		want map[string]string
	}{
		{
			name: "mounts",
			want: map[string]string{"/": regexp.MustCompile(`^[\d.]+[A-Z]?B/[\d.]+[A-Z]?B$`).String()},
		},
	}

	mounts := []ConfigFilesystem{
		{"/", "irrelevant"},
	}
	FilesystemMounts = mounts

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filesystemMounts(); !cmp.Equal(got, tt.want, regexpComparer) {
				t.Errorf("filesystemMounts() = %v, want %v", got, tt.want)
			}
		})
	}
}
