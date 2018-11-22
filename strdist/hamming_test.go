package strdist_test

import (
	"github.com/NickWells/golem/strdist"
	"testing"
)

// TestHamming ...
func TestHamming(t *testing.T) {
	testCases := []struct {
		a, b    string
		expDist int
	}{
		{
			a:       "a",
			b:       "b",
			expDist: 1,
		},
		{
			a:       "ab",
			b:       "b",
			expDist: 2,
		},
		{
			a:       "aaa",
			b:       "aba",
			expDist: 1,
		},
		{
			a:       "aaa",
			b:       "a¶a",
			expDist: 1,
		},
		{
			a:       "a§a",
			b:       "a¶a",
			expDist: 1,
		},
		{
			a:       "a§abc",
			b:       "a¶a",
			expDist: 3,
		},
		{
			a:       "a",
			b:       "a",
			expDist: 0,
		},
		{
			a:       "",
			b:       "",
			expDist: 0,
		},
		{
			a:       "abc",
			b:       "abc",
			expDist: 0,
		},
	}

	for i, tc := range testCases {
		for _, order := range []string{"a,b", "b,a"} {
			a, b := tc.a, tc.b
			if order == "b,a" {
				a, b = b, a
			}

			if dist := strdist.Hamming(a, b); dist != tc.expDist {
				t.Errorf("test %d (%s): Hamming('%s', '%s') should have been"+
					" %d but was %d",
					i, order, a, b, tc.expDist, dist)
			}
		}
	}
}
