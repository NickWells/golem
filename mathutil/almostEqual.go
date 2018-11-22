package mathutil

import "math"

// WithinNPercent returns true if a and b are within epsilon percent of one
// another.
// Strictly speaking the test is for whether the difference between a and b as
// a proportion of the larger value is less than epsilon. So an epsilon of 10
// will test for numbers within 10% of each other
// Numbers with differing sign are always considered different regardless of
// proximity
func WithinNPercent(a, b, epsilon float64) bool {
	if a == b {
		return true
	}

	absA := math.Abs(a)
	absB := math.Abs(b)

	isPosA := (a == absA)
	isPosB := (b == absB)

	if !((isPosA && isPosB) || (!isPosA && !isPosB)) {
		return false
	}

	diff := 100.0 * math.Abs(a-b)
	return (diff / math.Max(absA, absB)) < epsilon
}

// AlmostEqual returns true if a and b are within epsilon of one another.
func AlmostEqual(a, b, epsilon float64) bool {
	if a == b {
		return true
	}

	return math.Abs(a-b) < epsilon
}
