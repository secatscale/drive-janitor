package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// TestWhichOs tests the WhichOs function that returns the OS type
func TestWhichOs(t *testing.T) {
	t.Run("TestWhichOs is not null", func(t *testing.T) {
		osName := WhichOs()
		if osName == "" {
			t.Errorf("WhichOs() returned null")
		}
	})
	t.Run("TestWhichOs returns one of the expected values", func(t *testing.T) {
		osName := WhichOs()
		expected := []string{"windows", "linux", "darwin"}
		// fmt.Println(osName)
		// if osName is not in the expected list, then fail the test
		if !contains(expected, osName) {
			t.Errorf("WhichOs() = %v, want one of these : %v", osName, expected)
		}
	})
	t.Run("TestWhichOs returns the correct value", func(t *testing.T) {
		osName := WhichOs()
		expected := runtime.GOOS
		if osName != expected {
			t.Errorf("WhichOs() = %v, want %v", osName, expected)
		}
	})
}

// ngl the above test is overkill because the function is just returning runtime.GOOS

// TestWhereTrash tests the WhereTrash function that returns the path to the trash folder for an OS
func TestWhereTrash(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Errorf("os.UserHomeDir() returned an error: %v", err)
	}
	trashTests := []struct {
		testName  string
		osName    string
		trashPath string
	}{
		{testName: "Windows", osName: "windows", trashPath: "C:\\$Recycle.Bin"},
		{testName: "Linux", osName: "linux", trashPath: filepath.Join(home, "/.local/share/Trash")},
		{testName: "Darwin", osName: "darwin", trashPath: filepath.Join(home, "/.Trash")},
	}
	for _, tt := range trashTests {
		t.Run(tt.testName, func(t *testing.T) {
			if WhichOs() != tt.osName {
				t.Skip("Skipping test for OS: " + tt.osName)
			}
			trashPath, err := WhereTrash(tt.osName)
			if err != nil {
				t.Errorf("WhereTrash() returned an error: %v", err)
			}
			if trashPath == "" {
				t.Errorf("WhereTrash() returned null")
			}
			if trashPath != tt.trashPath {
				t.Errorf("WhereTrash() = %v, want %v", trashPath, tt.trashPath)
			}
		})
	}
}
