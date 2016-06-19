// +build windows

package cmdline // import "go.waywardcode.com/cmdline"

import (
	"os"
	"path/filepath"
)

// GlobArgs is a no-op on UNIX, but on
// windows it glob-expands all the arguments
// in the os.Args array.
func GlobArgs() {
	var exp = append(make([]string, 0, len(os.Args)), os.Args[0])

	var err error
	var ea []string
	for _, a := range os.Args[1:] {
		if ea, err = filepath.Glob(a); err != nil || len(ea) == 0 {
			exp = append(exp, a)
		} else {
			exp = append(exp, ea...)
		}
	}
	os.Args = exp // override default args...
}
