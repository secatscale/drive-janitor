package detection

import (
	"io"
	"os"
	"time"

	"github.com/h2non/filetype"
)

func (detection DetectionConfig) FileTypeMatching(path string) (bool, error) {
	// Call la function check type sur le path
	fileType, err := CheckType(path)
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
	age, err := CheckAge(path)
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

func CheckType(filePath string) (fileType string, err error) {

	fd, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer fd.Close()

	buf := make([]byte, 32)
	n, err := fd.Read(buf)
	if err == io.EOF {
		return "", nil
	}
	if err != nil {
		// Attention
		return "", nil
	}
	if len(buf) == 0 {
		return "empty", nil
	}
	buf = buf[:n]
	kind, err := filetype.Match(buf)
	if err != nil {
		return "", err
	}
	if kind == filetype.Unknown {
		ProbablyText := isProbablyText(buf)
		if ProbablyText {
			return "text", nil
		}
		return "unknown", nil
	}
	return kind.MIME.Value, nil
}

func isProbablyText(buf []byte) bool {
	for _, b := range buf {
		if b == 0 {
			return false
		}
		if (b < 32 && b != 9 && b != 10 && b != 13) || b > 126 {
			return false
		}
	}
	return true
}

// CheckAge is a function that takes a file path and returns the age of the file in days.
func CheckAge(filePath string) (int, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	age := time.Since(info.ModTime()).Hours() / 24
	// fmt.Println(age)
	return int(age), nil
}
