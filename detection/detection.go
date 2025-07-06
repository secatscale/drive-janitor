package detection

import (
	"drive-janitor/detection/checkage"
	"drive-janitor/detection/checktype"
	"drive-janitor/detection/checkyara"
	"path/filepath"
	"regexp"
)

// Look for a matching criteria on the file path for the current rules
func (detection *Detection) IsDetected(path string, detectionInfos *[]DetectionInfo) (bool, error) {
	// Call la function check type sur le path
	typeMatch, err := detection.FileTypeMatching(path)
	if err != nil {
		return false, err
	}

	// Call la function check age sur le path
	ageMatch, err := detection.FileAgeMatching(path)
	if err != nil {
		return false, err
	}
	filenameMatch, err := detection.FileNameMatching(path)
	if err != nil {
		return false, err
	}
	yaraMatch, err := detection.YaraMatching(path)
	if err != nil {
		return false, err
	}

	if typeMatch && ageMatch && filenameMatch && yaraMatch {
		detectionInfo := DetectionInfo{
			TypeMatch:     typeMatch,
			AgeMatch:      ageMatch,
			FilenameMatch: filenameMatch,
			Path:          path,
			YaraMatch:     yaraMatch,
			Detection:     detection,
		}

		*detectionInfos = append(*detectionInfos, detectionInfo)
	}

	return typeMatch && ageMatch && filenameMatch && yaraMatch, nil
}

func (detection Detection) FileTypeMatching(path string) (bool, error) {
	if detection.MimeType == "" {
		return true, nil
	}
	// Call la function check type sur le path
	fileType, err := checktype.GetType(path)
	if err != nil {
		return false, err
	}
	// print path and fileType
	matchingType := detection.MimeType
	if fileType == matchingType {
		return true, nil
	}
	return false, nil
}

func (detection Detection) FileAgeMatching(path string) (bool, error) {
	// Case where the age is not set
	if detection.Age == 0 || detection.Age == -1 {
		return true, nil
	}
	// Call la function check age sur le path
	age, err := checkage.GetAge(path)
	if err != nil {
		return false, err
	}
	machingAge := detection.Age
	if age < machingAge {
		return false, nil
	}
	// Call la function check age sur le path
	return true, nil
}

func (detection Detection) FileNameMatching(path string) (bool, error) {
	if detection.Filename == "" {
		return true, nil
	}
	filename := filepath.Base(path)
	match, err := regexp.MatchString(detection.Filename, filename)

	return match, err
}

func (detection Detection) YaraMatching(path string) (bool, error) {
	if detection.YaraRulesDir == "" {
		return true, nil
	}
	// Call la function check yara sur le path, pour les regles dans YaraRulesDir
	match, err := checkyara.CheckYara(path, detection.YaraRulesDir)
	if err != nil {
		return false, err
	}
	return match, nil
}
