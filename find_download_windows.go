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

		// contains the address of a pointer to a null-terminated Unicode string that specifies the path of the known folder.
		// SHGetKnownFolderPath retrieves the path of a known folder from the repository identifier (PIDL) of the folder.
		/*
			The calling process is responsible for freeing this resource once it is no longer needed by calling CoTaskMemFree,
			whether SHGetKnownFolderPath succeeds or not.
		*/
		path, err := windows.KnownFolderPath(rfid, 0)
		if err != nil {
			return "", fmt.Errorf("SHGetKnownFolderPath failed: %w", err)
		}
		fmt.Println(path)
		return path, nil
	}
	return "", fmt.Errorf("unsupported OS")
}
