package check

import (
	"strings"
	"testing"
)

func TestInvert(t *testing.T) {
	if e := Invert(Pass)(""); e == nil {
		t.Fatal("Invert failed to turn no error into an error")
	}
	if e := Invert(Fail)(""); e != nil {
		t.Fatalf("Invert failed to suppress and error and gave: %s", e)
	}
}

func TestAnd(t *testing.T) {
	if e := And(Pass, Pass, Pass)(""); e != nil {
		t.Fatalf("And should have passed but returned error: %s", e)
	}
	if e := And(Pass, Fail, Pass)(""); e == nil {
		t.Fatal("And should have failed with an error but gave none")
	}
}

func TestOr(t *testing.T) {
	if e := Or(Fail, Pass, Fail)(""); e != nil {
		t.Fatalf("Or should have passed but returned error: %s", e)
	}
	if e := Or(Fail, Fail, Fail)(""); e == nil {
		t.Fatal("Or should have failed with an error but gave none")
	}
}

func TestEquality(t *testing.T) {
	check := Equality(`"abc"`)
	if e := check(`"abc"`); e != nil {
		t.Fatalf("Equality should have passed but returned error: %s", e)
	}
	if e := check("abc"); e == nil {
		t.Fatal("Equality should have failed with an error but gave none")
	}
}

func TestRegexMatch(t *testing.T) {
	check := RegexMatch("^[0-9]+")
	if e := check(`123`); e != nil {
		t.Fatalf("RegexMatch should have passed but returned error: %s", e)
	}
	if e := check("abc"); e == nil {
		t.Fatal("RegexMatch should have failed with an error but gave none")
	}
}

func TestTransform(t *testing.T) {
	check := Transform(func(s string) string {
		return strings.Trim(s, `"`)
	}, Equality(`abc`))
	if e := check(`"abc"`); e != nil {
		t.Fatalf("Transform should have passed but returned error: %s", e)
	}
	if e := check(`'abc'`); e == nil {
		t.Fatal("Transform should have failed with an error but gave none")
	}
}
