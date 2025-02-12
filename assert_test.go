package test_test

import (
	"testing"

	"github.com/krostar/test"
	"github.com/krostar/test/check"
)

func TestFoo(t *testing.T) {
	test.Assert(check.Panics(t, func() {
	}, nil))
}
