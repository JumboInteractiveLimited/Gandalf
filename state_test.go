package gandalf

import (
	"net/http"
	"testing"
	"time"

	"github.com/JumboInteractiveLimited/Gandalf/pathing"
)

func TestToState(t *testing.T) {
	c := &Contract{
		Name: t.Name(),
		Request: &DummyRequester{
			&http.Response{
				StatusCode: 201,
			},
		},
		Check: &DummyChecker{},
		Export: &ToState{
			Key: t.Name(),
		},
	}
	c.Assert(t)
	r, ok := GetState().KV[t.Name()+".response"].(*http.Response)
	if !ok {
		t.Fatalf("Could not retrieve %s.response key and turn it into an *http.Response", t.Name())
	}
	if r.StatusCode != 201 {
		t.Fatalf("Stored response did not have the expected status code 201, instead it has %d", r.StatusCode)
	}
}

func TestState(t *testing.T) {
	s := GetState()
	s.KV["ABC"] = 42
	s.KV["EFG"] = 4242
	if v := GetState().KV["ABC"].(int); v != 42 {
		t.Fatalf("Value not shared across multiple uses of GetState")
	}
	GetState().Clear()
	if _, ok := GetState().KV["EFG"]; ok {
		t.Fatalf("Value not wiped after a state clear")
	}
	s.KV["ABC"] = 42
	s.KV["AEG"] = 43
	s.KV["BCD"] = 4242
	GetState().ClearRegex("A.*")
	if _, ok := GetState().KV["BCD"]; !ok {
		t.Fatalf("Key was wiped that did not match regex")
	}
	GetState().ClearKey("BCD")
	if _, ok := GetState().KV["BCD"]; ok {
		t.Fatalf("Single Key was not wiped")
	}
}

func ExampleToState() {
	_ = []*Contract{
		{Name: "ExampleStateContractRetrieve",
			Request: NewSimpleRequester("GET", "http://provider/token", "", nil, time.Second),
			Check: &SimpleChecker{
				HTTPStatus: 200,
			},
			Export: &ToState{
				Key: "ExampleStateContract",
			}},
		{Name: "ExampleStateContractUse",
			Request: &DynamicRequester{
				Builder: func(_ int) Requester {
					// get a token from the body of the last contract's response.
					body := GetResponseBody(GetState().GetResponse("ExampleStateContract"))
					found, err := pathing.GJSON(body, "result.token")
					if err != nil || len(found) == 0 {
						panic("Could not get the token from previous ExampleStateContract response")
					}
					token := found[0]
					// use the token in the header to get protected information
					return NewSimpleRequester("GET", "http://provider/info", "",
						http.Header{"Authorization": {"Bearer " + token}}, time.Second)
				},
			},
			Check: &SimpleChecker{
				HTTPStatus: 200,
			}},
	}
}
