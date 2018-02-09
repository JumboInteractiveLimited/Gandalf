package pathing

import (
	"errors"
	"strings"

	"github.com/JumboInteractiveLimited/Gandalf/check"
	"github.com/tidwall/gjson"
)

// GJSON (https://github.com/tidwall/gjson) for JSON extraction.
// This is a valid Extractor that can be used for JSON in place of
// JSON if desired. This function has two special cases where
// it does not behave the same as GJSON, if the path is and
// empty string when source is not then source is returned
// as is in a slice, the second case is when value extracted
// is a map/hash then it will skip iterating over keys and return
// the full json of the map for optional nesting of Checks.
// Produces an error when no value is extracted.
func GJSON(source, path string) (found []string, err error) {
	found = []string{}
	if path == "" {
		if source != "" {
			found = append(found, source)
		}
		return found, err
	}
	result := gjson.Get(source, path)
	if strings.HasPrefix(result.Raw, "{") {
		found = append(found, result.Raw)
		return found, err
	}
	result.ForEach(func(key, value gjson.Result) bool {
		found = append(found, value.Raw)
		return true
	})
	if len(found) == 0 {
		err = errors.New("no value extracted")
	}
	return found, err
}

// GJSONChecks is a preloaded version of Checks with GJSON as the extractor.
func GJSONChecks(pcs PathChecks) check.Func {
	return Checks(GJSON, pcs)
}
