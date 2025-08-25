package internal

import (
	"context"
	"time"
)

// TestingT is an interface for testing types.
// It mimics the standard library's *testing.T.
// It is placed here to avoid imports cycles.
type TestingT interface {
	Helper()
	Cleanup(f func())

	Fail()
	FailNow()

	Log(args ...any)
	Logf(format string, args ...any)

	Context() context.Context
	Deadline() (time.Time, bool)
}
