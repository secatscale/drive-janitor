package detection

type Detection struct {
	Name	 string // Nom de la detection
	MimeType string // A adapter au moment du parsing genre audio/mpeg si .mp3 etc
	Age      int    // Age du fichier en jour
	Filename string // Regex sur le nom du fichier
}

type DetectionArray []Detection

func (detectionArray DetectionArray) AsMatch(filepath string) (bool, error) {
	for _, detection := range detectionArray {
		match, _ := detection.IsDetected(filepath)
		if (match) {
			return true, nil
		} 
	}
	return false, nil
}