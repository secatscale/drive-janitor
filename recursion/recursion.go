package recursion

import (
	"drive-janitor/action"
	"drive-janitor/detection"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
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

func isInSkipDirectories(path string, skipDirectories []string) bool {
	//fmt.Println("Path: ", path, "Skip: ", skipDirectories)
	return slices.Contains(skipDirectories, path)
}

// Main recursion loop, from the initial path, we will recurse through all directories and files
func (config *Recursion) Recurse(detectionsArray detection.DetectionArray, action *action.Action) error {
	initialPathFs := os.DirFS(config.InitialPath)
	err := fs.WalkDir(initialPathFs, ".", func(path string, entry fs.DirEntry, err error) error {
		path = filepath.FromSlash(path)

		if err != nil {
			if os.IsPermission(err) {
				// skip si on a pas les permissions
				return nil
			}
			return err
		}
		badDirs := []string{"/dev/fd", "/proc", "fd/"}

		for _, bad := range badDirs {
			if strings.HasPrefix(path, bad) {
				return fs.SkipDir
			}
		}

		if isAboveMaxDepth(path, config.MaxDepth) || isInSkipDirectories(path, config.SkipDirectories) {
			return fs.SkipDir
		}
		if entry.Type().IsRegular() {
			absolutePath := filepath.Join(config.InitialPath, path)
			// Check if we need action on the current file
			detectionsMatch, needAction, err := detectionsArray.AsMatch(absolutePath)
			if err != nil {
				return err
			}
			if needAction {
				//				fmt.Println("File detected: ", absolutePath)
				// call the action
				action.TakeAction(absolutePath, detectionsMatch)
			}
			// Maybe useless TODO remove this and fix the tests
			config.BrowseFiles += 1
		}
		return nil
	})
	return err
}
