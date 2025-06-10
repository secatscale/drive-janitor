//go:build darwin
// +build darwin

package checkage

import (
	"os"
	"syscall"
	"time"
)

func GetAgeDarwin(info os.FileInfo) int {
	stat := info.Sys().(*syscall.Stat_t)
	age := int(time.Since(time.Unix(stat.Ctimespec.Sec, stat.Ctimespec.Nsec)).Hours()) / 24
	return age
}

func GetAgeLinux(info os.FileInfo) int {
	return -1
}

func GetAgeWindows(info os.FileInfo) int {
	return -1
}
