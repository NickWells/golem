package location

// Err wraps a message and a location into a single struct
// this allows for a more structured error message
type Err struct {
	Msg string
	Loc L
}

func (e Err) Error() string {
	return e.Msg + " (at " + e.Loc.String() + ")"
}
