package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func getType(filename string, line int) (*ast.File, ast.Node, string, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return nil, nil, "", err
	}

	fs := token.NewFileSet()
	fs.AddFile(filename, fs.Base(), int(fi.Size()))

	var res ast.Node
	var name string

	file, err := parser.ParseFile(fs, filename, nil, parser.AllErrors)
	ast.Inspect(file, func(node ast.Node) bool {
		if node != nil && node.Pos().IsValid() {
			if fs.Position(node.Pos()).Line == line+1 {
				switch n := node.(type) {
				case *ast.TypeSpec:
					if str, isStruct := n.Type.(*ast.StructType); isStruct {
						res = str
						name = n.Name.Name
					}
				case *ast.FuncDecl:
					res = n.Type
					name = n.Name.Name
				}
			}
		}
		return true
	})

	if res == nil {
		return nil, nil, "", fmt.Errorf("'%s' does not refer to a struct or function", name)
	}

	return file, res, name, nil
}
