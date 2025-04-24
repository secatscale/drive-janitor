package main

import (
	"drive-janitor/parsing-conf"
	"drive-janitor/rules"
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
	// Comment on check ? 1 pas d'erreur
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