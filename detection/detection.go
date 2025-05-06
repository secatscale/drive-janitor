package detection

import (
	"drive-janitor/detection/checkage"
	"drive-janitor/detection/checktype"
	"path/filepath"
	"regexp"
)

func (detection Detection) IsDetected(path string) (bool, error) {
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
	return typeMatch && ageMatch && filenameMatch, nil
}

func (detection Detection) FileTypeMatching(path string) (bool, error) {
	if detection.MimeType == "" {
		return true, nil
	}
	// Call la function check type sur le path
	fileType, err := checktype.CheckType(path)
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
	age, err := checkage.CheckAge(path)
	if err != nil {
		return false, err
	}
	machingAge := detection.Age
	if age > machingAge {
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
