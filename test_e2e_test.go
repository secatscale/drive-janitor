package main

import (
	"drive-janitor/parsing-conf"
	"drive-janitor/rules"
	"fmt"
	"os"
	"testing"
)

func TestEndToEnd(t *testing.T) {
	// Parse a config for the test
	// Then load the program from the parsed config

	t.Run("Test basic config log", func(t *testing.T) {
		pwd, err := os.Getwd()
		if (err != nil) {
			t.Fatalf("Error getting current working directory: %v", err)
		}
		configPath := pwd + "/config_test/config_basic.yml"
		rules, err := parsing.ParsingConfigFile(configPath)
		if err != nil {
			t.Fatalf("Error parsing config file: %v", err)
		}
		// Should not happen
		if len(rules) == 0 {
			t.Fatalf("Parsed rules are empty")
		}

		rules.Loop()
		// Check the log file and what it contains
		// For the moment just the file 
		checkLogs(rules, t)
	})

	t.Run("End to end recursion, with skip directories", func(t *testing.T) {
		pwd, err := os.Getwd()
		if (err != nil) {
			t.Fatalf("Error getting current working directory: %v", err)
		}
		configPath := pwd + "/config_test/config_recursion.yml"
		rules, err := parsing.ParsingConfigFile(configPath)
		if err != nil {
			t.Fatalf("Error parsing config file: %v", err)
		}
		// Should not happen
		if len(rules) == 0 {
			t.Fatalf("Parsed rules are empty")
		}
		rules.Loop()
		for _, rule := range rules {
			fmt.Println("total", rule.Recursion.BrowseFiles)
			// Abritray value, but could scale if needed
			// need to to count the number of files in the directory
			// and check if the number of files is correct
			// by subing the number of files in the skip directories
			if (rule.Recursion.BrowseFiles != 57) {
				t.Fatalf("Number of files browsed is not correct: %d\n We probaly didn't skipdir, or we changed samples repo", rule.Recursion.BrowseFiles);
			}
		}
	})

}

func assertLogFileExist(logFile string, t *testing.T) bool {
	if (logFile == "") {
		t.Fatalf("Log file path is empty")
		return false
	}
	_, err := os.Stat(logFile)
	if os.IsNotExist(err) {
		t.Fatalf("Log file %s does not exist", logFile)
		return false
	}
	return true
}

func checkLogs(rules rules.RulesArray, t *testing.T) {
	for _, rule := range rules {
		if rule.Action.Log && assertLogFileExist(rule.Action.LogConfig.LogRepository, t) { 
			// Check the content, just openeing for the moment
			logFile, err := os.Open(rule.Action.LogConfig.LogRepository)
			if err != nil {
				t.Fatalf("Error opening log file %s: %v", rule.Action.LogConfig.LogRepository, err)
			}
			defer logFile.Close()
		}
	}
}