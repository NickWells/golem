package paction_test

import (
	"fmt"
	"github.com/nickwells/golem/param"
	"github.com/nickwells/golem/param/paction"
	"github.com/nickwells/golem/param/paramset"
	"github.com/nickwells/golem/param/psetter"
)

// Example_count provides an example of how the paction.Counter might be
// used. Here we simply use it to check that only one of a pair of parameters
// has been set but it can be used to check for more complex rules. For
// instance, if one is set then all must be set or at least one of the set of
// parameters is set.
func Example_count() {
	// Declare the parameter values
	var (
		param1 bool
		param2 bool
	)

	// Declare the counter and make the associated action function. These two
	// lines and the extra arguments to the parameter creation are the core
	// of what you have to do
	var paramCounter paction.Counter
	af := (&paramCounter).MakeActionFunc()

	// Create the parameter set ...
	ps, _ := paramset.New()

	// ... and add the parameters. For each parameter we set the function to
	// be called after they have been set to the action function created
	// above. Note that we record the parameter we have created so that we
	// can report the name we gave it below
	p1 := ps.Add("p1", psetter.BoolSetter{Value: &param1},
		"the first flag (only set 1)",
		param.PostAction(af)) // This sets the action function
	p2 := ps.Add("p2", psetter.BoolSetter{Value: &param2},
		"the second flag (only set 1)",
		param.PostAction(af)) // This sets the action function

	// Now parse a set of supplied parameters. We force the arguments for the
	// purposes of the example, typically you would not pass anything and the
	// Parse function will use the command-line arguments.
	ps.Parse([]string{"-p1", "-p2"})

	// Now we can check the counter to see how many different parameters have
	// been set
	if paramCounter.Count() > 1 {
		fmt.Printf("Both of %s and %s have been set. At:\n",
			p1.Name(), p2.Name())

		// range over the parameters set and report them. Alternatively we
		// could use the SetBy function on the Counter and get a string
		// describing where the parameters were set
		for _, pSource := range paramCounter.ParamsSetAt {
			fmt.Println(pSource)
		}
	}

	// Output: Both of p1 and p2 have been set. At:
	// Param: p1 (at supplied parameters:1: -p1)
	// Param: p2 (at supplied parameters:2: -p2)
}
