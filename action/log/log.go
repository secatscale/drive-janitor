package log

import (
	"fmt"
	"os"
	"sync"

	"github.com/gofrs/flock"
)

// Mutex global pour proteger l'acces concurrentiel au fichier
var fileMutex sync.Mutex

func SaveToFile(data string, logPath string) error {
	// Use internal mutex to serialize access within THIS process
	fileMutex.Lock()
	defer fileMutex.Unlock()

	// Now use flock for inter-process synchronization
	return writeToFileWithLock(data, logPath)
}

func writeToFileWithLock(data string, logPath string) error {
	// Create lock for inter-process synchronization
	lock := flock.New(logPath + ".lock") // Use separate lock file

	// Try to acquire the lock (this will wait if another process has it)
	err := lock.Lock()
	if err != nil {
		return fmt.Errorf("failed to acquire file lock: %w", err)
	}
	defer func() {
		if unlockErr := lock.Unlock(); unlockErr != nil {
			fmt.Printf("Warning: failed to unlock file: %v\n", unlockErr)
		}
	}()

	// Now safely open and write to the actual file
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close file: %v\n", closeErr)
		}
	}()

	// Write data
	_, err = f.WriteString(data)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	// Sync to ensure data is written to disk
	err = f.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	return nil
}
