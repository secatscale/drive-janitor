package parsing

import (
	"fmt"
	"log"
	"os"

	"drive-janitor/action"
	"drive-janitor/detection"
	"drive-janitor/recursion"

	"gopkg.in/yaml.v2"
)

type Rules struct {
	Name	 string
	Action 	 action.Action
	Detection detection.Detection
	Recursion recursion.Recursion
}

type ConfigDetection struct {
	Name     string
	MimeType string
	Filename string
	Max_Age  int
}

type ConfigRecursion struct {
	Name           string
	Path           string
	MaxDepth       int
	Path_To_Ignore string
}

type ConfigAction struct {
	Name   string
	Delete bool
	Log    string
}

type ConfigRule struct {
	Name	 string
	Action 	 string
	Detection string
	Recursion string
}

type ConfigLog struct {
	Name           string
	Log_Repository string
}

type Config struct {
	Name       string
	Version    string
	Detections []ConfigDetection
	Recursions []ConfigRecursion
	Actions    []ConfigAction
	Rules      []ConfigRule
	Logs	   []ConfigLog
}

func ParseYAMLFile(filePath string) (Config, error) {
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
func MandatoryFieldsGave(cfg Config) bool {
	err := CheckRecursion(cfg)
	if err != nil {
		log.Printf("error in recursion: %v", err)
		return false;
	}
	err = CheckDetection(cfg)
	if err != nil {
		log.Printf("error in detection: %v", err)
		return false;
	}
	err = CheckAction(cfg)
	if err != nil {
		log.Printf("error in action: %v", err)
		return false;
	}
	err = CheckRules(cfg)
	if err != nil {
		log.Printf("error in rules: %v", err)
		return false;
	}
	return true
}

//Checker que les noms soit bien differents, pas possible d'avoir le meme nom pour deux actions, deux recursions, etc
func CheckUniqueNames(cfg Config) error {
	seen := make(map[string]bool)

	for _, recursion := range cfg.Recursions {
		if seen[recursion.Name] {
			return fmt.Errorf("duplicate recursion name: %s", recursion.Name)
		}
		seen[recursion.Name] = true
	}

	for _, detection := range cfg.Detections {
		if seen[detection.Name] {
			return fmt.Errorf("duplicate detection name: %s", detection.Name)
		}
		seen[detection.Name] = true
	}

	for _, action := range cfg.Actions {
		if seen[action.Name] {
			return fmt.Errorf("duplicate action name: %s", action.Name)
		}
		seen[action.Name] = true
	}

	for _, rule := range cfg.Rules {
		if seen[rule.Name] {
			return fmt.Errorf("duplicate rule name: %s", rule.Name)
		}
		seen[rule.Name] = true
	}

	return nil
}

func CheckRecursion(cfg Config) error {
	if len(cfg.Recursions) == 0 {
		return fmt.Errorf("at least one recursion is required")
	}
	for _, recursion := range cfg.Recursions {
		if recursion.Name == "" {
			return fmt.Errorf("name is required")
		}
		if recursion.Path == "" {
			return fmt.Errorf("path is required")
		}
	}
	return nil
}
func CheckDetection(cfg Config) error {
	if len(cfg.Detections) == 0 {
		return fmt.Errorf("at least one detection is required")
	}
	for _, detection := range cfg.Detections {
		if detection.Name == "" {
			return fmt.Errorf("name is required")
		}
		if detection.MimeType == "" && detection.Max_Age == 0 && detection.Filename == "" {
			return fmt.Errorf("at least one of mime type, max age or filename is required")
		}
	}
	return nil
}
func CheckAction(cfg Config) error {
	if len(cfg.Actions) == 0 {
		return fmt.Errorf("at least one action is required")
	}
	for _, action := range cfg.Actions {
		if action.Name == "" {
			return fmt.Errorf("name is required")
		}
	}
	return nil
}

func CheckRules(cfg Config) error {
	if len(cfg.Rules) == 0 {
		return fmt.Errorf("at least one rule is required")
	}
	for _, rule := range cfg.Rules {
		if rule.Name == "" {
			return fmt.Errorf("name is required")
		}
		if rule.Action == "" {
			return fmt.Errorf("action is required")
		}
		if rule.Detection == "" {
			return fmt.Errorf("detection is required")
		}
		if rule.Recursion == "" {
			return fmt.Errorf("recursion is required")
		}
	}
	return nil
}

func fillStructs(cfg Config) []Rules {
	var rules []Rules
	for _, rulesCfg := range cfg.Rules {
		localRules := Rules{}
		for _, recursion := range cfg.Recursions {
			if rulesCfg.Recursion == recursion.Name {
				rulesCfg.Recursion = recursion.Name
				// New recursion struct
				localRules.Recursion = recursion.Recursion{
					Name:             recursion.Name,
					InitialPath:         recursion.Path,
					MaxDepth:            recursion.MaxDepth,
					SkipDirectories:     []string{},
					// Will be deleted later i think
					BrowseFiles:         0,
				}
			}
		}
		for _, detection := range cfg.Detections {
			var detectionStruct []detection.Detection
			if rulesCfg.Detection == detection.Name {
				rulesCfg.Detection = detection.Name
				// New detection struct
				detectionStruct = append(detectionStruct, detection.Detection{
					Name:     detection.Name,
					MimeType: detection.MimeType,
					Filename: detection.Filename,
					Max_Age:  detection.Max_Age,
				})
			}
		}
		for _, action := range cfg.Actions {
			if rulesCfg.Action == action.Name {

				// New action struct
				localRules.Action = action.Action{
					Name:   action.Name,
					Delete: action.Delete,
//					Log:    action.Log,
//				LogConfig: action.Log{
//						Log_Repository: action.Log,
						// Will be deleted later i think
//						Files:         []string{},
	//					Format:        action.LogFormatText,
					}
			}
		}
		rules = append(rules, localRules)
	}
	return rules
}
