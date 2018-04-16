package gandalf

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
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

func Test_splitHeaders(t *testing.T) {
	type args struct {
		in http.Header
	}
	tests := []struct {
		name    string
		args    args
		wantOut http.Header
	}{
		{
			name: "Separate Values",
			args: args{http.Header{
				"Stuff": []string{"one", "two"},
			}},
			wantOut: http.Header{
				"Stuff": []string{"one", "two"},
			},
		},
		{
			name: "CSV Values",
			args: args{http.Header{
				"Stuff": []string{"one,two"},
			}},
			wantOut: http.Header{
				"Stuff": []string{"one", "two"},
			},
		},
		{
			name: "Mixed Values",
			args: args{http.Header{
				"Stuff": []string{"one", "two,three"},
			}},
			wantOut: http.Header{
				"Stuff": []string{"one", "two", "three"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := splitHeaders(tt.args.in); !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("splitHeaders() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
