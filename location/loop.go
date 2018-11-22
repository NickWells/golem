package location

// HasLoop checks to see if the name already appears in the LocChain which
// records the chain of sources getting to this point. If a loop is detected
// then this returns true and a description of the path taken to get to the
// loop, otherwise it returns false and an empty string.
func (chain LocChain) HasLoop(name string) (bool, string) {
	var loopDesc string
	incl := " --> "
	var loopDetected bool
	for _, l := range chain {
		if l.name == name {
			loopDetected = true
			incl = " ==> "
		}
		loopDesc += l.String() + incl
	}
	if loopDetected {
		loopDesc += "Back to " + name
		return true, loopDesc
	}
	return false, ""
}
