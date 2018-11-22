package strdist_test

import (
	"fmt"
	"github.com/NickWells/golem/mathutil"
	"github.com/NickWells/golem/strdist"
	"testing"
)

// TestNGrams tests the NGrams function
func TestNGrams(t *testing.T) {
	testCases := []struct {
		name              string
		s                 string
		n                 int
		expDistinctNGrams int
		expErr            bool
	}{
		{
			name:              "some Repeats",
			s:                 "helloello",
			n:                 3,
			expDistinctNGrams: 5,
		},
		{
			name:              "short string",
			s:                 "hell",
			n:                 4,
			expDistinctNGrams: 1,
		},
		{
			name:              "too short string",
			s:                 "hel",
			n:                 4,
			expDistinctNGrams: 0,
		},
		{
			name:              "bad n - zero",
			s:                 "hel",
			n:                 0,
			expDistinctNGrams: 0,
			expErr:            true,
		},
		{
			name:              "bad n - negative",
			s:                 "hel",
			n:                 -1,
			expDistinctNGrams: 0,
			expErr:            true,
		},
	}

	for i, tc := range testCases {
		testDesc := fmt.Sprintf("test: %d: %s: NGrams('%s', %d): ",
			i, tc.name, tc.s, tc.n)
		m, err := strdist.NGrams(tc.s, tc.n)

		if tc.expErr {
			if err == nil {
				t.Error(testDesc + "should return an error but didn't")
			}
			continue
		} else if err != nil {
			t.Errorf(testDesc+"shouldn't return an error but did: %s", err)
		}

		if len(m) != tc.expDistinctNGrams {
			t.Errorf(testDesc+"should return %d n-grams but returned %d",
				tc.expDistinctNGrams, len(m))
		}

		totNGrams := 0
		for k, v := range m {
			if len(k) != tc.n {
				t.Errorf(testDesc+"some n-grams are not of length %d eg: '%s'",
					tc.n, k)
				break
			}
			totNGrams += v
		}

		expTotNGrams := len(tc.s) - tc.n + 1
		if expTotNGrams < 0 {
			expTotNGrams = 0
		}
		if totNGrams != expTotNGrams {
			t.Errorf(testDesc+"the string should contain %d n-grams not %d",
				expTotNGrams, totNGrams)
		}
	}

}

// ExampleNGrams demonstrates the use of NGrams(...)
func ExampleNGrams() {
	m, err := strdist.NGrams("bigbig", 2)
	if err != nil {
		fmt.Println("Unexpected error:", err)
	}
	fmt.Printf(
		"number of distinct n-grams: %d, number of instances of 'bi': %d\n",
		len(m), m["bi"])
	// Output: number of distinct n-grams: 3, number of instances of 'bi': 2
}

// TestNGramUnion tests the functions for constructing unions of n-grams
func TestNGramUnion(t *testing.T) {
	testCases := []struct {
		name           string
		m1, m2         map[string]int
		expLen         int
		expWeightedLen int
		expUnion       map[string]int
	}{
		{
			name: "two the same",
			m1: map[string]int{
				"ab": 1,
				"bc": 99,
			},
			m2: map[string]int{
				"ab": 1,
				"bc": 99,
			},
			expLen:         2,
			expWeightedLen: 200,
			expUnion: map[string]int{
				"ab": 2,
				"bc": 198,
			},
		},
		{
			name: "two different",
			m1: map[string]int{
				"ab": 1,
				"bc": 99,
			},
			m2: map[string]int{
				"cd": 1,
				"ef": 99,
			},
			expLen:         4,
			expWeightedLen: 200,
			expUnion: map[string]int{
				"ab": 1,
				"bc": 99,
				"cd": 1,
				"ef": 99,
			},
		},
		{
			name: "one empty",
			m1:   map[string]int{},
			m2: map[string]int{
				"cd": 1,
				"ef": 99,
			},
			expLen:         2,
			expWeightedLen: 100,
			expUnion: map[string]int{
				"cd": 1,
				"ef": 99,
			},
		},
		{
			name: "one nil",
			m2: map[string]int{
				"cd": 1,
				"ef": 99,
			},
			expLen:         2,
			expWeightedLen: 100,
			expUnion: map[string]int{
				"cd": 1,
				"ef": 99,
			},
		},
	}

	for i, tc := range testCases {
		testID := fmt.Sprintf("test %d: %s", i, tc.name)
		m1, m2 := tc.m1, tc.m2
		for _, order := range []string{"m1/m2", "m2/m1"} {
			if order == "m2/m1" {
				m1, m2 = m2, m1
			}

			m := strdist.NGramUnion(m1, m2)
			if len(m) != tc.expLen {
				t.Errorf("%s (%s): the length should have been %d but was %d",
					testID, order, tc.expLen, len(m))
			}

			calcLen := strdist.NGramLenUnion(m1, m2)
			if len(m) != calcLen {
				t.Errorf("%s (%s):"+
					" the length from NGramLenUnion should have been"+
					" %d but was %d",
					testID, order, tc.expLen, calcLen)
			}

			calcLen = strdist.NGramWeightedLenUnion(m1, m2)
			if tc.expWeightedLen != calcLen {
				t.Errorf("%s (%s):"+
					" the length from NGramWeightedLenUnion should have been"+
					" %d but was %d",
					testID, order, tc.expWeightedLen, calcLen)
			}

			if !strdist.NGramsEqual(m, tc.expUnion) {
				t.Errorf("%s (%s): the union was not as expected: %v",
					testID, order, m)
			}
		}
	}
}

