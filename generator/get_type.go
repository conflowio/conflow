package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func getType(dir string, packageName string, name string) (ast.Expr, *ast.File, error) {
	var err error

	packages, err := loadPackages(dir)
	if err != nil {
		return nil, nil, err
	}

	pkg := packages[packageName]
	var file *ast.File
	var res ast.Expr

	for _, f := range pkg.Files {
		ast.Inspect(f, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.TypeSpec:
				if n.Name.Name == name {
					file = f
					var isStruct bool
					res, isStruct = n.Type.(*ast.StructType)
					if !isStruct {
						err = fmt.Errorf("'%s' does not refer to a struct", name)
					}
					return true
				}
			case *ast.FuncDecl:
				if n.Name.Name == name {
					file = f
					res = n.Type
					return true
				}
			default:
				return true
			}
			return false
		})

		if res != nil {
			break
		}
	}

	if res == nil {
		return nil, nil, fmt.Errorf("'%s' does not refer to a struct or function", name)
	}

	return res, file, err
}

func loadPackages(dir string) (map[string]*ast.Package, error) {
	return parser.ParseDir(token.NewFileSet(), dir, nil, parser.AllErrors)
}
