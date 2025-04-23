package rules

import (
	"fmt"
	"os"
)

func (r RulesArray) Loop() {
	for _, rules := range r {
		err := rules.Recursion.Recurse(rules.Detection, &rules.Action)
		if err != nil {
			panic(err)
		}
		err = rules.Action.GetLogFileName()
		if err != nil {
			fmt.Println("Error getting log file name:", err)
			os.Exit(1)
		}
		fmt.Printf("Log file: %s\n", rules.Action.LogConfig.LogRepository)
		fmt.Printf("Log FilesInfo: %s\n", rules.Action.LogConfig.FilesInfo)
		err = rules.Action.SaveToFile()
		if err != nil {
			fmt.Println("Error saving log file:", err)
			os.Exit(1)
		}
	}
}
