// Package comparegotesttools_test provides a comparison between common testing approaches
// of the gotest.tools framework and the krostar/test package.
//
// This package demonstrates how to implement equivalent assertions using both
// frameworks. Each test case shows comparable functionality using each approach.
//
// The test cases cover the following assertion categories:
// - Boolean assertions
// - Equality assertions (Equal, DeepEqual, Contains, etc.)
// - Error assertions (Error, NoError, ErrorIs, ErrorType, etc.)
// - Numeric assertions (Less, Greater, etc.)
// - Skip assertions (Skip conditions)
// - Advanced polling/waiting
package comparegotesttools_test

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"math"
	"slices"
	"strings"
	"testing"
	"time"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
	"gotest.tools/v3/poll"
	"gotest.tools/v3/skip"

	"github.com/krostar/test"
	"github.com/krostar/test/check"
)

func Test_Comparison(t *testing.T) {
	/*
		Boolean assertions
	*/
	t.Run("Boolean Assertions", func(t *testing.T) {
		{
			assert.Assert(t, true)
			test.Assert(t, true)
		}

		{
			assert.Assert(t, !false)
			test.Assert(t, !false)
		}
	})

	/*
		Equality assertions:
		- Equal asserts that two values are equal
		- DeepEqual asserts that two values are deeply equal
		- Contains asserts that a string/array/slice/map contains a value
	*/
	t.Run("Equality", func(t *testing.T) {
		t.Run("Equal", func(t *testing.T) {
			type typ struct {
				A int
				B string
			}

			a := typ{
				A: 1,
				B: "abc",
			}
			b := typ{
				A: 1,
				B: "abc",
			}
			c := typ{
				A: 2,
				B: "notabc",
			}

			{
				assert.Equal(t, "abc", "abc")
				test.Assert(t, "abc" == "abc")
			}

			{
				assert.Equal(t, a, b)
				test.Assert(t, a == b)
			}

			{
				assert.Assert(t, a != c)
				test.Assert(t, a != c)
			}
		})

		t.Run("DeepEqual", func(t *testing.T) {
			a := map[string]int{"a": 1, "b": 2}
			b := map[string]int{"b": 2, "a": 1}

			assert.DeepEqual(t, a, b)
			test.Assert(check.Compare(t, a, b))
		})

		t.Run("Contains", func(t *testing.T) {
			{
				s := "hello world"
				assert.Assert(t, cmp.Contains(s, "world"))
				test.Assert(t, strings.Contains(s, "world"))
			}

			{
				slice := []string{"apple", "banana", "orange"}
				assert.Assert(t, cmp.Contains(slice, "banana"))
				test.Assert(t, slices.Contains(slice, "banana"))
			}

			{
				m := map[string]int{"a": 1, "b": 2, "c": 3}
				assert.Assert(t, cmp.Contains(m, "b"))
				test.Assert(t, slices.Contains(slices.Collect(maps.Keys(m)), "b"))
			}
		})
	})

	/*
		Error assertions:
		- Error asserts that a function returned an error
		- NoError asserts that a function returned no error
		- ErrorContains asserts that an error contains a specific substring
		- ErrorIs is a wrapper for errors.Is
		- ErrorType asserts that an error is of a specific type
	*/
	t.Run("Error Handling", func(t *testing.T) {
		t.Run("Error - NoError", func(t *testing.T) {
			err := errors.New("boom")
			var noErr error

			assert.Assert(t, err != nil)
			test.Assert(t, err != nil)

			assert.NilError(t, noErr)
			test.Assert(t, noErr == nil)
		})

		t.Run("ErrorContains", func(t *testing.T) {
			err := errors.New("boom")

			assert.ErrorContains(t, err, "boom")
			test.Assert(t, err != nil && strings.Contains(err.Error(), "boom"))
		})

		t.Run("ErrorIs", func(t *testing.T) {
			err := errors.New("boom")
			errw := fmt.Errorf("%w", err)
			errv := fmt.Errorf("%v", err)

			assert.Assert(t, cmp.ErrorIs(errw, err))
			test.Assert(t, errors.Is(errw, err))

			assert.Assert(t, !errors.Is(errv, err))
			test.Assert(t, !errors.Is(errv, err))
		})

		t.Run("ErrorType", func(t *testing.T) {
			var pathErr *fs.PathError
			err := &fs.PathError{
				Op:   "open",
				Path: "file.txt",
				Err:  errors.New("file not found"),
			}

			assert.ErrorType(t, err, &fs.PathError{})
			test.Assert(t, errors.As(err, &pathErr))
		})
	})

	/*
		Numeric assertions:
		- Assert with custom comparisons
	*/
	t.Run("Numeric Comparisons", func(t *testing.T) {
		{
			assert.Assert(t, 21 < 31)
			test.Assert(t, 21 < 31)
		}

		{
			assert.Assert(t, 21 <= 21)
			test.Assert(t, 21 <= 21)
		}

		{
			assert.Assert(t, 31 > 21)
			test.Assert(t, 31 > 21)
		}

		{
			assert.Assert(t, 31 >= 31)
			test.Assert(t, 31 >= 31)
		}

		{
			delta := 0.2
			assert.Assert(t, math.Abs(1.0-1.1) <= delta)
			test.Assert(t, math.Abs(1.0-1.1) <= delta)
		}
	})

	/*
		Skip assertions:
		- Skip allows conditionally skipping tests
	*/
	t.Run("Skip Assertions", func(t *testing.T) {
		// Using gotest.tools - simple condition
		skip.If(t, false, "This test would be skipped if condition was true")

		// Using krostar/test
		if false {
			t.Skip("This test would be skipped if condition was true")
		}
	})

	/*
		Custom assertions with cmp
	*/
	t.Run("Custom Comparisons", func(t *testing.T) {
		// Using gotest.tools with custom comparator
		customComparison := func() cmp.Comparison {
			return func() cmp.Result {
				if 2+2 == 4 {
					return cmp.ResultSuccess
				}
				return cmp.ResultFailure("math is broken")
			}
		}
		assert.Assert(t, customComparison())

		// Using krostar/test
		test.Assert(t, func() bool { return 2+2 == 4 }())
	})

	/*
		Wait assertions (poll package)
	*/
	t.Run("Wait/Poll Assertions", func(t *testing.T) {
		t.Run("Basic Polling", func(t *testing.T) {
			// gotest.tools approach with poll package
			counter := 0

			// Using WaitOn with check functions
			poll.WaitOn(t, func(t poll.LogT) poll.Result {
				counter++
				if counter >= 3 {
					return poll.Success()
				}
				return poll.Continue("waiting for counter to reach 3, current: %d", counter)
			}, poll.WithTimeout(time.Second), poll.WithDelay(10*time.Millisecond))

			// Reset counter for krostar/test approach
			counter = 0

			// krostar/test approach
			ctx, cancel := context.WithTimeout(test.Context(t), time.Second)
			defer cancel()

			check.Eventually(ctx, t, func(context.Context) error {
				counter++
				if counter >= 3 {
					return nil
				}
				return errors.New("lets continue")
			}, 10*time.Millisecond)
		})
	})
}
