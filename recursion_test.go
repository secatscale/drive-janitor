package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"testing"
)

var errFileNotFound = errors.New("File not found")

func lookingForFile(name string, path string) (string, error) {
	start := os.DirFS(path);
	matches, err := fs.Glob(start, name)
	if (err != nil) {
		return "", err
	}
	fmt.Println(matches)
	if (len(matches) == 0) {
		return "", errFileNotFound
	}
	return "", nil
}

func TestSearchFile(t *testing.T) {
	t.Run("file found", func(t *testing.T) {
		name := "file.txt"
		path := "./"
		_, err := os.Create(path + name)
		defer os.Remove(path + name)
		if (err != nil) {
			t.Errorf("Erreur creating file for the test : %q", err.Error())
		}
		_, err = lookingForFile(name, path)
		assertNoError(t, err, nil)
	})
	t.Run("file not found", func(t *testing.T) {
		name := "file.txt"
		path := "./"
		_, err := lookingForFile(name, path)
		assertError(t, err, errFileNotFound)
	})
}