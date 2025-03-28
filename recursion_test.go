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

// Test 1 -> checker la max depth
// Test 2 -> checker qu'on parcours bien tous les fichiers


// cette function parcours toues les fichiers et dossiers

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

// find a way to generate test files easily
func TestRecursion(t *testing.T) {
	dir := "test"
	os.MkdirAll(dir, 0755)
	dir_one := "1"
	os.MkdirAll(filepath.Join(dir, dir_one), 0755)
	name := "file.txt"
	path, err :=  os.Getwd()
	if (err != nil) {
		t.Errorf("Error while getting the path: %v", err)
	}
	_, err = os.Create(filepath.Join(path, dir, name))
	if (err != nil) {
		t.Errorf("Error while creating file: %v", err)
	}
	_, err = os.Create(filepath.Join(path, dir, dir_one, "file2.txt"))
	if (err != nil) {
		t.Errorf("Error while creating file: %v", err)
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
			dir_two := "2"
			os.MkdirAll(filepath.Join(path, dir, dir_one, dir_two), 0755)
			_, err = os.Create(filepath.Join(path, dir, dir_one, dir_two, "file3.txt"))
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
		defer os.RemoveAll(path + dir)
	})
}