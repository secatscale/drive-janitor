package main

import (
	"io/ioutil"

	"github.com/h2non/filetype"
)

func CheckType(filePath string) (fileType string, err error) {
	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	if len(buf) == 0 {
		return "empty", nil
	}
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
