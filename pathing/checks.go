// Pathing is a collection of functions that are able to
// extract values from data using on a path/query.
package pathing

import (
	"fmt"

	"github.com/JumboInteractiveLimited/Gandalf/check"
)

// An Extractor is any func that takes two strings and returns a slice of strings and an error.
// These functions take a source string and a path string. The source string is the structured
// text document that this extractor will try to extract data from. The path string is an
// address/expression in a pathing DSL (eg. XPath or JSonPath) that describes what data should be
// extracted from the source. All Extractors should return the extracted data as a slice of string
// along with an error to indicate failure or success.
type Extractor func(source, path string) ([]string, error)

// A collection of path expressions for a given Extrator and
// check.Func's to check the value(s) with.
type PathChecks map[string]check.Func

// Extracts paths (keys of the given PathChecks and passes the value
// add that path to a check.Func for assertion.
func Checks(ex Extractor, pcs PathChecks) check.Func {
	cname := "PathingChecks"
	return func(found string) error {
		for path, check := range pcs {
			extracts, err := ex(found, path)
			if err != nil {
				return throw(cname, found,
					fmt.Errorf("check for path %s failed due to error:\n%s", path, err))
			}
			for _, s := range extracts {
				if e := check(s); e != nil {
					return throw(cname, found,
						fmt.Errorf("check for path %s failed due to error:\n%s", path, e))
				}
			}
		}
		return nil
	}
}

func throw(check, found string, err error) error {
	return fmt.Errorf("%s Failed when checking found value '%s':\n%s", check, found, err)
}
