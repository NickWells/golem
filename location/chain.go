package location

// LocChain represents a chain of locations. This might be used to record a
// string of included files.
type LocChain []L

// NewChain returns a new (empty but non-nil) instance of a LocChain
func NewChain() LocChain {
	return make(LocChain, 0)
}

// String provides a string representing a LocChain
func (lc LocChain) String() string {
	var chainDesc string
	sep := ""
	for _, l := range lc {
		chainDesc += sep + l.String()
		sep = " --> "
	}
	return chainDesc
}
