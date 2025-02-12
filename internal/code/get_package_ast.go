// Package code provides utilities for working with Go source code and ASTs.
package code

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/tools/go/packages"
)

//nolint:gochecknoglobals // those variable are required to keep a global cache
var (
	// _astLock provides synchronization for the package AST cache.
	_astLock sync.Mutex

	// _astPkgPathToPkg is a global cache of parsed package ASTs.
	// The first key is the package directory, and the second key is the package path.
	// This allows for efficient reuse of parsed ASTs across multiple assertions.
	_astPkgPathToPkg map[string]map[string]*packages.Package
)

// InitPackageASTCache initializes the package AST cache.
// It is usually called from a TestMain function.
// It parses and caches the AST for the package located at pkgDir.
// It panics if the package cannot be parsed.
func InitPackageASTCache(pkgDir string) {
	if _, err := GetPackageAST(pkgDir); err != nil {
		panic(fmt.Errorf("fail to init package cache: %v", err))
	}
}

// GetPackageAST retrieves the parsed AST for a given package directory.
// It returns a map from package paths to parsed packages.
// The function uses a global cache to avoid reparsing the same package multiple times.
// If the package is not already cached it attempts to parse the package, caches it and
// returns the result.
// It returns an error if the package cannot be parsed.
func GetPackageAST(pkgDir string) (map[string]*packages.Package, error) {
	_astLock.Lock()
	defer _astLock.Unlock()

	if found, ok := _astPkgPathToPkg[pkgDir]; ok {
		return found, nil
	}

	pkgPathToPkg, err := ParsePackageAST(context.Background(), pkgDir)
	if err != nil {
		return nil, fmt.Errorf("unable to parse caller package %q: %w", pkgDir, err)
	}

	if _astPkgPathToPkg == nil {
		_astPkgPathToPkg = make(map[string]map[string]*packages.Package)
	}

	if _astPkgPathToPkg[pkgDir] == nil {
		_astPkgPathToPkg[pkgDir] = make(map[string]*packages.Package)
	}

	_astPkgPathToPkg[pkgDir] = pkgPathToPkg
	return pkgPathToPkg, nil
}
