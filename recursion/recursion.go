package recursion

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func countSeparator(path string) int {
	if (!strings.ContainsRune(path, os.PathSeparator)) {
		return 0;
	}
	return strings.Count(path, string(os.PathSeparator))
}

func getDepth(path string) int {
	return countSeparator(path)
}

func isAboveMaxDepth(path string, maxDepth int) bool {
	if (maxDepth < 0) {
		return false
	}
	return getDepth(path) > maxDepth
}

func (config *RecursionConfig) Recurse(/* May take dectection and action struct*/) error {
	initialPathFs := os.DirFS(config.InitialPath);
	err := fs.WalkDir(initialPathFs, ".", func(path string, entry fs.DirEntry, err error) error {
		path = filepath.FromSlash(path)
		//fmt.Println(getDepth(path), config.InitialPath, path, entry.Type().IsDir())
		if (isAboveMaxDepth(path, config.MaxDepth)) {
			return fs.SkipDir;
		}
		if (entry.Type().IsRegular()) {
			config.BrowseFiles += 1;
		}
		if (err != nil) {
			return err
		}
		return nil
	})
	if (err != nil) {
		return err
	}
	return nil
}