package main

import (
	"runtime"
)

// WhichOs returns the OS type
// overkill function
func WhichOs() string {
	return runtime.GOOS
}