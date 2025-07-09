package detection

import (
	"fmt"

	"github.com/h2non/filetype"
)

// Predefined map of extensions to MIME types for better coverage.
var SupportedMimeTypes = map[string]bool{
	"text/markdown": true,
	"text/csv":      true,
	//"application/json",
	//"application/x-yaml",
	//"application/x-yaml",
	//"text/plain",
	//"text/plain",
	//"application/toml",
	//"image/svg+xml",
	//"image/webp",
	//"application/vnd.rar",
	//"application/x-7z-compressed",
	//"application/x-tar",
	//"application/gzip",
	//"application/x-xz",
	//"audio/mpeg",
	//"audio/wav",
	//"audio/flac",
	//"audio/ogg",
	//"video/quicktime",
	//"video/x-msvideo",
	//"video/x-ms-wmv",
	//"video/x-flv",
	//"application/x-httpd-php",
	//"text/x-python",
	//"text/x-java-source",
	//"application/x-ruby",
	//"text/x-c",
	//"text/x-c++",
	//"text/plain",
	//"text/x-swift",
	//"text/x-rust",
	//"text/x-kotlin",
	//"application/x-sh",
	//"application/x-msdos-program",
	//"application/x-powershell",
	//"application/sql",
	//"font/ttf",
	//"font/woff",
	//"font/woff2",
	//"text/x-go",
}

// Check if the provided file extension is supported
func SupportType(mimeType string) (bool, error) {
	// Check custom MIME types first
	if filetype.IsMIMESupported(mimeType) {
		fmt.Println("MIME type is supported:", mimeType)
		return true, nil
	} else if mimeTypeIsSupported, exists := SupportedMimeTypes[mimeType]; exists {
		return mimeTypeIsSupported, nil
	}
	return false, nil
}
