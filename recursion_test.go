/*
	config recursion:
		* max depth
		* skip directories (regex)
		* priority directories (regex)
*/

/*
	// struct recursion with info path depth etc
	// struct de detection
	// struct d'action
	recursion(struct recursion, )

	on a la recursion qui loop, et en fonction de la config, on aura la detection, donc idealement la recursion prend une structure avec la config de detection
*/

package main

import (
	"drive-janitor/testhelper"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

type RecursionConfig struct {
	InitialPath string
	MaxDepth int
	SkipDirectories []string
	PriorityDirectories []string
	BrowseFiles int
}

func countSeparator(path string) int {
	count := strings.Split(path, string(os.PathSeparator))
	return len(count)
}

func getDepth(path string) int {
	if (path == ".") {
		return 0;
	}
	return countSeparator(path)
}

func isAboveMaxDepth(path string, maxDepth int) bool {
	return getDepth(path) > maxDepth
}

func (config *RecursionConfig) recurse(/* May take dectection and action struct*/) error {
	initialPathFs := os.DirFS(config.InitialPath);
	err := fs.WalkDir(initialPathFs, ".", func(path string, entry fs.DirEntry, err error) error {
		path = filepath.FromSlash(path)
		if (isAboveMaxDepth(path, config.MaxDepth)) {
			return fs.SkipDir;
		}
		//fmt.Println(getDepth(path), config.InitialPath, path, entry.Type().IsDir())
		if (!entry.Type().IsDir()) {
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

func	generateTestFS(layers int, filesInfo map[int][]string) {
	path, err := os.Getwd()
	if (err != nil) {
		panic(err)
	}
	for i := range layers {
		path = filepath.Join(path, strconv.Itoa(i))
		err = os.MkdirAll(path, 0755)
		if (err != nil) {
			panic(err)
		}
		filename, layerExist := filesInfo[i]
		if (layerExist) {
			for _, file := range filename {
				_, err = os.Create(filepath.Join(path, file))
				if (err != nil) {
					panic(err)
				}
			}
		}
	}
}

// find a way to generate test files easily
func TestRecursion(t *testing.T) {
	generateTestFS(3, map[int][]string{
		0: {"file.txt"},
		1: {"file1.txt"},
		2: {"file2.txt"},
		3: {"file3.txt"},
	})
	dir := "0"
	path, err :=  os.Getwd()
	if (err != nil) {
		t.Fatalf("Error getting current directory: %v", err)
	}
	t.Run("Test Browsering", func(t *testing.T) {
		testhelper.RunOSDependentTest(t, "Test Browsering", func(t *testing.T) {
			config := RecursionConfig{
				InitialPath: filepath.Join(path, dir),
				MaxDepth: 2,
				SkipDirectories: []string{},
				PriorityDirectories: []string{},
				BrowseFiles: 0,
			}
			err = config.recurse()
			if (err != nil) {
				t.Errorf("Error while browsing files: %v", err)
			}
			if (config.BrowseFiles != 2) {
				t.Errorf("Expected 2files, got %d", config.BrowseFiles)
			}
		}, map[string]bool{"linux": true, "darwin": true, "windows": true})
	});

	t.Run("Test max depth", func(t *testing.T) {
		testhelper.RunOSDependentTest(t, "Test max depth", func(t *testing.T) {
			config := RecursionConfig{
				InitialPath: filepath.Join(path, dir),
				MaxDepth: 2,
				SkipDirectories: []string{},
				PriorityDirectories: []string{},
				BrowseFiles: 0,
			}
			err = config.recurse()
			if (err != nil) {
				t.Errorf("Error while browsing files: %v", err)
			}
			if (config.BrowseFiles != 2) {
				t.Errorf("Expected 2files, got %d", config.BrowseFiles)
			}
		}, map[string]bool{"linux": true, "darwin": true, "windows": true})
	})


	t.Cleanup(func() {
		defer os.RemoveAll(filepath.Join(path, dir))
	})
}