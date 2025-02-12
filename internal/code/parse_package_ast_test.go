package code

import (
	"context"
	"strings"
	"testing"
)

func Test_ParsePackageAST(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		pkgs, err := ParsePackageAST(t.Context(), "./testdata/ok")
		if err != nil {
			t.Fatalf("unable to parse package: %v", err)
		}

		if len(pkgs) == 0 {
			t.Fatalf("expected at least one package")
		}

		pkg, exists := pkgs["github.com/krostar/test/internal/code/testdata/ok"]
		if pkg == nil || !exists {
			t.Fatalf("package testdata not found")
		}

		if pkg.Name != "ok" {
			t.Fatalf("package name mismatch: expected %q, found %q", "ok", pkg.Name)
		}

		if pkg.TypesInfo == nil {
			t.Fatalf("TypesInfo should not be nil")
		}

		if pkg.TypesInfo.Defs == nil {
			t.Fatalf("TypesInfo.Defs should not be nil")
		}

		if pkg.TypesInfo.Uses == nil {
			t.Fatalf("TypesInfo.Uses should not be nil")
		}
	})

	t.Run("pkg loaded with errors", func(t *testing.T) {
		pkgs, err := ParsePackageAST(t.Context(), "./testdata/404")
		if err == nil || pkgs != nil {
			t.Fatalf("pkgs should be nil && err should be not nil: %v", err)
		}
		if !strings.Contains(err.Error(), "loaded packages contained errors") {
			t.Errorf("unexpected error message %q", err.Error())
		}
	})

	t.Run("pkg not loaded", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		pkgs, err := ParsePackageAST(ctx, "./testdata/404")
		if err == nil || pkgs != nil {
			t.Fatalf("pkgs should be nil && err should be not nil: %v", err)
		}
		if !strings.Contains(err.Error(), "unable to load packages") {
			t.Errorf("unexpected error message %s", err.Error())
		}
	})
}
