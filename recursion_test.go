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
	if (!strings.ContainsRune(path, os.PathSeparator)) {
		return 0;
	}
	return strings.Count(path, string(os.PathSeparator))
}

func getDepth(path string) int {
	return countSeparator(path)
}

func isAboveMaxDepth(path string, maxDepth int) bool {
	return getDepth(path) > maxDepth
}

func (config *RecursionConfig) recurse(/* May take dectection and action struct*/) error {
	initialPathFs := os.DirFS(config.InitialPath);
	err := fs.WalkDir(initialPathFs, ".", func(path string, entry fs.DirEntry, err error) error {
		path = filepath.FromSlash(path)
	//	fmt.Println(getDepth(path), config.InitialPath, path, entry.Type().IsDir())
		if (isAboveMaxDepth(path, config.MaxDepth)) {
			return fs.SkipDir;
		}
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
func TestRecursionComplex(t *testing.T) {
	// Create a complex file structure with multiple files and deeper nesting
	generateTestFS(9, map[int][]string{
		0: {"root1.txt", "root2.txt", "root3.txt"},
		1: {"level1_1.txt", "level1_2.txt", "level1_3.txt"},
		2: {"level2_1.txt", "level2_2.txt", "level2_3.txt", "level2_4.txt"},
		3: {"level3_1.txt", "level3_2.txt"},
		4: {"level4_1.txt", "level4_2.txt", "level4_3.txt"},
		5: {"level5_1.txt", "level5_2.txt"},
		6: {"level6_1.txt", "level6_2.txt", "level6_3.txt", "level6_4.txt"},
		7: {"level7_1.txt"},
		8: {"level8_1.txt", "level8_2.txt"},
	})
	dir := "0"
	path, err := os.Getwd()
	if (err != nil) {
		t.Fatalf("Error getting current directory: %v", err)
	}
	t.Run("Test with max depth 3", func(t *testing.T) {
		testhelper.RunOSDependentTest(t, "Test with max depth 3", func(t *testing.T) {
			config := RecursionConfig{
				InitialPath: filepath.Join(path, dir),
				MaxDepth: 3,
				SkipDirectories: []string{},
				PriorityDirectories: []string{},
				BrowseFiles: 0,
			}
			err = config.recurse()
			if (err != nil) {
				t.Errorf("Error while browsing files: %v", err)
			}
			// 3 files in root + 3 files in level 1 + 4 files in level 2 + 2 files in level 3
			expectedFiles := 12
			if (config.BrowseFiles != expectedFiles) {
				t.Errorf("Expected %d files, got %d", expectedFiles, config.BrowseFiles)
			}
		}, map[string]bool{"linux": true, "darwin": true, "windows": true})
	})
	t.Run("Test with max depth 5", func(t *testing.T) {
		testhelper.RunOSDependentTest(t, "Test with max depth 5", func(t *testing.T) {
			config := RecursionConfig{
				InitialPath: filepath.Join(path, dir),
				MaxDepth: 5,
				SkipDirectories: []string{},
				PriorityDirectories: []string{},
				BrowseFiles: 0,
			}
			err = config.recurse()
			if (err != nil) {
				t.Errorf("Error while browsing files: %v", err)
			}
			// 12 from previous test + 3 in level 4 + 2 in level 5
			expectedFiles := 17
			if (config.BrowseFiles != expectedFiles) {
				t.Errorf("Expected %d files, got %d", expectedFiles, config.BrowseFiles)
			}
		}, map[string]bool{"linux": true, "darwin": true, "windows": true})
	})
	t.Run("Test with max depth 8", func(t *testing.T) {
		testhelper.RunOSDependentTest(t, "Test all depths", func(t *testing.T) {
			config := RecursionConfig{
				InitialPath: filepath.Join(path, dir),
				MaxDepth: 8,
				SkipDirectories: []string{},
				PriorityDirectories: []string{},
				BrowseFiles: 0,
			}
			err = config.recurse()
			if (err != nil) {
				t.Errorf("Error while browsing files: %v", err)
			}
			// All files from all levels (24 total)
			expectedFiles := 24
			if (config.BrowseFiles != expectedFiles) {
				t.Errorf("Expected %d files, got %d", expectedFiles, config.BrowseFiles)
			}
		}, map[string]bool{"linux": true, "darwin": true, "windows": true})
	})
	t.Cleanup(func() {
		defer os.RemoveAll(filepath.Join(path, dir))
	})
}