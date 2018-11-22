package fileparser

import (
	"github.com/NickWells/golem/location"
	"io"
	"os"
)

// LineParser is an interface defining the ParseLine function. This is used to
// support various different line parsers which can be used by the FileParser
// Parse method
type LineParser interface {
	ParseLine(line string, loc *location.L) error
}

// NullParser does nothing - it just discards it's inputs
//
// It could be useful to check the files to be parsed - it will still
// check for inaccessible files and include directives that aren't
// followed by a filename
type NullParser struct{}

// ParseLine for the NullParser does nothing
func (np NullParser) ParseLine(_ string, _ *location.L) error {
	return nil
}

// EchoParser will just write everything it is passed to the Writer as a
// byte slice. The Writer will default to the standard output if it is
// not set.
//
// It could be useful to print the post-processed contents of the parsed
// files. The parser will follow any include directives and will strip
// comments and white space. Blank lines are ignored
type EchoParser struct {
	Writer io.Writer
}

// ParseLine for the EchoParser writes its line to the Writer
func (ep EchoParser) ParseLine(line string, _ *location.L) error {
	if ep.Writer == nil {
		ep.Writer = os.Stdout
	}
	_, err := ep.Writer.Write([]byte(line))
	return err
}
