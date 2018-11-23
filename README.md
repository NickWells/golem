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
