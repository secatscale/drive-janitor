package rules

import (
	"drive-janitor/action"
	"fmt"
	"log"
	"os"
)

func saveLog(action *action.Action) {
	err := action.GetLogFileName()
	if err != nil {
		fmt.Println("Error getting log file name:", err)
		// Temporary exiting
		os.Exit(1)
	}
	fmt.Printf("Log file: %s\n", action.LogConfig.LogRepository)

	err = action.SaveToFile()
	if err != nil {
		fmt.Println("Error saving log file:", err)
		// Temporary exiting
		os.Exit(1)
	}
}

func (r RulesArray) Loop() {
	for _, rules := range r {
		// Running each rules in a separate goroutine
		go func(rules Rules) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Recovered panic in goroutine: %v", r)
				}
			}()
			err := rules.Recursion.Recurse(rules.Detection, rules.Action)
			if err != nil {
				panic(err) // be careful using panic in goroutines
			}
			if rules.Action.Log {
				saveLog(rules.Action)
			}
		}(rules)
	}
}
