package parsing

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"

	"drive-janitor/action" // Ensure this package contains the definition for Action
	"drive-janitor/detection"
	"drive-janitor/recursion"
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
	err = checkUniqueNames(cfg)
	if err != nil {
		log.Printf("error in unique names: %v", err)
		return false
	}
	return true
}

// Checker que les noms soit bien differents, pas possible d'avoir le meme nom pour deux actions, deux recursions, etc
func checkUniqueNames(cfg Config) error {
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

func checkRulesAsValidSubRulesName(cfg Config, rule ConfigRule) error {
	// Check that action rules name is valid
	if (slices.IndexFunc(cfg.Actions, func(a ConfigAction) bool {
		return a.Name == rule.Action
	}) == -1) {
		return fmt.Errorf("action %s not found", rule.Action)
	}
	// Check all the detection rules name are valid
	for _, detection := range rule.Detection {
		if (slices.IndexFunc(cfg.Detections, func(a ConfigDetection) bool {
			return a.Name == detection
		}) == -1) {
		}
	}
	// Check if the recursion rule exists
	if (slices.IndexFunc(cfg.Recursions, func(a ConfigRecursion) bool {
		return a.Name == rule.Recursion
	}) == -1) {
		return fmt.Errorf("action %s not found", rule.Action)
	}
	return nil
}

func checkRecursion(cfg Config) error {
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
		// Check if the path exists
		if _, err := os.Stat(recursion.Path); os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s", recursion.Path)
		}
	}
	return nil
}

func checkDetection(cfg Config) error {
	if len(cfg.Detections) == 0 {
		return fmt.Errorf("at least one detection is required")
	}
	for _, d := range cfg.Detections {
		if d.Name == "" {
			return fmt.Errorf("name is required")
		}
		if d.MimeType == "" && d.Max_Age == 0 && d.Filename == "" {
			return fmt.Errorf("at least one of mime type, max age or filename is required")
		}
		if d.MimeType != "" {
			mimeIsSupported, err := detection.SupportType(d.MimeType)
			if err != nil {
				return fmt.Errorf("error getting MIME type: %v", err)
			}
			if !mimeIsSupported {
				return fmt.Errorf("unsupported MIME type: %s", d.MimeType)
			}

		}
	}
	return nil
}
func checkAction(cfg Config) error {
	if len(cfg.Actions) == 0 {
		return fmt.Errorf("at least one action is required")
	}
	for _, action := range cfg.Actions {
		if action.Name == "" {
			return fmt.Errorf("name is required")
		}
		if action.Log != "" {
			// Check the log rules exist
			if (slices.IndexFunc(cfg.Logs, func(a ConfigLog) bool {
				return a.Name == action.Log
			}) == -1) {
				return fmt.Errorf("log rules: %s in action %s not found", action.Log, action.Name)
			}
		}
	}
	return nil
}

func checkLog(cfg Config) error {
	for _, log := range cfg.Logs {
		if log.Name == "" {
			return fmt.Errorf("name is required")
		}
		if log.Log_Repository == "" {
			return fmt.Errorf("log repository is required")
		}
	}
	return nil
}

func checkLog(cfg Config) error {
	if len(cfg.Logs) == 0 {
		return fmt.Errorf("at least one log is required")
	}
	for _, log := range cfg.Logs {
		if log.Name == "" {
			return fmt.Errorf("name is required")
		}
		if log.Log_Repository == "" {
			return fmt.Errorf("log repository is required")
		}
	}
	return nil
}

func checkRules(cfg Config) error {
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
		if (slices.Contains(rule.Detection, "")) {
			return fmt.Errorf("detection is required")
		}
		if rule.Recursion == "" {
			return fmt.Errorf("recursion is required")
		}
		err := checkRulesAsValidSubRulesName(cfg, rule)
		if err != nil {
			return fmt.Errorf("error in rule %s: %v", rule.Name, err)
		}
	}
	return nil
}


func getRelativePath(path string, pathToIgnore []string) []string {
	var relativePaths []string
	for _, p := range pathToIgnore {
		relative, err := filepath.Rel(path, p)
		if err != nil {
			log.Printf("error getting relative path: %v", err)
			continue
		}
		relativePaths = append(relativePaths, relative)
	}
	return relativePaths
}

func fillStructs(cfg Config) rules.RulesArray {
	var rules_array rules.RulesArray
	for _, rulesCfg := range cfg.Rules {
		localRules := rules.Rules{}
		for _, r := range cfg.Recursions {
			if rulesCfg.Recursion == r.Name {
				rulesCfg.Recursion = r.Name
				// New recursion struct
				localRules.Recursion = &recursion.Recursion{
					Name:             r.Name,
					InitialPath:         r.Path,
					MaxDepth:            r.Max_Depth,
					SkipDirectories:     getRelativePath(r.Path, r.Path_To_Ignore),
					// Will be deleted later i think
					BrowseFiles: 0,
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
						Age:      d.Max_Age,
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
					Delete:    a.Delete,
					LogConfig: actionLog,
					Log:       err == nil,
				}
			}
		}
		rules_array = append(rules_array, localRules)
	}
	return rules_array
}

func getLogRules(logsRules []ConfigLog, logRuleName string) (action.Log, error) {
	for _, log := range logsRules {
		if log.Name == logRuleName {
			return action.Log{
				Format:        action.LogFormatText,
				LogRepository: log.Log_Repository,
			}, nil
		}
	}
	return action.Log{}, fmt.Errorf("log rule not found")
}
