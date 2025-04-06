package action

// Structure log pour la config + le tableau de fichier a logger
// On va log a la fin de la recursion pour faire 1 seul write
type LogConfig struct {
	Files []string
}

type ActionConfig struct {
	Delete bool
	Log bool
	LogConfig LogConfig
}