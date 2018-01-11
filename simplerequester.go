package gandalf

import (
	"bytes"
	"net/http"
	"time"
)

// A Requester that executes the stored Request each time.
type SimpleRequester struct {
	Request  *http.Request
	Timeout  time.Duration
	lastRun  int
	lastResp *http.Response
}

// Wrapper to easily create a SimpleRequester.
func NewSimpleRequester(method, url, body string, headers http.Header, timeout time.Duration) *SimpleRequester {
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		panic(err)
	}
	maybeOverrideHost(req)
	req.Header = headers
	return &SimpleRequester{
		Request: req,
		Timeout: timeout,
	}
}

// Call the Request. The last response is stored to be given on multiple calls for the same run.
func (r *SimpleRequester) Call(run int) (*http.Response, error) {
	if r.lastResp != nil && run == r.lastRun {
		return r.lastResp, nil
	}
	if r.Timeout == 0 {
		r.Timeout = time.Second
	}
	client := &http.Client{Timeout: r.Timeout}
	res, err := client.Do(r.Request)
	if err == nil {
		r.lastResp = res
	}
	r.lastRun = run
	return res, err
}

// Get the Request.
func (r *SimpleRequester) GetRequest() *http.Request {
	return r.Request
}
