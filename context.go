package test

import (
	"context"
	"time"
)

// Context returns a context.Context derived from the test's context.
//
// This function extracts the context from the TestingT and handles timeout management:
// - If the test has a deadline, it adjusts the deadline to allow for clean shutdown
// - It reserves a small portion of time before the test deadline to allow for proper cleanup
// - The returned context will be automatically canceled when the test completes
func Context(t TestingT) context.Context {
	ctx := t.Context()

	if deadline, isset := t.Deadline(); isset {
		// deadline is in the future, timeout duration is expected to be negative
		panicTimeout := -time.Since(deadline)

		// because test panics after deadline, we anticipate 1% of timeout duration
		// to properly quit, it helps to have test failure messages instead of stack traces
		cleanDuration := time.Duration(int64(float64(panicTimeout) * 0.01))
		// we don't need to reserve much time, one second should be more than enough
		if cleanDuration > time.Second {
			cleanDuration = time.Second
		}

		// add negative value, meaning subtract time to deadline
		deadline = deadline.Add(-cleanDuration)

		var cancel context.CancelFunc

		ctx, cancel = context.WithDeadline(ctx, deadline)
		t.Cleanup(cancel)
	}

	return ctx
}
