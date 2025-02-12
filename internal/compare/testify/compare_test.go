// Package comparetestify_test provides a comparison between common testing approaches
// of the testify framework and the krostar/test package.
//
// This package demonstrates how to implement equivalent assertions using both
// frameworks. Each test case shows comparable functionality using each approach.
//
// The test cases cover the following assertion categories:
// - Boolean assertions (True/False)
// - Emptiness and length assertions (Empty, NotEmpty, Len, Zero, NotZero)
// - Nil and NotNil assertions
// - Numeric comparison assertions (Greater, GreaterOrEqual, Less, LessOrEqual, Positive, Negative)
// - Approximate numeric comparison (InDelta, InEpsilon)
// - Equality assertions (Equal, NotEqual, EqualValues, Same, NotSame)
// - Contains and NotContains assertions
// - ElementsMatch assertions
// - Subset assertions
// - Error assertions (Error, NoError, EqualError, ErrorContains, ErrorIs, ErrorAs)
// - Type-related assertions (IsType, Implements)
// - File-related assertions (FileExists, NoFileExists, DirExists, NoDirExists)
// - Regular expression assertions (Regexp, NotRegexp)
// - Time-related assertions (WithinDuration)
// - Runtime behavior assertions (Panics, NotPanics, PanicsWithValue, Eventually)
// - Condition assertions
//
// The krostar/test package does not provide equivalents for:
// - EqualValues
// - Subset / NotSubset
// - JSONEq / YamlEq
package comparetestify_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"maps"
	"math"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/krostar/test"
	"github.com/krostar/test/check"
)

