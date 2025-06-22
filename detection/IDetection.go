package detection

type Detection struct {
	Name         string // Nom de la detection
	MimeType     string // A adapter au moment du parsing genre audio/mpeg si .mp3 etc
	Age          int    // Age du fichier en jour
	Filename     string // Regex sur le nom du fichier
	YaraRulesDir string // Chemin vers le dossier contenant les règles Yara
}

type DetectionArrayInfo struct {
	Detections DetectionArray
	DetectionInfo *[]DetectionInfo
}

type DetectionArray []Detection

type DetectionInfo struct {
	TypeMatch bool
	AgeMatch  bool
	FilenameMatch bool
	YaraMatch bool
	Path string
	Detection *Detection
}
