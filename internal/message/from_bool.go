// Package message provides functions for generating formatted messages for assertions.
package message

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"path/filepath"
	"runtime"

	"golang.org/x/tools/go/packages"

	"github.com/krostar/test/internal/code"
)

// FromBool generates a customized message string based on a boolean result and caller information.
//
// `callerStackIndex` specifies the depth in the call stack to retrieve the caller information.
// This is used to identify the source code location of the assertion.
// `result` is the boolean value for which to generate the message.
//
// It returns a formatted message string and an error if one occurred during the process.
// The message string will be tailored based on the expression used in the assertion.
func FromBool(callerStackIndex int, result bool) (string, error) {
	_, callerFile, callerLine, ok := runtime.Caller(callerStackIndex + 1)
	if !ok {
		return "", errors.New("no caller information available")
	}

	pkgPathToPkg, err := code.GetPackageAST(filepath.Clean(filepath.Dir(callerFile)))
	if err != nil {
		return "", fmt.Errorf("unable to get package AST: %v", err)
	}

	expr, _, pkg, err := code.GetCallerCallExpr(pkgPathToPkg, callerFile, callerLine)
	if err != nil {
		return "", fmt.Errorf("unable to get call expr from caller: %v", err)
	}

	var arg ast.Expr
	switch l := len(expr.Args); {
	case l == 1: // interpret as custom checker like Assert(checker(t, ...))
		arg = expr.Args[0]
	case l >= 2: // interpret as regular call like Assert(t, bool, msg...)
		arg = expr.Args[1]
	default:
		return "", fmt.Errorf("unexpected call expr arguments number %d", l)
	}

	msg, err := customizeASTExprRepr(pkg, result, arg)
	if err != nil {
		return genericASTExprToString(pkg, expr), fmt.Errorf("unable to get arg repr: %v", err)
	}

	return msg, nil
}

