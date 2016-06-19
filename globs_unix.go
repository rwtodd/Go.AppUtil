// +build !windows

package cmdline

// GlobArgs is a no-op on UNIX, but on
// windows it glob-expands all the arguments
// in the os.Args array.
func GlobArgs() {
	return
}
