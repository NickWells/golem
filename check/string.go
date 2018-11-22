package check

import (
	"fmt"
	"regexp"
)

// String is the type of a check function for a string. It takes a string as
// a parameter and returns an error or nil if the check passes
type String func(s string) error

// StringLenEQ returns a function that will check that the
// length of the string is equal to the limit
func StringLenEQ(limit int) String {
	return func(s string) error {
		if len(s) == limit {
			return nil
		}
		return fmt.Errorf("the length of the value (%d) must equal %d",
			len(s), limit)
	}
}

// StringLenLT returns a function that will check that the
// length of the string is less than the limit
func StringLenLT(limit int) String {
	return func(s string) error {
		if len(s) < limit {
			return nil
		}
		return fmt.Errorf(
			"the length of the value (%d) must be less than %d",
			len(s), limit)
	}
}

// StringLenGT returns a function that will check that the
// length of the string is less than the limit
func StringLenGT(limit int) String {
	return func(s string) error {
		if len(s) > limit {
			return nil
		}
		return fmt.Errorf(
			"the length of the value (%d) must be greater than %d",
			len(s), limit)
	}
}

// StringLenBetween returns a function that will check that the
// length of the string is between the two limits (inclusive)
func StringLenBetween(low, high int) String {
	if low >= high {
		panic(fmt.Sprintf("Impossible checks passed to StringLenBetween:"+
			" the lower limit (%d) should be less than the upper limit (%d)",
			low, high))
	}

	return func(s string) error {
		if len(s) < low {
			return fmt.Errorf(
				"the length of the value (%d) must be between %d and %d"+
					" - too short",
				len(s), low, high)
		}
		if len(s) > high {
			return fmt.Errorf(
				"the length of the value (%d) must be between %d and %d"+
					" - too long",
				len(s), low, high)
		}
		return nil
	}
}

// StringMatchesPattern returns a function that checks that the
// string matches the supplied regexp
func StringMatchesPattern(re *regexp.Regexp, reDesc string) String {
	return func(v string) error {
		if !re.MatchString(v) {
			return fmt.Errorf("%s does not match the pattern: %s",
				v, reDesc)
		}
		return nil
	}
}

// StringOr returns a function that will check that the value, when
// passed to each of the check funcs in turn, passes at least one of them
func StringOr(chkFuncs ...String) String {
	return func(s string) error {
		compositeErr := ""
		sep := "("

		for _, cf := range chkFuncs {
			err := cf(s)
			if err == nil {
				return nil
			}

			compositeErr += sep + err.Error()
			sep = " OR "
		}
		return fmt.Errorf("%s)", compositeErr)
	}
}

// StringAnd returns a function that will check that the value, when
// passed to each of the check funcs in turn, passes all of them
func StringAnd(chkFuncs ...String) String {
	return func(s string) error {
		for _, cf := range chkFuncs {
			err := cf(s)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
