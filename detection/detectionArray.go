package detection

func (detectionArray DetectionArray) AsMatch(filepath string) (bool, error) {
	for _, detection := range detectionArray {
		match, _ := detection.IsDetected(filepath)
		if match {
			return true, nil
		}
	}
	return false, nil
}