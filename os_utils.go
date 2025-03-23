package main

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

// WhichOs returns the OS type
func WhichOs() string {
	return runtime.GOOS
}

// overkill function, but it's here for testing purposes

func WhereTrash(osName string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	switch osName {
	case "windows":
		return filepath.Join("C:", "$Recycle.Bin"), nil
	case "linux":
		return filepath.Join(home, ".local", "share", "Trash"), nil
	case "darwin":
		return filepath.Join(home, ".Trash"), nil
	default:
		return "", errors.New("OS not supported")
	}
}
