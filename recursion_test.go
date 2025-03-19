package main

import (
	"errors"
	"io/fs"
	"os"
	"testing"
)

//		var args = []string{*name, *path}
// Next step is to step with relative path i think from arguments

var errFileNotFound = errors.New("File not found")

// Using WalkDir
func lookingForFileStepByStep(name string, start_path string) ([] string, error) {
	var paths []string
	err := fs.WalkDir(os.DirFS(start_path), ".", func(path string, d fs.DirEntry, err error) error {
		if (err != nil) {
			return err
		}
		if (d.Name() == name && d.Type().IsRegular()) {
			paths = append(paths, path)
		}
		return nil
	})
	if (err != nil) {
		return nil, err
	}
	if (len(paths) == 0) {
		return nil, errFileNotFound
	}
	return paths, nil
}

// Using Glob
func lookingForFileByPattern(pattern string, start_path string) ([]string, error) {
	// Better creating well the path before calling this function
	start := os.DirFS(start_path);

	// To match in sub directory
	matches, err := fs.Glob(start, pattern)
	if (err != nil) {
		// Case fs.Glob crashed
		return nil, err
	}
	if (len(matches) == 0) {
		return nil, errFileNotFound
	}
	var paths []string
    for _, v := range matches {
       paths = append(paths, start_path + v)
    }
	return paths, nil
}


func TestSearchFile(t *testing.T) {
	t.Run("file found by pattern", func(t *testing.T) {
		name := "file.txt"
		path := "./"
		_, err := os.Create(path + name)
		defer os.Remove(path + name)
		if (err != nil) {
			t.Errorf("Erreur creating file for the test : %q", err.Error())
		}
		// Kind of regex, will search only at the root of the path
		pattern := "file.txt"
		_, err = lookingForFileByPattern(pattern, path)
		assertNoError(t, err)
	})
	t.Run("file not found by pattern", func(t *testing.T) {
		name := "file.txt"
		path := "./"
		_, err := lookingForFileByPattern(name, path)
		assertError(t, err, errFileNotFound)
	})

	t.Run("file found into sub dir by pattern", func(t *testing.T) {
		dir := "test"
		os.MkdirAll(dir, 0755)
		defer os.RemoveAll(dir)
		name := "file.txt"
		path := "./"
		_, err := os.Create(path + dir + "/" + name)
		if (err != nil) {
			t.Errorf("Erreur creating file for the test : %q", err.Error())
		}
		// will look only into sub directory
		pattern := "**/file.txt"
		_, err = lookingForFileByPattern(pattern, path)
		assertNoError(t, err)
	})
	t.Run("file found", func(t *testing.T) {
		name := "file.txt"
		path := "./"
		_, err := os.Create(path + name)
		defer os.Remove(path + name)
		if (err != nil) {
			t.Errorf("Erreur creating file for the test : %q", err.Error())
		}
		// Kind of regex, will search only at the root of the path
		_, err = lookingForFileStepByStep(name, path)
		assertNoError(t, err)
	})

	t.Run("file not found", func(t *testing.T) {
		name := "file.txt"
		path := "./"
		_, err := lookingForFileStepByStep(name, path)
		assertError(t, err, errFileNotFound)
	})

	t.Run("file found into sub dir", func(t *testing.T) {
		dir := "test"
		os.MkdirAll(dir, 0755)
		defer os.RemoveAll(dir)
		name := "file.txt"
		path := "./"
		_, err := os.Create(path + dir + "/" + name)
		if (err != nil) {
			t.Errorf("Erreur creating file for the test : %q", err.Error())
		}
		// will look only into sub directory
		_, err = lookingForFileStepByStep(name, path)
		assertNoError(t, err)
	})
}