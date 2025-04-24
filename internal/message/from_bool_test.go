package message

import (
	"bytes"
	"context"
	"errors"
	"go/ast"
	"maps"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"testing"

	"golang.org/x/tools/go/packages"

	"github.com/krostar/test/internal/code"
)

func TestMain(m *testing.M) {
	code.InitPackageASTCache(".")
	m.Run()
}

func Test_FromBool(t *testing.T) {
	tests := map[string]struct {
		getResult       func() (string, error)
		expectedMessage string
		expectedError   string
	}{
		"ok": {
			getResult: func() (string, error) {
				var err error
				return FromBool(0, err == nil)
			},
			expectedMessage: "err is nil",
		},
		"ko": {
			getResult: func() (string, error) {
				return FromBool(100, true)
			},
			expectedError: "no caller information available",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch msg, err := tt.getResult(); {
			case tt.expectedError == "" && err != nil:
				t.Errorf("expected success but got error: %v", err)
			case tt.expectedError != "" && err == nil:
				t.Errorf("expected failure with following error: %s", tt.expectedError)
			case tt.expectedError != "" && err != nil && !strings.Contains(err.Error(), tt.expectedError):
				t.Errorf("expected failure with message: %s, got %s", tt.expectedError, err.Error())
			case tt.expectedError == "" && err == nil && !strings.Contains(msg, tt.expectedMessage):
				t.Errorf("expected %q to contain %q", msg, tt.expectedMessage)
			}
		})
	}
}

