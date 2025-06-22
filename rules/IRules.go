package rules

import (
	"drive-janitor/action"
	"drive-janitor/detection"
	"drive-janitor/recursion"
	"sync"
)

type Rules struct {
	Name      string
	Action    *action.Action
	Detection *detection.DetectionArrayInfo
	Recursion *recursion.Recursion
}

type RulesArray []Rules

type RulesInfo struct {
	RulesArray RulesArray
	// Used to wait all goroutine after the exectuion
	WaitGroup *sync.WaitGroup
	// Channel to pass data out of the go routine
	InfoLoop chan recursion.Recursion
}
