package checktype

import (
	"bufio"
	"io"
	"os"

	"github.com/h2non/filetype"
)

func CheckType(filePath string) (fileType string, err error) {
	//Stat rapide pour eviter les fichiers spéciaux
	info, err := os.Lstat(filePath)
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() {
		// Skip sockets, symlinks, devices, etc.
		return "", nil
	}

	fd, err := os.Open(filePath)
	// if permission denied, return empty string
	if os.IsPermission(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	defer fd.Close()

	reader := bufio.NewReader(fd)
	buf := make([]byte, 24)
	n, err := reader.Read(buf)
	// if file is empty, return empty string
	if err == io.EOF || n == 0 {
		return "empty", nil
	}
	if err != nil {
		return "", err
	}
	buf = buf[:n]
	kind, err := filetype.Match(buf)
	if err != nil {
		return "", err
	}
	if kind == filetype.Unknown {
		if isProbablyText(buf) {
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
		if (b < 32 && b != 9 && b != 10 && b != 13) || (b > 126 && b < 160) {
			return false
		}
	}
	return true
}
