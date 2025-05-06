package action

// Structure log pour la config + le tableau de fichier a logger
// On va log a la fin de la recursion pour faire 1 seul write
type LogFormat string


type FileInfo map[string]string

const (
	LogFormatText LogFormat = "text"
	LogFormatJSON LogFormat = "json"
	LogFormatCSV  LogFormat = "csv"
)

type Log struct {
	Format        LogFormat
	LogRepository string
	FilesInfo     []FileInfo
}

type Action struct {
	Delete    bool
	Log       bool
	LogConfig Log
}
