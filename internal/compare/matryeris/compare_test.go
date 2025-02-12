// Package comparematryeris_test provides a comparison between common testing approaches
// of the matryer/is framework and the krostar/test package.
//
// This package demonstrates how to implement equivalent assertions using both
// frameworks. Each test case shows comparable functionality using each approach.
//
// The matryer/is package is designed with minimalism in mind - it has very few methods
// but provides high-quality, informative failure messages with expected vs actual values.
package comparematryeris_test

import (
	"maps"
	"slices"
	"strings"
	"testing"

	"github.com/matryer/is"

	"github.com/krostar/test"
)

func Test_Comparison(t *testing.T) {
	/*
		Equality assertions
	*/
	t.Run("Equality", func(t *testing.T) {
		t.Run("Equal", func(t *testing.T) {
			is := is.New(t)

			is.Equal("hello", "hello")
			test.Assert(t, "hello" == "hello")

			type person struct {
				Name string
				Age  int
			}

			p1 := person{Name: "Alice", Age: 30}
			p2 := person{Name: "Alice", Age: 30}

			is.Equal(p1, p2)
			test.Assert(t, p1 == p2)

			s1 := []string{"a", "b", "c"}
			s2 := []string{"a", "b", "c"}

			is.Equal(s1, s2)
			test.Assert(t, slices.Equal(s1, s2))

			m1 := map[string]int{"a": 1, "b": 2}
			m2 := map[string]int{"a": 1, "b": 2}

			is.Equal(m1, m2)
			test.Assert(t, maps.Equal(m1, m2))
		})

		t.Run("True", func(t *testing.T) {
			is := is.New(t)

			is.True(1 == 1)
			is.True(strings.Contains("hello world", "world"))

			test.Assert(t, 1 == 1)
			test.Assert(t, strings.Contains("hello world", "world"))
		})

		t.Run("False", func(t *testing.T) {
			is := is.New(t)

			is.True(!(1 == 2))
			is.True(!strings.Contains("hello world", "xyz"))

			test.Assert(t, !(1 == 2))
			test.Assert(t, !strings.Contains("hello world", "xyz"))
		})
	})

	/*
		Error handling
	*/
	t.Run("Error Handling", func(t *testing.T) {
		is := is.New(t)

		var nilErr error = nil

		is.NoErr(nilErr)
		test.Assert(t, nilErr == nil)
	})
}
