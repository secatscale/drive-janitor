//go:build !windows
// +build !windows

// go build !windows
package os_utils

import (
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