package components

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"syscall"
)

var (
	FilesystemPath       = "/proc/mounts"
	FilesystemRootPath   = "/"
	FilesystemHomePath   = "/home"
	FilesystemSDCardPath = "/mnt/sdcard"
	FilesystemIconSDCard = "\uf7c2"
	FilesystemIconHDD    = "\uf0a0"
	FilesystemIconHome   = "\uf015"
	FilesystemUnits      = []string{
		"B", "KB", "MB", "GB", "TB",
	}
)

func filesystemData(path string) string {
	var statFs syscall.Statfs_t
	check(syscall.Statfs(path, &statFs))

	blockSize := uint64(statFs.Bsize)
	filledSize := statFs.Blocks - statFs.Bfree
	totalSize := float64(blockSize*(filledSize+statFs.Bavail))
	usedSize := float64(blockSize*filledSize)

	unitUsed, unitTotal := calculateUnit(&usedSize, FilesystemUnits), calculateUnit(&totalSize, FilesystemUnits)

	usedSizeString := strconv.FormatFloat(math.Round(usedSize*100)/100, 'f', 2, 64)
	totalSizeString := strconv.FormatFloat(math.Round(totalSize), 'f', 0, 64)

	return fmt.Sprintf("%s%s/%s%s", usedSizeString, unitUsed, totalSizeString, unitTotal)
}

func filesystemMounts() map[string]bool {
	mounts := make(map[string]bool)

	data, err := ioutil.ReadFile(FilesystemPath)
	check(err)

	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		mount := strings.Split(line, " ")[1]

		switch mount {
		case FilesystemRootPath, FilesystemHomePath, FilesystemSDCardPath:
			mounts[mount] = true
		}
	}

	return mounts
}

func filesystemStatus(path string, out []string, mounts map[string]bool, position int) {
	if _, ok := mounts[path]; ok {
		out[position] = filesystemData(path)
	} else {
		out[position] = "REMOVED"
	}
}

func init() {
	conf := Config["filesystem"]
	FilesystemRootPath = conf["root"]
	FilesystemHomePath = conf["home"]
	FilesystemSDCardPath = conf["sdcard"]
}

func Filesystem(_ uint64) string {
	mounts := filesystemMounts()

	output := []string{
		FilesystemIconHDD, "",
		FilesystemIconHome, "",
		FilesystemIconSDCard, "",
	}

	filesystemStatus(FilesystemRootPath, output, mounts, 1)
	filesystemStatus(FilesystemHomePath, output, mounts, 3)
	filesystemStatus(FilesystemSDCardPath, output, mounts, 5)

	return strings.Join(output, " ")
}
