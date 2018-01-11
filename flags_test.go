package gandalf

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
)

func TestHostOverrideFlags(t *testing.T) {
	cases := []struct {
		urlin   string
		ohost   string
		osuffix string
		urlout  string
	}{
		{"http://target", "", "", "http://target"},
		{"http://target", "", ".lan", "http://target.lan"},
		{"http://target/foo", "", ".lan", "http://target.lan/foo"},
		{"http://target", "mock", "", "http://mock"},
		{"http://target/foo", "mock", "", "http://mock/foo"},
		{"http://target/foo", "mock", ".lan", "http://mock.lan/foo"},
		{"http://target", "mock:8080", "", "http://mock:8080"},
		{"http://target/foo", "mock:8080", "", "http://mock:8080/foo"},
		{"http://target", "mock:8080", ".lan", "http://mock.lan:8080"},
		{"http://target/foo", "mock:8080", ".lan", "http://mock.lan:8080/foo"},
		{"http://target/foo", "mock:8080", "lan", "http://mock.lan:8080/foo"},
	}
	overrides := func(host, suffix string) {
		OverrideHost = host
		OverrideHostSuffix = suffix
	}
	for i, tt := range cases {
		t.Run(fmt.Sprintf("Case %d", i+1), func(st *testing.T) {
			defer overrides(OverrideHost, OverrideHostSuffix)
			overrides(tt.ohost, tt.osuffix)

			req, err := http.NewRequest("GET", tt.urlin, bytes.NewBuffer([]byte{}))
			if err != nil {
				st.Fatal(err)
			}
			if req.URL.String() != tt.urlin {
				st.Fatalf("Request URL %s does not match input %s", req.URL.String(), tt.urlin)
			}
			maybeOverrideHost(req)
			if req.URL.String() != tt.urlout {
				st.Fatalf("Request URL post override %s does not match expected %s",
					req.URL.String(),
					tt.urlout)
			}
		})
	}
}
