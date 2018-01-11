package gandalf

import (
	"net/http"
	"strings"
	"testing"
)

// Checks nothing.
type DummyChecker struct {
}

// Returns nil for great justice. always passes.
func (c *DummyChecker) Assert(_ *http.Response) error {
	return nil
}

// Return sane empty response.
func (c *DummyChecker) GetResponse() *http.Response {
	return SaneResponse()
}

// Exports to the nothing at all.
type DummyExporter struct {
}

// Returns nil.
func (m *DummyExporter) Save(_ *Contract) error {
	return nil
}

// Returns the same response forever.
type DummyRequester struct {
	Response *http.Response
}

// Get that response (again?).
func (r *DummyRequester) Call(run int) (*http.Response, error) {
	return r.Response, nil
}

// Get a generic http request.
func (r *DummyRequester) GetRequest() *http.Request {
	req, err := http.NewRequest("GET", "/", strings.NewReader("A"))
	if err != nil {
		panic(err)
	}
	return req
}

func ExampleDummyRequester() {
	_ = &Contract{
		Name: "DummyContract",
		Request: &DummyRequester{
			&http.Response{},
		},
		Check:  &DummyChecker{},
		Export: &DummyExporter{},
	}
	// Output:
}

func ExampleDummyChecker() {
	_ = Contract{
		Name:  "DummyContract",
		Check: &DummyChecker{},
	}
	// Output:
}

func ExampleDummyExporter() {
	_ = Contract{
		Name:   "DummyContract",
		Export: &DummyExporter{},
	}
	// Output:
}

func getDummyContract() *Contract {
	return &Contract{
		Name: "DummyContract",
		Request: &DummyRequester{
			&http.Response{},
		},
		Check:  &DummyChecker{},
		Export: &DummyExporter{},
	}
}

func TestDummyContractValidate(t *testing.T) {
	getDummyContract().Validate(t)
}

func TestDummyContract(t *testing.T) {
	getDummyContract().Assert(t)
}

func BenchmarkDummyContractAsserts(b *testing.B) {
	for n := 0; n < b.N; n++ {
		getDummyContract().Assert(b)
	}
}

func BenchmarkDummyContractInSequence(b *testing.B) {
	for n := 0; n < b.N; n++ {
		getDummyContract().BenchmarkInSequence(b)
	}
}
