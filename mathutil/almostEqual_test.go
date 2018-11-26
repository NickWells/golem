package mathutil_test

import (
	"github.com/nickwells/golem/mathutil"
	"math"
	"testing"
)

// TestAlmostEqual tests the AlmostEqual function
func TestAlmostEqual(t *testing.T) {
	testCases := []struct {
		name          string
		a, b, epsilon float64
		expResult     bool
	}{
		{
			name:      "both positive, true",
			a:         1.2345,
			b:         1.23456,
			epsilon:   0.0001,
			expResult: true,
		},
		{
			name:      "both positive, false",
			a:         1.2345,
			b:         1.23456,
			epsilon:   0.00001,
			expResult: false,
		},
		{
			name:      "both negative, true",
			a:         -1.2345,
			b:         -1.23456,
			epsilon:   0.0001,
			expResult: true,
		},
		{
			name:      "both negative, false",
			a:         -1.2345,
			b:         -1.23456,
			epsilon:   0.00001,
			expResult: false,
		},
		{
			name:      "different sign, true",
			a:         -0.00000001,
			b:         0.00000001,
			epsilon:   0.00001,
			expResult: true,
		},
		{
			name:      "identical",
			a:         1.23456789,
			b:         1.23456789,
			epsilon:   0.00000001,
			expResult: true,
		},
		{
			name:      "max positive",
			a:         math.MaxFloat64,
			b:         math.MaxFloat64 - 1,
			epsilon:   0.00000000000001,
			expResult: true,
		},
		{
			name:      "max negative",
			a:         -math.MaxFloat64,
			b:         -math.MaxFloat64 + 1,
			epsilon:   0.00000000000001,
			expResult: true,
		},
	}

	for i, tc := range testCases {
		res := mathutil.AlmostEqual(tc.a, tc.b, tc.epsilon)
		if res != tc.expResult {
			t.Logf("test %d: %s :\n", i, tc.name)
			t.Errorf("\t: AlmostEqual(%.9f, %.9f, %.9f)"+
				" should have returned %v but didn't",
				tc.a, tc.b, tc.epsilon, tc.expResult)
		}
	}
}

func TestWithinNPercent(t *testing.T) {
	testCases := []struct {
		name      string
		a, b, pct float64
		expResult bool
	}{
		{
			name:      "both positive, true",
			a:         100.1,
			b:         100.0,
			pct:       0.15,
			expResult: true,
		},
		{
			name:      "both positive, false",
			a:         100.3,
			b:         100.0,
			pct:       0.15,
			expResult: false,
		},
		{
			name:      "both negative, true",
			a:         -100.1,
			b:         -100.0,
			pct:       0.15,
			expResult: true,
		},
		{
			name:      "both negative, false",
			a:         -100.3,
			b:         -100.0,
			pct:       0.15,
			expResult: false,
		},
		{
			name:      "identical positive",
			a:         100.0,
			b:         100.0,
			pct:       0.0,
			expResult: true,
		},
		{
			name:      "identical negative",
			a:         -100.0,
			b:         -100.0,
			pct:       0.0,
			expResult: true,
		},
	}

	for i, tc := range testCases {
		res := mathutil.WithinNPercent(tc.a, tc.b, tc.pct)
		if res != tc.expResult {
			t.Logf("test %d: %s :\n", i, tc.name)
			t.Logf("\t: testing whether %11.7f is within %5.1f%% of %11.7f\n",
				tc.a, tc.pct, tc.b)
			t.Logf("should have returned %v but didn't\n", tc.expResult)
			t.Errorf("\t: Failed\n")
		}
	}

}
