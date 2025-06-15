package detection

// Checking if the file matches any detection criteria in the array
func (detectionArray DetectionArray) AsMatch(filepath string) (bool, error) {
	for _, detection := range detectionArray {
		match, _ := detection.IsDetected(filepath)
		if match {
			return true, nil
		}
	}
	return false, nil
}
