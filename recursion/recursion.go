package recursion

import (
	"drive-janitor/action"
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

func (config *Recursion) Recurse(detection detection.DetectionArray, action *action.Action) error {
	initialPathFs := os.DirFS(config.InitialPath)
	err := fs.WalkDir(initialPathFs, ".", func(path string, entry fs.DirEntry, err error) error {
		path = filepath.FromSlash(path)
	//	fmt.Println(path, err, entry, entry.Type().IsRegular(), isAboveMaxDepth(path, config.MaxDepth))

		if err != nil {
			if os.IsPermission(err) {
				// skip si on a pas les permissions
				return nil
			}
			return err
		}
		if strings.HasPrefix(path, "proc") {
			// On ignore le dossier proc
			return fs.SkipDir
		}

		//fmt.Println(isAboveMaxDepth(path, config.MaxDepth), config.InitialPath, path, entry.Type().IsDir())
		if isAboveMaxDepth(path, config.MaxDepth) {
			return fs.SkipDir
		}
		if entry.Type().IsRegular() {
			// We should check if the file should be detected or not
			// If it is, then we do the action

			absolutePath := filepath.Join(config.InitialPath, path)
			needAction, err := detection.AsMatch(absolutePath)
			if err != nil {
				return err
			}
			if needAction {
				fmt.Println("Detected file: ", path)
				// call the action
				action.TakeAction(absolutePath)
	//		}
			}
			config.BrowseFiles += 1
		}
		return nil
	})
	return err
}
