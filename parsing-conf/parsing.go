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
	pathsTrashToReplace, pathsDownloadToReplace, err := findPathToReplace(cfg)
	if err != nil {
		log.Printf("error finding paths to replace: %v", err)
		return err
	}
	// Recherche le repertoire de la corbeille uniquement si il y a des chemins corbeille à remplacer
	if len(pathsTrashToReplace) > 0 {
		trashPath, err := os_utils.WhereTrash(os_utils.WhichOs())
		if err != nil {
			log.Printf("error getting trash path: %v", err)
			trashPath = ""
		}
		for _, path := range pathsTrashToReplace {
			*path = expandSinglePath(*path, trashPath, "")
		}
	}
	// Recherche le repertoire de download uniquement si il y a des chemins download à remplacer
	if len(pathsDownloadToReplace) > 0 {
		downloadPath, err := os_utils.GetDownloadPath()
		if err != nil {
			log.Printf("error getting download path: %v", err)
			downloadPath = ""
		}
		for _, path := range pathsDownloadToReplace {
			*path = expandSinglePath(*path, "", downloadPath)
		}
	}
	return nil
}

// fonction pour parcourir la config, et detecter les paths qui nécéssite d'etre remplacé
func findPathToReplace(cfg *Config) ([]*string, []*string, error) {
	var pathsTrashToReplace []*string
	var pathsDownloadToReplace []*string
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
	}
	return pathsTrashToReplace, pathsDownloadToReplace, nil
}

// expandSinglePath remplace $TRASH et $DOWNLOAD dans un seul chemin
func expandSinglePath(path, trashPath, downloadPath string) string {
	pathsplit := strings.Split(path, "/")
	for i, part := range pathsplit {
		if strings.Contains(strings.ToLower(part), "$trash") {
			pathsplit[i] = trashPath
		}
		if strings.Contains(strings.ToLower(part), "$download") {
			pathsplit[i] = downloadPath
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
