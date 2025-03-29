package testhelper

import (
	"runtime"
	"testing"
)

func isOSSupported(osNames map[string]bool) bool {
	// Check if the current OS is in the map
	_, exists := osNames[runtime.GOOS]
	return exists
}

// RunOSDependentTest filters and runs tests based on the OS.
func RunOSDependentTest(t *testing.T, testName string, testFunc func(t *testing.T), osNames map[string]bool) {
	t.Run(testName + " OS: "+runtime.GOOS, func(t *testing.T) {
		if (!isOSSupported(osNames)) {
			t.Skipf("Skipping test as it's not applicable for %s", runtime.GOOS)
		} else {
			t.Run(testName, testFunc)
		}
	})
}
