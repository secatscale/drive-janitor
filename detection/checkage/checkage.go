package checkage

import (
	"drive-janitor/os_utils"
	"fmt"
	"os"
)

// CheckAge is a function that takes a file path and returns the age of the file in days.
func GetAge(filePath string) (int, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, fmt.Errorf("error getting file info for %s: %v", filePath, err)
	}
	var age int
	switch os_utils.WhichOs() {
	case "darwin":
		age = GetAgeDarwin(info)
	case "linux":
		age = GetAgeLinux(info)
	case "windows":
		age = GetAgeWindows(info)
	default:
		return 0, fmt.Errorf("unsupported OS: %s", os_utils.WhichOs())
	}
	// ModTime returns the last modification time of the file
	// It's cross-platform compatible
	// age := time.Since(info.ModTime()).Hours() / 24
	return age, nil
}
