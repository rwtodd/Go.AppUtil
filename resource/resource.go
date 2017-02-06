// The resource package helps an application
// locate files important to the program.
// In this early state, the package searches
// $GOPATH, but eventually more features will
// be added.
package resource

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// A Locator knows how to find resources
type Locator interface {
	// Path translates a resource name into
	// a full filesystem path.
	Path(rsc string) (string, error)
}

type pathsLocator struct {
	paths []string
}

// NewPathLocator creates a Locator which uses
// a fixed set of base paths to find a given
// resource.  Set a non-empty goPathRepo to append
// $GOPATH/src/$goPathRepo to the list of paths. It is
// not considered an error if $GOPATH isn't
// defined.
func NewPathLocator(paths []string, goPathRepo string) Locator {
	if len(goPathRepo) > 0 {
		gpth := os.Getenv("GOPATH")
		if len(gpth) > 0 {
			for _, base := range strings.Split(gpth, ":") {
				paths = append(paths, filepath.Join(base, "src", goPathRepo))
			}
		}
	}
	return &pathsLocator{paths}
}

func (l *pathsLocator) Path(rsc string) (string, error) {
	for _, root := range l.paths {
		attempt := filepath.Join(root, rsc)
		_, err := os.Stat(attempt)
		if err == nil {
			return attempt, nil
		}
	}
	// if we got here, rsc was not found in any of the paths
	return "", errors.New(rsc + ": resource not found")
}
