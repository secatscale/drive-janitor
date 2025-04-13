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
