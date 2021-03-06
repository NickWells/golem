package psetter

import (
	"errors"
	"fmt"
	"github.com/nickwells/golem/check"
	"github.com/nickwells/golem/param"
	"strconv"
)

// Float64Setter allows you to specify a parameter that can be used to set an
// float64 value. You can also supply a check function that will validate
// the Value. There are some helper functions given below (called
// Float64Check...) which will return functions that can perform a few
// common checks. For instance you can ensure that the value is positive by
// setting one of the Checks to the value returned by
// Float64CheckGT(0)
type Float64Setter struct {
	Value  *float64
	Checks []check.Float64
}

// ValueReq returns param.Mandatory indicating that some value must follow
// the parameter
func (s Float64Setter) ValueReq() param.ValueReq { return param.Mandatory }

// Set (called when there is no following value) returns an error
func (s Float64Setter) Set(_ string) error {
	return errors.New("no number given (it should be followed by '=num')")
}

// SetWithVal (called when a value follows the parameter) checks that the value
// can be parsed to a float, if it cannot be parsed successfully it returns an
// error. If there is a check and the check is violated it returns an
// error. Only if the value is parsed successfully and the check is not
// violated is the Value set.
func (s Float64Setter) SetWithVal(_ string, paramVal string) error {
	v, err := strconv.ParseFloat(paramVal, 64)
	if err != nil {
		return fmt.Errorf("could not parse '%s' as a float value: %s",
			paramVal, err)
	}

	if len(s.Checks) != 0 {
		for _, check := range s.Checks {
			if check == nil {
				continue
			}

			err := check(v)
			if err != nil {
				return err
			}
		}
	}

	*s.Value = v
	return nil
}

// AllowedValues returns a string describing the allowed values
func (s Float64Setter) AllowedValues() string {
	rval := "any value that can be read as a number with a decimal place"
	if len(s.Checks) != 0 {
		rval += " subject to checks"
	}
	return rval
}

// CurrentValue returns the current setting of the parameter value
func (s Float64Setter) CurrentValue() string {
	return fmt.Sprintf("%v", *s.Value)
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s Float64Setter) CheckSetter(name string) {
	if s.Value == nil {
		panic(name + ": Float64Setter Check failed: the Value to be set is nil")
	}
}
