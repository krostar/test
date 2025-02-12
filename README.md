# krostar/test

This library offers helper functions for writing cleaner tests in Go.
Using it encourages using regular go code instead of multiplying specific methods.

## Usage

### Assert and Require

There are only two functions exposed, `Assert` and `Require`, they both takes a boolean and make the test fail if that boolean is false.

```go
import (
    "testing"

    "github.com/krostar/test"
)

func TestSomething(t *testing.T) {
    got := 42
    want := 24

    test.Assert(t, got == want)
}
```

### Extending with custom checks

Assertions function deal with boolean. To create more complex assertions, you can simply write a function that perform that check.

```go
func myCustomCheck(t test.TestingT, foo *Object) (test.TestingT, bool) {
  if foo.bar != 42 {
    return t, false
  }

  return t, true
}

test.Assert(myCustomCheck(t, foo))
```

Some example of such functions are exposed in the check package:

- `check.Compare`: peform deep comparison using go/cmp, and return false if they have differences
- `check.Panics`: checks whether the provided function panics
- `check.ZeroValue`: checks whether the provided value is the zero value of its type

### Comparison with other libs:

#### stdlib

#### stretchr/testify

#### gotestyourself/gotest.tools

#### gotestyourself/gotest.tools

#### earthboundkid/be

#### matryer/is
