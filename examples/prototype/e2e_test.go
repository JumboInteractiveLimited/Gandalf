package prototype

import (
	"net/http"
	"testing"
	"time"

	. "github.com/JumboInteractiveLimited/Gandalf"
	c "github.com/JumboInteractiveLimited/Gandalf/check"
	p "github.com/JumboInteractiveLimited/Gandalf/pathing"
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
		{Name: "Healthz_Green",
			Request: NewSimpleRequester("GET", "http://provider/healthz", "", nil, time.Second),
			Check: &SimpleChecker{
				HTTPStatus: 200,
				Headers: http.Header{
					"Content-Type":                []string{"application/json"},
					"Access-Control-Allow-Origin": []string{"*"},
				},
				ExampleBody: `{"status":"green"}`,
				BodyCheck: p.JSONChecks(p.PathChecks{
					"$.status+": c.Equality(`"green"`),
				}),
			},
			Export: &ToMMock{
				Scenario:      "health",
				TriggerStates: []string{"not_started"},
				NewState:      "checked",
			},
		},
		{Name: "Healthz_Orange",
			Request: NewSimpleRequester("GET", "http://provider/healthz", "", nil, time.Second),
			Check: &SimpleChecker{
				HTTPStatus: 200,
				Headers: http.Header{
					"Content-Type":                []string{"application/json"},
					"Access-Control-Allow-Origin": []string{"*"},
				},
				ExampleBody: `{"status":"orange"}`,
				BodyCheck: p.JSONChecks(p.PathChecks{
					"$.status+": c.Equality(`"orange"`),
				}),
			},
			Export: &ToMMock{
				Scenario:      "health",
				TriggerStates: []string{"checked"},
				NewState:      "not_started",
			},
		},
	}
}

func TestStatelessContracts(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	for _, tc := range statelessContracts() {
		t.Run(tc.Name, func(st *testing.T) {
			tc.Assert(st)
		})
	}
}

func BenchmarkStatelessContracts(b *testing.B) {
	if testing.Short() {
		b.Skip()
	}
	for _, bc := range statelessContracts() {
		b.Run(bc.Name, func(sb *testing.B) {
			bc.Benchmark(sb)
		})
	}
}

func TestMain(m *testing.M) {
	Main(m)
}
