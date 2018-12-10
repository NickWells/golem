# golem
This is a collection of useful go packages. The packages are all documented
and mostly have comprehensive tests.

## packages
* param is a replacement for the go flag package that adds the ability to check the supplied values and much more
* check collects some standard checks that can be applied to params
* filecheck offers some checks specific to files
* fileparser provides a standard means of parsing a file (params can be specified in files too)
* location is a common way of recording where a parameter was set
* mathutil provides some missing mathematical functions
* strdist provides some string distance functions which are used to suggest alternative parameters in error messages
* testhelper offers some functions commonly used when testing
* xdg supports the XDG Base Directory Specification

## How to use the param package

Define a ParamSet and populate it, then parse the command line arguments.

```go
var param1 int64
var param2 bool

func main() {
	ps, err := paramset.New(addParams,
		param.SetProgramDescription("this program will do cool stuff"))
	if err != nil {
		log.Fatal("Couldn't construct the param set: ", err)
	}
	ps.Parse()
```

The work is done mostly in the addParam function which should take a pointer to a
param.ParamSet and return an error.

```go
func addParams(ps *param.ParamSet) error {
	ps.Add("param-1",
		psetter.Int64Setter{
			Value:  &param1,
			Checks: []check.Int64{check.Int64LT(42)},
		},
		"this sets the value of param1",
		param.AltName("p1"))
		
	ps.Add("param-2", psetter.BoolSetter{Value: &param2},
		"this sets the value of param2")
		
	return nil
}
```

This illustrates the simplest use of the param package but you can specify
the behaviour much more precisely.

Additionally you can have positional parameters as well as named parameters.

You can specify a terminal parameter (by default `--`) and the remaining
parameters will be available for further processing without being parsed.

## Standard parameters
The default behaviour of the package is to add some standard
parameters. These allow the user to see a help message which is automatically
generated from the parameters added above. This can be in varying levels of
detail. If you just pass the `-help` param you will get the standard help
message but the `-help-full` parameter shows some hidden parameters and the
`-help-short` parameter gives a summary of the parameters.

Additionally the standard parameters offer the chance to examine where
parameters have been set and to control the parsing behaviour.

## The help message
The standard help message generated if the user passes the -help parameter
will show the program description and the non-hidden parameters. For each
parameter it will show:
* the parameter description
* any alternative names
* the initial value of the parameter
* the allowed values and whether there are any additional constraints

Additionally if there are any configuration files that have been specified
(use the `SetConfigFile` and `AddConfigFile` functions on the ParamSet) or
any environment variable prefixes have been given (use the `SetEnvPrefix` and
`AddEnvPrefix` functions on the ParamSet) these will be reported at the end
of the help message.

## Parameter Groups
Parameters can be grouped together so that they are reported together rather
than in alphabetical order. This is to allow logically related parameters to
be reported together. The standard parameters offer a way of showing the help
message just for specified parameter groups. You can add a description to be
shown for the parameter group and you can have configuration files and
environment variable prefixes which are specific to just the parameters in
the group (use corresponding `SetGroup...` and `AddGroup...` functions on the
ParamSet).
