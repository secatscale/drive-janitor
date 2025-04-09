//go:build windows
// +build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/windows"
)

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
	if (osName == "windows") {
		// GUID for Downloads folder: {374DE290-123F-4565-9164-39C4925E467B}
		var rfid = windows.FOLDERID_Downloads

		// contains the address of a pointer to a null-terminated Unicode string that specifies the path of the known folder.
		var pathPtr uintptr
		// SHGetKnownFolderPath retrieves the path of a known folder from the repository identifier (PIDL) of the folder.
		/*
		The calling process is responsible for freeing this resource once it is no longer needed by calling CoTaskMemFree,
		whether SHGetKnownFolderPath succeeds or not.
		*/
		err := windows.SHGetKnownFolderPath(&rfid, 0, 0, &pathPtr)
		defer windows.CoTaskMemFree(unsafe.Pointer(pathPtr))
		if err != nil {
			return "", fmt.Errorf("SHGetKnownFolderPath failed: %w", err)
		}

		// Convert the pointer to a string
		path := windows.UTF16ToString((*[1 << 16]uint16)(unsafe.Pointer(pathPtr))[:])
		return path, nil
	}
	return "", fmt.Errorf("unsupported OS")
}