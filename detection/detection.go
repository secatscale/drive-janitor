package detection

func (detection DetectionConfig) FileTypeMatching(path string) (bool, error){
	// Call la function check type sur le path
	return true, nil
}

func (detection DetectionConfig) FileAgeMatching(path string) (bool, error){
	// Case where the age is not set
	if (detection.Age < 0) {
		return true, nil
	}
	// Call la function check age sur le path
	return true, nil
}