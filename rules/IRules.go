package rules

import (
	"drive-janitor/action"
	"drive-janitor/detection"
	"drive-janitor/recursion"
)

type Rules struct {
	Name	 string
	Action 	 action.Action
	Detection []detection.Detection
	Recursion recursion.Recursion
}

type RulesArray []Rules;
