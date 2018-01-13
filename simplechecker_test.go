package gandalf

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/JumboInteractiveLimited/Gandalf/check"
)

func ExampleSimpleChecker() {
	_ = Contract{
		Name: "SimpleCheckerContract",
		Check: &SimpleChecker{
			HTTPStatus: 200,
			Headers: http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			},
			ExampleBody: "{}",
			BodyCheck:   check.Equality("{}"),
		},
	}
	// Output:
}

func getSimpleCheckerContract() *Contract {
	body := "{}"
	return &Contract{
		Name: "SimpleCheckerContract",
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
					"Content-Type": []string{"application/json; charset=utf-8", "application/javascript"},
					"Allow":        []string{"application/json; charset=utf-8"},
				},
			},
		},
		Check: &SimpleChecker{
			HTTPStatus: 200,
			Headers: http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			},
			ExampleBody: "{}",
			BodyCheck:   check.Equality("{}"),
		},
	}
}

func TestSimpleChecker(t *testing.T) {
	getSimpleCheckerContract().Assert(t)
}

func BenchmarkSimpleCheckerAssert(b *testing.B) {
	for n := 0; n < b.N; n++ {
		getSimpleCheckerContract().Assert(b)
	}
}
