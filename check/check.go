package check

import (
	"context"
	"errors"
	"fmt"
	"time"

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

// Eventually repeatedly executes a check function until it succeeds or the context expires.
//
// The function continuously retries the check until one of the following occurs:
// - The check function returns nil (success)
// - The provided context is canceled or times out
//
// This is typically used for asynchronous tests that may take time to reach the desired state.
//
//	Example: test.Assert(check.Eventually(ctx, test.Context(t), func(ctx context.Context) error {
//		// ...
//	}, time.Millisecond*100))
func Eventually(ctx context.Context, t test.TestingT, check func(context.Context) error, timeBetweenRetries time.Duration) (test.TestingT, bool, string) {
	startedAt := time.Now()
	ticker := time.NewTimer(0)
	tryC := make(chan struct{}, 1)

	var (
		errs    [2]error
		retries uint
	)

	for {
		select {
		case <-ctx.Done():
			return t, false, fmt.Sprintf("check did not pass in %s with %d retries and now context is expired, last two errors: %s", time.Since(startedAt).String(), retries, errors.Join(errs[0], errs[1]))

		case <-tryC:
			if err := check(ctx); err != nil {
				errs[retries%2] = err
			} else {
				return t, true, fmt.Sprintf("check passed in %s with %d retries", time.Since(startedAt).String(), retries)
			}

			retries++

			ticker.Reset(timeBetweenRetries)

		case <-ticker.C:
			select {
			case tryC <- struct{}{}:
			default:
			}
		}
	}
}

// Not inverts the result of a boolean test check.
//
// This function is typically used with other check functions to negate their results.
// When used with other check functions, it will invert their pass/fail result.
//
// Example:
//
//	test.Assert(check.Not(check.Panics(t, func() { /* code that should not panic */ }, nil)))
func Not(t test.TestingT, result bool, msgAndArgs ...any) (test.TestingT, bool, string) {
	t.Helper()

	if len(msgAndArgs) > 0 && msgAndArgs[0] == "" {
		msgAndArgs = msgAndArgs[1:]
	}

	var msg string
	if len(msgAndArgs) > 0 {
		if str, ok := msgAndArgs[0].(string); ok {
			msg = fmt.Sprintf(str+"; and the result was inverted", msgAndArgs[1:]...)
		} else {
			msg = "NOT " + fmt.Sprint(msgAndArgs...)
		}
	}

	return t, !result, msg
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

	tt = t

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
