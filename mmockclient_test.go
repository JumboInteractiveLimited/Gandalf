package gandalf

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/jmartin82/mmock/definition"
	"gopkg.in/h2non/gock.v1"
)

func urlOrFail(t *testing.T, rawurl string) *url.URL {
	t.Helper()
	u, err := url.Parse(rawurl)
	if err != nil {
		t.Fatalf("Could not parse %s to url due to error: %s", rawurl, err)
	}
	return u
}

func Test_clientMMock_constructURL(t *testing.T) {
	type fields struct {
		base   *url.URL
		client *http.Client
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantURL *url.URL
		wantErr bool
	}{
		{
			name: "clean",
			fields: fields{
				base:   urlOrFail(t, "http://test:8082"),
				client: &http.Client{Timeout: time.Second},
			},
			args:    args{""},
			wantURL: urlOrFail(t, "http://test:8082"),
			wantErr: false,
		},
		{
			name: "api",
			fields: fields{
				base:   urlOrFail(t, "http://test:8082"),
				client: &http.Client{Timeout: time.Second},
			},
			args:    args{"/api"},
			wantURL: urlOrFail(t, "http://test:8082/api"),
			wantErr: false,
		},
		{
			name: "api relative",
			fields: fields{
				base:   urlOrFail(t, "http://test:8082"),
				client: &http.Client{Timeout: time.Second},
			},
			args:    args{"api"},
			wantURL: urlOrFail(t, "http://test:8082/api"),
			wantErr: false,
		},
		{
			name: "api trailing",
			fields: fields{
				base:   urlOrFail(t, "http://test:8082"),
				client: &http.Client{Timeout: time.Second},
			},
			args:    args{"api/"},
			wantURL: urlOrFail(t, "http://test:8082/api/"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &clientMMock{
				base:   tt.fields.base,
				client: tt.fields.client,
			}
			gotURL, err := self.constructURL(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("clientMMock.constructURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotURL, tt.wantURL) {
				t.Errorf("clientMMock.constructURL() = %v, want %v", gotURL, tt.wantURL)
			}
		})
	}
}

func Test_clientMMock_getDefinitions(t *testing.T) {
	defer gock.Off()
	mocks := []definition.Mock{
		{
			URI:         "AAA",
			Description: "BBB",
			Request: definition.Request{
				Method: "GET",
				Path:   "/healthz",
			},
			Response: definition.Response{
				StatusCode: 200,
				Body:       "OK",
			},
		},
	}
	gock.New("http://test:8082").
		Get("/api/mapping").
		Reply(http.StatusOK).
		JSON(mocks)
	type fields struct {
		base   *url.URL
		client *http.Client
	}
	tests := []struct {
		name    string
		fields  fields
		wantOut []definition.Mock
		wantErr bool
	}{
		{
			name: "get single",
			fields: fields{
				base:   urlOrFail(t, "http://test:8082"),
				client: &http.Client{Timeout: time.Second},
			},
			wantOut: mocks,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &clientMMock{
				base:   tt.fields.base,
				client: tt.fields.client,
			}
			gotOut, err := self.getDefinitions()
			if (err != nil) != tt.wantErr {
				t.Errorf("clientMMock.getDefinitions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("clientMMock.getDefinitions() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func Test_clientMMock_sendDefinition(t *testing.T) {
	defer gock.Off()
	mock := definition.Mock{
		URI:         "AAA",
		Description: "BBB",
		Request: definition.Request{
			Method: "GET",
			Path:   "/healthz",
		},
		Response: definition.Response{
			StatusCode: 200,
			Body:       "OK",
		},
	}
	gock.New("http://test:8082").
		Post("/api/mapping/AAA").
		MatchType("json").
		JSON(mock).
		Reply(http.StatusCreated)
	gock.New("http://test:8082").
		Put("/api/mapping/AAA").
		MatchType("json").
		JSON(mock).
		Reply(http.StatusOK)
	type fields struct {
		base   *url.URL
		client *http.Client
	}
	type args struct {
		method string
		mock   definition.Mock
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "POST",
			fields: fields{
				base:   urlOrFail(t, "http://test:8082"),
				client: &http.Client{Timeout: time.Second},
			},
			args: args{
				method: "POST",
				mock:   mock,
			},
			wantErr: false,
		},
		{
			name: "PUT",
			fields: fields{
				base:   urlOrFail(t, "http://test:8082"),
				client: &http.Client{Timeout: time.Second},
			},
			args: args{
				method: "PUT",
				mock:   mock,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &clientMMock{
				base:   tt.fields.base,
				client: tt.fields.client,
			}
			if err := self.sendDefinition(tt.args.method, tt.args.mock); (err != nil) != tt.wantErr {
				t.Errorf("clientMMock.sendDefinition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
