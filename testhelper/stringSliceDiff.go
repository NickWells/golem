package testhelper

func StringSliceDiff(a, b []string) bool {
	if a == nil && b == nil {
		return false
	}

	if a == nil || b == nil {
		return true
	}

	if len(a) != len(b) {
		return true
	}

	for i, v := range a {
		if b[i] != v {
			return true
		}
	}
	return false
}
