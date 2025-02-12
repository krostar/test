package code

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/packages"
)

// GetCallerCallExpr retrieves the *ast.CallExpr at a specific location in the caller's source code.
//
// `pkgs` is a map of package paths to *packages.Package, representing the parsed ASTs.
// `callerFile` is the filename of the caller's source file.
// `callerLine` is the line number in the caller's source file where the call expression is located.
//
// It returns the *ast.CallExpr, the *ast.File containing the expression, the *packages.Package
// to which the file belongs, and an error if any occurred during the process.
// Returns nil values if the package, file or expression is not found.
func GetCallerCallExpr(pkgs map[string]*packages.Package, callerFile string, callerLine int) (*ast.CallExpr, *ast.File, *packages.Package, error) {
	pkg, file := findCallerPackageAndASTFile(pkgs, callerFile)
	if pkg == nil || file == nil {
		return nil, nil, nil, fmt.Errorf("unable to find ast file and package for %s", callerFile)
	}

	expr, err := getASTCallExprAtLine(pkg.Fset, file, callerLine)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to get call expression: %v", err)
	}

	return expr, file, pkg, nil
}

// findCallerPackageAndASTFile searches the provided map of packages for a specific file.
// It returns the package and file if found, or nil if not.
func findCallerPackageAndASTFile(pkgs map[string]*packages.Package, callerFile string) (*packages.Package, *ast.File) {
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			if pkg.Fset.Position(file.Pos()).Filename == callerFile {
				return pkg, file
			}
		}
	}
	return nil, nil
}

// getASTCallExprAtLine retrieves the *ast.CallExpr at a specified line within an *ast.File.
//
// `fset` is the *token.FileSet used for position information.
// `file` is the *ast.File to search within.
// `line` is the target line number.
//
// Returns the *ast.CallExpr if found on the specified line, an error otherwise.
func getASTCallExprAtLine(fset *token.FileSet, file *ast.File, line int) (*ast.CallExpr, error) {
	var callExpr *ast.CallExpr

	ast.Inspect(file, func(node ast.Node) bool {
		if node == nil || callExpr != nil {
			return false
		}

		if fset.Position(node.Pos()).Line != line {
			return true
		}

		if call, ok := node.(*ast.CallExpr); ok {
			callExpr = call
			return false
		}

		return true
	})

	if callExpr == nil {
		return nil, errors.New("ast inspection did not return a node")
	}

	return callExpr, nil
}
