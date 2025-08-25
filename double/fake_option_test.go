package double

import (
	"testing"
	"time"
)

func Test_FakeWithContext(t *testing.T) {
	o := new(fakeOptions)

	FakeWithContext(t.Context())(o)

	if o.context == nil {
		t.Error("o.context should not be nil")
	}
}

func Test_FakeWithDeadline(t *testing.T) {
	o := new(fakeOptions)

	FakeWithDeadline(time.Now())(o)

	if o.deadline.IsZero() {
		t.Error("o.deadline should not be zero")
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
