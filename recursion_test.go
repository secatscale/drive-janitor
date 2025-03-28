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
	"fmt"
	"io/fs"
	"os"
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

func getDepth(originalPath string, path string) int {
	return countSeparator(path) - countSeparator(originalPath) + 1
}

func isAboveMaxDepth(initialPath string, path string, maxDepth int) bool {
	return getDepth(initialPath, path) > maxDepth
}

func (config *RecursionConfig) recurse(/* May take dectection and action struct*/) error {
	initialPathFs := os.DirFS(config.InitialPath);
	err := fs.WalkDir(initialPathFs, ".", func(path string, entry fs.DirEntry, err error) error {
//		fmt.Println(getDepth(config.InitialPath, path), config.InitialPath, path, entry.Type().IsDir())
		if (isAboveMaxDepth(config.InitialPath, path, config.MaxDepth)) {
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

// find a way to generate test files easily
func TestRecursion(t *testing.T) {
	dir := "test"
	os.MkdirAll(dir, 0755)
	os.MkdirAll("test/1", 0755)
	name := "file.txt"
	path := "./"
	_, err := os.Create(path + dir + "/" + name)
	if (err != nil) {
		t.Errorf("Error while creating file: %v", err)
	}
	_, err = os.Create(path + dir + "/1/" + "file2.txt")
	if (err != nil) {
		t.Errorf("Error while creating file: %v", err)
	}
	testhelper.RunOSDependentTest(t, "Test Browsering", func(t *testing.T) {
		config := RecursionConfig{
			InitialPath: "./test",
			MaxDepth: 1,
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
	}, map[string]bool{"linux": true, "darwin": true});

	testhelper.RunOSDependentTest(t, "Test max depth", func(t *testing.T) {
		os.MkdirAll("test/1/2", 0755)
		_, err = os.Create(path + dir + "/1/2/" + "file3.txt")
		config := RecursionConfig{
			InitialPath: "./test",
			MaxDepth: 1,
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
	}, map[string]bool{"linux": true, "darwin": true});

	testhelper.RunOSDependentTest(t, "Test Browsering", func(t *testing.T) {
		path, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
		config := RecursionConfig{
			InitialPath: path,
			MaxDepth: 1,
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
	}, map[string]bool{"windows": true});

	testhelper.RunOSDependentTest(t, "Test max depth", func(t *testing.T) {
		path, err := os.Getwd()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		os.MkdirAll(path + "\\1\\2", 0755)
		_, err = os.Create(path + "\\1\\2\\" + "file3.txt")
		config := RecursionConfig{
			InitialPath: "path",
			MaxDepth: 1,
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
	}, map[string]bool{"windows": true});

	t.Cleanup(func() {
		path, err := os.Getwd()
		if (err != nil) {
			fmt.Println("Error:", err)
		}
		if (WhichOs() == "windows") {
			defer os.RemoveAll(path)

		} else {
			defer os.RemoveAll(path + dir + "/1/2")
			defer os.RemoveAll(dir)
		}
	})
}