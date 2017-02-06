# Resource Package

This is a simple Go package to help a utility find its resources.  In the java world, I'd
typically package files into my Jar, but have no equivalent for Go other than compiling
the files into the executable as []bytes. 

## Status

This is a first cut an what I need to get by.  An interface, `Locator`, is provided, which
translates between resource names and full path names.  A singe kind of Locator is 
provided, which will search a list of provided root paths for a resource.

Example use:

	// second arg is appended to $GOPATH/src, and added after first arg
	loc := resource.NewPathLocator([]string{"/usr/local/share/go/mypackage"}, "/github.com/rwtodd/mypackage")

	// look up the licence file in our resources
	license, err := loc.Path("LICENSE") 

## Future Plans

Other types of `Locator` are possible: one that hashes all the available files in
a set of directories so the lookup is very quick, etc.  I plan to provide these
as the need arises.


