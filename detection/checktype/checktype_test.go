package checktype

import (
	"path/filepath"
	"testing"
)

func TestCheckType(t *testing.T) {
	// Test cases
	var tests = []struct {
		testName string
		fileName string
		fileType string
	}{
		// images
		{"Test jpg", "sample.jpg", "image/jpeg"},
		{"Test png", "sample.png", "image/png"},
		{"Test gif", "sample.gif", "image/gif"},
		{"Test webp", "sample.webp", "image/webp"},

		// videos
		{"Test webm", "sample.webm", "video/webm"},
		{"Test mp4", "sample.mp4", "video/mp4"},
		{"Test mov", "sample.mov", "video/quicktime"},
		{"Test flv", "sample.flv", "video/x-flv"},

		// audio
		{"Test mp3", "sample.mp3", "audio/mpeg"},
		{"Test wav", "sample.wav", "audio/x-wav"},
		{"Test flac", "sample.flac", "audio/x-flac"},
		{"Test aac", "sample.aac", "audio/aac"},

		// documents
		{"Text xlsx", "sample.xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"},
		// {"Test xlsm", "sample.xlsm", "application/vnd.ms-excel.sheet.macroEnabled.12"},
		// {"Test xlsb", "sample.xlsb", "application/vnd.ms-excel.sheet.binary.macroEnabled.12"},
		// {"Test xltm", "sample.xltm", "application/vnd.ms-excel.template.macroEnabled.12"},
		// {"Test xls", "sample.xls", "application/vnd.ms-excel"},
		// {"Test xlt", "sample.xlt", "application/vnd.ms-excel"},
		// {"Test xml", "sample.xml", "application/xml"},
		{"Test doc", "sample.doc", "application/msword"},
		// {"Test dot", "sample.dot", "application/msword"},
		// {"Test docm", "sample.docm", "application/vnd.ms-word.document.macroEnabled.12"},
		// {"Test dotm", "sample.dotm", "application/vnd.ms-word.template.macroEnabled.12"},
		{"Test docx", "sample.docx", "application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
		{"Test pptx", "sample.pptx", "application/vnd.openxmlformats-officedocument.presentationml.presentation"},
		// {"Test pptm", "sample.pptm", "application/vnd.ms-powerpoint.presentation.macroEnabled.12"},
		// {"Test ppt", "sample.ppt", "application/vnd.ms-powerpoint"},
		// {"Test pot", "sample.pot", "application/vnd.ms-powerpoint"},
		// {"Test potm", "sample.potm", "application/vnd.ms-powerpoint.template.macroEnabled.12"},
		// {"Test potx", "sample.potx", "application/vnd.openxmlformats-officedocument.presentationml.template"},
		// {"Test pps", "sample.pps", "application/vnd.ms-powerpoint"},
		// {"Test ppsm", "sample.ppsm", "application/vnd.ms-powerpoint.slideshow.macroEnabled.12"},
		// {"Test ppsx", "sample.ppsx", "application/vnd.openxmlformats-officedocument.presentationml.slideshow"},
		// {"Test rtf", "sample.rtf", "application/rtf"},

		// archive / application
		{"Test pdf", "sample.pdf", "application/pdf"},
		{"Test zip", "sample.zip", "application/zip"},
		{"Test tar", "sample.tar", "application/x-tar"},
		// {"Test tgz", "sample.tgz", "application/x-gzip"},
		// {"Test gz", "sample.gz", "application/gzip"},
		// {"Test bz2", "sample.bz2", "application/x-bzip2"},
		{"Test 7z", "sample.7z", "application/x-7z-compressed"},
		{"Test rar", "sample.rar", "application/vnd.rar"},
		// {"Test exe", "sample.exe", "application/x-msdownload"},
		// {"Test msi", "sample.msi", "application/x-msdownload"},
		// {"Test dmg", "sample.dmg", "application/x-apple-diskimage"},
		// {"Test iso", "sample.iso", "application/x-iso9660-image"},
		// {"Test deb", "sample.deb", "application/x-debian-package"},
		// {"Test rpm", "sample.rpm", "application/x-rpm"},
		// {"Test apk", "sample.apk", "application/vnd.android.package-archive"},
		// {"Test jar", "sample.jar", "application/java-archive"},
		// {"Test war", "sample.war", "application/java-archive"},
		// {"Test ear", "sample.ear", "application/java-archive"},
		// {"Test html", "sample.html", "text/html"},

		// text
		{"Test txt", "sample.txt", "text"},
		{"Test md", "sample.md", "text"},
		{"Test js", "sample.js", "text"},
		{"Test json", "sample.json", "text"},
		{"Test xml", "sample.xml", "text"},
		{"Test csv", "sample.csv", "text"},
		{"Test html", "sample.html", "text"},
		{"Test css", "sample.css", "text"},

		// empty
		{"Test empty", "empty", "empty"},
	}

	// Test loop
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			filePath, err := filepath.Abs(filepath.Join("../../samples", test.fileName))
			if err != nil {
				t.Errorf("Error getting file path: %s", err)
			}
			fileType, err := CheckType(filePath)
			if err != nil {
				t.Errorf("Error checking file type: %s", err)
			}
			if fileType != test.fileType {
				t.Errorf("Expected %s, got %s", test.fileType, fileType)
			}
		})
	}
}
