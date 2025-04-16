package os_utils

import (
	"drive-janitor/testhelper"
	"os"
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
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Fatalf("Download path does not exist: %v", path)
		}
		t.Logf("Download path: %s", path)
	},map[string]bool{"windows": true, "darwin": true, "linux": true})
}
