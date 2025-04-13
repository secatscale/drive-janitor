//go:build windows
// +build windows

package main

import (
	"fmt"

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
