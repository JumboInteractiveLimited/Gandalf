package gandalf

import (
	"errors"
	"net/http"
	"regexp"
	"sync"
)

// State is an in memory repository that can be used to perform stateful
// requests and response checks. This uses a thread safe singleton pattern and
// should not be instantiated anywhere other than GetState.
type State struct {
	KV map[string]interface{}
}

// Clear wipes all state clean.
func (s *State) Clear() {
	s.KV = map[string]interface{}{}
}

// ClearKey wipes a single key.
func (s *State) ClearKey(key string) {
	delete(s.KV, key)
}

// ClearRegex wipes all keys that match expr.
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

// GetResponse returns data stored at "{{Key}}.response".
func (s *State) GetResponse(key string) *http.Response {
	return s.KV[key+".response"].(*http.Response)
}

var stateInstance *State
var stateInstanceOnce sync.Once

// GetState return the thread safe global State instance.
func GetState() *State {
	stateInstanceOnce.Do(func() {
		stateInstance = &State{}
		stateInstance.Clear()
	})
	return stateInstance
}

// ToState is an exporter that will store the
// response for later usage.
type ToState struct {
	Key     string
	lastRun int
}

// Save the response of the current Requester run to
// a key in State.KV of the format (go tmpl style)
// "{{ ToState.Key }}.response" each time. It is expected
// that the Requester implements debouncing/caching so that
// Requester.Call can be rexecuted in the same run.
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
