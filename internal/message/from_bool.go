// Package message provides functions for generating formatted messages for assertions.
package message

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/constant"
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
//nolint:goconst // some hardcoded values could be consts, but for readability reasons it seems better to keep them raw
func customizeASTExprRepr(pkg *packages.Package, result bool, expr ast.Expr) (string, error) {
	typ := pkg.TypesInfo.TypeOf(expr)

	switch expr := expr.(type) {
	case *ast.BinaryExpr:
		x, y := genericASTExprToString(pkg, expr.X), genericASTExprToString(pkg, expr.Y)

		switch {
		case expr.Op == token.LAND || expr.Op == token.LOR:
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
				if result {
					sep = ", and "
				} else {
					sep = ", or "
				}
			case token.LOR:
				if result {
					sep = ", or "
				} else {
					sep = ", and "
				}
			default:
				return "", fmt.Errorf("unhandled binary operator %v", expr.Op)
			}

			return x + sep + y, nil
		case expr.Op == token.EQL || expr.Op == token.NEQ:
			resultIsEqual := (expr.Op == token.EQL && result) || (expr.Op == token.NEQ && !result)
			xIsFunc := isExprFunc(pkg, expr.X)
			xIsFuncRetuningError := isExprFuncReturningOnlyError(pkg, expr.X)
			yIsNil := isExprNil(pkg, expr.Y)
			yIsBool := isExprBool(pkg, expr.Y)

			switch {
			case xIsFunc && xIsFuncRetuningError && yIsNil && resultIsEqual:
				return x + " returned no error", nil
			case xIsFunc && xIsFuncRetuningError && yIsNil && !resultIsEqual:
				return x + " returned an error", nil
			case xIsFunc && yIsNil && resultIsEqual:
				return x + " returned nil", nil
			case xIsFunc && yIsNil && !resultIsEqual:
				return x + " returned not nil", nil
			case xIsFunc && yIsBool:
				return fmt.Sprintf("%s returned %t", x, result), nil
			case !xIsFunc && yIsNil && resultIsEqual:
				return x + " is nil", nil
			case !xIsFunc && yIsNil && !resultIsEqual:
				return x + " is not nil", nil
			case !xIsFunc && yIsBool && resultIsEqual:
				return fmt.Sprintf("%s is %t", x, mustGetExprBoolValue(pkg, expr.Y)), nil
			case !xIsFunc && yIsBool && !resultIsEqual:
				return fmt.Sprintf("%s is not %t", x, mustGetExprBoolValue(pkg, expr.Y)), nil
			default:
				if resultIsEqual {
					return x + " is equal to " + y, nil
				}
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

	case *ast.CallExpr:
		var p, t string
		switch fun := expr.Fun.(type) {
		case *ast.FuncLit:
			return fmt.Sprintf("%s returned %t", genericASTExprToString(pkg, expr), result), nil
		case *ast.Ident:
			var err error
			if p, t, err = getIdentSelector(pkg, fun); err != nil {
				return "", fmt.Errorf("unable to get func ident selector from %T: %v", err, expr)
			}
		case *ast.SelectorExpr:
			var err error
			if p, t, err = getIdentSelector(pkg, fun.Sel); err != nil {
				return "", fmt.Errorf("unable to get func.Sel selector from %T: %v", err, expr)
			}
		default:
			return "", fmt.Errorf("unhandled function type %T", fun)
		}

		switch {
		case (p == "slices" && t == "Contains") || (p == "strings" && t == "Contains"):
			if result {
				return fmt.Sprintf("%s contains %s", genericASTExprToString(pkg, expr.Args[0]), genericASTExprToString(pkg, expr.Args[1])), nil
			}
			return fmt.Sprintf("%s does not contain %s", genericASTExprToString(pkg, expr.Args[0]), genericASTExprToString(pkg, expr.Args[1])), nil
		case (p == "bytes" && t == "Equal") || (p == "maps" && t == "Equal") || (p == "reflect" && t == "DeepEqual") || (p == "slices" && t == "Equal"):
			if result {
				return fmt.Sprintf("%s is equal to %s", genericASTExprToString(pkg, expr.Args[0]), genericASTExprToString(pkg, expr.Args[1])), nil
			}
			return fmt.Sprintf("%s is not equal to %s", genericASTExprToString(pkg, expr.Args[0]), genericASTExprToString(pkg, expr.Args[1])), nil
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

	case *ast.Ident:
		obj := pkg.TypesInfo.ObjectOf(expr)
		switch obj := obj.(type) {
		case *types.Var:
			return fmt.Sprintf("var %s is %t", obj.Name(), result), nil
		case *types.Const:
			if typ, ok := typ.(*types.Basic); ok && (typ.Kind() == types.Bool || typ.Kind() == types.UntypedBool) && obj.Parent() == types.Universe {
				return "literal " + obj.Name(), nil
			}
			return fmt.Sprintf("const %s is %t", obj.Name(), result), nil
		default:
			return "", fmt.Errorf("unexpected ident obj of type %T", obj)
		}

	case *ast.ParenExpr:
		return customizeASTExprRepr(pkg, result, expr.X)

	case *ast.SelectorExpr:
		return fmt.Sprintf("%s is %t", genericASTExprToString(pkg, expr), result), nil

	case *ast.UnaryExpr:
		switch op := expr.Op; op {
		case token.NOT:
			switch expr.X.(type) {
			case *ast.CallExpr, *ast.Ident, *ast.ParenExpr, *ast.UnaryExpr:
				return customizeASTExprRepr(pkg, !result, expr.X)
			default:
				return "", fmt.Errorf("unhandled unary expr operator %T", expr.X)
			}
		case token.ARROW:
			return customizeASTExprRepr(pkg, result, expr.X)
		default:
			return "", fmt.Errorf("unhandled unary operator %T", expr.Op)
		}

	default:
		return "", fmt.Errorf("unhandled expr type %T", expr)
	}
}

func genericASTExprToString(pkg *packages.Package, expr ast.Expr) string {
	var buf bytes.Buffer
	if err := format.Node(&buf, pkg.Fset, expr); err != nil {
		return fmt.Sprintf("<error formatting expression: %v>", err)
	}
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

func isExprFuncReturningOnlyError(pkg *packages.Package, expr ast.Expr) bool {
	if expr == nil {
		return false
	}

	cae, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}

	typ := pkg.TypesInfo.TypeOf(cae.Fun)
	if typ == nil {
		return false
	}

	sig, ok := typ.(*types.Signature)
	if !ok {
		return false
	}

	t := sig.Results()
	if t.Len() != 1 {
		return false
	}

	return types.Identical(t.At(0).Type(), types.Universe.Lookup("error").Type())
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

func mustGetExprBoolValue(pkg *packages.Package, expr ast.Expr) bool {
	if expr == nil {
		panic("expr is nil")
	}

	if tv, ok := pkg.TypesInfo.Types[expr]; ok && tv.IsValue() {
		if tv.Value != nil && tv.Value.Kind() == constant.Bool {
			return constant.BoolVal(tv.Value)
		}
	}

	panic(fmt.Sprintf("type %v not found or value is not a bool", expr))
}
