package detection

type Detection struct {
	Name         string // Nom de la detection
	MimeType     string // A adapter au moment du parsing genre audio/mpeg si .mp3 etc
	Age          int    // Age du fichier en jour
	Filename     string // Regex sur le nom du fichier
	YaraRulesDir string // Chemin vers le dossier contenant les règles Yara
}

type DetectionArray []Detection
