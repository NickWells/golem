package mathutil_test

import (
	"fmt"
	"github.com/NickWells/golem/mathutil"
	"testing"
)

func TestMinMaxOf(t *testing.T) {
	testCases := []struct {
		name     string
		vals     []float64
		panicExp bool
		expMin   float64
		expMax   float64
	}{
		{
			name:   "all good - one val",
			vals:   []float64{1.0},
			expMin: 1.0,
			expMax: 1.0,
		},
		{
			name:   "all good - multiple vals",
			vals:   []float64{-1.0, 2.0, 3.0},
			expMin: -1.0,
			expMax: 3.0,
		},
		{
			name:   "all good - multiple vals, out of order",
			vals:   []float64{2.0, 3.0, -1.0},
			expMin: -1.0,
			expMax: 3.0,
		},
		{
			name:     "panic expected",
			vals:     []float64{},
			panicExp: true,
		},
	}

	for i, tc := range testCases {
		testID := fmt.Sprintf("test %d: %s", i, tc.name)

		v, panicked, panicVal := panicSafeFloat64(mathutil.MinOf, tc.vals)
		if panicOK(t, testID+": MinOf", panicked, tc.panicExp, panicVal) {
			if v != tc.expMin {
				t.Logf("%s :\n", testID)
				t.Errorf("\t: min value should have been %v, not %v\n",
					tc.expMin, v)
			}
		}

		v, panicked, panicVal = panicSafeFloat64(mathutil.MaxOf, tc.vals)
		if panicOK(t, testID+": MaxOf", panicked, tc.panicExp, panicVal) {
			if v != tc.expMax {
				t.Logf("%s :\n", testID)
				t.Errorf("\t: max value should have been %v, not %v\n",
					tc.expMax, v)
			}
		}
	}

}

func TestMinMaxOfInt(t *testing.T) {
	testCases := []struct {
		name     string
		vals     []int
		panicExp bool
		expMin   int
		expMax   int
	}{
		{
			name:   "all good - one val",
			vals:   []int{1},
			expMin: 1,
			expMax: 1,
		},
		{
			name:   "all good - multiple vals",
			vals:   []int{-1, 2, 3},
			expMin: -1,
			expMax: 3,
		},
		{
			name:   "all good - multiple vals, out of order",
			vals:   []int{2, 3, -1},
			expMin: -1,
			expMax: 3,
		},
		{
			name:     "panic expected",
			vals:     []int{},
			panicExp: true,
		},
	}

	for i, tc := range testCases {
		testID := fmt.Sprintf("test %d: %s", i, tc.name)

		v, panicked, panicVal := panicSafeInt(mathutil.MinOfInt, tc.vals)
		if panicOK(t, testID+": MinOfInt", panicked, tc.panicExp, panicVal) {
			if v != tc.expMin {
				t.Logf("%s :\n", testID)
				t.Errorf("\t: min value should have been %v, not %v\n",
					tc.expMin, v)
			}
		}

		v, panicked, panicVal = panicSafeInt(mathutil.MaxOfInt, tc.vals)
		if panicOK(t, testID+": MaxOfInt", panicked, tc.panicExp, panicVal) {
			if v != tc.expMax {
				t.Logf("%s :\n", testID)
				t.Errorf("\t: max value should have been %v, not %v\n",
					tc.expMax, v)
			}
		}
	}

}

// panicOK will check that the panic was as expected and return true if all was as expected, otherwise it will report the problem and return false
func panicOK(t *testing.T, name string, panicked, expected bool, pVal interface{}) bool {
	t.Helper()

	if panicked && !expected {
		t.Logf("%s :\n", name)
		t.Errorf("\t: unexpected panic: %v\n", pVal)
		return false
	}

	if !panicked && expected {
		t.Logf("%s :\n", name)
		t.Errorf("\t: panic expected but not seen\n")
		return false
	}
	return true
}

// panicSafeFloat64 calls the passed function with the supplied vals, it
// recovers from any panic and returns whether or not a panic was detected
// and the panic value
func panicSafeFloat64(f func(...float64) float64, vals []float64) (result float64, panicked bool, panicVal interface{}) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicVal = r
		}
	}()

	result = f(vals...)
	return
}

// panicSafeInt calls the passed function with the supplied vals, it
// recovers from any panic and returns whether or not a panic was detected
// and the panic value
func panicSafeInt(f func(...int) int, vals []int) (result int, panicked bool, panicVal interface{}) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicVal = r
		}
	}()

	result = f(vals...)
	return
}
