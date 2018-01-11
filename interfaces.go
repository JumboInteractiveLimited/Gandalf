package gandalf

import (
	"net/http"
)

// A Checker is an object that can assert that a given HTTP response
// meets some kind of criteria. A Checker should also be able
// to provide an HTTP response that would meet its checks to act as as a
// basis for; examples, mocks, or validation.
type Checker interface {
	// If the given response satisfies the Checker's criteria no error will be returned,
	// otherwise an error describing what check failed on the given response.
	Assert(*http.Response) error
	// Get a new response that would satisfy this Checker's criteria.
	GetResponse() *http.Response
}

// An Exporter represents a method of transforming and exporting the given contract.
type Exporter interface {
	Save(contract *Contract) error
}

// Requester knows how to reliably call a provider service to get a HTTP response
// that may later be used in a Checker. A Requester should also be able
// to provide an HTTP request to act as a basis for; examples, mocks, or self testing.
type Requester interface {
	Call(run int) (*http.Response, error)
	GetRequest() *http.Request
}
