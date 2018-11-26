package check_test

import (
	"fmt"
	"github.com/nickwells/golem/check"
	"github.com/nickwells/golem/testhelper"
	"strings"
	"testing"
)

func TestInt64Slice(t *testing.T) {
	testCases := []struct {
		name           string
		checkFunc      check.Int64Slice
		val            []int64
		errExpected    bool
		errMustContain []string
	}{
		{
			name:        "NoDups - none - empty",
			checkFunc:   check.Int64SliceNoDups,
			val:         []int64{},
			errExpected: false,
		},
		{
			name:        "NoDups - none - 1 entry",
			checkFunc:   check.Int64SliceNoDups,
			val:         []int64{1},
			errExpected: false,
		},
		{
			name:        "NoDups - none",
			checkFunc:   check.Int64SliceNoDups,
			val:         []int64{1, 2, 3, 4, 5},
			errExpected: false,
		},
		{
			name:        "NoDups - has duplicates",
			checkFunc:   check.Int64SliceNoDups,
			val:         []int64{1, 2, 3, 3, 5},
			errExpected: true,
			errMustContain: []string{
				"list entries:",
				"are duplicates, both are: ",
			},
		},
		{
			name:        "Int64Check - LT - all good",
			checkFunc:   check.Int64SliceInt64Check(check.Int64LT(10)),
			val:         []int64{1, 2, 3, 3, 5},
			errExpected: false,
		},
		{
			name:        "Int64Check - LT - all bad",
			checkFunc:   check.Int64SliceInt64Check(check.Int64GT(10)),
			val:         []int64{1, 2, 3, 3, 5},
			errExpected: true,
			errMustContain: []string{
				"list entry:",
				"does not satisfy the check: ",
				"greater than",
			},
		},
		{
			name:        "LenLT: 1 < 2",
			checkFunc:   check.Int64SliceLenLT(2),
			val:         []int64{1},
			errExpected: false,
		},
		{
			name:        "LenLT: 1 !< 1",
			checkFunc:   check.Int64SliceLenLT(1),
			val:         []int64{1},
			errExpected: true,
			errMustContain: []string{
				"the length of the list",
				"must be less than",
			},
		},
		{
			name:        "LenLT: 2 !< 1",
			checkFunc:   check.Int64SliceLenLT(1),
			val:         []int64{1, 2},
			errExpected: true,
			errMustContain: []string{
				"the length of the list",
				"must be less than",
			},
		},
		{
			name:        "LenEQ: 1 == 1",
			checkFunc:   check.Int64SliceLenEQ(1),
			val:         []int64{1},
			errExpected: false,
		},
		{
			name:        "LenEQ: 1 != 2",
			checkFunc:   check.Int64SliceLenEQ(1),
			val:         []int64{1, 2},
			errExpected: true,
			errMustContain: []string{
				"the length of the list",
				"must equal",
			},
		},
		{
			name:        "LenEQ: 0 !=1",
			checkFunc:   check.Int64SliceLenEQ(1),
			val:         []int64{},
			errExpected: true,
			errMustContain: []string{
				"the length of the list",
				"must equal",
			},
		},
		{
			name:        "LenGT: 2 > 1",
			checkFunc:   check.Int64SliceLenGT(1),
			val:         []int64{1, 2},
			errExpected: false,
		},
		{
			name:        "LenGT: 1 !< 1",
			checkFunc:   check.Int64SliceLenGT(1),
			val:         []int64{1},
			errExpected: true,
			errMustContain: []string{
				"the length of the list",
				"must be greater than",
			},
		},
		{
			name:        "LenGT: 2 !< 1",
			checkFunc:   check.Int64SliceLenGT(2),
			val:         []int64{1},
			errExpected: true,
			errMustContain: []string{
				"the length of the list",
				"must be greater than",
			},
		},
		{
			name:        "LenBetween: 1 <= 2 <= 3",
			checkFunc:   check.Int64SliceLenBetween(1, 3),
			val:         []int64{1, 2},
			errExpected: false,
		},
		{
			name:        "LenBetween: 1 <= 1 <= 3",
			checkFunc:   check.Int64SliceLenBetween(1, 3),
			val:         []int64{1},
			errExpected: false,
		},
		{
			name:        "LenBetween: 1 <= 3 <= 3",
			checkFunc:   check.Int64SliceLenBetween(1, 3),
			val:         []int64{1, 2, 3},
			errExpected: false,
		},
		{
			name:        "LenBetween: 1 !<= 0 <= 3",
			checkFunc:   check.Int64SliceLenBetween(1, 3),
			val:         []int64{},
			errExpected: true,
			errMustContain: []string{
				"the length of the list",
				"must be between",
				"too short",
			},
		},
		{
			name:        "LenBetween: 1 <= 4 !<= 3",
			checkFunc:   check.Int64SliceLenBetween(1, 3),
			val:         []int64{1, 2, 3, 4},
			errExpected: true,
			errMustContain: []string{
				"the length of the list",
				"must be between",
				"too long",
			},
		},
		{
			name: "Or: 1 > 2 , 1 > 3, 1 < 3",
			checkFunc: check.Int64SliceOr(
				check.Int64SliceLenGT(2),
				check.Int64SliceLenGT(3),
				check.Int64SliceLenLT(3),
			),
			val:         []int64{1},
			errExpected: false,
		},
		{
			name: "Or: 7 > 8, 7 < 6, 7 divides 60",
			checkFunc: check.Int64SliceOr(
				check.Int64SliceLenGT(8),
				check.Int64SliceLenLT(6),
				check.Int64SliceLenEQ(60),
			),
			val:         []int64{1, 2, 3, 4, 5, 6, 7},
			errExpected: true,
			errMustContain: []string{
				"must be greater than",
				"must be less than",
				"must equal",
				" OR ",
			},
		},
		{
			name: "And: 5 > 2 , 5 > 3, 5 < 6",
			checkFunc: check.Int64SliceAnd(
				check.Int64SliceLenGT(2),
				check.Int64SliceLenGT(3),
				check.Int64SliceLenLT(6),
			),
			val:         []int64{1, 2, 3, 4, 5},
			errExpected: false,
		},
		{
			name: "And: 11 > 8, 11 < 10",
			checkFunc: check.Int64SliceAnd(
				check.Int64SliceLenGT(8),
				check.Int64SliceLenLT(10),
			),
			val:         []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
			errExpected: true,
			errMustContain: []string{
				"must be less than",
			},
		},
	}

	for i, tc := range testCases {

		err := tc.checkFunc(tc.val)
		if err != nil {
			if !tc.errExpected {
				t.Logf("test %d: %s :\n", i, tc.name)
				t.Errorf("\t: there was an unexpected err: %s\n", err)
			} else {
				emsg := err.Error()
				reported := false
				for _, s := range tc.errMustContain {
					if !strings.Contains(emsg, s) {
						if !reported {
							t.Logf("test %d: %s :\n", i, tc.name)
							t.Logf("\t: Error: '%s'\n", emsg)
							reported = true
						}
						t.Errorf("\t: the error should have contained: '%s'\n",
							s)
					}
				}
			}
		} else if err == nil && tc.errExpected {
			t.Logf("test %d: %s :\n", i, tc.name)
			t.Errorf("\t: an error was expected but none was returned\n")
		}
	}

}

func panicSafeTestInt64SliceLenBetween(t *testing.T, lowerVal, upperVal int) (panicked bool, panicVal interface{}) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicVal = r
		}
	}()
	check.Int64SliceLenBetween(lowerVal, upperVal)
	return panicked, panicVal
}

func TestInt64SliceLenBetweenPanic(t *testing.T) {
	testCases := []struct {
		name             string
		lower            int
		upper            int
		panicExpected    bool
		panicMustContain []string
	}{
		{
			name:          "LenBetween: 1, 3",
			lower:         1,
			upper:         3,
			panicExpected: false,
		},
		{
			name:          "LenBetween: 4, 3",
			lower:         4,
			upper:         3,
			panicExpected: true,
			panicMustContain: []string{
				"Impossible checks passed to Int64SliceLenBetween: ",
				"the lower limit",
				"should be less than the upper limit",
			},
		},
	}

	for i, tc := range testCases {
		testName := fmt.Sprintf("%d: %s", i, tc.name)
		panicked, panicVal := panicSafeTestInt64SliceLenBetween(t, tc.lower, tc.upper)
		testhelper.PanicCheckString(t, testName,
			panicked, tc.panicExpected,
			panicVal, tc.panicMustContain)
	}

}
