package gandalf

import (
	"errors"
	"net/http"
	"regexp"
	"sync"
)

// This is a state repository that can be used to perform
// stateful requests and response checks. This uses a thread
// safe singleton pattern and should not be instantiated
// anywhere other than GetState.
type State struct {
	KV map[string]interface{}
}

// Wipe all state clean.
func (s *State) Clear() {
	s.KV = map[string]interface{}{}
}

// Wipe a single key.
func (s *State) ClearKey(key string) {
	delete(s.KV, key)
}

// Wipe all keys that match expr.
func (s *State) ClearRegex(expr string) error {
	ex, err := regexp.Compile(expr)
	if err != nil {
		return err
	}
	for key := range s.KV {
		if ex.MatchString(key) {
			delete(s.KV, key)
		}
	}
	return nil
}

// Get a response object stored at "{{Key}}.response".
func (s *State) GetResponse(key string) *http.Response {
	return s.KV[key+".response"].(*http.Response)
}

var stateInstance *State
var stateInstanceOnce sync.Once

// Return the thread safe global State instance.
func GetState() *State {
	stateInstanceOnce.Do(func() {
		stateInstance = &State{}
		stateInstance.Clear()
	})
	return stateInstance
}

// Every time the contract is called this will store the
// response for later usage.
type ToState struct {
	Key     string
	lastRun int
}

// Exports the response of the current Requester run to
// a key in State.KV of the format (go tmpl style)
// "{{ ToState.Key }}.response" each time.
func (e *ToState) Save(c *Contract) error {
	if e.Key == "" {
		return errors.New("exporter ToState requires a non empty Key for storage")
	}
	key := e.Key + ".response"
	s := GetState()
	_, defined := s.KV[key]
	if e.lastRun != c.Run || !defined {
		r, e := c.Request.Call(c.Run)
		s.KV[key] = r
		return e
	}
	return nil
}
