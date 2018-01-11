package pathing

import (
	"fmt"
	"testing"

	"github.com/JumboInteractiveLimited/Gandalf/check"
)

func compareStringSlices(a, b []string) error {
	if len(a) != len(b) {
		return fmt.Errorf("A has a length of %d yet B has a length of %d", len(a), len(b))
	}
	if (a == nil) != (b == nil) {
		return fmt.Errorf("A and/or B is nil")
	}
	for i, v := range a {
		if v != b[i] {
			return fmt.Errorf("A has %#v at index %d but B has %#v", v, i, b[i])
		}
	}
	return nil
}

func testStringSlices(t *testing.T, expected, result []string) {
	t.Helper()
	if t.Skipped() {
		return
	}
	if e := compareStringSlices(expected, result); e != nil {
		t.Fatalf("result slice does not match the expected: %s\n", e)
	}
}

func testError(t *testing.T, expected bool, err error) {
	t.Helper()
	if (err != nil) != expected {
		if expected {
			t.Fatalf("expected an error but one was not returned")
		} else {
			t.Fatalf("expected no error but got: %s", err)
		}
	}
}

func TestChecks(t *testing.T) {
	cases := []struct {
		input string
		pcs   PathChecks
		err   bool
	}{
		{`{}`, PathChecks{"": check.Pass}, false},
		{`{}`, PathChecks{"": check.Fail}, true},
		{`{"obj":{"name":"thing"}}`, PathChecks{"obj": GJSONChecks(PathChecks{"name": check.Equality(`"thing"`)})}, false},
		{`{"obj":[1,1,1]}`, PathChecks{"obj": check.Equality("1")}, false},
		{`{"obj":[1,1,2]}`, PathChecks{"obj": check.Equality("1")}, true},
	}
	for i, tt := range cases {
		t.Run(fmt.Sprintf("Case %d", i+1), func(st *testing.T) {
			testError(st, tt.err, GJSONChecks(tt.pcs)(tt.input))
		})
	}
}
