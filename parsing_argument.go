package main

import "errors"

var errMissingArgs = errors.New("Missing arguments");

func checkNameSetError(args []string, name *string, err *error) {
	*name = args[0]
	if (*name == "") {
		*err = errMissingArgs
	}
}

func parse_arguments(args []string) (string, string, error) {
	lenArgs := len(args)
	path := "./"
	var name string = ""
	var err error = nil
	switch (lenArgs) {
		case 1:
			checkNameSetError(args, &name, &err)
			return name, path, err
		case 2:
			checkNameSetError(args, &name, &err)
			if (args[1] != "") {
				path = args[1]
			}
			return name, path, err
		default:
			return "", "", errMissingArgs
	}
}