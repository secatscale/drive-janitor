package checkage

import (
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
			got, err := CheckAge(filePath)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}
			if got != tt.age {
				t.Errorf("got %d, want %d", got, tt.age)
			}
		})
	}
}
