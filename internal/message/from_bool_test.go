package message

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	. "strings" //nolint:revive // duplicated-imports: on-purpose duplicated import to check ast parsing
	"testing"

	"github.com/krostar/test/internal/code"
)

func TestMain(m *testing.M) {
	code.InitPackageASTCache(".")
	m.Run()
}

func Test_FromBool(t *testing.T) {
	for group, tests := range map[string][]struct {
		enableBreakpoint bool // configure IDE to break when enabled

		getResult       func() (string, error)
		expectedMessage string
		expectedError   string
	}{
		"basic types aka bool": {
			{
				getResult:       func() (string, error) { return FromBool(0, true) },
				expectedMessage: "literal true",
			},
			{
				getResult:       func() (string, error) { return FromBool(0, false) },
				expectedMessage: "literal false",
			},
			{
				getResult:       func() (string, error) { return FromBool(0, !false) },
				expectedMessage: "!false",
			},
			{
				getResult:       func() (string, error) { return FromBool(0, !true) },
				expectedMessage: "!true",
			},
			{
				getResult: func() (string, error) {
					var b bool
					return FromBool(0, b)
				},
				expectedMessage: "var b is false",
			},
			{
				getResult: func() (string, error) {
					b := true
					return FromBool(0, b)
				},
				expectedMessage: "var b is true",
			},
			{
				getResult: func() (string, error) {
					b := false
					return FromBool(0, !b)
				},
				expectedMessage: "!b",
			},
		},

		"binary operations": {
			{
				getResult:       func() (string, error) { return FromBool(0, 1 == 1) },
				expectedMessage: "1 is equal to 1",
			},
			{
				getResult:       func() (string, error) { return FromBool(0, 1 == 2) },
				expectedMessage: "1 is not equal to 2",
			},
			{
				getResult:       func() (string, error) { return FromBool(0, 2 > 3) },
				expectedMessage: "2 is less or equal than 3",
			},
			{
				getResult:       func() (string, error) { return FromBool(0, 2 < 3) },
				expectedMessage: "2 is less than 3",
			},
			{
				getResult:       func() (string, error) { return FromBool(0, 2 >= 3) },
				expectedMessage: "2 is less than 3",
			},
			{
				getResult:       func() (string, error) { return FromBool(0, 2 <= 3) },
				expectedMessage: "2 is less or equal than 3",
			},
			{
				getResult:       func() (string, error) { return FromBool(0, 3 > 2) },
				expectedMessage: "3 is greater than 2",
			},
			{
				getResult:       func() (string, error) { return FromBool(0, 3 < 2) },
				expectedMessage: "3 is greater or equal than 2",
			},
			{
				getResult:       func() (string, error) { return FromBool(0, 3 >= 2) },
				expectedMessage: "3 is greater or equal than 2",
			},
			{
				getResult:       func() (string, error) { return FromBool(0, 3 <= 2) },
				expectedMessage: "3 is greater than 2",
			},
			{
				getResult: func() (string, error) {
					a := 2
					b := 3
					return FromBool(0, a >= b)
				},
				expectedMessage: "a is less than b",
			},
			{
				getResult: func() (string, error) {
					a := 2
					b := 3
					return FromBool(0, a <= b)
				},
				expectedMessage: "a is less or equal than b",
			},
			{
				getResult: func() (string, error) {
					var err error
					return FromBool(0, err == nil)
				},
				expectedMessage: "err is nil",
			},
			{
				getResult: func() (string, error) {
					var err error
					return FromBool(0, err != nil)
				},
				expectedMessage: "err is nil",
			},
			{
				getResult: func() (string, error) {
					err := errors.New("boom")
					return FromBool(0, err != nil)
				},
				expectedMessage: "err is not nil",
			},
			{
				getResult:       func() (string, error) { return FromBool(0, true && !false) },
				expectedMessage: "literal true && !false",
			},
			{
				getResult: func() (string, error) {
					var err error
					return FromBool(0, err == nil && strings.Contains("boom", "boo"))
				},
				expectedMessage: `err is nil && "boom" contains "boo"`,
			},
			{
				getResult: func() (string, error) {
					var err error
					return FromBool(0, err != nil && errors.Is(err, os.ErrDeadlineExceeded))
				},
				expectedMessage: `err is nil && os.ErrDeadlineExceeded is not in the error tree of err`,
			},
			{
				getResult: func() (string, error) {
					var err error
					return FromBool(0, err != nil || !errors.Is(err, os.ErrDeadlineExceeded))
				},
				expectedMessage:  `err is not nil || os.ErrDeadlineExceeded is not in the error tree of err`,
				enableBreakpoint: true,
			},
		},

		"foo": {
			{
				getResult: func() (string, error) {
					f := func() bool { return true }
					return FromBool(0, f() == true)
				},
				expectedMessage: "fezfze",
			},
		},

		"functions": {
			{
				getResult: func() (string, error) {
					f := func() bool { return true }
					return FromBool(0, f())
				},
				expectedMessage: "function f() returned true",
			},
			{
				getResult:       func() (string, error) { return FromBool(0, Contains("abc", "a")) },
				expectedMessage: `"abc" contains "a"`,
			},
			{
				getResult:       func() (string, error) { return FromBool(0, strings.Contains("abc", "fuu")) },
				expectedMessage: `"abc" does not contain "fuu"`,
			},
			{
				getResult: func() (string, error) {
					err := errors.New("boom")
					return FromBool(0, errors.Is(err, os.ErrDeadlineExceeded))
				},
				expectedMessage: `os.ErrDeadlineExceeded is not in the error tree of err`,
			},
			{
				getResult: func() (string, error) {
					err := fmt.Errorf("boom: %w", os.ErrDeadlineExceeded)
					return FromBool(0, errors.Is(err, os.ErrDeadlineExceeded))
				},
				expectedMessage: `err's error tree contains os.ErrDeadlineExceeded`,
			},
			{
				getResult: func() (string, error) {
					err := errors.New("boom")
					type boomErr interface{ Boom() }
					var bErr boomErr
					return FromBool(0, errors.As(err, &bErr))
				},
				expectedMessage: `err cannot be defined as *github.com/krostar/test/internal/message.boomErr`,
			},
			{
				getResult: func() (string, error) {
					err := fmt.Errorf("boom: %w", new(os.PathError))
					type timeoutErr interface {
						Timeout() bool
					}
					var tErr timeoutErr
					return FromBool(0, errors.As(err, &tErr))
				},
				expectedMessage: `err can be defined as *github.com/krostar/test/internal/message.timeoutErr`,
			},
		},
	} {
		t.Run(group, func(t *testing.T) {
			for idx, test := range tests {
				t.Run(strconv.Itoa(idx), func(t *testing.T) {
					switch msg, err := test.getResult(); {
					case test.expectedError == "" && err != nil:
						t.Errorf("expected success but got error: %v", err)
					case test.expectedError != "" && err == nil:
						t.Errorf("expected failure with following error: %s", test.expectedError)
					case test.expectedError != "" && err != nil && !strings.Contains(err.Error(), test.expectedError):
						t.Errorf("expected failure with message: %s, got %s", test.expectedError, err.Error())
					case test.expectedError == "" && err == nil && !strings.Contains(msg, test.expectedMessage):
						t.Errorf("expected %q to contain %q", msg, test.expectedMessage)
					}
				})
			}
		})
	}
}
