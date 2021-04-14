package utils

// IsContained returns: a \subset b
func IsContained(a, b []string) bool {
	res := true
	for _, ai := range a {
		isInB := false
		for _, bi := range b {
			if ai == bi {
				isInB = true
				break
			}
		}
		if !isInB {
			res = false
			break
		}
	}
	return res
}