// TestNGramIntersection tests the functions for constructing intersections
// of n-grams
func TestNGramIntersection(t *testing.T) {
	testCases := []struct {
		name            string
		m1, m2          map[string]int
		expLen          int
		expWeightedLen  int
		expIntersection map[string]int
	}{
		{
			name: "two the same",
			m1: map[string]int{
				"ab": 1,
				"bc": 99,
			},
			m2: map[string]int{
				"ab": 2,
				"bc": 98,
			},
			expLen:         2,
			expWeightedLen: 99,
			expIntersection: map[string]int{
				"ab": 1,
				"bc": 98,
			},
		},
		{
			name: "one in common",
			m1: map[string]int{
				"ab": 1,
				"bc": 99,
			},
			m2: map[string]int{
				"ab": 2,
				"cb": 99,
			},
			expLen:         1,
			expWeightedLen: 1,
			expIntersection: map[string]int{
				"ab": 1,
			},
		},
		{
			name: "two different",
			m1: map[string]int{
				"ab": 1,
				"bc": 99,
			},
			m2: map[string]int{
				"cd": 1,
				"ef": 99,
			},
			expLen:          0,
			expWeightedLen:  0,
			expIntersection: map[string]int{},
		},
		{
			name: "one empty",
			m1:   map[string]int{},
			m2: map[string]int{
				"cd": 1,
				"ef": 99,
			},
			expLen:          0,
			expWeightedLen:  0,
			expIntersection: map[string]int{},
		},
		{
			name: "one nil",
			m2: map[string]int{
				"cd": 1,
				"ef": 99,
			},
			expLen:          0,
			expWeightedLen:  0,
			expIntersection: map[string]int{},
		},
	}

	for i, tc := range testCases {
		testID := fmt.Sprintf("test %d: %s", i, tc.name)
		m1, m2 := tc.m1, tc.m2
		for _, order := range []string{"m1/m2", "m2/m1"} {
			if order == "m2/m1" {
				m1, m2 = m2, m1
			}

			m := strdist.NGramIntersection(m1, m2)
			if len(m) != tc.expLen {
				t.Logf("%s (%s) :\n", testID, order)
				t.Errorf("\t:: the length should have been %d but was %d",
					tc.expLen, len(m))
			}

			calcLen := strdist.NGramLenIntersection(m1, m2)
			if len(m) != calcLen {
				t.Logf("%s (%s) :\n", testID, order)
				t.Errorf("\t: NGramLenIntersection: expected len: %d got: %d",
					tc.expLen, calcLen)
			}

			calcLen = strdist.NGramWeightedLenIntersection(m1, m2)
			if tc.expWeightedLen != calcLen {
				t.Logf("%s (%s) :\n", testID, order)
				t.Errorf(
					"\t: NGramWeightedLenIntersection expected len: %d got: %d",
					tc.expWeightedLen, calcLen)
			}

			if !strdist.NGramsEqual(m, tc.expIntersection) {
				t.Logf("%s (%s) :\n", testID, order)
				t.Errorf("\t: bad intersection: expected: %v got: %v",
					tc.expIntersection, m)
			}
		}
	}
}

// TestNGramsEqual tests the NGramsEqual function
func TestNGramsEqual(t *testing.T) {
	testCases := []struct {
		name     string
		m1, m2   map[string]int
		expEqual bool
	}{
		{
			name: "both identical",
			m1: map[string]int{
				"ab": 1,
				"bc": 2,
			},
			m2: map[string]int{
				"ab": 1,
				"bc": 2,
			},
			expEqual: true,
		},
		{
			name: "count differs",
			m1: map[string]int{
				"ab": 1,
				"bc": 2,
			},
			m2: map[string]int{
				"ab": 1,
				"bc": 1,
			},
		},
		{
			name: "length differs",
			m1: map[string]int{
				"ab": 1,
				"bc": 2,
			},
			m2: map[string]int{
				"ab": 1,
				"bc": 2,
				"cd": 3,
			},
		},
		{
			name: "keys differ",
			m1: map[string]int{
				"ab": 1,
				"bc": 2,
			},
			m2: map[string]int{
				"ba": 1,
				"bc": 2,
			},
		},
	}

	for i, tc := range testCases {
		m1, m2 := tc.m1, tc.m2
		for _, order := range []string{"m1/m2", "m2/m1"} {
			if order == "m2/m1" {
				m1, m2 = m2, m1
			}

			if strdist.NGramsEqual(m1, m2) != tc.expEqual {
				t.Errorf("test %d: %s (%s): failed",
					i, tc.name, order)
			}
		}
	}
}

