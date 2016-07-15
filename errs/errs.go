// The errs package has a few utility
// functions to make handling errors
// a little more concise in certain cases.
package errs

import (
	"strings"
)

// the errors in this package will remember their
// ultimate root cause, as best as they can
type rootedError interface {
	error
	Root() error
}

// RootCause gets the ultimate root cause error
// out of a wrapped error.  If the given 
// error is not a wrapped error, the input
// is returned as-is.
func RootCause(e error) error {
	switch re := e.(type) {
	case rootedError:
		e  = re.Root()
	}
	return e
}


type combinedError struct {
	desc   string
	errors []error
}

func (ce *combinedError) Error() string {
	var all = []string{ce.desc}
	for _, v := range ce.errors {
		all = append(all, "\t"+v.Error())
	}
	return strings.Join(all, "\n")
}

func (ce *combinedError) Root() error {
	// just grab the root of the first error,
	// which is all we can really do here
	return RootCause(ce.errors[0])
}

// Combine returns the composite of all
// non-nil errors in the input, with the
// given text description added for context.
func Combine(desc string, e ...error) error {
	var ans *combinedError = nil
	for _, e := range e {
		if e != nil {
			if ans == nil {
				ans = &combinedError{desc, nil}
			}
			ans.errors = append(ans.errors, e)
		}
	}

	if ans == nil {
		return nil
	}
	return ans
}

type errorWithRoot struct {
	desc string
	base error
}

func (e *errorWithRoot) Error() string {
	return e.desc + " | Root Error: " + e.base.Error();
}

func (e *errorWithRoot) Root() error {
	return e.base
}

// Wrap simply wraps an existing error with some context.
func Wrap(desc string, e error) error {
	switch et := e.(type) {
	case *errorWithRoot:
		// preserve the ultimate root cause
		e = &errorWithRoot{desc + " | " + et.desc, et.Root()}
	case nil:
		e = nil
	default:
		e = &errorWithRoot{desc,e}
	}
	return e 
}

// First returns the first non-nil error
// given in the input, wrapped with the given
// description string.  If there isn't
// one, First returns nil.
func First(desc string, e ...error) error {
	for _, e := range e {
		if e != nil {
			return Wrap(desc, e)
		}
	}
	return nil
}


