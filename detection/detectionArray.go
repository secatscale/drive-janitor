package detection

func (detectionArray DetectionArrayInfo) AsMatch(filepath string) (bool, error) {
	for _, detection := range detectionArray.Detections {
		match, _ := detection.IsDetected(filepath, detectionArray.DetectionInfo)
		if match {
			return true, nil
		}
	}
	return false, nil
}
