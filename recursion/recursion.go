package recursion

import (
	"drive-janitor/detection"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func countSeparator(path string) int {
	if !strings.ContainsRune(path, os.PathSeparator) {
		return 0
	}
	return strings.Count(path, string(os.PathSeparator))
}

func getDepth(path string) int {
	return countSeparator(path)
}

func isAboveMaxDepth(path string, maxDepth int) bool {
	if maxDepth < 0 {
		return false
	}
	return getDepth(path) > maxDepth
}

func isDetected(path string, detectionConfig detection.DetectionConfig) (bool, error) {
	// Call la function check type sur le path
	typeMatch, err := detectionConfig.FileTypeMatching(path)
	if err != nil {
		return false, err
	}
	// Call la function check age sur le path
	ageMatch, err := detectionConfig.FileAgeMatching(path)
	if err != nil {
		return false, err
	}
	return typeMatch && ageMatch, nil
}

func (config *RecursionConfig) Recurse(detectionConfig detection.DetectionConfig /* May take dectection and action struct*/) error {
	initialPathFs := os.DirFS(config.InitialPath)
	err := fs.WalkDir(initialPathFs, ".", func(path string, entry fs.DirEntry, err error) error {
		path = filepath.FromSlash(path)
		//fmt.Println(isAboveMaxDepth(path, config.MaxDepth), config.InitialPath, path, entry.Type().IsDir())
		if isAboveMaxDepth(path, config.MaxDepth) {
			return fs.SkipDir
		}
		if entry.Type().IsRegular() {
			// We should check if the file should be detected or not
			// If it is, then we do the action
			needAction, err := isDetected(filepath.Join(config.InitialPath, path), detectionConfig)
			if err != nil {
				return err
			}
			if needAction {
				fmt.Println("Detected file: ", path)
				// call the action
			}
			config.BrowseFiles += 1
		}
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
