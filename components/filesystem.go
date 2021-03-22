package components

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"syscall"
)

var (
	FilesystemPath  = "/proc/mounts"
	FilesystemUnits = []string{
		"B", "KB", "MB", "GB", "TB",
	}
)

func filesystemData(path string) string {
	var statFs syscall.Statfs_t
	check(syscall.Statfs(path, &statFs))

	blockSize := uint64(statFs.Bsize)
	filledSize := statFs.Blocks - statFs.Bfree
	totalSize := float64(blockSize * (filledSize + statFs.Bavail))
	usedSize := float64(blockSize * filledSize)

	unitUsed, unitTotal := calculateUnit(&usedSize, FilesystemUnits), calculateUnit(&totalSize, FilesystemUnits)

	usedSizeString := strconv.FormatFloat(math.Round(usedSize*100)/100, 'f', 2, 64)
	totalSizeString := strconv.FormatFloat(math.Round(totalSize), 'f', 0, 64)

	return fmt.Sprintf("%s%s/%s%s", usedSizeString, unitUsed, totalSizeString, unitTotal)
}

func filesystemMounts() map[string]string {
	mounts := make(map[string]string)

	data, err := os.ReadFile(FilesystemPath)
	check(err)

	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		mount := strings.Split(line, " ")[1]

		for _, config := range FilesystemMounts {
			if config.path == mount {
				mounts[mount] = filesystemData(config.path)
			}
		}
	}

	return mounts
}

func filesystemStatus(path string, mounts map[string]string) string {
	if value, ok := mounts[path]; ok {
		return value
	} else {
		return "REMOVED"
	}
}

func Filesystem(_ uint64) string {
	mounts := filesystemMounts()

	var output []string

	for _, config := range FilesystemMounts {
		output = append(output, config.icon, filesystemStatus(config.path, mounts))
	}

	return strings.Join(output, " ")
}
