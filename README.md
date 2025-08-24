# krostar/test

[![Go Reference](https://pkg.go.dev/badge/github.com/krostar/test.svg)](https://pkg.go.dev/github.com/krostar/test)
[![Go Report Card](https://goreportcard.com/badge/github.com/krostar/test)](https://goreportcard.com/report/github.com/krostar/test)
![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8.svg)

A lightweight, dependency-minimal testing library for Go that provides detailed assertion messages through source code analysis.
Designed to make your tests as clear, readable, and maintainable as the code being tested.

## Why

- **Simplicity**: Write standard Go expressions in your tests without learning a complex API
- **Clarity**: Get error messages that directly relate to your test code through source analysis
- **Minimalism**: Enjoy a small, focused API instead of numerous specialized assertions
- **Readability**: Keep your test code looking like regular Go code
- **Composability**: Easily create custom test checks

## Installation

```bash
go get github.com/krostar/test
```

## Quick start

### Assert and Require

The library exposes two main functions:

- `Assert(t, condition, [msg...])`: Reports test failure if condition is false but continues execution
- `Require(t, condition, [msg...])`: Reports test failure and stops execution immediately if condition is false

```go
package foo

import (
    "context"
    "errors"
    "testing"

    "github.com/krostar/test"
)

func Test_Something(t *testing.T) {
    got := 42
    want := 24

    test.Assert(t, got == want)

    err := errors.New("test error")
    test.Assert(t, errors.Is(err, context.Canceled), "context should have canceled and produced an error")
}
```

Output:

```sh
go test -v ./...
=== RUN   Test_Something
    test_test.go:13: Error: got is not equal to want
    test_test.go:17: Error: context.Canceled is not in the error tree of err [context should have canceled and produced an error]
--- FAIL: Test_Something (0.36s)

FAIL
```

### Automatic error messages

The library generates detailed error messages by analyzing your test expressions:

```go
// With simple expressions
test.Assert(t, user.Name == "John")
// Error: user.Name is not equal to "John"

// With compound expressions
test.Assert(t, err == nil && user != nil)
// Error: err is not nil, or user is nil

// With function calls
test.Assert(t, strings.Contains(response, "success"))
// Error: response does not contain "success"
```

### Built-in checks

The `check` package provides several built-in checks for common testing scenarios:

```go
import (
    "testing"
    "time"
    "context"

    "github.com/krostar/test"
    "github.com/krostar/test/check"
)

func TestChecks(t *testing.T) {
    // Deep comparison
    got := map[string]int{"a": 1, "b": 2}
    want := map[string]int{"b": 2, "a": 1}
    test.Assert(check.Compare(t, got, want))

    // Wait for async condition
    test.Assert(check.Eventually(test.Context(t), t, func(ctx context.Context) error {
        // Check some condition that may take time to become true
        return nil
    }, time.Millisecond*100))

    // Verify a function panics
    test.Assert(check.Panics(t, func() { panic("boom") }, nil))

    // Invert a check
    test.Assert(check.Not(check.Panics(t, func() { panic("boom") }, nil)))

    // Verify zero values
    test.Assert(check.ZeroValue(t, 0))
}
```

### Custom Checks

You can easily create custom checks to extend functionality:

```go
func IsValidUser(t test.TestingT, user *User) (test.TestingT, bool, string) {
    if user == nil {
        return t, false, "user is nil"
    }

    if user.ID == 0 {
        return t, false, "user ID should not be zero"
    }

    if user.Name == "" {
        return t, false, "user name should not be empty"
    }

    return t, true, "user is valid"
}

// Usage
test.Assert(IsValidUser(t, user))
```

## Test Doubles

The `double` package provides test doubles for testing code that uses a `testing.T` interface:

```go
import (
    "testing"

    "github.com/krostar/test"
    "github.com/krostar/test/double"
)

func TestWithSpy(t *testing.T) {
    // Create a spy on a fake to record calls
    spyT := double.NewSpy(double.NewFake())

    // Use the spy in your test
    functionUnderTest(spyT)

    // Verify calls were made
    spyT.ExpectRecords(t, false,
        double.SpyTestingTRecord{Method: "Logf", Inputs: []any{"Success: %s", double.SpyTestingTRecordIgnoreParam}},
    )

    // Verify test passed/failed
    spyT.ExpectTestToPass(t)
}
```

## Comparison with Other Testing Libraries

| Library | API Design | Implementation | Error Messages | Maintenance |
|-------------------------|------------------------------------------------|--------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------|-----------------------------------------------------------|
| **Go Standard Library** | ✅ Core functionality only | ✅ Standard go, no external dependencies | ❌ Basic messages requiring manual formatting | ✅ Part of go core |
| **stretchr/testify** | ⚠️ Very large API with many similar assertions | ❌ Heavy use of reflection<br/>⚠️ Encourages using non native (suite), or often misused (mock) packages | ✅ Detailed error messages<br/>✅ Good diffs for complex types | ✅ Actively maintained<br/>✅ Well-established in community |
| **gotest.tools/v3** | ✅ Minimalist and clear | ✅ Thin layer over go testing lib | ✅ Detailed error messages<br/>✅ Good diffs for complex types | ✅ Actively maintained |
| **matryer/is** | ⚠️ Extremely minimalist API | ✅ Thin layer over go testing lib | ⚠️ Basic but clear error messages | ❌ No longer actively maintained |
| **krostar/test** | ✅ Minimalist and clear | ✅ Thin layer over go testing lib | ✅ Detailed error messages<br/>⚠️ Diffs via go-cmp integration when using Compare check | ✅ Actively maintained |
