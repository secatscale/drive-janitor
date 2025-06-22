package rules

import (
	"drive-janitor/action"
	"drive-janitor/detection"
	"fmt"
	"log"
	"os"
)

func saveLog(action *action.Action, detectionInfo *[]detection.DetectionInfo) {
	err := action.GetLogFileName()
	if err != nil {
		fmt.Println("Error getting log file name:", err)
		// Temporary exiting
		os.Exit(1)
	}
	fmt.Printf("Log file: %s\n", action.LogConfig.LogRepository)

	err = action.SaveToFile(*detectionInfo)
	if err != nil {
		fmt.Println("Error saving log file:", err)
		// Temporary exiting
		os.Exit(1)
	}
}

// Main function to loop through the rules and execute them
func (r RulesInfo) Loop() {
	for _, rules := range r.RulesArray {
		// Running each rules in a separate goroutine
		go func(rules Rules) {
			defer func() {

				// To handle a panic error in a goroutine
				if r := recover(); r != nil {
					log.Printf("Recovered panic in goroutine: %v", r)
				}
				r.WaitGroup.Done()
			}()
			// Add a goroutine to wait
			r.WaitGroup.Add(1)
			err := rules.Recursion.Recurse(rules.Detection, rules.Action)
			// Send recursion info into the channel
			r.InfoLoop <- *rules.Recursion
			if err != nil {
				panic(err) // be careful using panic in goroutines
			}
			if rules.Action.Log {
				saveLog(rules.Action, rules.Detection.DetectionInfo)
			}
		}(rules)
		// This way we make sure to get info from the goroutine
		<-r.InfoLoop
	}
}
