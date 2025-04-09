package detection

type DetectionConfig struct {
	MimeType string // A adapter au moment du parsing genre audio/mpeg si .mp3 etc
	Age      int    // Age du fichier en jour
}
