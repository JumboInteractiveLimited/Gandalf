package pathing

import (
	"fmt"

	"github.com/JumboInteractiveLimited/Gandalf/check"
	"github.com/NodePrime/jsonpath"
)

// JsonPath function that meets the Extractor type. This extractor uses
// http://github.com/NodePrime/jsonpath which describes the supported
// jsonpath expressions/paths/addresses that you may use.
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
		if extract, ok := eval.Next(); err == nil && ok {
			if eval.Error != nil {
				err = fmt.Errorf("failed to extract jsonpath %#v on iteration %d due to error:\n%s", path, excnt, err)
			}
			found = append(found, string(extract.Value))
			excnt++
		} else {
			break
		}
	}
	return found, err
}

// This is a preloaded version of Checks with JSON as an extractor.
func JSONChecks(pcs PathChecks) check.Func {
	return Checks(JSON, pcs)
}
