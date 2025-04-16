//go:build windows
// +build windows

package os_utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

func GetDownloadPath() (string, error) {
	osName := WhichOs()
	if osName == "windows" {
		// GUID for Downloads folder: {374DE290-123F-4565-9164-39C4925E467B}
		var rfid = windows.FOLDERID_Downloads

		// It call SHGetKnownFolderPath
		path, err := windows.KnownFolderPath(rfid, 0)
		if err != nil {
			return "", fmt.Errorf("SHGetKnownFolderPath failed: %w", err)
		}
		fmt.Println(path)
		return path, nil
	}
	return "", fmt.Errorf("unsupported OS")
}

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

