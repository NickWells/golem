package fileparser

import "fmt"

// Stats records various details of the operation of the FileParser
type Stats struct {
	filesVisited int
	linesRead    int
	linesParsed  int
}

// String reports the contents of a Stats object
func (s Stats) String() string {
	return fmt.Sprintf("files: %3d lines read: %5d parsed: %5d",
		s.filesVisited, s.linesRead, s.linesParsed)
}

// FilesVisited returns the number of files visited by the FileParser during
// the last call to Parse
func (s Stats) FilesVisited() int { return s.filesVisited }

// LinesRead returns the number of lines resd by the FileParser during the last
// call to Parse
func (s Stats) LinesRead() int { return s.linesRead }

// LinesParsed returns the number of lines parsed by the FileParser (the number
// of lines for which ParseLine was called) during the last call to Parse
func (s Stats) LinesParsed() int { return s.linesParsed }
