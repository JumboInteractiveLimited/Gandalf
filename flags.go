package gandalf

import (
	"flag"
	"net/http"
	"strings"
)

// OverrideColour will force coloured cli output regardless of being a TTY or
// not. This can be set using the `-gandalf.colour` switch.
var OverrideColour bool

// OverrideHost rewrites the target provider api to be targeted when making
// real outbound via requesters that are correctly written to use this such as
// SimpleRequester. can be overridden globally with the
// `-gandalf.provider-host` cli switch.
var OverrideHost string

// OverrideHostSuffix rewrites the target provider hostname.  This can be
// useful if your contracts reference different hosts for various services,
// then setting OverrideHostSuffix to your dev instances domain to retarget at
// runtime.  This can be done using the `-gandalf.provider-suffix` cli switch.
var OverrideHostSuffix string

// OverrideWebroot gets prepended to all requests URI's. This can be useful
// when targeting an environment that uses webroot routing to the service to be
// tested. This can be done using the `-gandalf.provider-webroot` cli switch.
var OverrideWebroot string

// OverrideHTTPS if true will make all external requests use HTTPS. This may be
// required when targeting a production environment. This can be done using the
// `-gandalf.provider-https` cli switch.
var OverrideHTTPS bool

// OverrideChaos enables MMock definitions support chaos testing with random
// 5xx responses by setting the ChaoticEvil switch in ToMMock exporters. You
// can also override this in all definitions with the `-gandalf.mmock-chaos`
// cli switch.
var OverrideChaos bool

// MockSkip when set to true will not write mock definitions to disk. You can
// also override this wth the `-gandalf.mmock-skip` cli switch.
var MockSkip bool

// MockSleep sets the sleep period after exporting a mock definition set this
// to the number of milliseconds to sleep or use the `-gandalf.mock-sleep` cli
// switch.
var MockSleep int

// Gandalf can be configured with custom flags given
// to the `go test` command or be setting the respective
// global variables.
//
// MockSavePath tells exporters where to write generated mock should they have
// that functionality, eg. for mmock ingestion.  use the `-gandalf.mock-dest`
// cli switch to specify where
var MockSavePath string

func init() {
	flag.BoolVar(&OverrideChaos, "gandalf.mmock-chaos", false,
		"Force enable chaos testing in all output mmock definitions.")
	flag.BoolVar(&OverrideHTTPS, "gandalf.provider-https", false,
		"Force all requests to use HTTPS.")
	flag.BoolVar(&MockSkip, "gandalf.mmock-skip", false,
		"Skip exporting contract definitions to mmock.")
	flag.BoolVar(&OverrideColour, "gandalf.colour", false,
		"Override tty detection and force colour output.")
	flag.IntVar(&MockSleep, "gandalf.mock-sleep", 250,
		"Override milliseconds to wait after exporting a mock definition.")
	flag.StringVar(&MockSavePath, "gandalf.mock-dest", "./",
		"Destination to use when saving mocks.")
	flag.StringVar(&OverrideHost, "gandalf.provider-host", "",
		"if set to a non empty string all http requests for calls will be rewritten to use this address as the hostname and optional port.")
	flag.StringVar(&OverrideHostSuffix, "gandalf.provider-suffix", "",
		"when provided, this will be appended to the hostname of any and all external requests.")
	flag.StringVar(&OverrideWebroot, "gandalf.provider-webroot", "",
		"when provided, this will be prepended to the request URI of any and all external requests.")
	if !flag.Parsed() {
		flag.Parse()
	}
}

// Takes into account override flags.
func maybeOverrideHost(req *http.Request) {
	parts := strings.Split(req.URL.Host, ":")
	host := parts[0]
	if OverrideHost != "" {
		parts = strings.Split(OverrideHost, ":")
		host = parts[0]
	}
	if OverrideHTTPS {
		req.URL.Scheme = "https"
	}
	if OverrideWebroot != "" {
		req.URL.Path = strings.Join([]string{
			strings.TrimRight(OverrideWebroot, "/"),
			strings.TrimLeft(req.URL.Path, "/"),
		}, "/")
	}
	if OverrideHostSuffix != "" {
		if !strings.HasSuffix(host, ".") {
			host += "."
		}
		if !strings.HasPrefix(OverrideHostSuffix, ".") {
			host += OverrideHostSuffix
		} else {
			host += OverrideHostSuffix[1:]
		}
	}
	if len(parts) > 1 {
		host += ":" + parts[len(parts)-1]
	}
	req.Host = host
	req.URL.Host = host
}
