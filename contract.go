package gandalf

import (
	"testing"
)

// Contract is at the core of Gandalf, it represents the contract between a
// consumer and provider in two main parts, the Request, and the Check.
// The Request object is responsible for geting information into Gandalf
// from the provider for testing. Then the response is given to the Check
// object to  that the response meets whatever criteria the Checker supports.
type Contract struct {
	// Unique identifier for this contract.
	Name    string
	Check   Checker
	Request Requester
	Export  Exporter
	// Run stores the number of times that Assert has been run and executed all parts of the contract.
	// This allows for some information such as the request to differ per call if desired.
	Run int
	// If an Optional Contract fails its checks it will not fail the whole test run.
	Optional bool
	// Set to true after this contract is tested the first time, pass or fail.
	Tested     bool // internal state to mark if a contract has already been tested
	notHonored bool // internal state to mark when a test for this contract failed
}

// Checks the error and skips or fails the Testable based on the Optional switch.
func (c *Contract) honorCheck(t Testable, err error) {
	t.Helper()
	if c.notHonored = err != nil; c.notHonored {
		if c.Optional {
			t.Skipf("Contract %#v was not honored, but does not have to be, due to error:\n%s", c.Name, err)
		} else {
			t.Fatalf("Contract %#v was not honored, and it must be, due to error:\n%s", c.Name, err)
		}
	}
}

// Testable is the common interface between tests and benchmarks required to handle them interchangeably.
type Testable interface {
	Helper()
	Fatalf(format string, args ...interface{})
	Skipf(format string, args ...interface{})
}

func (c *Contract) export(t Testable) {
	if c.Export != nil {
		if e := c.Export.Save(c); e != nil {
			t.Fatalf("Export save failed due to error: %s", e)
		}
	}
}

// Assert runs a request and checks response on the contract causing a pass or fail.
// This executes the Exporter before and after calling the Requester, allowing
// for pre and post exporters.
func (c *Contract) Assert(t Testable) {
	defer func() { c.Run++ }()
	defer c.export(t)
	c.export(t)

	resp, err := c.Request.Call(c.Run)
	c.honorCheck(t, err)
	c.honorCheck(t, c.Check.Assert(resp))
	c.Tested = true
}

// Benchmark just the Requester's ability to provider responses in sequence.
// This uses the benchmark run counter instead of the Contract.Run field.
func (c *Contract) Benchmark(b *testing.B) {
	if c.Tested && c.notHonored {
		b.Skipf("Contract %#v benchmark skipped as the contract was not honored by passing its test")
	}
	errors := 0
	for n := 0; n < b.N; n++ {
		_, err := c.Request.Call(n)
		if err != nil {
			errors++
		}
	}
	if errors > 0 {
		b.Logf("%d calls out of %d (%.2f%%) resulted in an error",
			errors, b.N, float32(errors)*float32(100)/float32(b.N))
	}
}

// BenchmarkInOrder takes a list of contracts and benchmarks the time it takes
// to call each of their requests in sequence before starting the next run.
// This may be useful if a list of contracts defined, for example, a common
// customer journey to be benchmarked.
func BenchmarkInOrder(b *testing.B, contracts []*Contract) {
	errors := 0
	for n := 0; n < b.N; n++ {
		for _, c := range contracts {
			_, err := c.Request.Call(n)
			if err != nil {
				errors++
			}
		}
	}
	if errors > 0 {
		b.Logf("%d calls out of %d (%.2f%%) resulted in an error",
			errors, b.N, float32(errors)*float32(100)/float32(b.N))
	}
}
