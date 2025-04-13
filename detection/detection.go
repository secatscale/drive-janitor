package detection

import (
	"drive-janitor/detection/checkage"
	"drive-janitor/detection/checktype"
)

func (detection DetectionConfig) FileTypeMatching(path string) (bool, error) {
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

func (detection DetectionConfig) FileAgeMatching(path string) (bool, error) {
	// Case where the age is not set
	if detection.Age < 0 {
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
