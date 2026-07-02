package double

import (
	"context"
)

// FakeOption is a function that configures a Fake instance.
// It follows the functional options pattern for configuring the Fake test double.
type FakeOption func(o *fakeOptions)

// FakeWithContext sets a specific context for a Fake.
// This replaces the default background context returned by Context().
func FakeWithContext(ctx context.Context) FakeOption {
	return func(o *fakeOptions) { o.context = ctx }
}

// FakeWithRegisterCleanup configures the cleanup registration function for a Fake.
// This allows tests to capture or control the behavior of cleanup registrations.
func FakeWithRegisterCleanup(f func(func())) FakeOption {
	return func(o *fakeOptions) { o.registerCleanup = f }
}

type fakeOptions struct {
	registerCleanup func(func())
	context         context.Context //nolint:containedctx // we store a context so fake can return it
}
