package strdist

import (
	"github.com/nickwells/golem/mathutil"
	"unicode/utf8"
)

// Levenshtein calculates the Levenshtein distance between strings a and b
func Levenshtein(a, b string) int {
	aLen := utf8.RuneCountInString(a)
	bLen := utf8.RuneCountInString(b)
	d := make([][]int, aLen+1)
	for i := range d {
		d[i] = make([]int, bLen+1)
		d[i][0] = i
	}

	for i := 1; i <= bLen; i++ {
		d[0][i] = i
	}

	for j, bRune := range b {
		for i, aRune := range a {
			var subsCost int
			if aRune != bRune {
				subsCost = 1
			}

			del := d[i][j+1] + 1
			ins := d[i+1][j] + 1
			sub := d[i][j] + subsCost

			d[i+1][j+1] = mathutil.MinOfInt(del, ins, sub)
		}
	}

	return d[aLen][bLen]
}
