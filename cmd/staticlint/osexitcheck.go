package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var OSExitCheckAnalizer = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc:  "check for calls to os.Exit in main.go",
	Run:  run,
}

// run inspects the abstract syntax tree of the source code and checks for calls
// to os.Exit in the main.go file.
func run(pass *analysis.Pass) (interface{}, error) {
	// Iterate over all files in the package.
	for _, file := range pass.Files {
		// Use the ast.Inspect function to traverse the abstract syntax tree of
		// each file.
		ast.Inspect(file, func(node ast.Node) bool {
			// Only inspect the main.go file.
			if file.Name.Name != "main" {
				return true
			}
			// Only inspect CallExpr nodes.
			callExpr, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			// Check whether the CallExpr node is a call to os.Exit.
			if fun, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := fun.X.(*ast.Ident); ok && ident.Name == "os" && fun.Sel.Name == "Exit" {
					// If the CallExpr is a call to os.Exit, report the error.
					pass.Reportf(ident.NamePos, "call to os.Exit in main.go")
					return false
				}
			}
			return true
		})
	}
	return nil, nil
}
