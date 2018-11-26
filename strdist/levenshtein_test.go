package strdist_test

import (
	"github.com/nickwells/golem/strdist"
	"testing"
)

// TestLevenshtein ...
func TestLevenshtein(t *testing.T) {
	testCases := []struct {
		name    string
		a, b    string
		expDist int
	}{
		{
			name:    "zero char same",
			a:       "",
			b:       "",
			expDist: 0,
		},
		{
			name:    "single char same",
			a:       "a",
			b:       "a",
			expDist: 0,
		},
		{
			name:    "single char differ",
			a:       "a",
			b:       "b",
			expDist: 1,
		},
		{
			name:    "differ 2",
			a:       "aa",
			b:       "ab",
			expDist: 1,
		},
		{
			name:    "Kitten/Sitting",
			a:       "Kitten",
			b:       "Sitting",
			expDist: 3,
		},
		{
			name:    "Saturday/Sunday",
			a:       "Saturday",
			b:       "Sunday",
			expDist: 3,
		},
	}

	for i, tc := range testCases {
		for _, order := range []string{"a,b", "b,a"} {
			a, b := tc.a, tc.b
			if order == "b,a" {
				a, b = b, a
			}
			dist := strdist.Levenshtein(a, b)
			if dist != tc.expDist {
				t.Logf("test %d: %s :\n", i, tc.name)
				t.Errorf(
					"\t: Levenshtein('%s', '%s') expected distance: %d got: %d",
					a, b, tc.expDist, dist)
			}
		}
	}
}
