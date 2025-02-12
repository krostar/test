package check

import (
	"fmt"

	gocmp "github.com/google/go-cmp/cmp"

	"github.com/krostar/test"
)

// Compare checks if two values are equal using go-cmp.
// This is usually used like test.Assert(check.Compare(t, got, want)).
func Compare[T any](t test.TestingT, got, want T, gocmpOpts ...gocmp.Option) (test.TestingT, bool, string) {
	if diff := gocmp.Diff(got, want, gocmpOpts...); diff != "" {
		return t, false, "comparison differs: \n" + diff
	}
	return t, true, "no differences"
}

// Panics checks if a function panics.
// The `f` argument is the function to be tested for panic, `assertReason` is an optional function that can be used to assert on the recovered panic value.
// If `f` panics, and `assertReason` is provided and returns an error, Panics will return false and the error message.
// This is usually used like test.Assert(check.Panics(t, func(){panic("boom")}, nil)).
//
//nolint:revive // bare-return: we need to do this due because of defer
func Panics(t test.TestingT, f func(), assertReason func(reason any) error) (tt test.TestingT, result bool, msg string) { //nolint:nonamedreturns // by design of how panics works named return are required
	if f == nil {
		return t, false, "function to test for panic must not be nil"
	}

	tt = t //
	defer func() {
		reason := recover()

		if reason == nil {
			msg = "expected function to panic"
			return
		}

		if assertReason != nil {
			if reasonErr := assertReason(reason); reasonErr != nil {
				msg = fmt.Sprintf("function panicked like expected, but reason assertion failed: %v", reasonErr)
				return
			}
		}

		result = true
		msg = "function panicked like expected"
	}()

	f()

	return
}

// ZeroValue checks if a value is equal to the zero value of its type.
// This is usually used like test.Assert(check.ZeroValue(t, 0, nil)).
func ZeroValue[T comparable](t test.TestingT, v T) (test.TestingT, bool, string) {
	var zero T
	if v != zero {
		return t, false, fmt.Sprintf("expected %v (%T's zero value), got %v", zero, v, v)
	}
	return t, true, fmt.Sprintf("%#v is the zero value of type %T", v, v)
}
