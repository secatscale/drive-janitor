package detection

// Checking if the file matches any detection criteria in the array
func (detectionArray DetectionArray) AsMatch(filepath string) ([]string, bool, error) {
	detectionsMatch := []string{}
	for _, detection := range detectionArray {
		match, _ := detection.IsDetected(filepath)
		if match {
			detectionsMatch = append(detectionsMatch, detection.Name)
		}
	}
	return detectionsMatch, len(detectionsMatch) > 0, nil
}
