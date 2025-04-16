package os_utils

import (
	"runtime"
	"slices"
	"testing"
)

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
		if !slices.Contains(expected, osName) {
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
