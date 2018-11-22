package mathutil

import "math"

// MinOf returns the lesser of the slice of values; it will panic if the
// slice is empty
func MinOf(vals ...float64) float64 {
	min := vals[0]
	for _, v := range vals[1:] {
		min = math.Min(min, v)
	}
	return min
}

// MaxOf returns the greater of the slice of values; it will panic if the
// slice is empty
func MaxOf(vals ...float64) float64 {
	max := vals[0]
	for _, v := range vals[1:] {
		max = math.Max(max, v)
	}
	return max
}

// MinOfInt returns the lesser of the slice of values; it will panic if the
// slice is empty
func MinOfInt(vals ...int) int {
	min := vals[0]
	for _, v := range vals[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

// MaxOfInt returns the greater of the slice of values; it will panic if the
// slice is empty
func MaxOfInt(vals ...int) int {
	max := vals[0]
	for _, v := range vals[1:] {
		if v > max {
			max = v
		}
	}
	return max
}
