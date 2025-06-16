package detection

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
