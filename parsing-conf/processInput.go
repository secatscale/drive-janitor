package parsing

import (
	"drive-janitor/action"
	"drive-janitor/detection"
	"drive-janitor/recursion"
	"drive-janitor/rules"
	"fmt"
	"log"
	"path/filepath"
	"slices"
)

func fillRecursion(currentRecursionName string, cfgRecursion ConfigRecursion, currentRules *rules.Rules) {
	if currentRecursionName == cfgRecursion.Name {
		// New recursion struct
		currentRules.Recursion = &recursion.Recursion{
			Name:            cfgRecursion.Name,
			InitialPath:     cfgRecursion.Path,
			MaxDepth:        cfgRecursion.Max_Depth,
			SkipDirectories: getRelativePath(cfgRecursion.Path, cfgRecursion.Path_To_Ignore),
			// Will be deleted later i think
			BrowseFiles: 0,
		}
	}
}

func fillDetection(cfgDetectionList []ConfigDetection, currentDetectionListName []string, currentRules *rules.Rules) {
	detectionArray := detection.DetectionArray{}

	for _, name := range currentDetectionListName {
		matchingIndex := slices.IndexFunc(cfgDetectionList, func(detection ConfigDetection) bool {
			return detection.Name == name
		})
		if matchingIndex != -1 {
			{
				detectionArray = append(detectionArray, detection.Detection{
					Name:     cfgDetectionList[matchingIndex].Name,
					MimeType: cfgDetectionList[matchingIndex].MimeType,
					Filename: cfgDetectionList[matchingIndex].Filename,
					Age:      cfgDetectionList[matchingIndex].Max_Age,
				})
			}
		}
		currentRules.Detection = detectionArray
	}
}

func fillAction(cfgActionList []ConfigAction, currentActionName string, currentRules *rules.Rules, logsRulesList []ConfigLog) {
	for _, a := range cfgActionList {
		if a.Name == currentActionName {
			actionLog, err := getLogRules(logsRulesList, a.Log)
			asLog := err == nil
			if !asLog {
				log.Printf("error getting log rules: %v", err)
			}
			currentRules.Action = action.Action{
				Delete:    a.Delete,
				LogConfig: actionLog,
				Log:       asLog,
			}
		}
	}
}

func fillStructs(cfg Config) rules.RulesArray {
	var rules_array rules.RulesArray
	for _, rulesCfg := range cfg.Rules {
		localRules := rules.Rules{}
		for _, r := range cfg.Recursions {
			fillRecursion(rulesCfg.Recursion /* Name of the recursion rules */, r, &localRules)
		}
		fillDetection(cfg.Detections, rulesCfg.Detection /* List of the detection rules name */, &localRules)
		fillAction(cfg.Actions, rulesCfg.Action /* Name of the action rules */, &localRules, cfg.Logs)

		rules_array = append(rules_array, localRules)
	}
	return rules_array
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
