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
		https   bool
		root    string
		ohost   string
		osuffix string
		urlout  string
	}{
		{"http://target", false, "", "", "", "http://target"},
		{"http://target", false, "", "", ".lan", "http://target.lan"},
		{"http://target/foo", false, "", "", ".lan", "http://target.lan/foo"},
		{"http://target", false, "", "mock", "", "http://mock"},
		{"http://target/foo", false, "", "mock", "", "http://mock/foo"},
		{"http://target/foo", false, "", "mock", ".lan", "http://mock.lan/foo"},
		{"http://target", false, "", "mock:8080", "", "http://mock:8080"},
		{"http://target/foo", false, "", "mock:8080", "", "http://mock:8080/foo"},
		{"http://target", false, "", "mock:8080", ".lan", "http://mock.lan:8080"},
		{"http://target/foo", false, "", "mock:8080", ".lan", "http://mock.lan:8080/foo"},
		{"http://target/foo", false, "", "mock:8080", "lan", "http://mock.lan:8080/foo"},
		{"http://target/foo", true, "", "mock:8080", ".lan", "https://mock.lan:8080/foo"},
		{"https://target/foo", true, "", "mock:8080", ".lan", "https://mock.lan:8080/foo"},
		{"https://target/foo", false, "", "mock:8080", ".lan", "https://mock.lan:8080/foo"},
		{"http://target/foo", false, "api", "mock:8080", "lan", "http://mock.lan:8080/api/foo"},
		{"http://target/foo", false, "/api", "mock:8080", "lan", "http://mock.lan:8080/api/foo"},
		{"http://target/foo", false, "/api/", "mock:8080", "lan", "http://mock.lan:8080/api/foo"},
		{"http://target/foo", false, "api/", "mock:8080", "lan", "http://mock.lan:8080/api/foo"},
	}
	overrides := func(host, suffix, root string, https bool) {
		OverrideHost = host
		OverrideHostSuffix = suffix
		OverrideHTTPS = https
		OverrideWebroot = root
	}
	for i, tt := range cases {
		t.Run(fmt.Sprintf("Case %d", i+1), func(st *testing.T) {
			defer overrides(OverrideHost, OverrideHostSuffix, OverrideWebroot, OverrideHTTPS)
			overrides(tt.ohost, tt.osuffix, tt.root, tt.https)

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
