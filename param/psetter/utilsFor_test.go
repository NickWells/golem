package psetter_test

import "github.com/NickWells/golem/param"

// panicSafeCheck runs the CheckSetter and catches any panic, returning true
// if a panic was caught
func panicSafeCheck(s param.Setter) (panicked bool, panicVal interface{}) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicVal = r
		}
	}()
	s.CheckSetter("test")
	return
}
