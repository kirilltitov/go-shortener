package noexitinmain

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/types/typeutil"
)

// Analyzer представляет собою анализатор, запрещающий использование функции [os.Exit] в функции main.main
var Analyzer = &analysis.Analyzer{
	Name: "noexitinmain",
	Doc:  "Осуществляет проверку на присутствие вызова функции os.Exit в функции main.main",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	isExit := func(expr *ast.CallExpr) bool {
		const pkg = "os"
		const exitFuncName = "Exit"

		f := typeutil.StaticCallee(pass.TypesInfo, expr)

		if f == nil {
			return false
		}
		if f.Pkg() == nil || f.Pkg().Path() != pkg {
			return false
		}
		if f.Type().(*types.Signature).Recv() != nil {
			return false
		}

		return f.Name() == exitFuncName
	}

	const reportText = `os.Exit in main.main function is forbidden`

	for _, file := range pass.Files {
		inMain := 0
		ast.Inspect(file, func(node ast.Node) bool {
			if inMain == 0 {
				switch x := node.(type) {
				case *ast.FuncDecl:
					if x.Name.Name == "main" {
						inMain++
					} else {
						return false
					}
				}
				return true
			}

			if node == nil {
				inMain--
				return true
			}
			inMain++

			switch x := node.(type) {
			case *ast.ExprStmt:
				if call, ok := x.X.(*ast.CallExpr); ok {
					if isExit(call) {
						pass.Reportf(call.Pos(), reportText)
					}
				}
			case *ast.DeferStmt:
				if isExit(x.Call) {
					pass.Reportf(x.Call.Pos(), reportText)
				}
			case *ast.GoStmt:
				if isExit(x.Call) {
					pass.Reportf(x.Call.Pos(), reportText)
				}
			}

			return true
		})
	}

	return nil, nil
}
