package check

import (
	"errors"
	"strings"
	"testing"

	gocmp "github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/krostar/test"
)

func assertCheck(t *testing.T, expectedResult bool, msgContains string, tt test.TestingT, result bool, msg string) {
	if result != expectedResult {
		t.Errorf("expected check to return %t, got %t", expectedResult, result)
	}

	if !strings.Contains(msg, msgContains) {
		t.Errorf("expected check message %s to contain %s", msg, msgContains)
	}

	if t != tt.(*testing.T) {
		t.Errorf("expected check to return the same testingT as provided")
	}
}

func Test_Compare(t *testing.T) {
	type c struct {
		a int
		B string
	}

	t.Run("ok", func(t *testing.T) {
		tt, result, msg := Compare(t, t, t, cmpopts.IgnoreUnexported(testing.T{}))
		assertCheck(t, true, "no differences", tt, result, msg)

		tt, result, msg = Compare(t, c{
			a: 42,
			B: "hello",
		}, c{
			a: 42,
			B: "hello",
		}, gocmp.AllowUnexported(c{}))
		assertCheck(t, true, "no differences", tt, result, msg)
	})

	t.Run("ko", func(t *testing.T) {
		tt, result, msg := Compare(t, 42, 21)
		assertCheck(t, false, "comparison differs", tt, result, msg)

		tt, result, msg = Compare(t, []string{"a, b, c"}, []string{"a", "c", "b"})
		assertCheck(t, false, "comparison differs", tt, result, msg)

		tt, result, msg = Compare(t, c{
			a: 42,
			B: "hello",
		}, c{
			a: 21,
			B: "hello",
		}, gocmp.AllowUnexported(c{}))
		assertCheck(t, false, "comparison differs", tt, result, msg)
	})
}

func Test_Panics(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		tt, result, msg := Panics(t, func() { panic(42) }, nil)
		assertCheck(t, true, "function panicked like expected", tt, result, msg)

		tt, result, msg = Panics(t, func() { panic(42) }, func(reason any) error {
			if reason != any(42) {
				return errors.New(reason.(string))
			}
			return nil
		})
		assertCheck(t, true, "function panicked like expected", tt, result, msg)
	})

	t.Run("ko", func(t *testing.T) {
		tt, result, msg := Panics(t, nil, nil)
		assertCheck(t, false, "function to test for panic must not be nil", tt, result, msg)

		tt, result, msg = Panics(t, func() {}, nil)
		assertCheck(t, false, "expected function to panic", tt, result, msg)

		tt, result, msg = Panics(t, func() { panic(42) }, func(any) error {
			return errors.New("boom")
		})
		assertCheck(t, false, "function panicked like expected, but reason assertion failed", tt, result, msg)
	})
}

func Test_ZeroValue(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		tt, result, msg := ZeroValue(t, 0)
		assertCheck(t, true, "0 is the zero value of type int", tt, result, msg)

		tt, result, msg = ZeroValue(t, "")
		assertCheck(t, true, `"" is the zero value of type string`, tt, result, msg)

		tt, result, msg = ZeroValue(t, [2]string{})
		assertCheck(t, true, `[2]string{"", ""} is the zero value of type [2]string`, tt, result, msg)
	})

	t.Run("ko", func(t *testing.T) {
		tt, result, msg := ZeroValue(t, 42)
		assertCheck(t, false, "expected 0 (int's zero value), got 42", tt, result, msg)
	})
}
