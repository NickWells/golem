package check_test

import (
	"fmt"
	"github.com/NickWells/golem/check"
	"github.com/NickWells/golem/testhelper"
	"regexp"
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	testCases := []struct {
		name           string
		checkFunc      check.String
		val            string
		errExpected    bool
		errMustContain []string
	}{
		{
			name:        "LenEQ: 2 == 2",
			checkFunc:   check.StringLenEQ(2),
			val:         "ab",
			errExpected: false,
		},
		{
			name:        "LenEQ: 1 != 2",
			checkFunc:   check.StringLenEQ(2),
			val:         "a",
			errExpected: true,
			errMustContain: []string{
				"the length of the value",
				"must equal",
			},
		}, {
			name:        "LenLT: 1 < 2",
			checkFunc:   check.StringLenLT(2),
			val:         "a",
			errExpected: false,
		},
		{
			name:        "LenLT: 1 !< 1",
			checkFunc:   check.StringLenLT(1),
			val:         "a",
			errExpected: true,
			errMustContain: []string{
				"the length of the value",
				"must be less than",
			},
		},
		{
			name:        "LenLT: 2 !< 1",
			checkFunc:   check.StringLenLT(1),
			val:         "ab",
			errExpected: true,
			errMustContain: []string{
				"the length of the value",
				"must be less than",
			},
		},
		{
			name:        "LenGT: 2 > 1",
			checkFunc:   check.StringLenGT(1),
			val:         "ab",
			errExpected: false,
		},
		{
			name:        "LenGT: 1 !< 1",
			checkFunc:   check.StringLenGT(1),
			val:         "a",
			errExpected: true,
			errMustContain: []string{
				"the length of the value",
				"must be greater than",
			},
		},
		{
			name:        "LenGT: 2 !< 1",
			checkFunc:   check.StringLenGT(2),
			val:         "a",
			errExpected: true,
			errMustContain: []string{
				"the length of the value",
				"must be greater than",
			},
		},
		{
			name:        "Between: 1 <= 2 <= 3",
			checkFunc:   check.StringLenBetween(1, 3),
			val:         "ab",
			errExpected: false,
		},
		{
			name:        "Between: 1 <= 1 <= 3",
			checkFunc:   check.StringLenBetween(1, 3),
			val:         "a",
			errExpected: false,
		},
		{
			name:        "Between: 1 <= 3 <= 3",
			checkFunc:   check.StringLenBetween(1, 3),
			val:         "abc",
			errExpected: false,
		},
		{
			name:        "Between: 1 !<= 0 <= 3",
			checkFunc:   check.StringLenBetween(1, 3),
			val:         "",
			errExpected: true,
			errMustContain: []string{
				"the length of the value",
				"must be between",
				" - too short",
			},
		},
		{
			name:        "Between: 1 <= 4 !<= 3",
			checkFunc:   check.StringLenBetween(1, 3),
			val:         "abcd",
			errExpected: true,
			errMustContain: []string{
				"the length of the value",
				"must be between",
				" - too long",
			},
		},
		{
			name: "Matches",
			checkFunc: check.StringMatchesPattern(
				regexp.MustCompile("^a[a-z]+d$"),
				"3 or more letters starting with an a and ending with d"),
			val:         "abcd",
			errExpected: false,
		},
		{
			name: "Matches",
			checkFunc: check.StringMatchesPattern(
				regexp.MustCompile("^a[a-z]+d$"),
				"3 or more letters starting with an a and ending with d"),
			val:         "xxx",
			errExpected: true,
			errMustContain: []string{
				"does not match the pattern",
			},
		},
		{
			name: "Or: len(\"a\") > 2 , len(\"a\") > 3, len(\"a\") < 3",
			checkFunc: check.StringOr(
				check.StringLenGT(2),
				check.StringLenGT(3),
				check.StringLenLT(3),
			),
			val:         "a",
			errExpected: false,
		},
		{
			name: "Or: len(\"ab\") > 3, len(\"ab\") < 1",
			checkFunc: check.StringOr(
				check.StringLenGT(3),
				check.StringLenLT(1),
			),
			val:         "ab",
			errExpected: true,
			errMustContain: []string{
				"must be greater than",
				"must be less than",
				" OR ",
			},
		},
		{
			name: "And: len(\"abcd\") > 2 , len(\"abcd\") > 3, len(\"abcd\") < 6",
			checkFunc: check.StringAnd(
				check.StringLenGT(2),
				check.StringLenGT(3),
				check.StringLenLT(6),
			),
			val:         "abcd",
			errExpected: false,
		},
		{
			name: "And: len(\"abcd\") > 1, len(\"abcd\") < 3",
			checkFunc: check.StringAnd(
				check.StringLenGT(1),
				check.StringLenLT(3),
			),
			val:         "abcd",
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

func panicSafeTestStringLenBetween(t *testing.T, lowerVal, upperVal int) (panicked bool, panicVal interface{}) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicVal = r
		}
	}()
	check.StringLenBetween(lowerVal, upperVal)
	return panicked, panicVal
}

func TestStringLenBetweenPanic(t *testing.T) {
	testCases := []struct {
		name             string
		lower            int
		upper            int
		panicExpected    bool
		panicMustContain []string
	}{
		{
			name:          "Between: 1, 3",
			lower:         1,
			upper:         3,
			panicExpected: false,
		},
		{
			name:          "Between: 4, 3",
			lower:         4,
			upper:         3,
			panicExpected: true,
			panicMustContain: []string{
				"Impossible checks passed to StringLenBetween: ",
				"the lower limit",
				"should be less than the upper limit",
			},
		},
	}

	for i, tc := range testCases {
		testName := fmt.Sprintf("%d: %s", i, tc.name)
		panicked, panicVal := panicSafeTestStringLenBetween(
			t, tc.lower, tc.upper)
		testhelper.PanicCheckString(t, testName,
			panicked, tc.panicExpected,
			panicVal, tc.panicMustContain)
	}

}
