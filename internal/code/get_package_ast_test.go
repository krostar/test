package code

import (
	"strings"
	"testing"
)

func Test_InitPackageASTCache(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		pkgDir := "./testdata/ok"

		_astPkgPathToPkg = nil
		InitPackageASTCache(pkgDir)

		if _astPkgPathToPkg == nil || _astPkgPathToPkg[pkgDir] == nil {
			t.Errorf("package should be in cache")
		}
	})

	t.Run("ko", func(t *testing.T) {
		defer func() {
			reason := recover()
			if reason == nil {
				t.Errorf("should have panicked")
			}
			if !strings.Contains(reason.(error).Error(), "fail to init package cache") {
				t.Errorf("unexpected panic message")
			}
		}()
		InitPackageASTCache("./testdata/ko")
	})
}

func Test_GetPackageAST(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		pkgDir := "./testdata/ok"
		pkgPath := "github.com/krostar/test/internal/code/testdata/ok"

		_astPkgPathToPkg = nil

		// not in cache
		pkgs, err := GetPackageAST(pkgDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if _, found := pkgs[pkgPath]; !found {
			t.Fatalf("package %s not found in pkgs", pkgPath)
		}

		// now in cache
		if _astPkgPathToPkg == nil || _astPkgPathToPkg[pkgDir] == nil {
			t.Errorf("package should be in cache")
		}

		if _, err = GetPackageAST(pkgDir); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("ko", func(t *testing.T) {
		pkgs, err := GetPackageAST("./testdata/ko")
		if err == nil || pkgs != nil {
			t.Fatalf("expected failure")
		}
	})
}
