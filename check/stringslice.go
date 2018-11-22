package check

import "fmt"

// StringSlice is the type of a check function for a slice of strings. It
// takes a slice of strings as a parameter and returns an error or nil if
// there is no error
type StringSlice func(v []string) error

// StringSliceStringCheck returns a check function that checks that every
// member of the list matches the supplied String check func
func StringSliceStringCheck(sc String) StringSlice {
	return func(v []string) error {
		for i, s := range v {
			if err := sc(s); err != nil {
				return fmt.Errorf(
					"list entry: %d (%s) does not satisfy the check: %s",
					i, s, err)
			}
		}
		return nil
	}
}

// StringSliceLenEQ returns a check function that checks that the length of
// the list equals the supplied value
func StringSliceLenEQ(i int) StringSlice {
	return func(v []string) error {
		if len(v) != i {
			return fmt.Errorf("the length of the list (%d) must equal %d",
				len(v), i)
		}
		return nil
	}
}

// StringSliceLenGT returns a check function that checks that the length of
// the list is greater than the supplied value
func StringSliceLenGT(i int) StringSlice {
	return func(v []string) error {
		if len(v) <= i {
			return fmt.Errorf(
				"the length of the list (%d) must be greater than %d",
				len(v), i)
		}
		return nil
	}
}

// StringSliceLenLT returns a check function that checks that the length of
// the list is less than the supplied value
func StringSliceLenLT(i int) StringSlice {
	return func(v []string) error {
		if len(v) >= i {
			return fmt.Errorf(
				"the length of the list (%d) must be less than %d",
				len(v), i)
		}
		return nil
	}
}

// StringSliceLenBetween returns a check function that checks that the length
// of the list is between the two supplied values (inclusive)
func StringSliceLenBetween(low, high int) StringSlice {
	if low >= high {
		panic(fmt.Sprintf("Impossible checks passed to StringSliceLenBetween: "+
			"the lower limit (%d) should be less than the upper limit (%d)",
			low, high))
	}

	return func(v []string) error {
		if len(v) < low {
			return fmt.Errorf(
				"the length of the list (%d) must be between %d and %d"+
					" - too short",
				len(v), low, high)
		}
		if len(v) > high {
			return fmt.Errorf(
				"the length of the list (%d) must be between %d and %d"+
					" - too long",
				len(v), low, high)
		}
		return nil
	}
}

// StringSliceNoDups checks that the list contains no duplicates
func StringSliceNoDups(v []string) error {
	dupMap := make(map[string]int)
	for i, s := range v {
		if dup, ok := dupMap[s]; ok {
			return fmt.Errorf(
				"list entries: %d and %d are duplicates, both are: %s",
				dup, i, s)
		}
		dupMap[s] = i
	}
	return nil
}

// StringSliceOr returns a function that will check that the value, when
// passed to each of the check funcs in turn, passes at least one of them
func StringSliceOr(chkFuncs ...StringSlice) StringSlice {
	return func(v []string) error {
		compositeErr := ""
		sep := "("

		for _, cf := range chkFuncs {
			err := cf(v)
			if err == nil {
				return nil
			}

			compositeErr += sep + err.Error()
			sep = " OR "
		}
		return fmt.Errorf("%s)", compositeErr)
	}
}

// StringSliceAnd returns a function that will check that the value, when
// passed to each of the check funcs in turn, passes all of them
func StringSliceAnd(chkFuncs ...StringSlice) StringSlice {
	return func(v []string) error {
		for _, cf := range chkFuncs {
			err := cf(v)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
