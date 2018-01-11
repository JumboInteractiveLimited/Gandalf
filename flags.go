package gandalf

import (
	"flag"
	"net/http"
	"strings"
)

// This can be set using the `-gandalf.colour` switch
// and will force coloured cli output regardless of
// being a TTY or not.
var OverrideColour bool

// The target provider api to be targeted when making real
// outbound requests can be overridden globally with
// the `-gandalf.provider-host` cli switch.
var OverrideHost string

// If your contracts access multiple services using
// different fully qualified domain names, you may just
// want to append some custom domain (eg. dev.local) to
// all outbound requests. This can be done using the
// `-gandalf.provider-suffix` cli switch.
var OverrideHostSuffix string

// MMock definitions support chaos testing with random 5xx
// responses if you enable the ChaoticEvil switch in MMock
// structs. You can also override this in all definitions
// with the `-gandalf.mmock-chaos` cli switch.
var OverrideChaos bool

// If you would like to skip saving MMock definitions set
// this to true or use the `-gandalf.mmock-skip` cli switch.
var MockSkip bool

// Gandalf can be configured with custom flags given
// to the `go test` command or be setting the respective
// global variables.
//
// MMock structs can be used to generate MMock definitions,
// use the `-gandalf.mock-dest` cli switch to specify where
// to save these definitions for mmock ingestion.
var MockSavePath string

func init() {
	flag.BoolVar(&OverrideChaos, "gandalf.mmock-chaos", false, "Force enable chaos testing in all output mmock definitions.")
	flag.BoolVar(&MockSkip, "gandalf.mmock-skip", false, "Skip exporting contract definitions to mmock.")
	flag.BoolVar(&OverrideColour, "gandalf.colour", false, "Override tty detection and force colour output.")
	flag.StringVar(&MockSavePath, "gandalf.mock-dest", "./", "Destination to use when saving mocks.")
	flag.StringVar(&OverrideHost, "gandalf.provider-host", "", "if set to a non empty string all http requests for calls will be rewritten to use this address as the hostname and optional port.")
	flag.StringVar(&OverrideHostSuffix, "gandalf.provider-suffix", "", "when provided, this will be appended to the hostname of any and all outbound http requests.")
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
