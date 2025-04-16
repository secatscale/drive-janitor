package main

import (
	"drive-janitor/action"
	"drive-janitor/detection"
	"drive-janitor/recursion"
	"flag"
	"fmt"
	"os"
)

func takeArguments(path *string, depth *int, extension *string /* , action *string */, age *int) {

	// Get the current working directory
	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		os.Exit(1)
	}

	// Define command-line flags
	flag.StringVar(path, "path", currentPath, "Path from where we should to check")
	flag.IntVar(depth, "depth", -1, "Maximum directory depth to search (negative for no limit)")
	flag.StringVar(extension, "type", "", "File mimeType to filter (required)")
	// flag.StringVar(action, "action", "list", "Action to perform on files (list, count, size, delete)")
	flag.IntVar(age, "age", -1, "Age of files to filter (in days, negative for no limit)")
	// flag.StringVar(config, "config", "", "Path to the config.yml file (optional)")

	// Parse flags
	flag.Parse()
}

func validateArguments(path string, depth int, extension string /* action string, */, age int) {
	// Validate that path is provided
	if path == "" {
		fmt.Println("Error: path must be provided with -path flag")
		flag.Usage()
		os.Exit(1)
	}

	if extension == "" {
		fmt.Println("Error: file extension must be provided with -type flag")
		flag.Usage()
		os.Exit(1)
	}
}

func main() {

	// Temporary taking arguments from command line
	// Define command-line flags
	var (
		path     string
		depth    int
		mimeType string
		// action   string
		age int
	)

	// Take arguments from command line
	takeArguments(&path, &depth, &mimeType /* , &action */, &age)
	// Validate arguments
	validateArguments(path, depth, mimeType /* , action */, age)

	recursion := recursion.Recursion{
		InitialPath:         path,
		MaxDepth:            depth,
		BrowseFiles:         0,
		SkipDirectories:     []string{},
	}

	mimeIsSupported, err := detection.SupportType(mimeType)
	if err != nil {
		fmt.Println("Error getting MIME type:", err)
		os.Exit(1)
	}
	if !mimeIsSupported {
		fmt.Println("Error: MIME type not supported")
		os.Exit(1)
	}

	detection := detection.Detection{
		MimeType: mimeType,
		Age:      age,
	}

	action := action.Action{
		Delete: false,
		Log:    true,
		LogConfig: action.Log{
			Format:        action.LogFormatText,
			LogRepository: "logs/",
			Files:         []string{},
		},
	}

	fmt.Println("MIME type:", detection.MimeType)

	err = recursion.Recurse(detection, &action)
	if err != nil {
		fmt.Println("Error while browsing files:", err)
		os.Exit(1)
	}
	fmt.Println(recursion.BrowseFiles, path)

	fmt.Printf("Scanning directory: %s\n", path)
	fmt.Printf("File extension: %s\n", mimeType)
	fmt.Printf("Max depth: %v\n", depth)
	// fmt.Printf("Action: %s\n", action)

	err = action.GetLogFileName()
	if err != nil {
		fmt.Println("Error getting log file name:", err)
		os.Exit(1)
	}
	fmt.Printf("Log file: %s\n", action.LogConfig.LogRepository)

	err = action.SaveToFile()
	if err != nil {
		fmt.Println("Error saving log file:", err)
		os.Exit(1)
	}
}
