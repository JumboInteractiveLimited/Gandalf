package gandalf

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestDynamicRequester(t *testing.T) {
	r := DynamicRequester{
		Builder: func(run int) Requester {
			return NewSimpleRequester("GET", "https://api.github.com", strconv.Itoa(run), nil, time.Second*5)
		},
	}
	if body := GetRequestBody(r.GetRequest()); body != "0" {
		t.Fatalf("Incorrect request generated, body of %s instead of 0", body)
	}
	rs, err := r.Call(1)
	if err != nil {
		t.Fatalf("Unexpected error on call: %s", err)
	}
	if rs.StatusCode != 200 {
		t.Fatalf("Got status code %d instead of the expected 200", rs.StatusCode)
	}
}

func ExampleDynamicRequester() {
	_ = Contract{
		Name: "SimpleContract",
		Request: &DynamicRequester{
			Builder: func(run int) Requester {
				// Each run will get a different post id
				return NewSimpleRequester("GET", fmt.Sprintf("http://provider/post/%d", run), "", nil, time.Second*5)
			},
		},
	}
}
