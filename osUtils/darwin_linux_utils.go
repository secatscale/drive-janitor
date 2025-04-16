//go:build !windows
// +build !windows

package os_utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// as a user
func GetDownloadPath() (string, error){
	osName := WhichOs()
	if (osName == "darwin" || osName == "linux") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		downloads := filepath.Join(home, "Downloads")
		return downloads, nil
	}
	return "", fmt.Errorf("unsupported OS")
}

func GetWindowsTrashPath() (string, error) {
	return "", errors.New("GetWindowsTrashPath is not supported on non-Windows platforms")
}

func GetCurrentUserSID() (string, error) {
	return "", errors.New("GetCurrentUserSID is not supported on non-Windows platforms")
}
