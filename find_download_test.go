package main

import (
	"drive-janitor/testhelper"
	"testing"
)

func TestGetDownloadPath(t *testing.T) {
	testhelper.RunOSDependentTest(t, "Getting download path on windows", func(t *testing.T) {
		path, err := GetDownloadPath()
		if err != nil {
			t.Fatalf("Error getting download path: %v", err)
		}
		if path == "" {
			t.Fatal("Download path is empty")
		}
		t.Logf("Download path: %s", path)
	},map[string]bool{"windows": true})

	testhelper.RunOSDependentTest(t, "Getting download path on darwin", func(t *testing.T) {
		path, err := GetDownloadPath()
		if err != nil {
			t.Fatalf("Error getting download path: %v", err)
		}
		if path == "" {
			t.Fatal("Download path is empty")
		}
		t.Logf("Download path: %s", path)
	}, map[string]bool{"darwin": true, "linux": true})
}
