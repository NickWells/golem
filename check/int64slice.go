package check

import "fmt"

// Int64Slice is the type of a check function for a slice of int64's. It
// takes a slice of int64's and returns an error or nil if the check passes
type Int64Slice func(v []int64) error

// Int64SliceNoDups checks that the list contains no duplicates
func Int64SliceNoDups(v []int64) error {
	dupMap := make(map[int64]int)
	for i, s := range v {
		if dup, ok := dupMap[s]; ok {
			return fmt.Errorf(
				"list entries: %d and %d are duplicates, both are: %d",
				dup, i, s)
		}
		dupMap[s] = i
	}
	return nil
}

// Int64SliceInt64Check returns a check function that checks that every member
// of the list matches the supplied Int64 check function
func Int64SliceInt64Check(sc Int64) Int64Slice {
	return func(v []int64) error {
		for i, s := range v {
			if err := sc(s); err != nil {
				return fmt.Errorf(
					"list entry: %d (%d) does not satisfy the check: %s",
					i, s, err)
			}
		}
		return nil
	}
}

// Int64SliceLenEQ returns a check function that checks that the length of
// the list equals the supplied value
func Int64SliceLenEQ(i int) Int64Slice {
	return func(v []int64) error {
		if len(v) != i {
			return fmt.Errorf("the length of the list (%d) must equal %d",
				len(v), i)
		}
		return nil
	}
}

// Int64SliceLenGT returns a check function that checks that the length of
// the list is greater than the supplied value
func Int64SliceLenGT(i int) Int64Slice {
	return func(v []int64) error {
		if len(v) <= i {
			return fmt.Errorf(
				"the length of the list (%d) must be greater than %d",
				len(v), i)
		}
		return nil
	}
}

// Int64SliceLenLT returns a check function that checks that the length of
// the list is less than the supplied value
func Int64SliceLenLT(i int) Int64Slice {
	return func(v []int64) error {
		if len(v) >= i {
			return fmt.Errorf(
				"the length of the list (%d) must be less than %d",
				len(v), i)
		}
		return nil
	}
}

// Int64SliceLenBetween returns a check function that checks that the length
// of the list is between the two supplied values (inclusive)
func Int64SliceLenBetween(low, high int) Int64Slice {
	if low >= high {
		panic(fmt.Sprintf("Impossible checks passed to Int64SliceLenBetween: "+
			"the lower limit (%d) should be less than the upper limit (%d)",
			low, high))
	}

	return func(v []int64) error {
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

// Int64SliceOr returns a function that will check that the value, when
// passed to each of the check funcs in turn, passes at least one of them
func Int64SliceOr(chkFuncs ...Int64Slice) Int64Slice {
	return func(i []int64) error {
		compositeErr := ""
		sep := "("

		for _, cf := range chkFuncs {
			err := cf(i)
			if err == nil {
				return nil
			}

			compositeErr += sep + err.Error()
			sep = " OR "
		}
		return fmt.Errorf("%s)", compositeErr)
	}
}

// Int64SliceAnd returns a function that will check that the value, when
// passed to each of the check funcs in turn, passes all of them
func Int64SliceAnd(chkFuncs ...Int64Slice) Int64Slice {
	return func(i []int64) error {
		for _, cf := range chkFuncs {
			err := cf(i)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
