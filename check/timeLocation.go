package check

import "time"

// TimeLocation is the type of a check function for a time.Location. It takes
// a pointer to a time.Location parameter and returns an error or nil if the
// check passes
type TimeLocation func(l *time.Location) error
