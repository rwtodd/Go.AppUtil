package errs

import (
	"errors"
	"testing"
)

// a helper function to compare two errors, which might be different types
func compareErrors(a error, b error) bool {
	if a == nil {
		return b == nil
	}
	if b == nil {
		return false // we already know a != nil
	}
	return a.Error() == b.Error()
}

func TestCombine(t *testing.T) {
	theError := errors.New("I'm an Error")
	cases := []*struct {
		inputs []error
		ans    error
	}{
		{[]error{nil, nil, nil}, nil},
		{[]error{nil, nil, theError}, &combinedError{"A Test Case", []error{theError}}},
		{[]error{nil, theError, nil}, &combinedError{"A Test Case", []error{theError}}},
		{[]error{theError, nil, theError}, &combinedError{"A Test Case", []error{theError, theError}}},
	}

	for _, c := range cases {
		e := Combine("A Test Case", c.inputs...)
		if !compareErrors(e, c.ans) {
			t.Fatalf("Combine got <%v> instead of <%v>!", e, c.ans)
		}
	}
}

func TestFirst(t *testing.T) {
	err1 := errors.New("I'm an Error")
	err2 := errors.New("I'm another Error")
	err1ans := Wrap("A Test Case", err1)
	err2ans := Wrap("A Test Case", err2)

	cases := []*struct {
		inputs []error
		ans    error
	}{
		{[]error{nil, nil, nil}, nil},
		{[]error{nil, nil, err1}, err1ans},
		{[]error{err2, nil, err1}, err2ans},
		{[]error{nil, err2, err1, err1}, err2ans},
	}

	for _, c := range cases {
		e := First("A Test Case", c.inputs...)
		if !compareErrors(e, c.ans) {
			t.Fatalf("First got <%v> instead of <%v>!", e, c.ans)
		}
	}
}

func TestWrap(t *testing.T) {
	err1 := errors.New("I'm an error")
	err2 := Wrap("A Test Case", err1)
	err3 := Wrap("A Test Case 2", err2)
	
	if RootCause(err2) != err1 {
		t.Fatalf("Root cause not preserved!")
	}
	if RootCause(err3) != err1 {
		t.Fatalf("Root cause not preserved 2 levels deep!")
	}
}