// TestJaccard tests the Jaccard functions
func TestJaccard(t *testing.T) {
	testCases := []struct {
		name           string
		m1, m2         map[string]int
		expVal         float64
		expWeightedVal float64
	}{
		{
			name: "2 in m1 and the same in m2 plus one more",
			m1: map[string]int{
				"ab": 1,
				"bc": 2,
			},
			m2: map[string]int{
				"ab": 1,
				"bc": 2,
				"cd": 3,
			},
			expVal:         0.666666667,
			expWeightedVal: 0.33333333,
		},
		{
			name:           "both empty",
			m1:             map[string]int{},
			m2:             map[string]int{},
			expVal:         1.0,
			expWeightedVal: 1.0,
		},
		{
			name:           "both nil",
			expVal:         1.0,
			expWeightedVal: 1.0,
		},
	}

	for i, tc := range testCases {
		testID := fmt.Sprintf("test %d: %s", i, tc.name)
		ji := strdist.JaccardIndex(tc.m1, tc.m2)
		if !mathutil.AlmostEqual(ji, tc.expVal, 0.00001) {
			t.Errorf("%s : the returned index should have been"+
				" %9.7f but was %9.7f",
				testID, tc.expVal, ji)
		}

		ji = strdist.Jaccard(tc.m1, tc.m2)
		if !mathutil.AlmostEqual(ji, 1.0-tc.expVal, 0.00001) {
			t.Errorf("%s : the returned index should have been"+
				" %9.7f but was %9.7f",
				testID, 1.0-tc.expVal, ji)
		}

		wji := strdist.WeightedJaccardIndex(tc.m1, tc.m2)
		if !mathutil.AlmostEqual(wji, tc.expWeightedVal, 0.00001) {
			t.Errorf("%s (weighted) : the returned index should have been"+
				" %9.7f but was %9.7f",
				testID, tc.expWeightedVal, wji)
		}

		wji = strdist.WeightedJaccard(tc.m1, tc.m2)
		if !mathutil.AlmostEqual(wji, 1.0-tc.expWeightedVal, 0.00001) {
			t.Errorf("%s (weighted) : the returned index should have been"+
				" %9.7f but was %9.7f",
				testID, 1.0-tc.expWeightedVal, wji)
		}
	}

}

// TestOverlapCoefficient tests the OverlapCoefficient functions
func TestOverlapCoefficient(t *testing.T) {
	testCases := []struct {
		name           string
		m1, m2         map[string]int
		expVal         float64
		expWeightedVal float64
	}{
		{
			name: "m1 is distinct from m2",
			m1: map[string]int{
				"ab": 1,
				"bc": 2,
				"xy": 3,
				"yz": 4,
			},
			m2: map[string]int{
				"cd": 1,
				"ef": 2,
				"gh": 3,
			},
			expVal:         0.0,
			expWeightedVal: 0.0,
		},
		{
			name: "m1 is a subset of m2",
			m1: map[string]int{
				"ab": 1,
				"bc": 2,
			},
			m2: map[string]int{
				"ab": 1,
				"bc": 2,
				"cd": 3,
			},
			expVal:         1.0,
			expWeightedVal: 1.0,
		},
		{
			name: "m1 and m2 overlap",
			m1: map[string]int{
				"ab": 1,
				"bc": 2,
				"cd": 4,
			},
			m2: map[string]int{
				"ab": 1,
				"bc": 2,
				"ef": 3,
			},
			expVal:         0.66666666667,
			expWeightedVal: 0.5,
		},
		{
			name:           "both empty",
			m1:             map[string]int{},
			m2:             map[string]int{},
			expVal:         1.0,
			expWeightedVal: 1.0,
		},
		{
			name:           "both nil",
			expVal:         1.0,
			expWeightedVal: 1.0,
		},
	}

	for i, tc := range testCases {
		testID := fmt.Sprintf("test %d: %s", i, tc.name)
		oc := strdist.OverlapCoefficient(tc.m1, tc.m2)
		if !mathutil.AlmostEqual(oc, tc.expVal, 0.00001) {
			t.Errorf("%s : the returned coefficient should have been"+
				" %9.7f but was %9.7f",
				testID, tc.expVal, oc)
		}
		woc := strdist.WeightedOverlapCoefficient(tc.m1, tc.m2)
		if !mathutil.AlmostEqual(woc, tc.expWeightedVal, 0.00001) {
			t.Errorf(
				"%s (weighted) : the returned coefficient should have been"+
					" %9.7f but was %9.7f",
				testID, tc.expWeightedVal, woc)
		}
	}

}
