package double

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/krostar/test/internal"
)

// TestingT is an interface for testing types, mirroring the standard library's *testing.T.
type TestingT = internal.TestingT

// Spy implements the TestingT interface and records all method calls for later verification.
// It's useful for testing assertions and other test utilities by wrapping another TestingT
// implementation and intercepting all method calls.
type Spy struct {
	m           sync.RWMutex // mutex to protect concurrent access
	underlyingT TestingT     // the wrapped TestingT implementation

	failed  bool                // tracks whether Fail or FailNow was called
	logs    []string            // stores all messages logged with Logf
	records []SpyTestingTRecord // stores all method calls with their inputs and outputs
}

// NewSpy creates a new Spy that wraps the provided TestingT implementation.
// All method calls on the returned Spy will be recorded and also delegated
// to the underlying TestingT instance.
//
// This allows test code to verify the behavior of code that uses a TestingT
// without failing the actual test unless explicitly checked.
func NewSpy(underlyingT TestingT) *Spy {
	return &Spy{underlyingT: underlyingT}
}

// Helper implements the TestingT interface.
func (spy *Spy) Helper() {
	spy.m.Lock()
	defer spy.m.Unlock()

	spy.underlyingT.Helper()
	spy.records = append(spy.records, SpyTestingTRecord{Method: "Helper"})
}

// Cleanup implements the TestingT interface.
func (spy *Spy) Cleanup(cleanupFunc func()) {
	spy.m.Lock()
	defer spy.m.Unlock()

	spy.underlyingT.Cleanup(cleanupFunc)
	spy.records = append(spy.records, SpyTestingTRecord{
		Method:  "Cleanup",
		Inputs:  []any{cleanupFunc},
		Outputs: nil,
	})
}

// Fail implements the TestingT interface.
func (spy *Spy) Fail() {
	spy.m.Lock()
	defer spy.m.Unlock()

	spy.underlyingT.Fail()
	spy.records = append(spy.records, SpyTestingTRecord{Method: "Fail"})
	spy.failed = true
}

// FailNow implements the TestingT interface.
// Warning: the goroutine is not stopped.
func (spy *Spy) FailNow() {
	spy.m.Lock()
	defer spy.m.Unlock()

	spy.underlyingT.FailNow()
	spy.records = append(spy.records, SpyTestingTRecord{Method: "FailNow"})
	spy.failed = true
}

// Logf implements the TestingT interface.
func (spy *Spy) Logf(format string, args ...any) {
	spy.m.Lock()
	defer spy.m.Unlock()

	spy.underlyingT.Logf(format, args...)
	spy.records = append(spy.records, SpyTestingTRecord{
		Method:  "Logf",
		Inputs:  []any{format, args},
		Outputs: nil,
	})
	spy.logs = append(spy.logs, fmt.Sprintf(format, args...))
}

// Context implements the TestingT interface.
func (spy *Spy) Context() context.Context {
	spy.m.Lock()
	defer spy.m.Unlock()

	ctx := spy.underlyingT.Context()
	spy.records = append(spy.records, SpyTestingTRecord{
		Method:  "Context",
		Inputs:  nil,
		Outputs: []any{ctx},
	})

	return ctx
}

// Deadline implements the TestingT interface.
func (spy *Spy) Deadline() (time.Time, bool) {
	spy.m.Lock()
	defer spy.m.Unlock()

	deadline, isset := spy.underlyingT.Deadline()
	spy.records = append(spy.records, SpyTestingTRecord{
		Method:  "Deadline",
		Inputs:  nil,
		Outputs: []any{deadline, isset},
	})

	return deadline, isset
}
