package parsing

import (
	"fmt"
	"log"
	"os"

	"drive-janitor/action" // Ensure this package contains the definition for Action
	"drive-janitor/detection"
	"drive-janitor/recursion"

	"gopkg.in/yaml.v2"
)

type Rules struct {
	Name	 string
	Action 	 action.Action
	Detection []detection.Detection
	Recursion recursion.Recursion
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

	for _, log := range cfg.Logs {
		if seen[log.Name] {
			return fmt.Errorf("duplicate log name: %s", log.Name)
		}
		seen[log.Name] = true
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
		if len(rule.Detection) == 0 {
			return fmt.Errorf("at least one detection is required")
		}
		for _, detection := range rule.Detection {
			if detection == "" {
				return fmt.Errorf("detection is required")
			}
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
		for _, r := range cfg.Recursions {
			if rulesCfg.Recursion == r.Name {
				rulesCfg.Recursion = r.Name
				// New recursion struct
				localRules.Recursion = recursion.Recursion{
					Name:             r.Name,
					InitialPath:         r.Path,
					MaxDepth:            r.MaxDepth,
					SkipDirectories:     []string{},
					// Will be deleted later i think
					BrowseFiles:         0,
				}
			}
		}
		var detectionStruct []detection.Detection
		for _, cfd_detection := range cfg.Detections {
			for _, d := range cfg.Detections {
				if cfd_detection.Name == d.Name {
					// New detection struct
					detectionStruct = append(detectionStruct, detection.Detection{
						Name:     d.Name,
						MimeType: d.MimeType,
						Filename: d.Filename,
						Age:  d.Max_Age,
					})
				}
			}
		}
		localRules.Detection = detectionStruct
		for _, a := range cfg.Actions {
			if rulesCfg.Action == a.Name {
				// New action struct
				// Faut aller checker le nom de la regle de log et l'ajouter dans notre struct action
				actionLog, err := getLogRules(cfg.Logs, a.Log)
				if err != nil {
					log.Printf("error getting log rules: %v", err)
				}
					localRules.Action = action.Action{
					Delete: a.Delete,
					LogConfig:	actionLog,
					Log: err != nil,
					}
			}
		}
		rules = append(rules, localRules)
	}
	return rules
}

func getLogRules(logsRules []ConfigLog, logRuleName string) (action.Log, error) {
	for _, log := range logsRules {
		if (log.Name == logRuleName) {
			return action.Log{
				Format:			action.LogFormatText,
				LogRepository: log.Log_Repository,
			}, nil
		}
	}
	return action.Log{}, fmt.Errorf("log rule not found");
}