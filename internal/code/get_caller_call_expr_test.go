package code

import (
	"go/ast"
	"strings"
	"testing"
)

func Test_GetCallerCallExpr(t *testing.T) {
	pkgs, err := ParsePackageAST(t.Context(), "./testdata/ok")
	if err != nil {
		t.Fatalf("failed to parse package AST: %v", err)
	}
	ok := pkgs["github.com/krostar/test/internal/code/testdata/ok"]

	t.Run("ok", func(t *testing.T) {
		expr, file, pkg, err := GetCallerCallExpr(pkgs, ok.CompiledGoFiles[0], 13)
		if err != nil {
			t.Fatalf("failed to get caller expr: %v", err)
		}
		if expr == nil || file == nil || pkg == nil {
			t.Fatalf("expected expr, file, and pkg to be non-nil")
		}

		if pkg.PkgPath != ok.PkgPath {
			t.Errorf("expected pkg to be the ok pkg")
		}

		if fileName := pkg.Fset.Position(file.Pos()).Filename; fileName != ok.CompiledGoFiles[0] {
			t.Errorf("expected file to be %s, got %s", ok.CompiledGoFiles[0], fileName)
		}

		if fun := expr.Fun.(*ast.Ident).Name; fun != "launch" {
			t.Errorf("expected function to be launch, got %s", fun)
		}
	})

	t.Run("ko", func(t *testing.T) {
		t.Run("pkg not found", func(t *testing.T) {
			expr, file, pkg, err := GetCallerCallExpr(pkgs, "./notexisting.go", 1043)
			if err == nil {
				t.Errorf("expected failure")
			}

			if expr != nil || file != nil || pkg != nil {
				t.Errorf("expected expr, file, and pkg to be nil")
			}

			if !strings.Contains(err.Error(), "unable to find ast file and package") {
				t.Errorf("unexpected error message, got %s", err.Error())
			}
		})

		t.Run("expr not found", func(t *testing.T) {
			expr, file, pkg, err := GetCallerCallExpr(pkgs, ok.CompiledGoFiles[0], 5)
			if err == nil {
				t.Errorf("expected failure")
			}

			if expr != nil || file != nil || pkg != nil {
				t.Errorf("expected expr, file, and pkg to be nil")
			}

			if !strings.Contains(err.Error(), "unable to get call expression") {
				t.Errorf("unexpected error message, got %s", err.Error())
			}
		})
	})
}
