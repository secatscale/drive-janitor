package main

import (
	"drive-janitor/detection"
	"drive-janitor/recursion"
	"flag"
	"fmt"
	"os"
)

func takeArguments(path *string, depth *int, extension *string, action *string) {

	// Get the current working directory
	currentPath, err := os.Getwd()
	if (err != nil) {
		fmt.Println("Error getting current working directory:", err)
		os.Exit(1)
	}

	// Define command-line flags
	flag.StringVar(path, "path", currentPath, "Path from where we should to check")
	flag.IntVar(depth, "depth", -1, "Maximum directory depth to search (negative for no limit)")
	flag.StringVar(extension, "ext", "", "File extension to filter (required)")
	flag.StringVar(action, "action", "list", "Action to perform on files (list, count, size, delete)")

	// Parse flags
	flag.Parse()
}

func validateArguments(path string, depth int, extension string, action string) {
	// Validate that path is provided
	if path == "" {
		fmt.Println("Error: path must be provided with -path flag")
		flag.Usage()
		os.Exit(1)
	}

	if (extension == "") {
		fmt.Println("Error: file extension must be provided with -ext flag")
		flag.Usage()
		os.Exit(1)
	}
}

func main() {

	// Temporary taking arguments from command line
	// Define command-line flags
	var (
		path	  string
		depth     int
		extension string
		action    string
	)

	// Take arguments from command line
	takeArguments(&path, &depth, &extension, &action)
	// Validate arguments
	validateArguments(path, depth, extension, action)

	recursion := recursion.RecursionConfig{
		InitialPath: path,
		MaxDepth: depth,
		BrowseFiles: 0,
		SkipDirectories: []string{},
		PriorityDirectories: []string{},
	}

	mime, err := detection.GetMimeType(extension)
	if err != nil {
		fmt.Println("Error getting MIME type:", err)
		os.Exit(1)
	}

	detection := detection.DetectionConfig{
		MimeType: mime,
		Age: -1,
	}

	fmt.Println("MIME type:", detection.MimeType)
	
	recursion.Recurse()
	fmt.Println(recursion.BrowseFiles, path)

	fmt.Printf("Scanning directory: %s\n", path)
	fmt.Printf("File extension: %s\n", extension)
	fmt.Printf("Max depth: %v\n", depth)
	fmt.Printf("Action: %s\n", action)
}
