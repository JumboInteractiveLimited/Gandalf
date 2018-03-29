package gandalf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/jmartin82/mmock/definition"
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

func getMMockContracts() []*Contract {
	body := `{"test_url":"https://example.com"}`
	return []*Contract{
		{
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
		},
		{
			Name: "MMockContractPath",
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
				Path:          "/people/:name",
			},
		},
	}
}

func TestMMockExporter(t *testing.T) {
	for i, tc := range getMMockContracts() {
		t.Run(fmt.Sprintf("Case %d", i), func(st *testing.T) {
			tc.Assert(t)
			mock, err := readMMockDefinition(fmt.Sprintf("./%s.json", tc.Name))
			if err != nil {
				st.Fatalf("Could not read back MMock definition file due to error: %s", err)
			}
			if strings.Contains(tc.Name, "Path") {
				if !strings.Contains(mock.Request.Path, ":name") {
					st.Fatalf("Could not find ':name' in path that should be overridden '%s'\n", mock.Request.Path)
				}
			}
		})
	}
}

func readMMockDefinition(path string) (mock definition.Mock, err error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	return mock, json.Unmarshal(file, &mock)
}

func BenchmarkMMockContracts(b *testing.B) {
	if testing.Short() {
		b.Skipf("MMock saving to file system is slow and therefore skipped in short mode.")
	}
	for n := 0; n < b.N; n++ {
		for _, tc := range getMMockContracts() {
			tc.Assert(b)
		}
	}
}
