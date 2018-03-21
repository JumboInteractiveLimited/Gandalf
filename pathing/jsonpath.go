package pathing

import (
	"errors"
	"fmt"
	"strings"

	"github.com/JumboInteractiveLimited/Gandalf/check"
	"github.com/JumboInteractiveLimited/jsonpath"
)

// JSON meets the Extractor type for extracting JSON data with JSONPath. This
// extractor uses http://github.com/JumboInteractiveLimited/jsonpath which describes the
// supported jsonpath expressions/paths/addresses that you may use.
// Produces an error when a value cannot be extract for non wildcard paths.
func JSON(source, path string) (found []string, err error) {
	if path == "$+" {
		return []string{source}, nil
	}
	found = []string{}
	paths, err := jsonpath.ParsePaths(path)
	if err != nil {
		return
	}
	eval, err := jsonpath.EvalPathsInBytes([]byte(source), paths)
	if err != nil {
		return
	}
	excnt := 0
	for {
		if extract, ok := eval.Next(); ok {
			if eval.Error != nil {
				err = fmt.Errorf("failed to extract jsonpath %#v on iteration %d due to error:\n%s", path, excnt, err)
			}
			found = append(found, string(extract.Value))
			excnt++
		} else {
			break
		}
	}
	if len(found) == 0 && !strings.Contains(path, "*") {
		err = errors.New("no value extracted")
	}
	return found, err
}

// JSONChecks is a preloaded version of Checks with JSON as an extractor.
func JSONChecks(pcs PathChecks) check.Func {
	return Checks(JSON, pcs)
}