func Test_Comparison(t *testing.T) {
	/*
		Boolean assertions: True/False
	*/
	t.Run("Boolean Assertions", func(t *testing.T) {
		{
			assert.True(t, true)
			test.Assert(t, true)
		}

		{
			assert.False(t, false)
			test.Assert(t, !false)
		}
	})

	/*
		Emptiness and length assertions:
		- Empty asserts that the specified object is empty (nil, "", false, 0, len() == 0)
		- NotEmpty asserts that the specified object is not empty
		- Len asserts that the specified object has specific length
		- Zero asserts that the specified object is the zero value for its type
		- NotZero asserts that the specified object is not the zero value for its type
	*/
	t.Run("Emptiness and Length", func(t *testing.T) {
		t.Run("Empty - NotEmpty", func(t *testing.T) {
			{
				var m map[string]int
				assert.Empty(t, m)
				test.Assert(t, len(m) == 0)
			}

			{
				var str string
				assert.Empty(t, str)
				test.Assert(t, str == "")
			}

			{
				var a []string
				assert.Empty(t, a)
				test.Assert(t, len(a) == 0)
			}

			{
				c := make(chan struct{})
				assert.Empty(t, c)
				test.Assert(t, len(c) == 0)
			}

			{
				m := map[string]int{"a": 1}
				assert.NotEmpty(t, m)
				test.Assert(t, len(m) != 0)
			}

			{
				str := "abc"
				assert.NotEmpty(t, str)
				test.Assert(t, str != "")
			}

			{
				a := []string{"a", "b"}
				assert.NotEmpty(t, a)
				test.Assert(t, len(a) != 0)
			}

			{
				c := make(chan struct{}, 1)
				c <- struct{}{}
				assert.NotEmpty(t, c)
				test.Assert(t, len(c) != 0)
				close(c)
			}
		})

		t.Run("Len", func(t *testing.T) {
			{
				a := []string{"a", "b"}
				assert.Len(t, a, 2)
				test.Assert(t, len(a) == 2)
			}

			{
				a := map[string]int{"a": 1}
				assert.Len(t, a, 1)
				test.Assert(t, len(a) == 1)
			}
		})

		t.Run("Zero - NotZero", func(t *testing.T) {
			{
				var i int
				assert.Zero(t, i)
				test.Assert(t, i == 0)
			}

			{
				var m map[string]int
				assert.Zero(t, m)
				test.Assert(t, m == nil)
			}

			{
				i := 42
				assert.NotZero(t, i)
				test.Assert(t, i != 0)
			}

			{
				m := make(map[string]int)
				assert.NotZero(t, m)
				test.Assert(t, m != nil)
			}
		})
	})

	/*
		Nil and NotNil assertions:
		- Nil asserts that the specified object is nil
		- NotNil asserts that the specified object is not nil
	*/
	t.Run("Nil - NotNil", func(t *testing.T) {
		{
			var m map[string]int
			assert.Nil(t, m)
			test.Assert(t, m == nil)
		}

		{
			var p *int
			assert.Nil(t, p)
			test.Assert(t, p == nil)
		}

		{
			m := make(map[string]int)
			assert.NotNil(t, m)
			test.Assert(t, m != nil)
		}

		{
			p := new(int)
			assert.NotNil(t, p)
			test.Assert(t, p != nil)
		}
	})

	/*
		Numeric comparison assertions:
		- Greater asserts that the first element is greater than the second
		- GreaterOrEqual asserts that the first element is greater than or equal to the second
		- Less asserts that the first element is less than the second
		- LessOrEqual asserts that the first element is less than or equal to the second
		- Positive asserts that the specified element is positive
		- Negative asserts that the specified element is negative
	*/
	t.Run("Numeric Comparisons", func(t *testing.T) {
		t.Run("Greater - GreaterOrEqual - Less - LessOrEqual", func(t *testing.T) {
			{
				assert.Greater(t, 31, 21)
				test.Assert(t, 31 > 21)
			}

			{
				assert.GreaterOrEqual(t, 31, 21)
				test.Assert(t, 31 >= 21)
			}

			{
				assert.GreaterOrEqual(t, 31, 31)
				test.Assert(t, 31 == 31)
			}

			{
				assert.Less(t, 21, 31)
				test.Assert(t, 21 < 31)
			}

			{
				assert.LessOrEqual(t, 21, 31)
				test.Assert(t, 21 <= 31)
			}

			{
				assert.LessOrEqual(t, 21, 21)
				test.Assert(t, 21 == 21)
			}
		})

		t.Run("Positive - Negative", func(t *testing.T) {
			{
				assert.Positive(t, 1)
				test.Assert(t, 1 > 0)
			}

			{
				assert.Negative(t, -1)
				test.Assert(t, -1 < 0)
			}
		})
	})

	/*
		Approximate numeric comparison:
		- InDelta asserts that the two numerals are within delta of each other
		- InEpsilon asserts that the two numerals are within epsilon of their expected ratio
	*/
	t.Run("Approximate Numeric Comparison", func(t *testing.T) {
		assert.InDelta(t, 1.0, 1.1, 0.2)
		test.Assert(t, math.Abs(1.0-1.1) <= 0.2)

		assert.InEpsilon(t, 1.0, 1.1, 0.2)
		test.Assert(t, math.Abs(1.0-1.1)/1.0 <= 0.2)
	})

	/*
		Equality assertions:
		- Equal asserts that two objects are equal
		- NotEqual asserts that two objects are not equal
		- EqualValues asserts that two objects are equal or convertable to the same types and equal
		- Same asserts that two pointers reference the same object
		- NotSame asserts that two pointers do not reference the same object
	*/
	t.Run("Equality", func(t *testing.T) {
		t.Run("Equal - NotEqual", func(t *testing.T) {
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

			d := new(int)
			*d = 1
			e := new(int)
			*e = 1
			f := new(int)
			*f = 2

			type typ2 struct {
				a int
				B string
			}
			g := typ2{
				a: 1,
				B: "abc",
			}
			h := typ2{
				a: 1,
				B: "abc",
			}
			i := typ2{
				a: 0,
				B: "cba",
			}

			{
				assert.Equal(t, "abc", "abc")
				test.Assert(t, "abc" == "abc")
			}

			{
				assert.Equal(t, d, e)
				test.Assert(t, *d == *e)
			}

			{
				assert.Equal(t, a, b)
				test.Assert(t, a == b)
			}

			{
				assert.Equal(t, g, h)
				test.Assert(t, g == h)
			}

			{
				assert.NotEqual(t, "abc", "cba")
				test.Assert(t, "abc" != "cba")
			}

			{
				assert.NotEqual(t, d, f)
				test.Assert(t, *d != *f)
			}

			{
				assert.NotEqual(t, a, c)
				test.Assert(t, a != c)
			}

			{
				assert.NotEqual(t, g, i)
				test.Assert(t, g != i)
			}
		})

		t.Run("Same - NotSame", func(t *testing.T) {
			p1 := new(int)
			*p1 = 1
			p2 := p1
			p3 := new(int)
			*p3 = 1

			assert.Same(t, p1, p2)
			test.Assert(t, p1 == p2)

			assert.NotSame(t, p1, p3)
			test.Assert(t, p1 != p3)
		})
	})

	/*
		Contains asserts that the specified string, list(array, slice...) or map contains the specified substring or element
		NotContains asserts that the specified string, list(array, slice...) or map does not contain the specified substring or element
	*/
	t.Run("Contains - NotContains", func(t *testing.T) {
		{
			s := "hello world"
			assert.Contains(t, s, "world")
			test.Assert(t, strings.Contains(s, "world"))
		}

		{
			slice := []string{"apple", "banana", "orange"}
			assert.Contains(t, slice, "banana")
			test.Assert(t, slices.Contains(slice, "banana"))
		}

		{
			m := map[string]int{"a": 1, "b": 2, "c": 3}
			assert.Contains(t, m, "b")
			test.Assert(t, slices.Contains(slices.Collect(maps.Keys(m)), "b"))
		}

		{
			s := "hello world"
			assert.NotContains(t, s, "universe")
			test.Assert(t, !strings.Contains(s, "universe"))
		}

		{
			slice := []string{"apple", "banana", "orange"}
			assert.NotContains(t, slice, "grape")
			test.Assert(t, !slices.Contains(slice, "grape"))
		}

		{
			m := map[string]int{"a": 1, "b": 2, "c": 3}
			assert.NotContains(t, m, "d")
			test.Assert(t, !slices.Contains(slices.Collect(maps.Keys(m)), "d"))
		}
	})

	/*
		ElementsMatch asserts that the specified listA(array, slice...) is equal to specified listB(array, slice...) ignoring the order of the elements
	*/
	t.Run("ElementsMatch", func(t *testing.T) {
		a := []int{1, 2, 3}
		b := []int{3, 2, 1}
		assert.ElementsMatch(t, a, b)
		test.Assert(t, slices.Equal(slices.Sorted(slices.Values(a)), slices.Sorted(slices.Values(b))))
	})

	/*
		Error assertions:
		- Error asserts that a function returned an error (i.e. not nil)
		- NoError asserts that a function returned no error (i.e. nil)
		- EqualError asserts that a function returned an error with a specific message
		- ErrorContains asserts that a function returned an error containing a specific substring
		- ErrorIs is a wrapper for errors.Is
		- ErrorAs is a wrapper for errors.As
	*/
	t.Run("Error Handling", func(t *testing.T) {
		t.Run("Error - NoError", func(t *testing.T) {
			err := errors.New("boom")
			var noErr error

			assert.Error(t, err)
			test.Assert(t, err != nil)

			assert.NoError(t, noErr)
			test.Assert(t, noErr == nil)
		})

		t.Run("EqualError - ErrorContains", func(t *testing.T) {
			err := errors.New("boom")

			assert.EqualError(t, err, "boom")
			test.Assert(t, err != nil && err.Error() == "boom")

			assert.ErrorContains(t, err, "boom")
			test.Assert(t, err != nil && strings.Contains(err.Error(), "boom"))
		})

		t.Run("ErrorIs - NotErrorIs", func(t *testing.T) {
			err := errors.New("boom")
			errw := fmt.Errorf("%w", err)
			errv := fmt.Errorf("%v", err)

			assert.ErrorIs(t, errw, err)
			test.Assert(t, errors.Is(errw, err))

			assert.NotErrorIs(t, errv, err)
			test.Assert(t, !errors.Is(errv, err))
		})

		t.Run("ErrorAs - NotErrorAs", func(t *testing.T) {
			err := fmt.Errorf("%w", net.ErrClosed)

			assert.ErrorAs(t, err, new(interface{ Timeout() bool }))
			test.Assert(t, errors.As(err, new(interface{ Timeout() bool })))

			assert.NotErrorAs(t, err, new(interface{ Foo() bool }))
			test.Assert(t, !errors.As(err, new(interface{ Foo() bool })))
		})
	})

	/*
		Type-related assertions:
		- IsType asserts that an object is of the same type as the specified object
		- Implements asserts that an object implements an interface
	*/
	t.Run("Type Assertions", func(t *testing.T) {
		{
			var a int
			assert.IsType(t, 42, a)
			test.Assert(t, fmt.Sprintf("%T", a) == fmt.Sprintf("%T", 42))
		}

		{
			var reader io.Reader = strings.NewReader("hello")
			assert.Implements(t, (*io.Reader)(nil), reader)
			test.Assert(t, func() bool {
				_, ok := reader.(io.Reader)
				return ok
			}())
		}
	})

	/*
		File-related assertions:
		- FileExists asserts that a file exists
		- NoFileExists asserts that a file does not exist
		- DirExists asserts that a directory exists
		- NoDirExists asserts that a directory does not exist
	*/
	t.Run("File Operations", func(t *testing.T) {
		t.Run("FileExists - NoFileExists", func(t *testing.T) {
			file, err := os.Create(filepath.Join(t.TempDir(), "foo.txt"))
			test.Assert(t, err == nil)
			fileName := file.Name()
			test.Assert(t, file.Close() == nil)

			assert.FileExists(t, fileName)
			test.Assert(t, func() bool {
				_, err := os.Stat(fileName)
				return err == nil
			}())

			test.Assert(t, os.Remove(fileName) == nil)

			assert.NoFileExists(t, fileName)
			test.Assert(t, func() bool {
				_, err := os.Stat(fileName)
				return os.IsNotExist(err)
			}())
		})

		t.Run("DirExists - NoDirExists", func(t *testing.T) {
			dir := t.TempDir()

			assert.DirExists(t, dir)
			test.Assert(t, func() bool {
				info, err := os.Stat(dir)
				return err == nil && info.IsDir()
			}())

			test.Assert(t, os.Remove(dir) == nil)

			assert.NoDirExists(t, dir)
			test.Assert(t, func() bool {
				_, err := os.Stat(dir)
				return os.IsNotExist(err)
			}())
		})
	})

	/*
		Regular expression assertions:
		- Regexp asserts that a string matches a regular expression
		- NotRegexp asserts that a string does not match a regular expression
	*/
	t.Run("Regexp Assertions", func(t *testing.T) {
		{
			// Regular expression matching
			assert.Regexp(t, `^\d+$`, "12345")
			test.Assert(t, func() bool {
				matched, err := regexp.MatchString(`^\d+$`, "12345")
				return err == nil && matched
			}())

			assert.NotRegexp(t, `^\d+$`, "abc")
			test.Assert(t, func() bool {
				matched, err := regexp.MatchString(`^\d+$`, "abc")
				return err == nil && !matched
			}())
		}
	})

	/*
		Time-related assertions:
		- WithinDuration asserts that the two times are within duration delta of each other
	*/
	t.Run("Time Operations", func(t *testing.T) {
		{
			now := time.Now()
			fiveSecondsLater := now.Add(5 * time.Second)
			assert.WithinDuration(t, now, fiveSecondsLater, 10*time.Second)
			test.Assert(t, fiveSecondsLater.Sub(now) <= 10*time.Second)
		}
	})

	/*
		Runtime behavior assertions:
		- Panics asserts that the code inside the specified function panics
		- NotPanics asserts that the code inside the specified function does NOT panic
		- PanicsWithValue asserts that the code inside the specified function panics with the expected panic value
		- Eventually asserts that given condition will be met in waitFor time
	*/
	t.Run("Runtime Behavior", func(t *testing.T) {
		t.Run("Panics - NotPanics - PanicsWithValue", func(t *testing.T) {
			{
				assert.Panics(t, func() { panic("boom") })
				test.Assert(check.Panics(t, func() { panic("boom") }, nil))
			}

			{
				assert.NotPanics(t, func() { /* do nothing */ })
				test.Assert(check.Not(check.Panics(t, func() { /* do nothing */ }, nil)))
			}

			{
				assert.PanicsWithValue(t, "specific panic", func() { panic("specific panic") })
				test.Assert(check.Panics(t, func() { panic("specific panic") }, func(reason any) error {
					if str, ok := reason.(string); !ok || str != "specific panic" {
						return errors.New("not the panic we expected")
					}
					return nil
				}))
			}
		})

		t.Run("Eventually", func(t *testing.T) {
			{
				counter := 0
				assert.Eventually(t, func() bool {
					counter++
					return counter >= 3
				}, time.Second, 10*time.Millisecond)
			}

			{
				ctx, cancel := context.WithTimeout(test.Context(t), time.Second)
				defer cancel()

				counter := 0

				check.Eventually(ctx, t, func(context.Context) error {
					counter++
					if counter >= 3 {
						return nil
					}
					return errors.New("lets continue")
				}, 10*time.Millisecond)
			}
		})
	})

	t.Run("Condition Assertions", func(t *testing.T) {
		condition := func() bool { return true }

		assert.Condition(t, condition)
		test.Assert(t, condition())
	})
}
