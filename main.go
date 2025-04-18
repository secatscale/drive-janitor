package main

import (
	"drive-janitor/parsing-conf"
	"flag"
	"fmt"
	"os"
)

func checkConfigFileArgument(configPath *string) (error) {
	flag.StringVar(configPath, "config", "", "Path to the config file")
	flag.Parse()
	if *configPath == "" {
		fmt.Println("Error: config file path must be provided with -config flag")
		flag.Usage()
		return fmt.Errorf("config file path not provided")
	} else if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		fmt.Printf("Error: config file %s does not exist\n", *configPath)
		return fmt.Errorf("config file does not exist")
	}
	return nil
}

func main() {
	var configPath string
	if (checkConfigFileArgument(&configPath) != nil) {
//		fmt.Println("Error getting config file path")
		os.Exit(1)
	}
	rules, err := parsing.ParsingConfigFile(configPath)
	if (err != nil) {
		fmt.Println("Error parsing config file")
		os.Exit(1)
	}
	rules.Loop()

	if (err != nil) {
		fmt.Println("Error while looping on rules")
		os.Exit(1)
	}
}

