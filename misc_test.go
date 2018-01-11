package gandalf

import (
	"testing"
)

func TestGetResponseBody(t *testing.T) {
	r := SaneResponse()
	if b := GetResponseBody(r); b != "A" {
		t.Fatalf("Got %s instead of expected A", b)
	}
	if b := GetResponseBody(r); b != "A" {
		t.Fatalf("Got %s instead of expected A", b)
	}
	if b := GetResponseBody(r); b != "A" {
		t.Fatalf("Got %s instead of expected A", b)
	}
}

type captureExporter struct {
	seen map[int]bool
}

func (e *captureExporter) Save(c *Contract) error {
	if e.seen == nil {
		e.seen = map[int]bool{}
	}
	e.seen[c.Run] = true
	return nil
}

func TestToMultiple(t *testing.T) {
	ex1 := &captureExporter{}
	ex2 := &captureExporter{}
	c := &Contract{
		Name:    "ToMultipleContract",
		Request: &DummyRequester{SaneResponse()},
		Check:   &DummyChecker{},
		Export:  ExportToMultiple(ex1, ex2),
	}
	c.Assert(t)
	if len(ex1.seen)+len(ex2.seen) != 2 {
		t.Fatalf("Not all exporters were hit by ToMultiple:\n1: %#v\n2: %#v", ex1, ex2)
	}
}
