package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestWhereTrash tests the WhereTrash function that returns the path to the trash folder for an OS
func TestWhereTrash(t *testing.T) {
	t.Run("HomeDir is not null", func(t *testing.T) {
		home, err := os.UserHomeDir()
		if err != nil {
			t.Fatalf("os.UserHomeDir() returned an error: %v", err)

		}
		if home == "" {
			t.Fatalf("os.UserHomeDir() returned null")
		}
	})
	home, _ := os.UserHomeDir()
	trashTests := []struct {
		testName  string
		osName    string
		trashPath string
	}{
		{testName: "Windows", osName: "windows", trashPath: "C:\\$Recycle.Bin"},
		{testName: "Linux", osName: "linux", trashPath: filepath.Join(home, ".local", "share", "Trash")},
		{testName: "Darwin", osName: "darwin", trashPath: filepath.Join(home, ".Trash")},
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
	t.Run("Unsupported OS", func(t *testing.T) {
		osName := "SkibidiBopBopDopaMina"
		trashPath, err := WhereTrash(osName)
		if err == nil {
			t.Errorf("WhereTrash() did not return an error for unsupported OS")
		}
		if trashPath != "" {
			t.Errorf("WhereTrash() returned a path for unsupported OS")
		}
	})
}
