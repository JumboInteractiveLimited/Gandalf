// check functions are so simple there are many useful ways
// they can be generated.
package check

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Invert the output of the given Func, if it returned nil then return an
// error, if it returned error return nil.
func Invert(next Func) Func {
	cname := "Invert"
	return func(found string) error {
		if e := next(found); e != nil {
			return nil
		}
		return throw(cname, found,
			errors.New(`check was expected to fail but didn't`))
	}
}

// And combines multiple check Func's into a single one where all must pass.
// Short circuits on first failure.
func And(funcs ...Func) Func {
	cname := "And"
	return func(found string) error {
		var err error
		for i, f := range funcs {
			err = f(found)
			if err != nil {
				return throw(cname, found,
					fmt.Errorf(`check %d failed with error "%s"`, i, err))
			}
		}
		return nil
	}
}

// Or combines multiple check Func's into a single one where any one must pass, the rest may fail.
// Short circuits on first success.
func Or(funcs ...Func) Func {
	cname := "Or"
	return func(found string) error {
		var errs []error
		for _, f := range funcs {
			if e := f(found); e != nil {
				errs = append(errs, e)
			} else {
				return nil
			}
		}
		serrs := []string{}
		for _, e := range errs {
			serrs = append(serrs, e.Error())
		}
		return throw(cname, found,
			fmt.Errorf(`all checks failed with errors ["%s"]`, strings.Join(serrs, `","`)))
	}
}

// Equality dose a string equality check on the found value. This performs no
// deserialization of the value, this means you may need to include qoutes for
// strings.
func Equality(expected string) Func {
	cname := "Equality"
	return func(found string) error {
		if found != expected {
			return throw(cname, found,
				fmt.Errorf("found value %#v does not equal expected %#v", found, expected))
		}
		return nil
	}
}

// RegexMatch compiles the given regex and returns a patcheck Func that asserts
// values match the expression.
func RegexMatch(expr string) Func {
	cname := "RegexMatch"
	ex, err := regexp.Compile(expr)
	if err != nil {
		panic(err)
	}
	return func(found string) error {
		if !ex.MatchString(found) {
			return throw(cname, found,
				fmt.Errorf("expression %#v did not match %#v", expr, found))
		}
		return nil
	}
}

// Transform make it easy to alter a value before passing it to the next
// pathcheck Func. This might be useful for stripping whitespace or normalizing
// case. The func f will be run each time the Func next is evaluated.
func Transform(f func(string) string, next Func) Func {
	return func(found string) error {
		return next(f(found))
	}
}

func throw(check, found string, err error) error {
	return fmt.Errorf("Pathcheck '%s' Failed when checking found value '%s':\n%s", check, found, err)
}
