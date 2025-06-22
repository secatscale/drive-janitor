package action

import (
	"drive-janitor/action/log"
	"drive-janitor/detection"
	"drive-janitor/detection/checkage"
	"drive-janitor/detection/checktype"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func (action *Action) TakeAction(filePath string, detectionsMatch []string) {
	if action.Delete {
		// TODO: Implement delete action
		// os.Remove(filePath)
	}

	detectedBy := buildDetectedByString(detectionsMatch)

	if action.Log {
		FileInfo := FileInfo{"detectedBy": detectedBy, "path": filePath}
		action.LogConfig.FilesInfo = append(action.LogConfig.FilesInfo, FileInfo)
	}
}

func buildDetectedByString(detectionsMatch []string) string {
	var detectedBy string
	if len(detectionsMatch) > 0 {
		for _, detection := range detectionsMatch {
			if detectedBy != "" {
				detectedBy += " and "
			}
			detectedBy += detection
		}
	} else {
		detectedBy = "unknown"
	}
	return detectedBy
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

func (action *Action) SaveToFile(detectionInfo []detection.DetectionInfo) error {
	err := os.MkdirAll(filepath.Dir(action.LogConfig.LogRepository), 0755)
	if err != nil {
		return fmt.Errorf("error creating log directory: %w", err)
	}
	// Enrich logs with additional information
	action.EnrichLogs(detectionInfo)
	var logContent string
	switch action.LogConfig.Format {
	case LogFormatText:
		logContent = GenerateTXT(action.LogConfig.FilesInfo)
	case LogFormatJSON:
		logContent = GenerateJSON(action.LogConfig.FilesInfo)
	case LogFormatCSV:
		logContent = GenerateCSV(action.LogConfig.FilesInfo)
	default:
		return fmt.Errorf("unsupported log format: %s", action.LogConfig.Format)
	}
	err = log.SaveToFile(logContent, action.LogConfig.LogRepository)
	if err != nil {
		return fmt.Errorf("error saving log file: %w", err)
	}
	return nil
}

func GenerateTXT(FilesInfo []FileInfo) string {
	var logContent string
	for _, fileInfo := range FilesInfo {
		logContent += fmt.Sprintf("Detected By: %s\nPath: %s\nFile Type: %s\nFile Age: %s\n", fileInfo["detectedBy"], fileInfo["path"], fileInfo["file_type"], fileInfo["file_age"])
		logContent += "------------------------\n"
	}
	return logContent
}

func GenerateJSON(FilesInfo []FileInfo) string {
	var logContent string
	// logContent += "{"
	logContent += "["
	for _, fileInfo := range FilesInfo {
		logContent += fmt.Sprintf("{\"detected by\":\"%s\", \"path\": \"%s\", \"file_type\": \"%s\", \"file_age\": \"%s\"},", fileInfo["detectedBy"], fileInfo["path"], fileInfo["file_type"], fileInfo["file_age"])
	}
	// Remove the last comma and close the JSON object
	if len(logContent) > 1 {
		logContent = logContent[:len(logContent)-1]
	}
	logContent += "]"
	// logContent += "}"
	return logContent
}

func GenerateCSV(FilesInfo []FileInfo) string {
	var logContent string
	logContent += "detectedby,path,file_type,file_age\n"
	for _, fileInfo := range FilesInfo {
		logContent += fmt.Sprintf("%s,%s,%s,%s\n", fileInfo["detectedBy"], fileInfo["path"], fileInfo["file_type"], fileInfo["file_age"])
	}
	return logContent
}

// Enrich fileInfo with timestamp, file type, and file age in days
func (action *Action) EnrichLogs(detectionInfo []detection.DetectionInfo) {
	for _, detection := range detectionInfo {
		print("Detection: ", detection.Path, "\n")
		if (detection.TypeMatch) {
			print("We matched a file: ", detection.Path, "\n")
			print("On is type: ", detection.Detection.MimeType, "\n")
		}
	}
	for i, fileInfo := range action.LogConfig.FilesInfo {

		fileType, err := checktype.GetType(fileInfo["path"])
		if err != nil {
			fileType = "unknown"
		}

		fileInfo["file_type"] = string(fileType)
		fileAge, err := checkage.GetAge(fileInfo["path"])
		if err != nil {
			fileAge = -1 // Indicate an error in age calculation
		}
		fileInfo["file_age"] = fmt.Sprintf("%d days", fileAge)

		action.LogConfig.FilesInfo[i] = fileInfo
	}
}