func Test_customizeASTExprRepr(t *testing.T) {
	anError := errors.New("bim")
	errBoom := errors.New("boom")

	for name, tt := range map[string]map[string]struct {
		getResult       func(*testing.T) (string, error)
		expectedMessage string
		expectedError   string
	}{
		"BinaryExpr": {
			// AND / OR
			"AND_true": {
				getResult: func(t *testing.T) (string, error) {
					i, j := true, true
					pkg, expr := getTestingExpr[bool](t, i && j)
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "var i is true, and var j is true",
			},
			"AND_false": {
				getResult: func(t *testing.T) (string, error) {
					i, j := true, false
					pkg, expr := getTestingExpr[bool](t, i && j)
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "var i is false, or var j is false",
			},
			"OR_true": {
				getResult: func(t *testing.T) (string, error) {
					i, j := true, false
					pkg, expr := getTestingExpr[bool](t, i || j)
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "var i is true, or var j is true",
			},
			"OR_false": {
				getResult: func(t *testing.T) (string, error) {
					i, j := false, false
					pkg, expr := getTestingExpr[bool](t, i || j)
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "var i is false, and var j is false",
			},

			// EQ / NEQ
			"EQ-generic_true": {
				getResult: func(t *testing.T) (string, error) {
					a, b := "str", "str"
					pkg, expr := getTestingExpr[bool](t, a == b)
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "a is equal to b",
			},
			"EQ-generic_false": {
				getResult: func(t *testing.T) (string, error) {
					a, b := "str", "not str"
					pkg, expr := getTestingExpr[bool](t, a == b)
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "a is not equal to b",
			},
			"NEQ-generic_true": {
				getResult: func(t *testing.T) (string, error) {
					a, b := "str", "not str"
					pkg, expr := getTestingExpr[bool](t, a != b)
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "a is not equal to b",
			},
			"NEQ-generic_false": {
				getResult: func(t *testing.T) (string, error) {
					a, b := "str", "str"
					pkg, expr := getTestingExpr[bool](t, a != b)
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "a is equal to b",
			},
			"EQ-bool-func-compared-to-bool_false": {
				getResult: func(t *testing.T) (string, error) {
					f := func() bool { return true }
					pkg, expr := getTestingExpr[bool](t, f() == false)
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "f() returned false",
			},
			"NEQ-bool-func-compared-to-bool_false": {
				getResult: func(t *testing.T) (string, error) {
					f := func() bool { return true }
					pkg, expr := getTestingExpr[bool](t, f() != true)
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "f() returned false",
			},
			"EQ-bool-func-compared-to-bool_true": {
				getResult: func(t *testing.T) (string, error) {
					f := func() bool { return true }
					pkg, expr := getTestingExpr[bool](t, f() == true)
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "f() returned true",
			},
			"NEQ-bool-func-compared-to-bool_true": {
				getResult: func(t *testing.T) (string, error) {
					f := func() bool { return true }
					pkg, expr := getTestingExpr[bool](t, f() != false)
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "f() returned true",
			},
			"EQ-error-func-compared-to-nil_false": {
				getResult: func(t *testing.T) (string, error) {
					f := func() error { return anError }
					pkg, expr := getTestingExpr[bool](t, f() == nil)
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "f() returned an error",
			},
			"NEQ-error-func-compared-to-nil_false": {
				getResult: func(t *testing.T) (string, error) {
					f := func() error { return nil }
					pkg, expr := getTestingExpr[bool](t, f() != nil)
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "f() returned no error",
			},
			"EQ-error-func-compared-to-nil_true": {
				getResult: func(t *testing.T) (string, error) {
					f := func() error { return nil }
					pkg, expr := getTestingExpr[bool](t, f() == nil)
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "f() returned no error",
			},
			"NEQ-error-func-compared-to-nil_true": {
				getResult: func(t *testing.T) (string, error) {
					f := func() error { return anError }
					pkg, expr := getTestingExpr[bool](t, f() != nil)
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "f() returned an error",
			},
			"EQ-nullable-func-compared-to-nil_false": {
				getResult: func(t *testing.T) (string, error) {
					f := func() []string { return []string{"a", "b"} }
					pkg, expr := getTestingExpr[bool](t, f() == nil)
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "f() returned not nil",
			},
			"NEQ-nullable-func-compared-to-nil_false": {
				getResult: func(t *testing.T) (string, error) {
					f := func() []string { return nil }
					pkg, expr := getTestingExpr[bool](t, f() != nil)
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "f() returned nil",
			},
			"EQ-nullable-func-compared-to-nil_true": {
				getResult: func(t *testing.T) (string, error) {
					f := func() []string { return nil }
					pkg, expr := getTestingExpr[bool](t, f() == nil)
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "f() returned nil",
			},
			"NEQ-nullable-func-compared-to-nil_true": {
				getResult: func(t *testing.T) (string, error) {
					f := func() []string { return []string{"a", "b"} }
					pkg, expr := getTestingExpr[bool](t, f() != nil)
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "f() returned not nil",
			},
			"EQ-not-func-compared-to-nil_true": {
				getResult: func(t *testing.T) (string, error) {
					var a any
					pkg, expr := getTestingExpr[bool](t, a == nil)
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "a is nil",
			},
			"NEQ-not-func-compared-to-nil_true": {
				getResult: func(t *testing.T) (string, error) {
					a := new(bool)
					pkg, expr := getTestingExpr[bool](t, a != nil)
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "a is not nil",
			},
			"EQ-not-func-compared-to-nil_false": {
				getResult: func(t *testing.T) (string, error) {
					a := new(bool)
					pkg, expr := getTestingExpr[bool](t, a == nil)
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "a is not nil",
			},
			"NEQ-not-func-compared-to-nil_false": {
				getResult: func(t *testing.T) (string, error) {
					var a any
					pkg, expr := getTestingExpr[bool](t, a != nil)
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "a is nil",
			},
			"EQ-not-func-compared-to-bool_true-constant": {
				getResult: func(t *testing.T) (string, error) {
					var b bool
					pkg, expr := getTestingExpr[bool](t, b == false)
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "b is false",
			},
			"NEQ-not-func-compared-to-bool_true-constant": {
				getResult: func(t *testing.T) (string, error) {
					var b bool
					pkg, expr := getTestingExpr[bool](t, b != true)
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "b is not true",
			},
			"EQ-not-func-compared-to-bool_false-constant": {
				getResult: func(t *testing.T) (string, error) {
					var b bool
					pkg, expr := getTestingExpr[bool](t, b == true)
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "b is not true",
			},
			"NEQ-not-func-compared-to-bool_false-constant": {
				getResult: func(t *testing.T) (string, error) {
					var b bool
					pkg, expr := getTestingExpr[bool](t, b != false)
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "b is false",
			},
			"EQ-not-func-compared-to-bool-not-constant": {
				getResult: func(t *testing.T) (string, error) {
					var (
						b1 bool
						b2 bool
					)
					pkg, expr := getTestingExpr[bool](t, b1 == b2)
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "b1 is equal to b2",
			},
			"NEQ-not-func-compared-to-bool-not-constant": {
				getResult: func(t *testing.T) (string, error) {
					var b1 bool
					b2 := true
					pkg, expr := getTestingExpr[bool](t, b1 != b2)
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "b1 is not equal to b2",
			},
			"GTR_true": {
				getResult: func(t *testing.T) (string, error) {
					n1, n2 := 42, 3
					pkg, expr := getTestingExpr[bool](t, n1 > n2)
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "n1 is greater than n2",
			},
			"GTE_true": {
				getResult: func(t *testing.T) (string, error) {
					n1, n2 := 42, 3
					pkg, expr := getTestingExpr[bool](t, n1 >= n2)
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "n1 is greater than or equal to n2",
			},
			"GTR_false": {
				getResult: func(t *testing.T) (string, error) {
					n1, n2 := 3, 42
					pkg, expr := getTestingExpr[bool](t, n1 > n2)
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "n1 is less than or equal to n2",
			},
			"GTE_false": {
				getResult: func(t *testing.T) (string, error) {
					n1, n2 := 3, 42
					pkg, expr := getTestingExpr[bool](t, n1 >= n2)
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "n1 is less than n2",
			},
			"LSS_true": {
				getResult: func(t *testing.T) (string, error) {
					n1, n2 := 3, 42
					pkg, expr := getTestingExpr[bool](t, n1 < n2)
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "n1 is less than n2",
			},
			"LEQ_true": {
				getResult: func(t *testing.T) (string, error) {
					n1, n2 := 3, 42
					pkg, expr := getTestingExpr[bool](t, n1 <= n2)
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "n1 is less than or equal to n2",
			},
			"LSS_false": {
				getResult: func(t *testing.T) (string, error) {
					n1, n2 := 42, 3
					pkg, expr := getTestingExpr[bool](t, n1 < n2)
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "n1 is greater than or equal to n2",
			},
			"LEQ_false": {
				getResult: func(t *testing.T) (string, error) {
					n1, n2 := 42, 3
					pkg, expr := getTestingExpr[bool](t, n1 <= n2)
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "n1 is greater than n2",
			},
		},
		"CallExpr": {
			"FuncLit_true": {
				getResult: func(t *testing.T) (string, error) {
					pkg, expr := getTestingExpr[bool](t, func() bool { return true }())
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "func() bool { return true }() returned true",
			},
			"FuncLit_false": {
				getResult: func(t *testing.T) (string, error) {
					pkg, expr := getTestingExpr[bool](t, func() bool { return false }())
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "func() bool { return false }() returned false",
			},
			"Ident_true": {
				getResult: func(t *testing.T) (string, error) {
					i := func() bool { return true }
					pkg, expr := getTestingExpr[bool](t, i())
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "function i() returned true",
			},
			"Ident_false": {
				getResult: func(t *testing.T) (string, error) {
					pkg, expr := getTestingExpr[bool](t, errors.Is(anError, errBoom))
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "errBoom is not in the error tree of anError",
			},
			"SelectorExpr_true": {
				getResult: func(t *testing.T) (string, error) {
					pkg, expr := getTestingExpr[bool](t, os.IsExist(nil))
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "function os.IsExist(nil) returned true",
			},
			"SelectorExpr_false": {
				getResult: func(t *testing.T) (string, error) {
					pkg, expr := getTestingExpr[bool](t, os.IsExist(os.ErrExist))
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "function os.IsExist(os.ErrExist) returned false",
			},
			"strings.Contains_true": {
				getResult: func(t *testing.T) (string, error) {
					pkg, expr := getTestingExpr[bool](t, strings.Contains("foo", "bar"))
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: `"foo" contains "bar"`,
			},
			"strings.Contains_false": {
				getResult: func(t *testing.T) (string, error) {
					pkg, expr := getTestingExpr[bool](t, strings.Contains("foo", "bar"))
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: `"foo" does not contain "bar"`,
			},
			"slices.Contains_true": {
				getResult: func(t *testing.T) (string, error) {
					pkg, expr := getTestingExpr[bool](t, slices.Contains([]string{"foo", "bar"}, "bar"))
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: `[]string{"foo", "bar"} contains "bar"`,
			},
			"slices.Contains_false": {
				getResult: func(t *testing.T) (string, error) {
					pkg, expr := getTestingExpr[bool](t, slices.Contains([]string{"foo"}, "bar"))
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: `[]string{"foo"} does not contain "bar"`,
			},
			"bytes.Equal_true": {
				getResult: func(t *testing.T) (string, error) {
					pkg, expr := getTestingExpr[bool](t, bytes.Equal([]byte("str"), []byte("str")))
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: `[]byte("str") is equal to []byte("str")`,
			},
			"bytes.Equal_false": {
				getResult: func(t *testing.T) (string, error) {
					pkg, expr := getTestingExpr[bool](t, bytes.Equal([]byte("str"), []byte("abc")))
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: `[]byte("str") is not equal to []byte("abc")`,
			},
			"maps.Equal_true": {
				getResult: func(t *testing.T) (string, error) {
					pkg, expr := getTestingExpr[bool](t, maps.Equal(map[string]int{"a": 1}, map[string]int{"a": 1}))
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: `map[string]int{"a": 1} is equal to map[string]int{"a": 1}`,
			},
			"maps.Equal_false": {
				getResult: func(t *testing.T) (string, error) {
					pkg, expr := getTestingExpr[bool](t, maps.Equal(map[string]int{"a": 1}, map[string]int{"a": 2}))
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: `map[string]int{"a": 1} is not equal to map[string]int{"a": 2}`,
			},
			"reflect.DeepEqual_true": {
				getResult: func(t *testing.T) (string, error) {
					pkg, expr := getTestingExpr[bool](t, reflect.DeepEqual(map[string]int{"a": 1}, map[string]int{"a": 1}))
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: `map[string]int{"a": 1} is equal to map[string]int{"a": 1}`,
			},
			"reflect.DeepEqual_false": {
				getResult: func(t *testing.T) (string, error) {
					pkg, expr := getTestingExpr[bool](t, reflect.DeepEqual(map[string]int{"a": 1}, map[string]int{"a": 2}))
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: `map[string]int{"a": 1} is not equal to map[string]int{"a": 2}`,
			},
			"errors.Is_true": {
				getResult: func(t *testing.T) (string, error) {
					pkg, expr := getTestingExpr[bool](t, errors.Is(anError, errBoom))
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "anError's error tree contains errBoom",
			},
			"errors.Is_false": {
				getResult: func(t *testing.T) (string, error) {
					pkg, expr := getTestingExpr[bool](t, errors.Is(anError, errBoom))
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "errBoom is not in the error tree of anError",
			},
			"errors.As_true": {
				getResult: func(t *testing.T) (string, error) {
					type boomErr interface{ Boom() }
					var bErr boomErr
					pkg, expr := getTestingExpr[bool](t, errors.As(anError, &bErr))
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "anError can be defined as *github.com/krostar/test/internal/message.boomErr",
			},
			"errors.As_false": {
				getResult: func(t *testing.T) (string, error) {
					type boomErr interface{ Boom() }
					var bErr boomErr
					pkg, expr := getTestingExpr[bool](t, errors.As(anError, &bErr))
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "anError cannot be defined as *github.com/krostar/test/internal/message.boomErr",
			},
		},
		"Ident": {
			"literal_true": {
				getResult: func(t *testing.T) (string, error) {
					pkg, expr := getTestingExpr[bool](t, true)
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "literal true",
			},
			"literal_false": {
				getResult: func(t *testing.T) (string, error) {
					pkg, expr := getTestingExpr[bool](t, false)
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "literal false",
			},
			"var_true": {
				getResult: func(t *testing.T) (string, error) {
					var i bool
					pkg, expr := getTestingExpr[bool](t, i)
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "var i is true",
			},
			"var_false": {
				getResult: func(t *testing.T) (string, error) {
					var i bool
					pkg, expr := getTestingExpr[bool](t, i)
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "var i is false",
			},
			"const_true": {
				getResult: func(t *testing.T) (string, error) {
					const i = false
					pkg, expr := getTestingExpr[bool](t, i)
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "const i is true",
			},
			"const_false": {
				getResult: func(t *testing.T) (string, error) {
					const i = false
					pkg, expr := getTestingExpr[bool](t, i)
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "const i is false",
			},
		},
		"ParentExpr": {
			"bool": {
				getResult: func(t *testing.T) (string, error) {
					i := true
					pkg, expr := getTestingExpr[bool](t, (i))
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "i is true",
			},
			"gt": {
				getResult: func(t *testing.T) (string, error) {
					i := 21
					pkg, expr := getTestingExpr[bool](t, (i > 42))
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "i is less than or equal to 42",
			},
		},
		"SelectorExpr": {
			"foo": {
				getResult: func(t *testing.T) (string, error) {
					foo := struct {
						value bool
					}{
						value: true,
					}
					pkg, expr := getTestingExpr[bool](t, foo.value)
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "foo.value is true",
			},
		},
		"UnaryExpr": {
			"NOT-CallExpr": {
				getResult: func(t *testing.T) (string, error) {
					pkg, expr := getTestingExpr[bool](t, !errors.Is(anError, errBoom))
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "anError's error tree contains errBoom",
			},
			"NOT-Ident": {
				getResult: func(t *testing.T) (string, error) {
					i := true
					pkg, expr := getTestingExpr[bool](t, !i)
					return customizeASTExprRepr(pkg, false, expr)
				},
				expectedMessage: "var i is true",
			},
			"NOT-ParentExpr": {
				getResult: func(t *testing.T) (string, error) {
					pkg, expr := getTestingExpr[bool](t, !(21 > 42))
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "21 is less than or equal to 42",
			},
			"NOT-UnaryExpr": {
				getResult: func(t *testing.T) (string, error) {
					i := true
					pkg, expr := getTestingExpr[bool](t, !!i)
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "var i is true",
			},
			"ARROW-bool": {
				getResult: func(t *testing.T) (string, error) {
					b := make(chan bool, 1)
					b <- true
					pkg, expr := getTestingExpr[bool](t, <-b)
					return customizeASTExprRepr(pkg, true, expr)
				},
				expectedMessage: "var b is true",
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			for subName, test := range tt {
				t.Run(subName, func(t *testing.T) {
					switch msg, err := test.getResult(t); {
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

func Test_isExprNil(t *testing.T) {
	t.Run("nil expr", func(t *testing.T) {
		pkg, _ := rawGetTestingExpr(t, 0, "rawGetTestingExpr", 42)
		if isExprNil(pkg, nil) {
			t.Error("expected false but true")
		}
	})

	t.Run("non-nil expr", func(t *testing.T) {
		for i, tt := range []struct {
			got    bool
			expect bool
		}{
			{
				got:    isExprNil(getTestingExpr[*int](t, nil)),
				expect: true,
			},
			{
				got:    isExprNil(getTestingExpr[map[string]string](t, nil)),
				expect: true,
			},
			{
				got:    isExprNil(getTestingExpr[chan struct{}](t, nil)),
				expect: true,
			},
			{
				got:    isExprNil(getTestingExpr[error](t, nil)),
				expect: true,
			},
			{
				got:    isExprNil(getTestingExpr(t, 0)),
				expect: false,
			},
			{
				got:    isExprNil(getTestingExpr(t, "")),
				expect: false,
			},
		} {
			if tt.got != tt.expect {
				t.Errorf("[%d] expected %v, got %v", i, tt.expect, tt.got)
			}
		}
	})
}

func Test_isExprBool(t *testing.T) {
	t.Run("nil expr", func(t *testing.T) {
		pkg, _ := rawGetTestingExpr(t, 0, "rawGetTestingExpr", 42)
		if isExprBool(pkg, nil) {
			t.Error("expected false but true")
		}
	})

	t.Run("non-nil expr", func(t *testing.T) {
		const (
			ctrue  = true
			cfalse = false
		)
		var (
			vtrue  = true
			vfalse = false
		)

		for i, tt := range []struct {
			got    bool
			expect bool
		}{
			{
				got:    isExprBool(getTestingExpr(t, true)),
				expect: true,
			},
			{
				got:    isExprBool(getTestingExpr(t, false)),
				expect: true,
			},
			{
				got:    isExprBool(getTestingExpr(t, ctrue)),
				expect: true,
			},
			{
				got:    isExprBool(getTestingExpr(t, cfalse)),
				expect: true,
			},
			{
				got:    isExprBool(getTestingExpr(t, vtrue)),
				expect: true,
			},
			{
				got:    isExprBool(getTestingExpr(t, vfalse)),
				expect: true,
			},
			{
				got:    isExprBool(getTestingExpr(t, 42)),
				expect: false,
			},
			{
				got:    isExprBool(getTestingExpr(t, "42")),
				expect: false,
			},
			{
				got:    isExprBool(getTestingExpr(t, []int{42})),
				expect: false,
			},
		} {
			if tt.got != tt.expect {
				t.Errorf("[%d] expected %v, got %v", i, tt.expect, tt.got)
			}
		}
	})
}

func Test_isExprFunc(t *testing.T) {
	t.Run("nil expr", func(t *testing.T) {
		pkg, _ := rawGetTestingExpr(t, 0, "rawGetTestingExpr", 42)
		if isExprFunc(pkg, nil) {
			t.Error("expected false but true")
		}
	})

	t.Run("non-nil expr", func(t *testing.T) {
		f := func() error { return nil }

		for i, tt := range []struct {
			got    bool
			expect bool
		}{
			{
				got:    isExprFunc(getTestingExpr(t, f())),
				expect: true,
			},
			{
				got:    isExprFunc(getTestingExpr(t, func() error { return nil }())),
				expect: true,
			},
			{
				got:    isExprFunc(getTestingExpr(t, t.Context())),
				expect: true,
			},
			{
				got:    isExprFunc(getTestingExpr(t, f)),
				expect: false,
			},
			{
				got:    isExprFunc(getTestingExpr(t, 42)),
				expect: false,
			},
		} {
			if tt.got != tt.expect {
				t.Errorf("[%d] expected %v, got %v", i, tt.expect, tt.got)
			}
		}
	})
}

func Test_isExprFuncReturningOnlyError(t *testing.T) {
	t.Run("nil expr", func(t *testing.T) {
		pkg, _ := rawGetTestingExpr(t, 0, "rawGetTestingExpr", 42)
		if isExprFuncReturningOnlyError(pkg, nil) {
			t.Error("expected false but true")
		}
	})

	t.Run("non-nil expr", func(t *testing.T) {
		a := func() error { return nil }

		for i, tt := range []struct {
			got    bool
			expect bool
		}{
			{
				got:    isExprFuncReturningOnlyError(getTestingExpr(t, a())),
				expect: true,
			},
			{
				got:    isExprFuncReturningOnlyError(getTestingExpr(t, errors.New("boom"))),
				expect: true,
			},
			{
				got:    isExprFuncReturningOnlyError(getTestingExpr(t, func() error { return nil }())),
				expect: true,
			},
			{
				got:    isExprFuncReturningOnlyError(getTestingExpr(t, t.Context())),
				expect: false,
			},
			{
				got:    isExprFuncReturningOnlyError(getTestingExpr(t, 42)),
				expect: false,
			},
		} {
			if tt.got != tt.expect {
				t.Errorf("[%d] expected %v, got %v", i, tt.expect, tt.got)
			}
		}
	})
}

func Test_getIdentSelector(t *testing.T) {
	t.Run("bad expr", func(t *testing.T) {
		pkg, _ := rawGetTestingExpr(t, 0, "rawGetTestingExpr", 42)
		_, _, err := getIdentSelector(pkg, new(ast.Ident))
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "ident object is nil") {
			t.Errorf("expected different error reason: %v", err)
		}
	})

	t.Run("nil expr", func(t *testing.T) {
		pkg, _ := rawGetTestingExpr(t, 0, "rawGetTestingExpr", 42)
		p, i, err := getIdentSelector(pkg, nil)
		switch {
		case err != nil:
			t.Errorf("unexpected error: %v", err)
		case p != "" || i != "":
			t.Errorf("expected selector to be empty, got %s.%s", p, i)
		}
	})

	t.Run("errors.New", func(t *testing.T) {
		p, i, err := getIdentSelector(getTestingIdent(t, errors.New))
		switch {
		case err != nil:
			t.Errorf("unexpected error: %v", err)
		case p != "errors" || i != "New":
			t.Errorf("expected selector to be %s, got %s.%s", "errors/New", p, i)
		}
	})

	t.Run("context.Canceled", func(t *testing.T) {
		p, i, err := getIdentSelector(getTestingIdent(t, context.Canceled))
		switch {
		case err != nil:
			t.Errorf("unexpected error: %v", err)
		case p != "context" || i != "Canceled":
			t.Errorf("expected selector to be %s, got %s.%s", "context.Canceled", p, i)
		}
	})

	t.Run("local var", func(t *testing.T) {
		var b bool
		p, i, err := getIdentSelector(getTestingIdent(t, b))
		switch {
		case err != nil:
			t.Errorf("unexpected error: %v", err)
		case p != "github.com/krostar/test/internal/message" || i != "b":
			t.Errorf("expected selector to be %s, got %s.%s", "github.com/krostar/test/internal/message", p, i)
		}
	})
}

func Test_getExprBoolValue(t *testing.T) {
	if !*getExprBoolValue(getTestingExpr(t, true)) {
		t.Error("expected true")
	}

	if *getExprBoolValue(getTestingExpr(t, false)) {
		t.Error("expected false")
	}

	t.Run("nil value", func(t *testing.T) {
		if getExprBoolValue(nil, nil) != nil {
			t.Error("expected false")
		}
	})

	t.Run("non-bool value or not constant value", func(t *testing.T) {
		if getExprBoolValue(getTestingExpr(t, 42)) != nil {
			t.Error("expected nil")
		}

		var b bool
		if getExprBoolValue(getTestingExpr(t, b)) != nil {
			t.Error("expected nil")
		}
	})
}

func rawGetTestingExpr[T any](t *testing.T, stack int, funcName string, _ T) (*packages.Package, ast.Expr) {
	t.Helper()

	_, callerFile, callerLine, ok := runtime.Caller(stack + 1)
	if !ok {
		t.Fatalf("no caller information available")
	}

	pkgPathToPkg, err := code.GetPackageAST(filepath.Clean(filepath.Dir(callerFile)))
	if err != nil {
		t.Fatalf("unable to get package AST: %v", err)
	}

	_, file, pkg, err := code.GetCallerCallExpr(pkgPathToPkg, callerFile, callerLine)
	if err != nil {
		t.Fatalf("unable to get call expr from caller: %v", err)
	}

	var (
		prev *ast.CallExpr
		expr *ast.CallExpr
	)

	ast.Inspect(file, func(node ast.Node) bool {
		if node == nil || expr != nil {
			return false
		}

		if pkg.Fset.Position(node.Pos()).Line != callerLine {
			return true
		}

		if ident, ok := node.(*ast.Ident); prev != nil && ok {
			if ident.Name == funcName {
				expr = prev
				return false
			}
			prev = nil
			return true
		}

		if call, ok := node.(*ast.CallExpr); ok {
			prev = call
			return true
		}

		return true
	})

	if expr == nil {
		t.Fatalf("no call expression found")
	}

	return pkg, expr.Args[1]
}

func getTestingExpr[T any](t *testing.T, v T) (*packages.Package, ast.Expr) {
	t.Helper()

	return rawGetTestingExpr(t, 1, "getTestingExpr", v)
}

func getTestingIdent[T any](t *testing.T, v T) (*packages.Package, *ast.Ident) {
	t.Helper()

	pkg, expr := rawGetTestingExpr(t, 1, "getTestingIdent", v)
	if expr == nil {
		t.Fatalf("no call expression found")
		return nil, nil
	}

	switch expr := expr.(type) {
	case *ast.SelectorExpr:
		return pkg, expr.Sel
	case *ast.Ident:
		return pkg, expr
	default:
		t.Fatalf("unhandled type: %T", expr)
		return nil, nil
	}
}
