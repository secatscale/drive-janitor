//go:build windows
// +build windows

package checkage

import (
	"os"
	"syscall"
	"time"
)

func GetAgeWindows(info os.FileInfo) int {
	stat := info.Sys().(*syscall.Win32FileAttributeData)
	age := int(time.Since(time.Unix(0, stat.CreationTime.Nanoseconds())).Hours()) / 24
	return age
}

// will not be used, but is required for the build to succeed
func GetAgeDarwin(info os.FileInfo) int {
	return -1
}

// will not be used, but is required for the build to succeed
func GetAgeLinux(info os.FileInfo) int {
	return -1
}
