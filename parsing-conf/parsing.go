package parsing

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"drive-janitor/os_utils"
	"drive-janitor/recursion"
	"drive-janitor/rules"

	"gopkg.in/yaml.v2"
)

// We parse, checking for errors in the config file
// Also preparing our structs for the rules
func ParsingConfigFile(configPath string) (rules.RulesInfo, error) {
	cfg, err := parseYAMLFile(configPath)
	if err != nil {
		return rules.RulesInfo{}, err
	}
	err = expandPathsInConfig(&cfg)
	if err != nil {
		err = fmt.Errorf("error expanding paths in config file: %v", err)
		return rules.RulesInfo{}, err
	}
	if !mandatoryFieldsGave(cfg) {
		return rules.RulesInfo{}, fmt.Errorf("missing mandatory fields in config file")
	}
	rulesArray := fillStructs(cfg)
	var rulesInfo = rules.RulesInfo{
		RulesArray: rulesArray,
		WaitGroup:  &sync.WaitGroup{},
		InfoLoop:   make(chan recursion.Recursion),
	}
	return rulesInfo, nil
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

// expandPathsInConfig replaces $TRASH, $DOWNLOAD and $HOME in the config file
func expandPathsInConfig(cfg *Config) error {
	pathsTrashToReplace, pathsDownloadToReplace, pathsHomeToReplace, err := findPathToReplace(cfg)
	if err != nil {
		err = fmt.Errorf("error finding paths to replace: %w", err)
		return err
	}
	// Recherche le repertoire de la corbeille uniquement si il y a des chemins corbeille à remplacer
	if len(pathsTrashToReplace) > 0 {
		trashPath, err := os_utils.WhereTrash(os_utils.WhichOs())
		if err != nil {
			log.Printf("error getting trash path in expandPathsInConfig: %v", err)
			trashPath = ""
		}
		for _, path := range pathsTrashToReplace {
			*path = expandSinglePath(*path, trashPath, "", "")
		}
	}
	// Recherche le repertoire de download uniquement si il y a des chemins download à remplacer
	if len(pathsDownloadToReplace) > 0 {
		downloadPath, err := os_utils.GetDownloadPath()
		if err != nil {
			log.Printf("error getting download path in expandPathsInConfig: %v", err)
			downloadPath = ""
		}
		for _, path := range pathsDownloadToReplace {
			*path = expandSinglePath(*path, "", downloadPath, "")
		}
	}
	// Recherche le repertoire home uniquement si il y a des chemins home à remplacer
	if len(pathsHomeToReplace) > 0 {
		homePath, err := os.UserHomeDir()
		if err != nil {
			log.Printf("error getting home path in expandPathsInConfig: %v", err)
			homePath = ""
		}
		for _, path := range pathsHomeToReplace {
			*path = expandSinglePath(*path, "", "", homePath)
		}
	}
	return nil
}

// fonction pour parcourir la config, et detecter les paths qui nécéssite d'etre remplacé
func findPathToReplace(cfg *Config) ([]*string, []*string, []*string, error) {
	var pathsTrashToReplace []*string
	var pathsDownloadToReplace []*string
	var pathsHomeToReplace []*string
	// Parcours par indice pour eviter de faire des copies et pouvoir retourner des pointeurs
	// Paths dans ConfigRecursion :
	for i := range cfg.Recursions {
		rec := &cfg.Recursions[i]
		// Path
		lowerPath := strings.ToLower(rec.Path)
		if strings.Contains(lowerPath, "$trash") {
			pathsTrashToReplace = append(pathsTrashToReplace, &rec.Path)
		}
		if strings.Contains(lowerPath, "$download") {
			pathsDownloadToReplace = append(pathsDownloadToReplace, &rec.Path)
		}
		if strings.Contains(lowerPath, "$home") {
			pathsHomeToReplace = append(pathsHomeToReplace, &rec.Path)
		}
		// Path_To_Ignore
		for j := range rec.Path_To_Ignore {
			p := &rec.Path_To_Ignore[j]
			lower := strings.ToLower(*p)
			if strings.Contains(lower, "$trash") {
				pathsTrashToReplace = append(pathsTrashToReplace, p)
			}
			if strings.Contains(lower, "$download") {
				pathsDownloadToReplace = append(pathsDownloadToReplace, p)
			}
			if strings.Contains(lower, "$home") {
				pathsHomeToReplace = append(pathsHomeToReplace, p)
			}
		}
	}
	// Paths dans ConfigLog : LogRepository
	for i := range cfg.Logs {
		lg := &cfg.Logs[i]
		lowerLog := strings.ToLower(lg.Log_Repository)
		if strings.Contains(lowerLog, "$trash") {
			pathsTrashToReplace = append(pathsTrashToReplace, &lg.Log_Repository)
		}
		if strings.Contains(lowerLog, "$download") {
			pathsDownloadToReplace = append(pathsDownloadToReplace, &lg.Log_Repository)
		}
		if strings.Contains(lowerLog, "$home") {
			pathsHomeToReplace = append(pathsHomeToReplace, &lg.Log_Repository)
		}
	}
	return pathsTrashToReplace, pathsDownloadToReplace, pathsHomeToReplace, nil
}

// expandSinglePath remplace $TRASH et $DOWNLOAD dans un seul chemin
func expandSinglePath(path, trashPath, downloadPath, homePath string) string {
	pathsplit := strings.Split(path, "/")
	for i, part := range pathsplit {
		if strings.Contains(strings.ToLower(part), "$trash") {
			pathsplit[i] = trashPath
		}
		if strings.Contains(strings.ToLower(part), "$download") {
			pathsplit[i] = downloadPath
		}
		if strings.Contains(strings.ToLower(part), "$home") {
			pathsplit[i] = homePath
		}
	}
	path = strings.Join(pathsplit, "/")
	return path
}

// Parser pour avoir la config
// Ensuite on va remplir les autre structs en verifiants les champs depuis le yaml
func mandatoryFieldsGave(cfg Config) bool {
	err := checkRecursion(cfg)
	if err != nil {
		log.Printf("mandatoryFieldsGave: error in recursion: %v", err)
		return false
	}
	err = checkDetection(cfg)
	if err != nil {
		log.Printf("mandatoryFieldsGave: error in detection: %v", err)
		return false
	}
	err = checkAction(cfg)
	if err != nil {
		log.Printf("mandatoryFieldsGave: error in action: %v", err)
		return false
	}
	err = checkRules(cfg)
	if err != nil {
		log.Printf("mandatoryFieldsGave: error in rules: %v", err)
		return false
	}
	err = checkLog(cfg)
	if err != nil {
		log.Printf("mandatoryFieldsGave: error in logs: %v", err)
		return false
	}
	err = checkUniqueNames(cfg)
	if err != nil {
		log.Printf("mandatoryFieldsGave: error in unique names: %v", err)
		return false
	}
	return true
}
