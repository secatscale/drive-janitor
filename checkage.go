package main

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
	age := time.Since(info.ModTime()).Hours() / 24
	// fmt.Println(age)
	return int(age), nil
}
