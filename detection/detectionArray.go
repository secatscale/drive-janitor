package detection

func (detectionArray DetectionArrayInfo) AsMatch(filepath string) ([]string, bool, error) {
	detectionsMatch := []string{}
	for _, detection := range detectionArray.Detections {
		match, _ := detection.IsDetected(filepath, detectionArray.DetectionInfo)
		if match {
			detectionsMatch = append(detectionsMatch, detection.Name)
		}
	}
	return detectionsMatch, len(detectionsMatch) > 0, nil
}
