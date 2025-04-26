package parsing

import (
	"fmt"
	"log"
	"os"
	"strings"

	// Ensure this package contains the definition for Action
	"drive-janitor/os_utils"
	"drive-janitor/rules"

	"gopkg.in/yaml.v2"
)

func ParsingConfigFile(configPath string) (rules.RulesArray, error) {
	cfg, err := parseYAMLFile(configPath)
	if err != nil {
		fmt.Printf("Error parsing config file: %v\n", err)
		return rules.RulesArray{}, err
	}
	err = expandPathsInConfig(&cfg)
	if err != nil {
		fmt.Printf("Error expanding paths in config file: %v\n", err)
		return rules.RulesArray{}, err
	}
	if !mandatoryFieldsGave(cfg) {
		fmt.Println("Error: missing mandatory fields in config file")
		return rules.RulesArray{}, fmt.Errorf("missing mandatory fields in config file")
	}
	rulesArray := fillStructs(cfg)
	return rulesArray, nil
}

func parseYAMLFile(filePath string) (Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("filepath: %v\nerror unmarshaling yaml: %v", filePath, err)
	}
	return cfg, nil
}

// expandPathsInConfig replaces $TRASH et $DOWNLOAD in the config file
func expandPathsInConfig(cfg *Config) error {
	trashPath, err := os_utils.WhereTrash(os_utils.WhichOs())
	if err != nil {
		log.Printf("error getting trash path: %v", err)
		trashPath = ""
	}
	downloadPath, err := os_utils.GetDownloadPath()
	if err != nil {
		log.Printf("error getting download path: %v", err)
		downloadPath = ""
	}

	// Recursions
	for i := range cfg.Recursions {
		cfg.Recursions[i].Path = expandSinglePath(cfg.Recursions[i].Path, trashPath, downloadPath)
		for j := range cfg.Recursions[i].Path_To_Ignore {
			cfg.Recursions[i].Path_To_Ignore[j] = expandSinglePath(cfg.Recursions[i].Path_To_Ignore[j], trashPath, downloadPath)
		}
	}

	// Logs
	for i := range cfg.Logs {
		cfg.Logs[i].Log_Repository = expandSinglePath(cfg.Logs[i].Log_Repository, trashPath, downloadPath)
	}

	return nil
}

// expandSinglePath remplace $TRASH et $DOWNLOAD dans un seul chemin
func expandSinglePath(path, trashPath, downloadPath string) string {
	pathLower := strings.ToLower(path)
	if strings.Contains(pathLower, "$trash") && trashPath != "" {
		path = strings.ReplaceAll(path, "$TRASH", trashPath)
		path = strings.ReplaceAll(path, "$trash", trashPath)
		path = strings.ReplaceAll(path, "$Trash", trashPath)
	}
	if strings.Contains(pathLower, "$download") && downloadPath != "" {
		path = strings.ReplaceAll(path, "$DOWNLOAD", downloadPath)
		path = strings.ReplaceAll(path, "$download", downloadPath)
		path = strings.ReplaceAll(path, "$Download", downloadPath)
	}
	return path
}

// Parser pour avoir la config
// Ensuite on va remplir les autre structs en verifiants les champs depuis le yaml
func mandatoryFieldsGave(cfg Config) bool {
	err := checkRecursion(cfg)
	if err != nil {
		log.Printf("error in recursion: %v", err)
		return false
	}
	err = checkDetection(cfg)
	if err != nil {
		log.Printf("error in detection: %v", err)
		return false
	}
	err = checkAction(cfg)
	if err != nil {
		log.Printf("error in action: %v", err)
		return false
	}
	err = checkRules(cfg)
	if err != nil {
		log.Printf("error in rules: %v", err)
		return false
	}
	err = checkLog(cfg)
	if err != nil {
		log.Printf("error in logs: %v", err)
		return false
	}
	err = checkUniqueNames(cfg)
	if err != nil {
		log.Printf("error in unique names: %v", err)
		return false
	}
	return true
}
