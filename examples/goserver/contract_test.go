package goserver

import (
	"net/http"
	"testing"
	"time"

	. "github.com/JumboInteractiveLimited/Gandalf"
	c "github.com/JumboInteractiveLimited/Gandalf/check"
)

func statelessContracts() []*Contract {
	return []*Contract{
		{Name: "Robots",
			Request: NewSimpleRequester("GET", "http://provider/robots.txt", "", nil, time.Second),
			Check: &SimpleChecker{
				HTTPStatus: 200,
				Headers: http.Header{
					"Content-Type":                []string{"text/plain"},
					"Access-Control-Allow-Origin": []string{"*"},
				},
				ExampleBody: "User-agent: *\nDisallow: /",
				BodyCheck:   c.Equality("User-agent: *\nDisallow: /"),
			},
			Export: &ToMMock{},
		},
	}
}

func TestStatelessContracts(t *testing.T) {
	for _, tc := range statelessContracts() {
		t.Run(tc.Name, func(st *testing.T) {
			tc.Assert(st)
		})
	}
}

func BenchmarkStatelessContracts(b *testing.B) {
	for _, bc := range statelessContracts() {
		b.Run(bc.Name, func(sb *testing.B) {
			bc.Benchmark(sb)
		})
	}
}

func TestMain(m *testing.M) {
	MainWithHandler(m, setupServer())
}
