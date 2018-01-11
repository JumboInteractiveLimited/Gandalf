package gandalf

import (
	"testing"
	"time"
)

func ExampleSimpleRequester() {
	_ = Contract{
		Name:    "SimpleContract",
		Request: NewSimpleRequester("GET", "https://api.github.com", "", nil, time.Second*5),
	}
	// Output:
}

func getSimpleRequesterContract() *Contract {
	return &Contract{
		Name:    "SimpleContract",
		Request: NewSimpleRequester("GET", "https://api.github.com", "", nil, time.Second*5),
		Check:   &DummyChecker{},
	}
}

func TestSimpleRequester(t *testing.T) {
	getSimpleRequesterContract().Assert(t)
}

func TestSimpleRequesterCache(t *testing.T) {
	reqter := NewSimpleRequester("GET", "https://api.github.com", "", nil, time.Second*5)
	c := Contract{
		Name:    "SimpleContract",
		Request: reqter,
		Check:   &DummyChecker{},
	}
	if c.Run != 0 {
		t.Fatal("Contract run did not start at 0")
	}
	c.Assert(t)
	if c.Run != 1 {
		t.Fatalf("Contract run not incremented after assert, it is %d", c.Run)
	}
	c.Assert(t)
	if c.Run != 2 {
		t.Fatalf("Contract run not incremented after assert, it is %d", c.Run)
	}
	if reqter.lastRun != 1 {
		t.Fatalf("SimpleRequester last run is %d instead of the expected 1", reqter.lastRun)
	}
	if reqter.lastResp == nil {
		t.Fatal("SimpleRequester last response was not stored!")
	}
}

func BenchmarkSimpleRequester(b *testing.B) {
	if testing.Short() {
		b.Skipf("Skipped in short mode to avoid reaching github rate limit")
	}
	for n := 0; n < b.N; n++ {
		getSimpleRequesterContract().Assert(b)
	}
}
