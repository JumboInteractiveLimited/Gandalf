// Check is a collection of functions that assert various facts about the strings.
package check

import (
	"errors"
)

// A check function is defined here as a simple function
type Func func(contents string) error

// A path check that will always pass, never returns an error.
func Pass(_ string) error {
	return nil
}

// This will always return an error.
func Fail(_ string) error {
	return errors.New(`check always fails`)
}
