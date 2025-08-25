package check

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	gocmp "github.com/google/go-cmp/cmp"
	gocmpopts "github.com/google/go-cmp/cmp/cmpopts"

	"github.com/krostar/test"
)

func Test_Compare(t *testing.T) {
	type c struct {
		a int
		B string
	}

	t.Run("ok", func(t *testing.T) {
		tt, result, msg := Compare(t, t, t, gocmpopts.IgnoreUnexported(testing.T{}))
		assertCheck(t, tt, result, true, msg, "no differences")

		tt, result, msg = Compare(t, c{
			a: 42,
			B: "hello",
		}, c{
			a: 42,
			B: "hello",
		}, gocmp.AllowUnexported(c{}))
		assertCheck(t, tt, result, true, msg, "no differences")
	})

	t.Run("ko", func(t *testing.T) {
		tt, result, msg := Compare(t, 42, 21)
		assertCheck(t, tt, result, false, msg, "comparison differs")

		tt, result, msg = Compare(t, []string{"a, b, c"}, []string{"a", "c", "b"})
		assertCheck(t, tt, result, false, msg, "comparison differs")

		tt, result, msg = Compare(t, c{
			a: 42,
			B: "hello",
		}, c{
			a: 21,
			B: "hello",
		}, gocmp.AllowUnexported(c{}))
		assertCheck(t, tt, result, false, msg, "comparison differs")
	})
}

func Test_Eventually(t *testing.T) {
	t.Run("success after retries", func(t *testing.T) {
		retries := 0

		ctx, cancel := context.WithTimeout(t.Context(), 500*time.Millisecond)
		defer cancel()

		tt, result, msg := Eventually(ctx, t, func(context.Context) error {
			defer func() { retries++ }()

			if retries < 2 {
				return fmt.Errorf("boom %d", retries)
			}

			return nil
		}, time.Millisecond*10)

		assertCheck(t, tt, result, true, msg, "check passed")

		if retries < 2 {
			t.Errorf("expected at least 2 retries, got %d", retries)
		}
	})

	t.Run("timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
		defer cancel()

		tt, result, msg := Eventually(ctx, t, func(context.Context) error {
			return errors.New("always fails")
		}, time.Millisecond*10)

		assertCheck(t, tt, result, false, msg, "context is expired", "always fails")
	})

	t.Run("immediate success", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
		defer cancel()

		tt, result, msg := Eventually(ctx, t, func(context.Context) error {
			return nil
		}, time.Millisecond*10)

		assertCheck(t, tt, result, true, msg, "check passed", "0 retries")
	})
}

func Test_Not(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		tt, result, msg := Not(t, true, "foo")
		assertCheck(t, tt, result, false, msg, "foo; and the result was inverted")
	})

	t.Run("false", func(t *testing.T) {
		tt, result, msg := Not(t, false, "bar")
		assertCheck(t, tt, result, true, msg, "bar; and the result was inverted")
	})

	t.Run("check", func(t *testing.T) {
		tt, result, msg := Not(Compare(t, "foo", "bar"))
		assertCheck(t, tt, result, true, msg, "foo", "bar", "; and the result was inverted")
	})

	t.Run("empty message", func(t *testing.T) {
		tt, result, msg := Not(t, false, "")
		assertCheck(t, tt, result, true, msg)
	})

	t.Run("non string message", func(t *testing.T) {
		tt, result, msg := Not(t, false, 42)
		assertCheck(t, tt, result, true, msg, "42")
	})
}

func Test_Panics(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		tt, result, msg := Panics(t, func() { panic(42) }, nil)
		assertCheck(t, tt, result, true, msg, "function panicked like expected")

		tt, result, msg = Panics(t, func() { panic(42) }, func(reason any) error {
			if reason != any(42) {
				return errors.New(reason.(string))
			}
			return nil
		})

		assertCheck(t, tt, result, true, msg, "function panicked like expected")
	})

	t.Run("ko", func(t *testing.T) {
		tt, result, msg := Panics(t, nil, nil)
		assertCheck(t, tt, result, false, msg, "function to test for panic must not be nil")

		tt, result, msg = Panics(t, func() {}, nil)
		assertCheck(t, tt, result, false, msg, "expected function to panic")

		tt, result, msg = Panics(t, func() { panic(42) }, func(any) error {
			return errors.New("boom")
		})
		assertCheck(t, tt, result, false, msg, "function panicked like expected, but reason assertion failed")
	})
}

func Test_ZeroValue(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		tt, result, msg := ZeroValue(t, 0)
		assertCheck(t, tt, result, true, msg, "0 is the zero value of type int")

		tt, result, msg = ZeroValue(t, "")
		assertCheck(t, tt, result, true, msg, `"" is the zero value of type string`)

		tt, result, msg = ZeroValue(t, [2]string{})
		assertCheck(t, tt, result, true, msg, `[2]string{"", ""} is the zero value of type [2]string`)
	})

	t.Run("ko", func(t *testing.T) {
		tt, result, msg := ZeroValue(t, 42)
		assertCheck(t, tt, result, false, msg, "expected 0 (int's zero value), got 42")
	})
}

func assertCheck(t *testing.T, tt test.TestingT, result, expectedResult bool, msg string, msgContains ...string) {
	t.Helper()

	if t != tt.(*testing.T) {
		t.Error("expected check to return the same testingT as provided")
	}

	if result != expectedResult {
		t.Errorf("expected check to return %t, got %t with message %q", expectedResult, result, msg)
	}

	for _, m := range msgContains {
		if !strings.Contains(msg, m) {
			t.Errorf("expected check message %q to contain %q", msg, m)
		}
	}
}
