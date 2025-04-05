package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf16"
)

// TestWhereTrash tests the WhereTrash function that returns the path to the trash folder for an OS
func TestWhereTrash(t *testing.T) {
	t.Run("HomeDir is not null", func(t *testing.T) {
		home, err := os.UserHomeDir()
		if err != nil {
			t.Fatalf("os.UserHomeDir() returned an error: %v", err)
		}
		if home == "" {
			t.Fatalf("os.UserHomeDir() returned null")
		}
	})
	home, _ := os.UserHomeDir()
	windowsSID, _ := GetCurrentUserSID()
	trashTests := []struct {
		testName  string
		osName    string
		trashPath string
	}{
		{testName: "Windows", osName: "windows", trashPath: filepath.Join(os.Getenv("SystemDrive")+"\\", "$Recycle.Bin", windowsSID)},
		{testName: "Linux", osName: "linux", trashPath: filepath.Join(home, ".local", "share", "Trash")},
		{testName: "Darwin", osName: "darwin", trashPath: filepath.Join(home, ".Trash")},
	}
	for _, tt := range trashTests {
		t.Run(tt.testName, func(t *testing.T) {
			if WhichOs() != tt.osName {
				t.Skip("Skipping test for OS: " + tt.osName)
			}
			trashPath, err := WhereTrash(tt.osName)
			if err != nil {
				t.Errorf("WhereTrash() returned an error: %v", err)
			}
			if trashPath == "" {
				t.Errorf("WhereTrash() returned null")
			}
			if trashPath != tt.trashPath {
				t.Errorf("WhereTrash() = %v, want %v", trashPath, tt.trashPath)
			}
		})
	}
	t.Run("Unsupported OS", func(t *testing.T) {
		osName := "SkibidiBopBopDopaMina"
		trashPath, err := WhereTrash(osName)
		if err == nil {
			t.Errorf("WhereTrash() did not return an error for unsupported OS")
		}
		if trashPath != "" {
			t.Errorf("WhereTrash() returned a path for unsupported OS")
		}
	})
	t.Run("Is this actually where a trashed file go ?", func(t *testing.T) {
		trashPath, err := WhereTrash(WhichOs())
		if err != nil {
			t.Fatalf("WhereTrash() returned an error: %v", err)
		}
		isDeleteMeHere, err := isDeleteMeHere()
		if err != nil {
			t.Fatalf("isDeleteMeHere returned an error: %v", err)
		}
		if !isDeleteMeHere {
			err = respawnSample()
			if err != nil {
				t.Fatalf("respawnSample() returned an error: %v", err)
			}
		}
		isDeleteMeInTrash, err := isTrashedSampleInTrash(trashPath)
		if err != nil {
			t.Fatalf("isTrashedSampleInTrash() returned an error: %v", err)
		}
		if isDeleteMeInTrash {
			err = cleanTrashedSample()
			if err != nil {
				t.Fatalf("cleanTrashedSample() returned an error: %v", err)
			}
		}
		err = trashSample()
		if err != nil {
			t.Fatalf("trashSample() returned an error: %v", err)
		}

		isTrashedSampleInTrash, err := isTrashedSampleInTrash(trashPath)
		if err != nil {
			t.Fatalf("isTrashedSampleInTrash() returned an error: %v", err)
		}
		if !isTrashedSampleInTrash {
			t.Fatalf("Trashed sample file not found in trash folder")
		}
		t.Cleanup(func() {
			err = cleanTrashedSample()
			if err != nil {
				t.Fatalf("cleanTrashedSample() returned an error: %v", err)
			}
			err = respawnSample()
			if err != nil {
				t.Fatalf("respawnSample() returned an error: %v", err)
			}
		})
	})
}

// isDeleteMeHere checks if the sample file is in the current directory
func isDeleteMeHere() (bool, error) {
	_, err := os.Stat("delete_me.txt")
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("os.Stat() returned an error: %v", err)
	}
	return true, nil
}

// trashSample moves a sample file to the trash folder
func trashSample() error {
	absPath, err := filepath.Abs("delete_me.txt")
	if err != nil {
		return fmt.Errorf("filepath.Abs() returned an error: %v", err)
	}
	var cmd *exec.Cmd
	switch WhichOs() {
	case "windows":
		psScript := fmt.Sprintf(`
		Add-Type -AssemblyName Microsoft.VisualBasic
		[Microsoft.VisualBasic.FileIO.FileSystem]::DeleteFile("%s", 'OnlyErrorDialogs', 'SendToRecycleBin')`, absPath)
		fmt.Println("Powershell script:", psScript)
		cmd = exec.Command("powershell", "-Command", psScript)
	case "linux":
		//verify gio is installed
		_, err := exec.LookPath("gio")
		if err != nil {
			return fmt.Errorf("gio not found in PATH")
		}
		cmd = exec.Command("gio", "trash", absPath)
	case "darwin":
		//verify osascript is installed
		_, err := exec.LookPath("osascript")
		if err != nil {
			return fmt.Errorf("osascript not found in PATH")
		}
		cmd = exec.Command("osascript", "-e", fmt.Sprintf(`tell application "Finder" to delete POSIX file "%s"`, absPath))
	default:
		return fmt.Errorf("Unsupported OS %s", WhichOs())
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cmd.CombinedOutput() returned an error: %v, output: %s", err, string(out))
	}
	return nil
}

func decodeIFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	fmt.Println("data", data)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	if len(data) < 20 {
		return "", fmt.Errorf("invalid $I file: too short (%d bytes)", len(data))
	}

	// Longueur de la chaîne en UTF-16
	strLen := binary.LittleEndian.Uint32(data[16:20])
	if strLen == 0 || strLen > 4096 {
		return "", fmt.Errorf("invalid $I file: suspicious string length %d", strLen)
	}

	expectedSize := 20 + int(strLen)*2
	if expectedSize > len(data) {
		return "", fmt.Errorf("invalid $I file: declared string length (%d bytes) exceeds file size (%d bytes)", expectedSize, len(data))
	}

	rawUtf16 := data[20:expectedSize]

	u16 := make([]uint16, strLen)
	for i := 0; i < int(strLen); i++ {
		u16[i] = binary.LittleEndian.Uint16(rawUtf16[i*2 : i*2+2])
	}

	decoded := string(utf16.Decode(u16))
	return strings.TrimRight(decoded, "\x00"), nil
}

func GetWindowsTrashedFilePaths(trashPath string, originalFileName string) (filePath, metaPath string, err error) {
	entries, err := os.ReadDir(trashPath)
	if err != nil {
		return "", "", fmt.Errorf("os.ReadDir(%s) returned an error: %v", trashPath, err)
	}
	originalFileName = strings.ToLower(originalFileName)

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, "$I") {
			metaFile := filepath.Join(trashPath, name)
			originalPath, err := decodeIFile(metaFile)
			if err != nil {
				continue
			}
			if strings.HasSuffix(strings.ToLower(originalPath), originalFileName) {
				realFile := filepath.Join(trashPath, strings.Replace(name, "$I", "$R", 1))
				return realFile, metaFile, nil
			}

		}

	}
	return "", "", nil
}

// isTrashedSampleInTrash checks if the sample file is in the trash folder
func isTrashedSampleInTrash(trashPath string) (bool, error) {
	var err error
	switch WhichOs() {
	case "windows":
		filePath, _, err := GetWindowsTrashedFilePaths(trashPath, "delete_me.txt")
		if err != nil {
			return false, fmt.Errorf("GetWindowsTrashedSamplePaths() returned an error: %v", err)
		}
		if filePath == "" {
			return false, nil
		}
		_, err = os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				return false, nil
			}
			return false, fmt.Errorf("os.Stat() returned an error: %v", err)
		}
		return true, nil
	case "linux":
		_, err = os.Stat(filepath.Join(trashPath, "files", "delete_me.txt"))
	case "darwin":
		_, err = os.Stat(filepath.Join(trashPath, "delete_me.txt"))
	default:
		return false, fmt.Errorf("Unsupported OS %s", WhichOs())
	}
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("os.Stat() returned an error: %v", err)
	}
	return true, nil
}

func getTrashedSamplePaths() (filePath, metaPath string, err error) {
	trashPath, err := WhereTrash(WhichOs())
	if err != nil {
		return "", "", fmt.Errorf("WhereTrash() returned an error: %v", err)
	}
	switch WhichOs() {
	case "windows":
		filePath, metaPath, err = GetWindowsTrashedFilePaths(trashPath, "delete_me.txt")
		if err != nil {
			return "", "", fmt.Errorf("GetWindowsTrashedSamplePaths() returned an error: %v", err)
		}
		return filePath, metaPath, nil
	case "darwin":
		filePath = filepath.Join(trashPath, "delete_me.txt")
		return filePath, "", nil
	case "linux":
		filePath = filepath.Join(trashPath, "files", "delete_me.txt")
		metaPath = filepath.Join(trashPath, "info", "delete_me.txt.trashinfo")
		return filePath, metaPath, nil
	default:
		return "", "", fmt.Errorf("Unsupported OS %s", WhichOs())
	}
}

// cleanTrashedSample removes the sample file from the trash folder
func cleanTrashedSample() error {
	filePath, metaPath, err := getTrashedSamplePaths()
	if err != nil {
		return fmt.Errorf("getTrashedSamplePaths() returned an error: %v", err)
	}
	switch WhichOs() {
	case "windows":
		err = os.Remove(filePath)
		if err != nil {
			return fmt.Errorf("os.Remove() returned an error: %v", err)
		}
		if metaPath != "" {
			err = os.Remove(metaPath)
			if err != nil {
				return fmt.Errorf("os.Remove() returned an error: %v", err)
			}
		}
	case "linux":
		err = os.Remove(filePath)
		if err != nil {
			return fmt.Errorf("os.Remove() returned an error: %v", err)
		}
		err = os.Remove(metaPath)
		if err != nil {
			return fmt.Errorf("os.Remove() returned an error: %v", err)
		}
	case "darwin":
		cmd := exec.Command("osascript", "-e", fmt.Sprintf(`tell application "Finder" to delete POSIX file "%s"`, filePath))
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("cmd.CombinedOutput() returned an error: %v, output: %s", err, string(out))
		}
	default:
		return fmt.Errorf("Unsupported OS %s", WhichOs())
	}
	return nil
}

// respawnSample creates the sample file in the current directory
func respawnSample() error {
	file, err := os.Create("delete_me.txt")
	if err != nil {
		return fmt.Errorf("os.Create() returned an error: %v", err)
	}
	_, err = file.WriteString("This file is meant to be deleted during testing of the WhereTrash function\n")
	if err != nil {
		return fmt.Errorf("file.WriteString() returned an error: %v", err)
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("file.Close() returned an error: %v", err)
	}
	return nil
}
