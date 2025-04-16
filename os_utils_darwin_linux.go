//go:build !windows
// +build !windows

// go build !windows
package main

import "errors"

func GetWindowsTrashPath() (string, error) {
	return "", errors.New("GetWindowsTrashPath is not supported on non-Windows platforms")
}

func GetCurrentUserSID() (string, error) {
	return "", errors.New("GetCurrentUserSID is not supported on non-Windows platforms")
}
