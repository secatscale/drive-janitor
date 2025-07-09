package checkage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckAge(t *testing.T) {
	age := 0
	tests := []struct {
		testname string
		fileName string
		age      int // expected age in days
	}{
		{"Test1", "sample.txt", age},
		{"Test2", "sample.zip", age},
		{"Test3", "sample.csv", age},
		{"Test4", "empty", age},
	}
	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			filePath, err := filepath.Abs(filepath.Join("../../samples", tt.fileName))
			if err != nil {
				t.Fatalf("Error getting file path: %s", err)
			}

			// we need to touch the file to set the change time to 0
			f, err := os.OpenFile(filePath, os.O_WRONLY, 0644)
			if err != nil {
				t.Fatalf("Error opening file: %s", err)
			}
			defer f.Close()
			// Write zero bytes at offset 0 (no content change, but updates ctime!)
			if _, err := f.WriteAt([]byte{}, 0); err != nil {
				t.Fatalf("Error touching file: %s", err)
			}

			got, err := GetAge(filePath)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}
			if got != tt.age {
				t.Errorf("got %d, want %d", got, tt.age)
			}
		})
	}
}
