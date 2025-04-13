package main

import (
	"testing"

	"gopkg.in/yaml.v2"
)

func TestParseBasicYAML(t *testing.T) {
	yamlData := []byte(`
name: TestApp
version: 1.2.3
`)

	// Define a basic struct for testing
	type Config struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}
	// Create an instance of the Config struct
	var config Config

	
	// Unmarshal the YAML data into the struct
	err := yaml.Unmarshal(yamlData, &config)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	// Check if the strings are correctly parsed
	assertCorrectyamlstring(t, "TestApp", config.Name)
	
	assertCorrectyamlstring(t, "1.2.3", config.Version)
}

func TestParseYAMLFile1(t *testing.T) {
	t.Run("TestParseYAMLFile1", func(t *testing.T) {
		filePath := "test1.yaml"
		cfg, _ := ParseYAMLFile(filePath)
		assertCorrectyamlstring(t, "Test1", cfg.Name)	
		assertCorrectyamlstring(t, "1.2.3", cfg.Version)
		detection1 := cfg.Detections[0]
		assertCorrectyamlstring(t, "xxxx", detection1.Name)
		assertCorrectyamlstring(t, "image/png", detection1.MimeType)
		assertCorrectyamlstring(t, "*.exe", detection1.Filename)
		assertCorrectYAMLint(t, 2, detection1.Max_Age)
		detection2 := cfg.Detections[1]
		assertCorrectyamlstring(t, "xxxx2", detection2.Name)
		assertCorrectyamlstring(t, "application/pdf", detection2.MimeType)
		assertCorrectyamlstring(t, "pipi.prout", detection2.Filename)
		assertCorrectYAMLint(t, 10, detection2.Max_Age)
	})
	t.Run("TestParseYAMLFile2", func(t *testing.T) {
		filePath := "test2.yaml"
		cfg, _ := ParseYAMLFile(filePath)
		assertCorrectyamlstring(t, "Test2", cfg.Name)	
		assertCorrectyamlstring(t, "1.2.3", cfg.Version)
		recursion1 := cfg.Recursions[0]
		assertCorrectyamlstring(t, "recursion1", recursion1.Name)
		assertCorrectyamlstring(t, "/tmp", recursion1.Path)
		assertCorrectYAMLint(t, 5, recursion1.MaxDepth)
		assertCorrectyamlstring(t, "/tmp/ignore", recursion1.Path_To_Ignore)
	})
	t.Run("TestParseYAMLFile3", func(t *testing.T) {
		filePath := "test3.yaml"
		cfg, _ := ParseYAMLFile(filePath)
		assertCorrectyamlstring(t, "Test3", cfg.Name)	
		assertCorrectyamlstring(t, "1.2.3", cfg.Version)
		action1 := cfg.Actions[0]
		assertCorrectyamlstring(t, "xxxx", action1.Name)
		assertCorrectYAMLbool(t, true, action1.Delete)
//		log1 := action1.Log
	//	assertCorrectyamlstring(t, "/tmp/log", log1.Log_Repository)
//		assertCorrectYAMLbool(t, true, log1.Crypt)
	})
	t.Run("TestParseYAMLFile4", func(t *testing.T) {
		filePath := "test4.yaml"
		cfg, _ := ParseYAMLFile(filePath)
		assertCorrectyamlstring(t, "Test4", cfg.Name)	
		assertCorrectyamlstring(t, "1.2.3", cfg.Version)
		action1 := cfg.Actions[0]
		assertCorrectyamlstring(t, "xxxx", action1.Name)
		assertCorrectYAMLbool(t, false, action1.Delete)
	//	log1 := action1.Log
	//	assertCorrectyamlstring(t, "/tmp/log", log1.Log_Repository)
//		assertCorrectYAMLbool(t, false, log1.Crypt)
	})
}	

func assertCorrectYAMLint(t *testing.T, expected, got int) {
	if expected != got {
		t.Errorf("Expected %d, got %d", expected, got)
	}
}

func assertCorrectyamlstring(t *testing.T, expected, got string) {
	if expected != got {
		t.Errorf("Expected %s, got %s", expected, got)
	}
}

func assertCorrectYAMLbool(t *testing.T, expected, got bool) {
	if expected != got {
		t.Errorf("Expected %t, got %t", expected, got)
	}
}

