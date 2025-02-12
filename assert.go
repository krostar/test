// Package test provides a lightweight assertion library for Go tests.
//
// It offers simple assertion functions like Assert and Require which provide
// detailed error messages by analyzing the expression being tested.
package test

import (
	"flag"
	"fmt"

	"github.com/krostar/test/internal"
	"github.com/krostar/test/internal/message"
)

// TestingT is an interface for testing types.
// It mimics the standard library's *testing.T.
type TestingT internal.TestingT

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

	logResult(t, result, 1, msgAndArgs...)
	if !result {
		t.Fail()
	}

	return result
}

// Require stops the test execution immediately if `result` is false.
// Otherwise, it behaves the same as Assert.
func Require(t TestingT, result bool, msgAndArgs ...any) {
	t.Helper()

	logResult(t, result, 1, msgAndArgs...)
	if !result {
		t.FailNow()
	}
}

//nolint:gochecknoglobals // there is no clean way to deal with it, so global it is
var (
	// SuccessMessageEnabled controls whether to enable success messages logging in assert functions.
	SuccessMessageEnabled     = false
	_flagEnableSuccessMessage = flag.Bool("check.display-success-messages", false, "Whether to print messages in passing tests")
)

// logResult handles the logging of test results, with details about the assertion.
// It's used internally by Assert and Require functions.
//
// The function performs several tasks:
//   - Retrieves the source code expression that was evaluated from the caller's location
//   - Formats an appropriate message explaining what passed or failed
//   - Adds any custom messages provided by the caller
//   - Logs the resulting message as either a success or error message
func logResult(t TestingT, result bool, callerStackIndex int, msgAndArgs ...any) {
	t.Helper()

	// function that perform checks can return empty strings, don't display them
	if len(msgAndArgs) > 0 && msgAndArgs[0] == "" {
		msgAndArgs = msgAndArgs[1:]
	}

	var msg string

	if (result && (SuccessMessageEnabled || *_flagEnableSuccessMessage)) || !result {
		var err error
		msg, err = message.FromBool(callerStackIndex+1, result)
		if err != nil {
			t.Logf("krostar/test internal failure: unable to get assertion message: %v", err)
		}

		switch l := len(msgAndArgs); {
		case l == 1:
			msg = fmt.Sprintf("%s [%v]", msg, msgAndArgs[0])
		case l > 1:
			if format, ok := msgAndArgs[0].(string); ok {
				msg = fmt.Sprintf("%s [%s]", msg, fmt.Sprintf(format, msgAndArgs[1:]...))
			} else {
				msg = fmt.Sprintf("%s %v", msg, msgAndArgs)
			}
		}
	}

	if msg != "" {
		if result {
			t.Logf("Success: %s", msg)
		} else {
			t.Logf("Error: %s", msg)
		}
	}
}
