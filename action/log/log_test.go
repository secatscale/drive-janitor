package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{"test1", "test1, test1, test1\n"},
		{"test2", "test2, test2, test2\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SaveToFile(tt.data)
			if err != nil {
				t.Errorf("SaveToFile failed: %v", err)
			}
			assertDataIsInFile(t, tt.data)
		})
	}
}

// Test avec goroutines simultanées
func TestConcurrentLogging(t *testing.T) {
	count := 50
	dataPrefix := "concurrent-log-"

	done := make(chan bool)
	for i := 0; i < count; i++ {
		go func(i int) {
			msg := fmt.Sprintf("%s%s #%d\n", dataPrefix, time.Now().Format("15:04:05.000"), i)
			err := SaveToFile(msg)
			if err != nil {
				t.Errorf("SaveToFile failed in goroutine: %v", err)
			}
			done <- true
		}(i)
	}

	// Attend que toutes les goroutines aient fini
	for i := 0; i < count; i++ {
		<-done
	}
}

// Vérifie que la donnée est bien dans le fichier de log du jour
func assertDataIsInFile(t *testing.T, data string) {
	date := time.Now().Format("2006-01-02")
	logPath := filepath.Join("logs", "drive_janitor_logs_"+date+".log")

	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), data) {
		t.Errorf("Expected data not found in log file:\nExpected: %q\nFile Content: %s", data, string(content))
	}
}
