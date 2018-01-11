package gandalf

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func ExampleToMMock() {
	_ = &Contract{
		Name: "MMockContract",
		Export: &ToMMock{
			Scenario:      "happy_path",
			TriggerStates: []string{"not_started"},
			NewState:      "started",
			ChaoticEvil:   true,
		},
	}
	// Output:
}

func getMMockContract() *Contract {
	body := `{"test_url":"https://example.com"}`
	return &Contract{
		Name: "MMockContract",
		Request: &DummyRequester{
			&http.Response{
				Status:        "200 OK",
				StatusCode:    200,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Body:          ioutil.NopCloser(bytes.NewBufferString(body)),
				ContentLength: int64(len(body)),
				Request:       &http.Request{},
				Header: http.Header{
					"Content-Type": []string{"application/json; charset=utf-8"},
				},
			},
		},
		Check: &DummyChecker{},
		Export: &ToMMock{
			Scenario:      "thing",
			TriggerStates: []string{"not_started"},
			NewState:      "started",
		},
	}
}

func TestMMockExporter(t *testing.T) {
	getMMockContract().Assert(t)
}

func TestMMockValidate(t *testing.T) {
	getMMockContract().Validate(t)
}

func BenchmarkMMockContracts(b *testing.B) {
	if testing.Short() {
		b.Skipf("MMock saving to file system is slow and therefore skipped in short mode.")
	}
	for n := 0; n < b.N; n++ {
		getMMockContract().Assert(b)
	}
}
