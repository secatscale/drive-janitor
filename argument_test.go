package main

import (
	"flag"
	"testing"
)


var name = flag.String("name", "", "name of the file you search")
var path = flag.String("path", "./", "path from were you want to search")

func assertError(t *testing.T, got error, want error) {
	t.Helper()

	if (got != want) {
		t.Errorf("got %q, want %q", got.Error(), want.Error())
	}
}

func assertNoError(t *testing.T, got error, want error) {
	t.Helper()

	if (got != want) {
		t.Errorf("got %q, want %q", got.Error(), want.Error())
	}
}

//To run with go test -args -name "xx" -path "yy"
func TestArgument(t *testing.T) {
	t.Run("Two arguments no error", func(t *testing.T) {
		var args = []string{*name, *path}
		name, path, err := parse_arguments(args);
		if (name == "" || path == "") {
			assertError(t, err, errMissingArgs);
		}
		t.Logf("Name: %s, Path: %s", name, path)
		assertNoError(t, err, nil)
	})
	t.Run("One argument no error", func(t *testing.T) {
		var args = []string{*name}
		name, _, err := parse_arguments(args);
		if (name == "") {
			assertError(t, err, errMissingArgs);
		}
		t.Logf("Name: %s", name)
		assertNoError(t, err, nil);
	})
	t.Run("No arguments", func(t *testing.T) {
		var args = []string{}
		_, _, err := parse_arguments(args);
		assertError(t, err, errMissingArgs);
	})

}