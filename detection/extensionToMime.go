package detection

import (
	"errors"
	"mime"
)

// Predefined map of extensions to MIME types for better coverage.
var customMimeTypes = map[string]string{
	".md":   "text/markdown",
	".csv":  "text/csv",
	".json": "application/json",
	".yaml": "application/x-yaml",
	".yml":  "application/x-yaml",
	".log":  "text/plain",
	".ini":  "text/plain",
	".toml": "application/toml",
	".svg":  "image/svg+xml",
	".webp": "image/webp",
	".rar":  "application/vnd.rar",
	".7z":   "application/x-7z-compressed",
	".tar":  "application/x-tar",
	".gz":   "application/gzip",
	".xz":   "application/x-xz",
	".mp3":  "audio/mpeg",
	".wav":  "audio/wav",
	".flac": "audio/flac",
	".ogg":  "audio/ogg",
	".mov":  "video/quicktime",
	".avi":  "video/x-msvideo",
	".wmv":  "video/x-ms-wmv",
	".flv":  "video/x-flv",
	".go":   "text/x-go",
	".php":  "application/x-httpd-php",
	".py":   "text/x-python",
	".java": "text/x-java-source",
	".rb":   "application/x-ruby",
	".c":    "text/x-c",
	".cpp":  "text/x-c++",
	".cs":   "text/plain",
	".swift": "text/x-swift",
	".rs":   "text/x-rust",
	".kt":   "text/x-kotlin",
	".sh":   "application/x-sh",
	".bat":  "application/x-msdos-program",
	".ps1":  "application/x-powershell",
	".sql":  "application/sql",
	".ttf":  "font/ttf",
	".woff": "font/woff",
	".woff2": "font/woff2",
}

// GetMimeType returns the MIME type for a given file extension.
func GetMimeType(ext string) (string, error) {
	if ext == "" || ext[0] != '.' {
		ext = "." + ext // Ensure the extension starts with a dot
	}

	// Check custom MIME types first
	if mimeType, exists := customMimeTypes[ext]; exists {
		return mimeType, nil
	}

	// Fallback to Go's built-in detection
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		return "unknown/unknown", errors.New("unknown MIME type")
	}
	return mimeType, nil
}
