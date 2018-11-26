package check_test

import (
	"fmt"
	"github.com/nickwells/golem/check"
	"github.com/nickwells/golem/testhelper"
	"strings"
	"testing"
	"time"
)

func TestDuration(t *testing.T) {
	testCases := []struct {
		name           string
		checkFunc      check.Duration
		d              time.Duration
		errExpected    bool
		errMustContain []string
	}{
		{
			name:        "LT: 1 < 2",
			checkFunc:   check.DurationLT(2 * time.Second),
			d:           1 * time.Second,
			errExpected: false,
		},
		{
			name:           "LT: 1 !< 1",
			checkFunc:      check.DurationLT(1 * time.Second),
			d:              1 * time.Second,
			errExpected:    true,
			errMustContain: []string{"must be less than"},
		},
		{
			name:           "LT: 2 !< 1",
			checkFunc:      check.DurationLT(1 * time.Second),
			d:              2 * time.Second,
			errExpected:    true,
			errMustContain: []string{"must be less than"},
		},
		{
			name:        "GT: 2 > 1",
			checkFunc:   check.DurationGT(1 * time.Second),
			d:           2 * time.Second,
			errExpected: false,
		},
		{
			name:           "GT: 1 !> 1",
			checkFunc:      check.DurationGT(1 * time.Second),
			d:              1 * time.Second,
			errExpected:    true,
			errMustContain: []string{"must be greater than"},
		},
		{
			name:           "GT: 1 !> 2",
			checkFunc:      check.DurationGT(2 * time.Second),
			d:              1 * time.Second,
			errExpected:    true,
			errMustContain: []string{"must be greater than"},
		},
		{
			name:        "Between: 1 <= 2 <= 3",
			checkFunc:   check.DurationBetween(1*time.Second, 3*time.Second),
			d:           2 * time.Second,
			errExpected: false,
		},
		{
			name:        "Between: 1 <= 1 <= 3",
			checkFunc:   check.DurationBetween(1*time.Second, 3*time.Second),
			d:           1 * time.Second,
			errExpected: false,
		},
		{
			name:        "Between: 1 <= 3 <= 3",
			checkFunc:   check.DurationBetween(1*time.Second, 3*time.Second),
			d:           3 * time.Second,
			errExpected: false,
		},
		{
			name:        "Between: 1 !<= 0 <= 3",
			checkFunc:   check.DurationBetween(1*time.Second, 3*time.Second),
			d:           0 * time.Second,
			errExpected: true,
			errMustContain: []string{
				"the value",
				"must be between",
				" - too short",
			},
		},
		{
			name:        "Between: 1 <= 4 !<= 3",
			checkFunc:   check.DurationBetween(1*time.Second, 3*time.Second),
			d:           4 * time.Second,
			errExpected: true,
			errMustContain: []string{
				"the value",
				"must be between",
				" - too long",
			},
		},
	}

	for i, tc := range testCases {
		err := tc.checkFunc(tc.d)
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

func panicSafeTestDurationBetween(t *testing.T, ld, ud time.Duration) (panicked bool, panicVal interface{}) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicVal = r
		}
	}()
	check.DurationBetween(ld, ud)
	return panicked, panicVal
}

func TestDurationBetweenPanic(t *testing.T) {
	testCases := []struct {
		name             string
		lower            time.Duration
		upper            time.Duration
		panicExpected    bool
		panicMustContain []string
	}{
		{
			name:          "Between: 1, 3",
			lower:         1 * time.Second,
			upper:         3 * time.Second,
			panicExpected: false,
		},
		{
			name:          "Between: 4, 3",
			lower:         4 * time.Second,
			upper:         3 * time.Second,
			panicExpected: true,
			panicMustContain: []string{
				"Impossible checks passed to DurationBetween: ",
				"the lower limit",
				"should be less than the upper limit",
			},
		},
	}

	for i, tc := range testCases {
		testName := fmt.Sprintf("%d: %s", i, tc.name)
		panicked, panicVal := panicSafeTestDurationBetween(t, tc.lower, tc.upper)
		testhelper.PanicCheckString(t, testName,
			panicked, tc.panicExpected,
			panicVal, tc.panicMustContain)
	}

}
