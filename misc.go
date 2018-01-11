package gandalf

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// Allows using multiple Exporter structs in one contract.
type ToMultiple struct {
	Exporters []Exporter
}

// Convenience function for creating a ToMultiple
func ExportToMultiple(es ...Exporter) *ToMultiple {
	return &ToMultiple{
		Exporters: es,
	}
}

// Loops through Exporters and gives the Contract
// to each Save method, stopping on the first error.
func (m *ToMultiple) Save(c *Contract) error {
	for _, ex := range m.Exporters {
		if e := ex.Save(c); e != nil {
			return e
		}
	}
	return nil
}

// Reads the body from the request to be returned
// but also creates a new reader to put the body
// back into the response, allowing multiple reads.
func GetRequestBody(r *http.Request) string {
	if r == nil {
		return ""
	}
	bb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bb))
	return string(bb)
}

// Reads the body from the response to be returned
// but also creates a new reader to put the body
// back into the response, allowing multiple reads.
func GetResponseBody(r *http.Response) string {
	if r == nil {
		return ""
	}
	bb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bb))
	return string(bb)
}

// Returns a new HTTP response that should be sane;
// it has a 200 status code, body of "A", HTTP/1.1
// protocol, etc.
func SaneResponse() *http.Response {
	return &http.Response{
		StatusCode:    200,
		Status:        "200 OK",
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          ioutil.NopCloser(bytes.NewBufferString("A")),
		ContentLength: int64(len("A")),
	}
}
