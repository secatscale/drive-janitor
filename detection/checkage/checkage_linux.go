//go:build linux
// +build linux

package checkage

import (
	"os"
	"syscall"
	"time"
)

func GetAgeLinux(info os.FileInfo) int {
	stat := info.Sys().(*syscall.Stat_t)
	age := int(time.Since(time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)).Hours()) / 24
	return age
}

// will not be used, but is required for the build to succeed
func GetAgeDarwin(info os.FileInfo) int {
	return -1
}

// will not be used, but is required for the build to succeed
func GetAgeWindows(info os.FileInfo) int {
	return -1
}
