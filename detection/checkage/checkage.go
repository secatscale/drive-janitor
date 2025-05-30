package checkage

import (
	"os"
	"time"
)

// CheckAge is a function that takes a file path and returns the age of the file in days.
func CheckAge(filePath string) (int, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	// ModTime returns the last modification time of the file
	// It's cross-platform compatible
	age := time.Since(info.ModTime()).Hours() / 24
	return int(age), nil
}
