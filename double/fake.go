package double

import (
	"context"
	"time"
)

// Fake implements a minimal TestingT that does nothing.
// It's useful for tests that need a TestingT but don't need its real behavior.
type Fake struct {
	o *fakeOptions
}

// NewFake creates a new Fake test double.
func NewFake(opts ...FakeOption) *Fake {
	o := &fakeOptions{
		registerCleanup: func(func()) {},
		context:         context.Background(),
	}

	for _, opt := range opts {
		opt(o)
	}

	return &Fake{o: o}
}

// Helper implements the TestingT interface.
// This is a no-op implementation.
func (Fake) Helper() {}

// Cleanup implements the TestingT interface.
// Registers a function to be called when the test completes.
func (t Fake) Cleanup(f func()) { t.o.registerCleanup(f) }

// Fail implements the TestingT interface.
// This is a no-op implementation.
func (Fake) Fail() {}

// FailNow implements the TestingT interface.
// This is a no-op implementation.
func (Fake) FailNow() {}

// Logf implements the TestingT interface.
// This is a no-op implementation.
func (Fake) Logf(string, ...any) {}

// Context implements the TestingT interface.
// Returns the context specified during creation, or background context by default.
func (t Fake) Context() context.Context {
	return t.o.context
}

// Deadline implements the TestingT interface.
// Returns the deadline specified during creation, or no deadline by default.
func (t Fake) Deadline() (time.Time, bool) {
	return t.o.deadline, !t.o.deadline.IsZero()
}
