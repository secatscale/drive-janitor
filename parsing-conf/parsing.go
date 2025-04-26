package parsing

import (
	"fmt"
	"log"
	"os"

	// Ensure this package contains the definition for Action
	"drive-janitor/rules"

	"gopkg.in/yaml.v2"
)

func ParsingConfigFile(configPath string) (rules.RulesArray, error) {
	cfg, err := parseYAMLFile(configPath)
	if err != nil {
		fmt.Printf("Error parsing config file: %v\n", err)
		return rules.RulesArray{}, err
	}
	if !mandatoryFieldsGave(cfg) {
		fmt.Println("Error: missing mandatory fields in config file")
		return rules.RulesArray{}, fmt.Errorf("missing mandatory fields in config file")
	}
	rules := fillStructs(cfg)
	return rules, nil
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