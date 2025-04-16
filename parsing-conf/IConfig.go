package parsing

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
	Detection []string
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