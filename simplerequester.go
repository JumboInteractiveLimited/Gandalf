package gandalf

import (
	"bytes"
	"net/http"
	"time"
)

// SimpleRequester implements a Requester that executes the stored Request each
// time.
type SimpleRequester struct {
	Request  *http.Request
	Timeout  time.Duration
	lastRun  int
	lastResp *http.Response
}

// NewSimpleRequester is a wrapper to easily create a SimpleRequester given a
// limited set of common inputs.
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

func (r *SimpleRequester) getClient() *http.Client {
	return &http.Client{
		Timeout: r.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
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
	res, err := r.getClient().Do(r.Request)
	if err == nil {
		r.lastResp = res
	}
	r.lastRun = run
	return res, err
}

// GetRequest SimpleRequester.Request.
func (r *SimpleRequester) GetRequest() *http.Request {
	return r.Request
}
