package main

import (
	"drive-janitor/parsing-conf"
	"drive-janitor/rules"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"
)

func TestEndToEnd(t *testing.T) {
	// Parse a config for the test
	// Then load the program from the parsed config

	t.Run("Test basic config log", func(t *testing.T) {
		pwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("Error getting current working directory: %v", err)
		}
		configPath := pwd + "/config_test/config_basic.yml"
		rulesInfo, err := parsing.ParsingConfigFile(configPath)
		rules := rulesInfo.RulesArray
		if err != nil {
			t.Fatalf("Error parsing config file: %v", err)
		}
		// Should not happen
		if len(rules) == 0 {
			t.Fatalf("Parsed rules are empty")
		}

		rulesInfo.Loop()

		// we wait before checking the log file
		time.Sleep(1000 * 1000 * 1000) // 2 seconds

		// Check the log file and what it contains
		// For the moment just the file
		checkLogs(rules, t)
	})

	t.Run("End to end recursion, with skip directories", func(t *testing.T) {
		pwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("Error getting current working directory: %v", err)
		}
		configPath := pwd + "/config_test/config_recursion.yml"
		rulesInfo, err := parsing.ParsingConfigFile(configPath)
		rules := rulesInfo.RulesArray
		if err != nil {
			t.Fatalf("Error parsing config file: %v", err)
		}
		// Should not happen
		if len(rules) == 0 {
			t.Fatalf("Parsed rules are empty")
		}
		rulesInfo.Loop()
		for _, rule := range rules {
			fmt.Println("total", rule.Recursion.BrowseFiles)
			// Abritray value, but could scale if needed
			// need to to count the number of files in the directory
			// and check if the number of files is correct
			// by subing the number of files in the skip directories
			if rule.Recursion.BrowseFiles != 58 {
				t.Fatalf("Number of files browsed is not correct: %d\n We probaly didn't skipdir, or we changed samples repo", rule.Recursion.BrowseFiles)
			}
		}
	})

	t.Run("End to end recursion, checking regex filename", func(t *testing.T) {
		pwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("Error getting current working directory: %v", err)
		}
		configPath := pwd + "/config_test/config_regex_match.yml"
		rulesInfo, err := parsing.ParsingConfigFile(configPath)
		rules := rulesInfo.RulesArray
		if err != nil {
			t.Fatalf("Error parsing config file: %v", err)
		}
		// Should not happen
		if len(rules) == 0 {
			t.Fatalf("Parsed rules are empty")
		}
		rulesInfo.Loop()
		for i, rule := range rules {
			if i == 0 {
				for _, logInfo := range rule.Action.LogConfig.FilesInfo {
					assertMatchTestFilname(logInfo, t)
				}
			}
			fmt.Println("total", rule.Recursion.BrowseFiles)
		}
	})

	t.Run("End to end recursion, checking YARA rules", func(t *testing.T) {
		pwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("Error getting current working directory: %v", err)
		}
		configPath := pwd + "/config_test/config_yara_match.yml"
		rulesInfo, err := parsing.ParsingConfigFile(configPath)
		rules := rulesInfo.RulesArray
		if err != nil {
			t.Fatalf("Error parsing config file: %v", err)
		}
		// Should not happen
		if len(rules) == 0 {
			t.Fatalf("Parsed rules are empty")
		}
		rulesInfo.Loop()
		for i, rule := range rules {
			if i == 0 {
				for _, logInfo := range rule.Action.LogConfig.FilesInfo {
					assertMatchTestYara(logInfo, t)
				}
			}
			fmt.Println("total", rule.Recursion.BrowseFiles)
		}
	})

}

// Arbritary checking file matching in samples2
// This function only apply to the test : `End to end recursion, checking regex filename`
func assertMatchTestFilname(fileInfo map[string]string, t *testing.T) {
	allowed := []string{"samples2/Elephant.txt", "samples2/Kangourou.wav", "samples2/KangourouElephant.voc", "samples2/elephant.webp", "samples2/elkanelkangourou.tiff", "samples2/kangourou.ra"}
	for i, _ := range allowed {
		allowed[i] = filepath.FromSlash(allowed[i])
	}
	if !slices.Contains(allowed, fileInfo["path"]) {
		t.Fatalf("Match a file we should not match: %v", fileInfo["path"])
	}
}

// Arbritary checking YARA matching in samples2
func assertMatchTestYara(fileInfo map[string]string, t *testing.T) {
	allowed := []string{"samples/test_yara_match.txt", "detection/checkyara/yararules/malicious.yar"}
	for i, _ := range allowed {
		allowed[i] = filepath.FromSlash(allowed[i])
	}
	if !slices.Contains(allowed, fileInfo["path"]) {
		t.Fatalf("Match a file we should not match: %v", fileInfo["path"])
	}
}

func assertLogFileExist(logFile string, t *testing.T) bool {
	if logFile == "" {
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
