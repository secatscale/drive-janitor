package action

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (action *Action) TakeAction(filePath string) {
	if action.Delete {
		// TODO: Implement delete action
		// os.Remove(filePath)
	}
	if action.Log {
		action.LogConfig.Files = append(action.LogConfig.Files, filePath)
	}
}

func (action *Action) GetLogFileName() error {
	// Creer le chemin du fichier de log
	date := time.Now().Format("2006-01-02")
	var logFileName string
	switch action.LogConfig.Format {
	case LogFormatText:
		logFileName = fmt.Sprintf("drive_janitor_logs_%s.txt", date)
	case LogFormatJSON:
		logFileName = fmt.Sprintf("drive_janitor_logs_%s.json", date)
	case LogFormatCSV:
		logFileName = fmt.Sprintf("drive_janitor_logs_%s.csv", date)
	}
	logFileName = filepath.Join(action.LogConfig.LogRepository, logFileName)
	logFileName, err := filepath.Abs(logFileName)
	if err != nil {
		return err
	}
	action.LogConfig.LogRepository = logFileName
	return nil
}

func (action *Action) SaveToFile() error {
	err := os.MkdirAll(filepath.Dir(action.LogConfig.LogRepository), 0755)
	if err != nil {
		return fmt.Errorf("error creating log directory: %w", err)
	}
	switch action.LogConfig.Format {
	case LogFormatText:
		// Write log to text file
		err := os.WriteFile(action.LogConfig.LogRepository, []byte(strings.Join(action.LogConfig.Files, "\n")), 0644)
		if err != nil {
			return fmt.Errorf("error writing log file: %w", err)
		}
	default:
		//TODO: Implement other formats
		return fmt.Errorf("unsupported log format: %s", action.LogConfig.Format)
	}
	return nil
}
