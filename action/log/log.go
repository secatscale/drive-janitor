package log

import (
	"os"
	"sync"

	"github.com/gofrs/flock"
)

// Mutex global pour proteger l'acces concurentiel au fichier
var fileMutex sync.Mutex

func SaveToFile(data string, logPath string) error {

	// Verrou interne entre les goroutines
	fileMutex.Lock()
	defer fileMutex.Unlock()

	// Lock + ouverture du fichier
	f, lock, err := loadFileWithlock(logPath)
	if err != nil {
		return err
	}
	defer closeFileAndUnlock(f, lock)

	// Ecrire les donnees dans le fichier
	_, err = f.WriteString(data)
	if err != nil {
		return err
	}
	return nil
}

func loadFileWithlock(logPath string) (*os.File, *flock.Flock, error) {
	// Create a lock on the file using the new API
	lock := flock.New(logPath)

	// Try to lock the file
	err := lock.Lock()
	if err != nil {
		return nil, nil, err
	}

	// Open or create the file
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		lock.Unlock()
		return nil, nil, err
	}

	return f, lock, nil
}

func closeFileAndUnlock(f *os.File, lock *flock.Flock) error {
	err := f.Close()
	if err != nil {
		lock.Unlock()
		return err
	}
	err = lock.Unlock()
	if err != nil {
		return err
	}
	return nil
}
