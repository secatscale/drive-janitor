package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

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

type ConfigLog struct {
	Name           string
	Activated      bool
	Log_Repository string
	Crypt          bool
}

type ConfigAction struct {
	Name   string
	Delete bool
	Log    ConfigLog
}

type ConfigRule struct {
	FinalRuleName string
	Recursion     ConfigRecursion
	Detection     ConfigDetection
	Action        ConfigAction
}

type Config struct {
	Name       string
	Version    string
	Detections []ConfigDetection
	Recursions []ConfigRecursion
	Actions    []ConfigAction
	Rules      []ConfigRule
}

func ParseYAMLFile(filePath string) (Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("error unmarshaling yaml: %v", err)
	}
	return cfg, nil
}
