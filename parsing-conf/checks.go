package parsing

import (
	"drive-janitor/detection"
	"fmt"
	"os"
	"slices"
)

func checkRuleExists[T Named](list []T, ruleName string) error {
	if slices.IndexFunc(list, func(item T) bool {
		return item.GetName() == ruleName
	}) == -1 {
		return fmt.Errorf("rule %s not found", ruleName)
	}
	return nil
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
	err := checkRuleExists(cfg.Actions, rule.Action)
	if err != nil {
		return fmt.Errorf("action %s not found", rule.Action)
	}
	// Check all the detection rules name are valid
	for _, detection := range rule.Detection {
		err := checkRuleExists(cfg.Detections, detection)
		if err != nil {
			return fmt.Errorf("detection %s not found in rule %s", detection, rule.Name)
		}
	}
	// Check if the recursion rule exists
	err = checkRuleExists(cfg.Recursions, rule.Recursion)
	if err != nil {
		return fmt.Errorf("recursion %s not found in rule %s", rule.Recursion, rule.Name)
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

		if d.MimeType == "" && d.Max_Age == 0 && d.Filename == "" && d.Yara_Rules_Dir == "" {
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
			err := checkRuleExists(cfg.Logs, action.Log)
			if err != nil {
				return fmt.Errorf("log %s not found in action %s", action.Log, action.Name)
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
		if log.Log_Format == "" {
			return fmt.Errorf("log format is required")
		}
		if log.Log_Format != "text" && log.Log_Format != "json" && log.Log_Format != "csv" {
			return fmt.Errorf("unsupported log format: %s (only `text`, `json`, `csv` are supported)", log.Log_Format)
		}
		if log.Log_Repository == "" {
			return fmt.Errorf("log repository is required")
		}
	}
	return nil
}

func checkNotEmpty(str string) error {
	if str == "" {
		return fmt.Errorf("string is empty")
	}
	return nil
}

func checkRules(cfg Config) error {
	if len(cfg.Rules) == 0 {
		return fmt.Errorf("at least one rule is required")
	}
	for _, rule := range cfg.Rules {
		err := checkNotEmpty(rule.Name)
		if err != nil {
			return fmt.Errorf("name is required")
		}
		err = checkNotEmpty(rule.Action)
		if err != nil {
			return fmt.Errorf("action is required")
		}
		if len(rule.Detection) == 0 {
			return fmt.Errorf("at least one detection is required")
		}
		if slices.Contains(rule.Detection, "") {
			return fmt.Errorf("detection is required")
		}
		err = checkNotEmpty(rule.Recursion)
		if err != nil {
			return fmt.Errorf("recursion is required")
		}
		err = checkRulesAsValidSubRulesName(cfg, rule)
		if err != nil {
			return fmt.Errorf("error in rule %s: %v", rule.Name, err)
		}
	}
	return nil
}
