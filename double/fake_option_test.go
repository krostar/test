package double

import (
	"testing"
)

func Test_FakeWithContext(t *testing.T) {
	o := new(fakeOptions)

	FakeWithContext(t.Context())(o)

	if o.context == nil {
		t.Error("o.context should not be nil")
	}
}

func Test_FakeWithRegisterCleanup(t *testing.T) {
	o := new(fakeOptions)

	called := false
	FakeWithRegisterCleanup(func(f func()) {
		called = true

		f()
	})(o)

	o.registerCleanup(func() {})

	if !called {
		t.Error("registerCleanup was not set")
	}
}
