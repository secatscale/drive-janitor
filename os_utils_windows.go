//go:build windows
// +build windows

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

func GetWindowsTrashPath() (string, error) {
	drive := os.Getenv("SystemDrive")
	if drive == "" {
		return "", errors.New("SystemDrive not set")
	}
	sid, err := GetCurrentUserSID()
	if err != nil {
		return "", err
	}
	return filepath.Join(drive+"\\", "$Recycle.Bin", sid), nil
}

// GetCurrentUser returns the current user on windows
func GetCurrentUserSID() (string, error) {
	token, err := windows.OpenCurrentProcessToken()
	if err != nil {
		return "", fmt.Errorf("failed to open current process token: %w", err)
	}
	user, err := token.GetTokenUser()
	if err != nil {
		return "", fmt.Errorf("failed to get token user: %w", err)
	}
	sid := user.User.Sid.String()
	if sid == "" {
		return "", errors.New("failed to get SID")
	}
	return sid, nil
}
