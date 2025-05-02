//go:build !windows
// +build !windows

package os_utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// WhichOs returns the OS type
func WhichOs() string {
	return runtime.GOOS
}

// GetDownloadPath tries to find the user's download directory across Linux/macOS systems.
func GetDownloadPath() (string, error) {
	osName := WhichOs()
	if osName != "darwin" && osName != "linux" {
		return "", fmt.Errorf("unsupported OS: %s", osName)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// List of potential download paths to try, in order
	candidates := []string{}

	// Try XDG_DOWNLOAD_DIR environment variable first
	if xdgDownload := os.Getenv("XDG_DOWNLOAD_DIR"); xdgDownload != "" {
		candidates = append(candidates, xdgDownload)
	}

	// Common default folders
	candidates = append(candidates,
		filepath.Join(home, "Downloads"),
		filepath.Join(home, "Téléchargements"), // French
		filepath.Join(home, "Descargas"),       // Spanish
		filepath.Join(home, "Baixades"),        // Catalan
		filepath.Join(home, "Herunterladen"),   // German
	)

	// Try all candidates
	for _, path := range candidates {
		if stat, err := os.Stat(path); err == nil && stat.IsDir() {
			return path, nil
		}
	}

	return "", errors.New("no valid Downloads directory found")
}

func WhereTrash(osName string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %v", err)
	}
	switch osName {
	case "linux":
		return filepath.Join(home, ".local", "share", "Trash"), nil
	case "darwin":
		return filepath.Join(home, ".Trash"), nil
	default:
		return "", errors.New("OS not supported")
	}
}
