package test

import (
	"fmt"

	"github.com/krostar/test/internal/message"
)

// TestingT is an interface for testing types.
// It mimics the standard library's *testing.T.
type TestingT interface {
	Helper()
	Fail()
	FailNow()
	Logf(format string, args ...any)
}

// SuccessMessageEnabled controls whether to enable success messages logging in assert functions.
//
//nolint:gochecknoglobals // there is no clean way to deal with it, so global it is
var SuccessMessageEnabled = false

// Assert checks the provided boolean `result`.
//
// If `result` is false, it logs a detailed error message based on source code parsing
// and fails the test. The error message includes the expression that evaluated to false.
//
// Optionally, `msgAndArgs` can be provided to add custom messages to the error output.
//
// If check.SuccessMessageEnabled is true, it will log a success message even if `result` is true.
//
// Assert returns the same value as `result`.
//
// Example usage:
//
//	func Test_Something(t *testing.T) {
//		user, err := GetUser(context.Background(), "bob@example.com")
//		test.Require(t, err == nil && user != nil)
//		test.Assert(t, user.Name == "Bob" && user.Age == 42)
//	}
//
// -> Error: user.Name is not equal to "Bob", or user.Age is not equal to 42.
func Assert(t TestingT, result bool, msgAndArgs ...any) bool {
	t.Helper()

	// function that perform checks can return empty strings, don't display them
	if len(msgAndArgs) > 0 && msgAndArgs[0] == "" {
		msgAndArgs = msgAndArgs[1:]
	}

	return eval(t, result, 1, msgAndArgs...)
}

// Require is similar to Assert, but it stops the test execution immediately if `result` is false.
// Otherwise, it behaves the same as Assert.
func Require(t TestingT, result bool, msgAndArgs ...any) {
	t.Helper()

	// function that perform checks can return empty strings, don't display them
	if len(msgAndArgs) > 0 && msgAndArgs[0] == "" {
		msgAndArgs = msgAndArgs[1:]
	}

	if !eval(t, result, 1, msgAndArgs...) {
		t.FailNow()
	}
}

func eval(t TestingT, result bool, callerStackIndex int, msgAndArgs ...any) bool {
	t.Helper()

	var msg string

	if !result || (result && SuccessMessageEnabled) {
		var err error
		msg, err = message.FromBool(callerStackIndex+1, result)
		if err != nil {
			t.Logf("krostar/test internal failure: unable to get assertion message: %v", err)
		}

		switch l := len(msgAndArgs); {
		case l == 1:
			msg = fmt.Sprintf("%s [%v]", msg, msgAndArgs[0])
		case l > 1:
			msg = fmt.Sprintf("%s [%s]", msg, fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...))
		}
	}

	if len(msg) > 0 {
		if result {
			t.Logf("Success: %s", msg)
		} else {
			t.Logf("Error: %s", msg)
		}
	}

	return result
}
