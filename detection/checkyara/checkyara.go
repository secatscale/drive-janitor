package checkyara

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	yara "github.com/hillu/go-yara/v4"
)

var (
	rules     *yara.Rules
	rulesOnce sync.Once
	rulesErr  error
	rulesDir  string
)

func loadRules(rulesDir string) (*yara.Rules, error) {
	compiler, err := yara.NewCompiler()
	if err != nil {
		return nil, fmt.Errorf("failed to create compiler: %v", err)
	}

	err = filepath.Walk(rulesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing %s: %v", path, err)
		}
		if !info.IsDir() && strings.HasSuffix(path, ".yar") {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open rule file %s: %v", path, err)
			}
			defer file.Close()
			if err := compiler.AddFile(file, ""); err != nil {
				return fmt.Errorf("failed to compile rule %s: %v", path, err)
			}
			fmt.Printf("Loaded YARA rule: %s\n", path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	rules, err := compiler.GetRules()
	if err != nil {
		return nil, fmt.Errorf("failed to get compiled rules: %v", err)
	}
	return rules, nil
}

// Fonction d'initialisation unique
func initRules() {
	rules, rulesErr = loadRules(rulesDir)
}

// Utilisé dans chaque appel pour garantir un chargement unique
func getRules(yaraRulesDir string) (*yara.Rules, error) {
	// Si les ne sont pas déjà initialisées, on les charge
	if rules == nil {
		rulesDir = yaraRulesDir
	}
	rulesOnce.Do(initRules)
	return rules, rulesErr
}

func scanFile(filePath string, rules *yara.Rules) (bool, error) {
	// Prépare le slice qui recevra les matches
	var matches yara.MatchRules

	// 0 = flags par défaut, 0*time.Second = pas de timeout
	if err := rules.ScanFile(filePath, 0, 0*time.Second, &matches); err != nil {
		return false, fmt.Errorf("erreur lors du scan de %s : %v", filePath, err)
	}

	// Retourne true si au moins une règle a matché
	return len(matches) > 0, nil
}

func CheckYara(filePath string, yaraRulesDir string) (bool, error) {
	rules, err := getRules(yaraRulesDir)
	if err != nil {
		return false, fmt.Errorf("failed to get YARA rules: %v", err)
	}
	if rules == nil {
		return false, fmt.Errorf("no YARA rules loaded")
	}
	return scanFile(filePath, rules)
}
