package gandalf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/jmartin82/mmock/definition"
)

type clientMMock struct {
	base   *url.URL
	client *http.Client
}

var getMMockClient = func() func() *clientMMock {
	singleton := struct {
		client *clientMMock
		sync   sync.Once
	}{}
	return func() *clientMMock {
		if singleton.client == nil && strings.HasPrefix(MockSavePath, "http") {
			singleton.sync.Do(func() {
				c := &http.Client{Timeout: time.Duration(MockDelay) * time.Millisecond}
				u, err := url.Parse(MockSavePath)
				if err != nil {
					panic(fmt.Errorf("Could not construct mmock url due to error: %s", err))
				}
				singleton.client = &clientMMock{u, c}
			})
		}
		return singleton.client
	}
}()

func (self *clientMMock) constructURL(path string) (out *url.URL, err error) {
	uri, err := url.Parse(path)
	if err == nil {
		out = self.base.ResolveReference(uri)
	}
	return
}

func (self *clientMMock) call(method, path, body string) (*http.Response, error) {
	u, err := self.constructURL(path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, u.String(), bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}
	if method == http.MethodPut || method == http.MethodPost {
		req.Header.Add("Content-Type", "application/json")
	}
	return self.client.Do(req)
}

func (self *clientMMock) getDefinitions() (out []definition.Mock, err error) {
	resp, err := self.call("GET", "/api/mapping", "")
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		err = json.Unmarshal(body, &out)
	}
	return
}

func (self *clientMMock) sendDefinition(method string, mock definition.Mock) error {
	muri := getMMockDefURI(mock)
	bbody, err := json.Marshal(mock)
	if err != nil {
		return err
	}
	resp, err := self.call(method, "/api/mapping/"+muri, string(bbody))
	if err != nil {
		return err
	}
	if method == http.MethodPost && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("POST to MMock failed with status code %d", resp.StatusCode)
	}
	if method == http.MethodPut && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("PUT to MMock failed with status code %d", resp.StatusCode)
	}
	return nil
}

func (self *clientMMock) createDefinition(mock definition.Mock) error {
	return self.sendDefinition("POST", mock)
}

func (self *clientMMock) updateDefinition(mock definition.Mock) error {
	return self.sendDefinition("PUT", mock)
}

func (self *clientMMock) upsertDefinition(mock definition.Mock) error {
	all, err := self.getDefinitions()
	if err != nil {
		return err
	}
	exists := false
	wanted := getMMockDefURI(mock)
	for _, m := range all {
		if getMMockDefURI(m) == wanted {
			exists = true
		}
	}
	do := self.createDefinition
	if exists {
		do = self.updateDefinition
	}
	return do(mock)
}

func getMMockDefURI(mock definition.Mock) string {
	uri := mock.URI
	if uri == "" {
		uri = mock.Description
	}
	return uri
}
