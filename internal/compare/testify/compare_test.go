package compare_testify_test

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"testing"

	gocmp "github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"

	"github.com/krostar/test"
)

func Test_Oneliners(t *testing.T) {
	/*
		Documentation states that:
			- Empty asserts that the specified object is empty.
				I.e. nil, "", false, 0 or either a slice or a channel with len == 0.
			- Len asserts that the specified object has specific length.
				Len also fails if the object has a type that len() not accept.
	*/
	t.Run("Empty - NotEmpty - Len", func(t *testing.T) {
		{
			var m map[string]int
			assert.Empty(t, m)
			test.Assert(t, len(m) == 0)
		}

		{
			var str string
			assert.Empty(t, str)
			test.Assert(t, len(str) == 0)
		}

		{
			var b bool
			assert.Empty(t, b)
			test.Assert(t, b == false)
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
			test.Assert(t, len(str) != 0)
		}

		{
			b := true
			assert.NotEmpty(t, b)
			test.Assert(t, b)
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

		{
			a := []string{"a", "b"}
			assert.Len(t, a, 2)
			test.Assert(t, len(a) == 2)
		}
	})

	/*
		Documentation states that ErrorAs is a wrapper for errors.As.
	*/
	t.Run("ErrorAs - NotErrorAs", func(t *testing.T) {
		err := fmt.Errorf("%w", net.ErrClosed)
		{
			assert.ErrorAs(t, err, new(interface{ Timeout() bool }))
			test.Assert(t, errors.As(err, new(interface{ Timeout() bool })))
		}

		{
			assert.NotErrorAs(t, err, new(interface{ Foo() bool }))
			test.Assert(t, !errors.As(err, new(interface{ Foo() bool })))
		}
	})

	/*
		Documentation states that ErrorIs is a wrapper for errors.Is.
	*/
	t.Run("ErrorIs - NotErrorIs", func(t *testing.T) {
		err := errors.New("boom")
		errw := fmt.Errorf("%w", err)
		errv := fmt.Errorf("%v", err)

		{
			assert.ErrorIs(t, errw, err)
			test.Assert(t, errors.Is(errw, err))
		}

		{
			assert.NotErrorIs(t, errv, err)
			test.Assert(t, !errors.Is(errv, err))
		}
	})

	/*
		Documentation states that:
			- EqualError asserts that a function returned an error (i.e. not `nil`)
				and that it is equal to the provided error.
			- ErrorContains asserts that a function returned an error (i.e. not `nil`)
				and that the error contains the specified substring.
	*/

	t.Run("EqualErrors - ErrorContains", func(t *testing.T) {
		err := errors.New("boom")
		{
			assert.EqualError(t, err, "boom")
			test.Assert(t, err != nil && err.Error() == "boom")
		}

		{
			assert.ErrorContains(t, err, "boom")
			test.Assert(t, err != nil && strings.Contains(err.Error(), "boom"))
		}
	})

	/*
		Documentation states that True asserts that the specified value is true.
	*/
	t.Run("True - False", func(t *testing.T) {
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
		Documentation states that Nil asserts that the specified object is nil.
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
		Documentation states that:
			- Greater asserts that the first element is greater than the second.
			- GreaterOrEqual asserts that the first element is greater than or equal to the second.
			- Less asserts that the first element is less than the second.
			- LessOrEqual asserts that the first element is less than or equal to the second.
	*/
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

	/*
		Documentation states that Same asserts that two pointers reference the same object.
	*/
	t.Run("Same - NotSame", func(t *testing.T) {
		a := new(int)
		*a = 1
		b := a
		c := new(int)
		*c = 1

		{
			assert.Same(t, a, b)
			test.Assert(t, a == b)
		}

		{
			assert.NotSame(t, a, c)
			test.Assert(t, a != c)
		}
	})

	/*
		Documentation states that Positive asserts that the specified element is positive
	*/
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

	/*
		Documentation states that Error asserts that a function returned an error (i.e. not `nil`).
	*/
	t.Run("Error - NoError", func(t *testing.T) {
		err := errors.New("boom")
		var noErr error

		{
			assert.Error(t, err)
			test.Assert(t, err != nil)
		}

		{
			assert.NoError(t, noErr)
			test.Assert(t, noErr == nil)
		}
	})

	/*
		Documentation states that Equal asserts that two objects are equal.
			Pointer variable equality is determined based on the equality of the
			referenced values (as opposed to the memory addresses).
	*/
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
			test.Assert(t, gocmp.Equal(g, h, gocmp.AllowUnexported(typ2{})))
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
			test.Assert(t, !gocmp.Equal(g, i, gocmp.AllowUnexported(typ2{})))
		}
	})
}

/*
	Zero
	NotZero

	Contains
	NotContains

	Panics
	NotPanics
	PanicsWithValue


	ElementsMatch
	EqualExportedValues
	EqualValues
	Eventually
	Exactly
	FileExists
	Implements
	InDelta
	InEpsilon
	IsType
	JSONEq
	NoFileExists
	NotRegexp
	ObjectsAreEqual
	ObjectsAreEqualValues
	Regexp
	WithinDuration
	WithinRange
	YAMLEq
*/
