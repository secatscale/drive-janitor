//go:build !windows
// +build !windows

// go build !windows
package os_utils

import (
	"drive-janitor/os_utils"
	"fmt"
	"os"
	"path/filepath"
)

// as a user
func GetDownloadPath() (string, error){
	osName := os_utils.WhichOs()
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