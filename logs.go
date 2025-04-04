package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gofrs/flock"
)

// Mutex global pour proteger l'acces concurentiel au fichier
var fileMutex sync.Mutex

func SaveToFile(data string) error {

	// Creer le dossier logs s'il n'existe pas
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		return err
	}

	// Creer le chemin du fichier de log
	date := time.Now().Format("2006-01-02")
	logPath := fmt.Sprintf("drive_janitor_logs_%s.log", date)
	logPath = filepath.Join("logs", logPath)
	logPath, err = filepath.Abs(logPath)
	if err != nil {
		return err
	}

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
	// Creer un lock sur le fichier
	lock := flock.NewFlock(logPath)
	// Essayer de locker le fichier
	err := lock.Lock()
	if err != nil {
		return nil, nil, err
	}
	// Ouvrir ou creer le fichier
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
		return err
	}
	err = lock.Unlock()
	if err != nil {
		return err
	}
	return nil
}
