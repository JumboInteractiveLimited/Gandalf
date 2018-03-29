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

	"github.com/eapache/go-resiliency/retrier"
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

func (mmock *clientMMock) constructURL(path string) (*url.URL, error) {
	uri, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	return mmock.base.ResolveReference(uri), nil
}

func (mmock *clientMMock) call(method, path, body string) (*http.Response, error) {
	u, err := mmock.constructURL(path)
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
	return mmock.client.Do(req)
}

func (mmock *clientMMock) getDefinitions() (out []definition.Mock, err error) {
	resp, err := mmock.call("GET", "/api/mapping", "")
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		err = json.Unmarshal(body, &out)
	}
	return out, err
}

func (mmock *clientMMock) sendDefinition(method string, mock definition.Mock) error {
	muri := getMMockDefURI(mock)
	bbody, err := json.Marshal(mock)
	if err != nil {
		return err
	}
	resp, err := mmock.call(method, "/api/mapping/"+muri, string(bbody))
	if err != nil {
		return err
	}
	switch method {
	case http.MethodPost:
		if resp.StatusCode == http.StatusConflict {
			return mmock.sendDefinition("PUT", mock)
		} else if resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("POST to MMock failed with status code %d", resp.StatusCode)
		}
	case http.MethodPut:
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("PUT to MMock failed with status code %d", resp.StatusCode)
		}
	default:
		return fmt.Errorf("Cannot send definitons to MMock via HTTP Method %s", method)
	}
	return nil
}

func (mmock *clientMMock) getRetrier() *retrier.Retrier {
	return retrier.New(
		retrier.ConstantBackoff(3, time.Duration(MockDelay)*time.Millisecond),
		nil,
	)
}

func (mmock *clientMMock) createDefinition(mock definition.Mock) error {
	return mmock.getRetrier().Run(
		func() error {
			return mmock.sendDefinition("POST", mock)
		},
	)
}

func (mmock *clientMMock) updateDefinition(mock definition.Mock) error {
	return mmock.getRetrier().Run(
		func() error {
			return mmock.sendDefinition("PUT", mock)
		},
	)
}

func (mmock *clientMMock) upsertDefinition(mock definition.Mock) error {
	all, err := mmock.getDefinitions()
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
	do := mmock.createDefinition
	if exists {
		do = mmock.updateDefinition
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
