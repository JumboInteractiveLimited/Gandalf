package gandalf

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/JumboInteractiveLimited/Gandalf/check"
)

// SimpleChecker implements a Checker that asserts the expected HTTP status
// code, headers, and uses pathing.check.Func for checking the contents of the
// body.
type SimpleChecker struct {
	// HTTP Status code expected, ignored if left as default (0).
	HTTPStatus int
	// Assert that these headers have at least the values given. ignored if left as default.
	Headers http.Header
	// Uses a check.Func to assert the body is as expected.
	BodyCheck check.Func
	// Provide an example response body that should meet BodyCheck.
	ExampleBody string
}

// GetResponse returns a new HTTP response that should meet all checks.
func (c *SimpleChecker) GetResponse() *http.Response {
	r := SaneResponse()
	r.StatusCode = c.HTTPStatus
	r.Status = fmt.Sprintf("%d %s", c.HTTPStatus, http.StatusText(c.HTTPStatus))
	r.Header = c.Headers
	r.Body = ioutil.NopCloser(bytes.NewBufferString(c.ExampleBody))
	r.ContentLength = int64(len(c.ExampleBody))
	return r
}

// Assert the given HTTP response has the expected status code
func (c *SimpleChecker) assertStatus(res *http.Response) (err error) {
	if c.HTTPStatus != 0 && res.StatusCode != c.HTTPStatus {
		err = fmt.Errorf("HTTP Status code %d (%s) does not match expected %d (%s)",
			res.StatusCode, http.StatusText(res.StatusCode), c.HTTPStatus, http.StatusText(c.HTTPStatus))
	}
	return err
}

// Assert the given HTTP response has the expected headers, this check allows
// for additional headers to those that are expected without error but all expected
// headers must have (at least) the specified value(s).
func (c *SimpleChecker) assertHeaders(res *http.Response) (err error) {
	if len(c.Headers) == 0 {
		return nil
	}
	for expectedKey := range c.Headers {
		hasHeaderKey := false
		for k := range res.Header {
			if k == expectedKey {
				hasHeaderKey = true
			}
		}
		if !hasHeaderKey {
			return fmt.Errorf("Expected header key %#v not found in response", expectedKey)
		}
		for _, cv := range c.Headers[expectedKey] {
			hasHeaderValue := false
			for _, rv := range res.Header[expectedKey] {
				if rv == cv {
					hasHeaderValue = true
				}
			}
			if !hasHeaderValue {
				return fmt.Errorf("Expected header value %#v not found in header key %#v with %d values %v", cv, expectedKey, len(res.Header[expectedKey]), res.Header[expectedKey])
			}
		}
	}
	return err
}

// Check that the on the given responses body is as expected. The
// BodyCheck field on SimpleChecker can be either a single check.Func
// or you could use pathing.Checks to extract data from a structured
// text body. If BodyCheck is not set then it will use a string
// equality against the ExampleBody field.
func (c *SimpleChecker) assertBody(res *http.Response) (err error) {
	if c.BodyCheck == nil {
		c.BodyCheck = check.Equality(c.ExampleBody)
	}
	return c.BodyCheck(GetResponseBody(res))
}

// Assert the given HTTP response meets all checks.
// Executes methods in the following order:
//  1) SimpleChecker.assertStatus
//  2) SimpleChecker.assertHeaders
//  3) SimpleChecker.assertBody
func (c *SimpleChecker) Assert(res *http.Response) error {
	if e := c.assertStatus(res); e != nil {
		return fmt.Errorf("Status check failed, response body was:\n%s\n%s", GetResponseBody(res), e)
	}
	if e := c.assertHeaders(res); e != nil {
		return e
	}
	return c.assertBody(res)
}
