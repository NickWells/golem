package check_test

import (
	"fmt"
	"github.com/nickwells/golem/check"
	"github.com/nickwells/golem/testhelper"
	"strings"
	"testing"
)

func TestStringSlice(t *testing.T) {
	testCases := []struct {
		name           string
		checkFunc      check.StringSlice
		val            []string
		errExpected    bool
		errMustContain []string
	}{
		{
			name:        "NoDups - none - empty",
			checkFunc:   check.StringSliceNoDups,
			val:         []string{},
			errExpected: false,
		},
		{
			name:        "NoDups - none - 1 entry",
			checkFunc:   check.StringSliceNoDups,
			val:         []string{"a"},
			errExpected: false,
		},
		{
			name:        "NoDups - none",
			checkFunc:   check.StringSliceNoDups,
			val:         []string{"a", "b", "c", "d", "e"},
			errExpected: false,
		},
		{
			name:        "NoDups - has duplicates",
			checkFunc:   check.StringSliceNoDups,
			val:         []string{"a", "b", "b", "c", "d"},
			errExpected: true,
			errMustContain: []string{
				"list entries:",
				"are duplicates, both are: ",
			},
		},
		{
			name:        "StringCheck - LT - all good",
			checkFunc:   check.StringSliceStringCheck(check.StringLenLT(10)),
			val:         []string{"a", "b", "c", "c", "e"},
			errExpected: false,
		},
		{
			name:        "StringCheck - LT - all bad",
			checkFunc:   check.StringSliceStringCheck(check.StringLenGT(10)),
			val:         []string{"a", "b", "c", "c", "e"},
			errExpected: true,
			errMustContain: []string{
				"list entry:",
				"does not satisfy the check: ",
				"greater than",
			},
		},
		{
			name:        "LenLT: 1 < 2",
			checkFunc:   check.StringSliceLenLT(2),
			val:         []string{"a"},
			errExpected: false,
		},
		{
			name:        "LenLT: 1 !< 1",
			checkFunc:   check.StringSliceLenLT(1),
			val:         []string{"a"},
			errExpected: true,
			errMustContain: []string{
				"the length of the list",
				"must be less than",
			},
		},
		{
			name:        "LenLT: 2 !< 1",
			checkFunc:   check.StringSliceLenLT(1),
			val:         []string{"a", "b"},
			errExpected: true,
			errMustContain: []string{
				"the length of the list",
				"must be less than",
			},
		},
		{
			name:        "LenEQ: 1 == 1",
			checkFunc:   check.StringSliceLenEQ(1),
			val:         []string{"a"},
			errExpected: false,
		},
		{
			name:        "LenEQ: 1 != 2",
			checkFunc:   check.StringSliceLenEQ(1),
			val:         []string{"a", "b"},
			errExpected: true,
			errMustContain: []string{
				"the length of the list",
				"must equal",
			},
		},
		{
			name:        "LenEQ: 0 !=1",
			checkFunc:   check.StringSliceLenEQ(1),
			val:         []string{},
			errExpected: true,
			errMustContain: []string{
				"the length of the list",
				"must equal",
			},
		},
		{
			name:        "LenGT: 2 > 1",
			checkFunc:   check.StringSliceLenGT(1),
			val:         []string{"a", "b"},
			errExpected: false,
		},
		{
			name:        "LenGT: 1 !< 1",
			checkFunc:   check.StringSliceLenGT(1),
			val:         []string{"a"},
			errExpected: true,
			errMustContain: []string{
				"the length of the list",
				"must be greater than",
			},
		},
		{
			name:        "LenGT: 2 !< 1",
			checkFunc:   check.StringSliceLenGT(2),
			val:         []string{"a"},
			errExpected: true,
			errMustContain: []string{
				"the length of the list",
				"must be greater than",
			},
		},
		{
			name:        "LenBetween: 1 <= 2 <= 3",
			checkFunc:   check.StringSliceLenBetween(1, 3),
			val:         []string{"a", "b"},
			errExpected: false,
		},
		{
			name:        "LenBetween: 1 <= 1 <= 3",
			checkFunc:   check.StringSliceLenBetween(1, 3),
			val:         []string{"a"},
			errExpected: false,
		},
		{
			name:        "LenBetween: 1 <= 3 <= 3",
			checkFunc:   check.StringSliceLenBetween(1, 3),
			val:         []string{"a", "b", "c"},
			errExpected: false,
		},
		{
			name:        "LenBetween: 1 !<= 0 <= 3",
			checkFunc:   check.StringSliceLenBetween(1, 3),
			val:         []string{},
			errExpected: true,
			errMustContain: []string{
				"the length of the list",
				"must be between",
				"too short",
			},
		},
		{
			name:        "LenBetween: 1 <= 4 !<= 3",
			checkFunc:   check.StringSliceLenBetween(1, 3),
			val:         []string{"a", "b", "c", "4"},
			errExpected: true,
			errMustContain: []string{
				"the length of the list",
				"must be between",
				"too long",
			},
		},
		{
			name: "Or: 1 > 2 , 1 > 3, 1 < 3",
			checkFunc: check.StringSliceOr(
				check.StringSliceLenGT(2),
				check.StringSliceLenGT(3),
				check.StringSliceLenLT(3),
			),
			val:         []string{"a"},
			errExpected: false,
		},
		{
			name: "Or: 7 > 8, 7 < 6, 7 divides 60",
			checkFunc: check.StringSliceOr(
				check.StringSliceLenGT(8),
				check.StringSliceLenLT(6),
				check.StringSliceLenEQ(60),
			),
			val:         []string{"a", "b", "c", "d", "e", "f", "g"},
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
			checkFunc: check.StringSliceAnd(
				check.StringSliceLenGT(2),
				check.StringSliceLenGT(3),
				check.StringSliceLenLT(6),
			),
			val:         []string{"a", "b", "c", "d", "e"},
			errExpected: false,
		},
		{
			name: "And: 11 > 8, 11 < 10",
			checkFunc: check.StringSliceAnd(
				check.StringSliceLenGT(8),
				check.StringSliceLenLT(10),
			),
			val: []string{
				"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"},
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

func panicSafeTestStringSliceLenBetween(t *testing.T, lowerVal, upperVal int) (panicked bool, panicVal interface{}) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicVal = r
		}
	}()
	check.StringSliceLenBetween(lowerVal, upperVal)
	return panicked, panicVal
}

func TestStringSliceLenBetweenPanic(t *testing.T) {
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
				"Impossible checks passed to StringSliceLenBetween: ",
				"the lower limit",
				"should be less than the upper limit",
			},
		},
	}

	for i, tc := range testCases {
		testName := fmt.Sprintf("%d: %s", i, tc.name)
		panicked, panicVal := panicSafeTestStringSliceLenBetween(t, tc.lower, tc.upper)
		testhelper.PanicCheckString(t, testName,
			panicked, tc.panicExpected,
			panicVal, tc.panicMustContain)
	}

}
