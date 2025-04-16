package os_utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// WhichOs returns the OS type
// overkill function
func WhichOs() string {
	return runtime.GOOS
}

func WhereTrash(osName string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %v", err)
	}
	switch osName {
	case "windows":
		return GetWindowsTrashPath()
	case "linux":
		return filepath.Join(home, ".local", "share", "Trash"), nil
	case "darwin":
		return filepath.Join(home, ".Trash"), nil
	default:
		return "", errors.New("OS not supported")
	}
}