// customizeASTExprRepr generates a representation of an AST expression,
// customizing it based on the type and context of the expression.
//
// `pkg` provides type information for the expression.
// `result` is the result of the assertion, used to tailor the message.
// `expr` is the AST expression to represent.
//
// It returns the formatted string representation of the expression and
// an error if any occurred during the processing.
//
//nolint:exhaustive // type switch of ast.Expr lead to a lot of cases but many of them are un-needed here
func customizeASTExprRepr(pkg *packages.Package, result bool, expr ast.Expr) (string, error) {
	typ := pkg.TypesInfo.TypeOf(expr)

	switch expr := expr.(type) {
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s is %t", genericASTExprToString(pkg, expr), result), nil

	case *ast.Ident:
		obj := pkg.TypesInfo.ObjectOf(expr)
		switch obj := obj.(type) {
		case *types.Var:
			return fmt.Sprintf("var %s is %t", obj.Name(), result), nil
		case *types.Const:
			switch typ := typ.(type) {
			case *types.Basic:
				switch typ.Kind() {
				case types.UntypedBool, types.Bool:
					return "literal" + expr.String(), nil
				default:
					return "", fmt.Errorf("unhandled basic kind %v", typ.Kind())
				}
			default:
				return "", fmt.Errorf("unhandled ident type %v", typ)
			}
		default:
			return "", fmt.Errorf("unexpected ident obj of type %T", obj)
		}

	case *ast.BinaryExpr:
		x, y := genericASTExprToString(pkg, expr.X), genericASTExprToString(pkg, expr.Y)

		if expr.Op == token.LAND || expr.Op == token.LOR {
			var err error
			if x, err = customizeASTExprRepr(pkg, result, expr.X); err != nil {
				return "", fmt.Errorf("unable to get LAND/LOR %T.X custom repr: %v", expr, err)
			}
			if y, err = customizeASTExprRepr(pkg, result, expr.Y); err != nil {
				return "", fmt.Errorf("unable to get LAND/LOR %T.Y custom repr: %v", expr, err)
			}
			var sep string
			switch expr.Op {
			case token.LAND:
				sep = ", or "
			case token.LOR:
				sep = ", and "
			}
			return x + sep + y, nil
		}

		switch {
		case (expr.Op == token.EQL && result) || (expr.Op == token.NEQ && !result):
			switch {
			case isExprFunc(pkg, expr.X) && (isExprBool(pkg, expr.Y) || isExprNil(pkg, expr.Y)):
				return fmt.Sprintf("%s returned %s", x, y), nil
			case !isExprFunc(pkg, expr.X) && (isExprBool(pkg, expr.Y) || isExprNil(pkg, expr.Y)):
				return fmt.Sprintf("%s is %s", x, y), nil
			default:
				return fmt.Sprintf("%s is equal to %s", x, y), nil
			}
		case (expr.Op == token.EQL && !result) || (expr.Op == token.NEQ && result):
			switch {
			case isExprFunc(pkg, expr.X) && (isExprBool(pkg, expr.Y)):
				return x + " returned " + y, nil
			case isExprFunc(pkg, expr.X) && isExprNil(pkg, expr.Y):
				return x + " returned not nil", nil
			case !isExprFunc(pkg, expr.X) && (isExprBool(pkg, expr.Y) || isExprNil(pkg, expr.Y)):
				return x + " is not " + y, nil
			default:
				return x + " is not equal to " + y, nil
			}
		case (expr.Op == token.GTR && result) || (expr.Op == token.LEQ && !result):
			return x + " is greater than " + y, nil
		case (expr.Op == token.GEQ && !result) || (expr.Op == token.LSS && result):
			return x + " is less than " + y, nil
		case (expr.Op == token.GEQ && result) || (expr.Op == token.LSS && !result):
			return x + " is greater than or equal to " + y, nil
		case (expr.Op == token.GTR && !result) || (expr.Op == token.LEQ && result):
			return x + " is less than or equal to " + y, nil
		default:
			return "", fmt.Errorf("unhandled binary operator %v", expr.Op)
		}

	case *ast.UnaryExpr:
		switch op := expr.Op; op {
		case token.NOT:
			switch expr.X.(type) {
			case *ast.CallExpr, *ast.BinaryExpr:
				return customizeASTExprRepr(pkg, !result, expr.X)
			default:
				return fmt.Sprintf("%s is %t", genericASTExprToString(pkg, expr), result), nil
			}
		default:
			return "", fmt.Errorf("unhandled unary operator %T", expr.Op)
		}

	case *ast.CallExpr:
		var p, t string
		switch fun := expr.Fun.(type) {
		case *ast.SelectorExpr:
			var err error
			if p, t, err = getIdentSelector(pkg, fun.Sel); err != nil {
				return "", fmt.Errorf("unable to get func.Sel selector from %T: %v", err, expr)
			}
		case *ast.Ident:
			var err error
			if p, t, err = getIdentSelector(pkg, fun); err != nil {
				return "", fmt.Errorf("unable to get func ident selector from %T: %v", err, expr)
			}
		default:
			return "", fmt.Errorf("unhandled function type %T", fun)
		}

		if sig, ok := typ.(*types.Signature); ok {
			sigTy := sig.Underlying()
			_ = sigTy
			if result {
				return fmt.Sprintf("check %s passed", genericASTExprToString(pkg, expr)), nil
			}
			return fmt.Sprintf("check %s failed", genericASTExprToString(pkg, expr)), nil
		}

		switch {
		case p == "strings" && t == "Contains":
			if result {
				return fmt.Sprintf("%s contains %s", genericASTExprToString(pkg, expr.Args[0]), genericASTExprToString(pkg, expr.Args[1])), nil
			}
			return fmt.Sprintf("%s does not contain %s", genericASTExprToString(pkg, expr.Args[0]), genericASTExprToString(pkg, expr.Args[1])), nil
		case p == "errors" && t == "Is":
			if result {
				return fmt.Sprintf("%s's error tree contains %s", genericASTExprToString(pkg, expr.Args[0]), genericASTExprToString(pkg, expr.Args[1])), nil
			}
			return fmt.Sprintf("%s is not in the error tree of %s", genericASTExprToString(pkg, expr.Args[1]), genericASTExprToString(pkg, expr.Args[0])), nil
		case p == "errors" && t == "As":
			if result {
				return fmt.Sprintf("%s can be defined as %s", genericASTExprToString(pkg, expr.Args[0]), pkg.TypesInfo.TypeOf(expr.Args[1])), nil
			}
			return fmt.Sprintf("%s cannot be defined as %s", genericASTExprToString(pkg, expr.Args[0]), pkg.TypesInfo.TypeOf(expr.Args[1])), nil
		default:
			return fmt.Sprintf("function %s returned %t", genericASTExprToString(pkg, expr), result), nil
		}
	default:
		return "", fmt.Errorf("unhandled expr type %T", expr)
	}
}

func genericASTExprToString(pkg *packages.Package, expr ast.Expr) string {
	var buf bytes.Buffer
	_ = format.Node(&buf, pkg.Fset, expr) //nolint:errcheck // we parse AST from a package compile for a test, this should'nt fail
	return buf.String()
}

func isExprNil(pkg *packages.Package, expr ast.Expr) bool {
	if expr == nil {
		return false
	}

	ident, ok := expr.(*ast.Ident)
	if !ok {
		return false
	}

	_, ok = pkg.TypesInfo.ObjectOf(ident).(*types.Nil)
	return ok
}

func isExprBool(pkg *packages.Package, expr ast.Expr) bool {
	if expr == nil {
		return false
	}

	typ := pkg.TypesInfo.TypeOf(expr)
	basic, ok := typ.(*types.Basic)
	if !ok {
		return false
	}

	return basic.Kind() == types.Bool || basic.Kind() == types.UntypedBool
}

func isExprFunc(_ *packages.Package, expr ast.Expr) bool {
	if expr == nil {
		return false
	}

	_, ok := expr.(*ast.CallExpr)

	return ok
}

func getIdentSelector(pkg *packages.Package, expr *ast.Ident) (string, string, error) {
	if expr == nil {
		return "", "", nil
	}

	obj := pkg.TypesInfo.ObjectOf(expr)
	if obj == nil {
		return "", "", errors.New("ident object is nil")
	}

	return obj.Pkg().Path(), obj.Name(), nil
}
