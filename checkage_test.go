package main

import (
	"path/filepath"
	"testing"
)

func TestCheckAge(t *testing.T) {
	tests := []struct {
		testname string
		fileName string
		age      int // expected age in days
	}{
		{"Test1", "sample.txt", 5},
		{"Test2", "sample.zip", 5},
		{"Test3", "sample.csv", 5},
		{"Test4", "empty", 0},
	}
	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			filePath, err := filepath.Abs(filepath.Join("samples", tt.fileName))
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
