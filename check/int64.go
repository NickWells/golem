package check

import "fmt"

// Int64 is the type of the check function for int64 values. It takes an
// int64 parameter and returns an error
type Int64 func(i int64) error

// Int64EQ  returns a function that will check that the value is
// greater than the limit
func Int64EQ(limit int64) Int64 {
	return func(i int64) error {
		if i == limit {
			return nil
		}
		return fmt.Errorf("the value (%d) must be equal to %d", i, limit)
	}
}

// Int64GT  returns a function that will check that the value is
// greater than the limit
func Int64GT(limit int64) Int64 {
	return func(i int64) error {
		if i > limit {
			return nil
		}
		return fmt.Errorf("the value (%d) must be greater than %d", i, limit)
	}
}

// Int64LT  returns a function that will check that the value is less
// than the limit
func Int64LT(limit int64) Int64 {
	return func(i int64) error {
		if i < limit {
			return nil
		}
		return fmt.Errorf("the value (%d) must be less than %d", i, limit)
	}
}

// Int64Between  returns a function that will check that the value
// lies between the upper and lower limits (inclusive)
func Int64Between(low, high int64) Int64 {
	if low >= high {
		panic(fmt.Sprintf("Impossible checks passed to Int64Between:"+
			" the lower limit (%d) should be less than the upper limit (%d)",
			low, high))
	}

	return func(i int64) error {
		if i < low {
			return fmt.Errorf(
				"the value (%d) must be between %d and %d - too small",
				i, low, high)
		}
		if i > high {
			return fmt.Errorf(
				"the value (%d) must be between %d and %d - too big",
				i, low, high)
		}
		return nil
	}
}

// Int64Divides  returns a function that will check that the value
// is a divisor of d
func Int64Divides(d int64) Int64 {
	return func(i int64) error {
		if d%i == 0 {
			return nil
		}
		return fmt.Errorf("the value (%d) must be a divisor of %d", i, d)
	}
}

// Int64IsAMultiple returns a function that will check that the value
// is a multiple of d
func Int64IsAMultiple(d int64) Int64 {
	return func(i int64) error {
		if i%d == 0 {
			return nil
		}
		return fmt.Errorf("the value (%d) must be a multiple of %d", i, d)
	}
}

// Int64Or returns a function that will check that the value, when
// passed to each of the check funcs in turn, passes at least one of them
func Int64Or(chkFuncs ...Int64) Int64 {
	return func(i int64) error {
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

// Int64And returns a function that will check that the value, when
// passed to each of the check funcs in turn, passes all of them
func Int64And(chkFuncs ...Int64) Int64 {
	return func(i int64) error {
		for _, cf := range chkFuncs {
			err := cf(i)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
