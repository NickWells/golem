package location

import "fmt"

// L represents a location - a named source and the index into that source.
// For instance, this might be used to represent a file and a line number or
// a database table and a row number. You can also record the content of the
// location but this is not done routinely and so is only reported if it has
// been set. Similarly you can add a note to the location and this is only
// reported if it is not empty.
type L struct {
	name       string
	idx        int64
	hasContent bool
	content    string
	note       string
}

// New returns a new instance of a location
func New(name string) *L {
	return &L{name: name}
}

// String provides the string representing a location
func (l L) String() string {
	var s string
	if l.note != "" {
		s = "[ " + l.note + " ]: "
	}
	s += fmt.Sprintf("%s:%d", l.name, l.idx)
	if l.hasContent {
		s += ": " + l.content
	}

	return s
}

// Incr increments the location index and marks the location as having no
// content
func (l *L) Incr() {
	l.idx++
	l.hasContent = false
}

// Idx returns the current index value
func (l L) Idx() int64 {
	return l.idx
}

// SetContent sets the content and marks the location as having content
func (l *L) SetContent(s string) {
	l.content = s
	l.hasContent = true
}

// SetNote sets the notes field on the location
func (l *L) SetNote(s string) {
	l.note = s
}

// Errorf constructs a location Err setting the Msg to the results of
// formatting the arguments and the Loc to the location
func (l L) Errorf(s string, args ...interface{}) Err {
	return Err{
		Msg: fmt.Sprintf(s, args...),
		Loc: l,
	}
}

// Error constructs a location Err setting the Msg to the passed value and
// the Loc to the location
func (l L) Error(s string) Err {
	return Err{
		Msg: s,
		Loc: l,
	}
}
