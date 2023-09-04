package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "norawquery is a linter that checks for raw sql query function calls in your Go files"

var Analyzer = &analysis.Analyzer{
	Name: "norawquery",
	Doc:  doc,
	Run:  runAnalyzer,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func runAnalyzer(pass *analysis.Pass) (interface{}, error) {
	result := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	result.Preorder(nodeFilter, func(n ast.Node) {
		if n, ok := n.(*ast.CallExpr); ok {
			selector, ok := n.Fun.(*ast.SelectorExpr)
			if !ok {
				return
			}

			if pass.TypesInfo == nil || pass.TypesInfo.Uses[selector.Sel] == nil || pass.TypesInfo.Uses[selector.Sel].Pkg() == nil {
				return
			}

			if pass.TypesInfo.Uses[selector.Sel].Pkg().Path() != "gorm.io/gorm" {
				return
			}

			if !stringIncludes(selector.Sel.Name, "Raw", "Exec") {
				return
			}

			pass.Reportf(n.Fun.Pos(), "Please do not use raw sql query, use gorm ORM instead!")
		}
	})

	return nil, nil
}

func stringIncludes(str string, includes ...string) bool {
	for _, v := range includes {
		if str == v {
			return true
		}
	}
	return false
}
