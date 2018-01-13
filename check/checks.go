// Package check is a collection of functions that assert various facts about the strings.
package check

import (
	"errors"
)

// Func is a check function, defined here as a simple function that takes a
// string and returns an error explaining why the check failed or nil if the
// input passed the check.
type Func func(contents string) error

// Pass is a path check that will always pass, never returns an error.
func Pass(_ string) error {
	return nil
}

// Fail check will always return an error.
func Fail(_ string) error {
	return errors.New(`check always fails`)
}
