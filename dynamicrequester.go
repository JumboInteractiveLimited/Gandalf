package gandalf

import "net/http"

// DynamicRequester allows for generating requests using a a function that
// creates a requester each time it is called. This is useful for requests that
// should change at runtime based on, for example, State values. Caches a
// single requester based on the run for debouncing.
type DynamicRequester struct {
	Builder       func(run int) Requester
	lastRun       int
	lastRequester Requester
}

// Handles caching the last requester built.
func (r *DynamicRequester) getRequester(run int) Requester {
	if r.lastRequester == nil || r.lastRun != run {
		r.lastRun = run
		r.lastRequester = r.Builder(run)
	}
	return r.lastRequester
}

// Call executes the builder (or retrieve from cache if the run is the same as
// the last Call execution) then Call the requester passing on the run.
func (r *DynamicRequester) Call(run int) (*http.Response, error) {
	rs, err := r.getRequester(run).Call(run)
	return rs, err
}

// GetRequest passes down to the current Requester's GetRequest method.  This
// uses the last run given to Call (or 0) as the run to give to the builder.
func (r *DynamicRequester) GetRequest() *http.Request {
	return r.getRequester(r.lastRun).GetRequest()
}
