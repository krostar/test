package code

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/multierr"
	"golang.org/x/tools/go/packages"
)

// ParsePackageAST parses and loads the AST for a given package directory.
//
// `ctx` is the context for the operation, allowing cancellation.
// `pkgDir` is the directory of the package to parse.
//
// It returns a map of package paths to *packages.Package, and an error if the parsing fails
// or if any of the loaded packages contains errors.
func ParsePackageAST(ctx context.Context, pkgDir string) (map[string]*packages.Package, error) {
	// https://github.com/golang/go/issues/27556#issuecomment-419468978
	pkgs, err := packages.Load(&packages.Config{
		Context: ctx,
		Logf:    func(string, ...any) {},
		Mode:    packages.NeedCompiledGoFiles | packages.NeedName | packages.NeedSyntax | packages.NeedTypesInfo,
		Tests:   true,
	}, pkgDir)
	if err != nil {
		return nil, fmt.Errorf("unable to load packages: %w", err)
	}

	{ // check loaded package for errors
		var errs []error
		for _, pkg := range pkgs {
			for _, err := range pkg.Errors {
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			return nil, fmt.Errorf("loaded packages contained errors: %v", multierr.Combine(errs...))
		}
	}

	// allows to easily find a package based on its path
	pkgPathToPkg := make(map[string]*packages.Package)

	packages.Visit(pkgs, func(pkg *packages.Package) bool {
		pkgPathToPkg[pkg.PkgPath] = pkg
		if strings.HasPrefix(pkg.PkgPath, "vendor/") {
			pkgPathWithoutVendor := strings.TrimPrefix(pkg.PkgPath, "vendor/")
			pkgPathToPkg[pkgPathWithoutVendor] = pkg
		}
		return true
	}, nil)

	return pkgPathToPkg, nil
}
