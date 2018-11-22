package strdist

import "unicode/utf8"

// Hamming returns the Hamming distance of the two strings. if the two
// strings are of different length then the Hamming distance is increased
// by the difference in lengths. Note that it compares runes rather than
// characters or chars
func Hamming(a, b string) int {
	var d = utf8.RuneCountInString(b) - utf8.RuneCountInString(a)
	if d < 0 {
		d *= -1
		a, b = b, a // a is longer than b so swap
	}

	var offset int
	for _, aRune := range a {
		bRune, width := utf8.DecodeRuneInString(b[offset:])
		offset += width
		if bRune != aRune {
			d++
		}
	}

	return d
}
