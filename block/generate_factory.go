package block

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"text/template"

	"github.com/opsidian/basil/util"
	"github.com/pkg/errors"
)

func GenerateFactory(dir string, name string, pkgName string) error {
	params, err := generateTemplateParams(dir, name, pkgName)
	if err != nil {
		return err
	}

	tmpl := template.New("block_factory")
	tmpl.Funcs(map[string]interface{}{
		"trimPrefix": func(s string, prefix string) string {
			return strings.TrimPrefix(s, prefix)
		},
	})
	if _, err := tmpl.Parse(factoryTemplate); err != nil {
		return err
	}

	res := &bytes.Buffer{}
	err = tmpl.Execute(res, params)
	if err != nil {
		return err
	}

	formatted, err := format.Source(res.Bytes())
	if err != nil {
		return err
	}

	filename := regexp.MustCompile("[A-Z][a-z0-9_]*").ReplaceAllStringFunc(name, func(str string) string {
		return "_" + strings.ToLower(str)
	})
	filename = strings.TrimLeft(filename, "_") + ".basil.go"
	filePath := path.Join(dir, filename)

	err = ioutil.WriteFile(filePath, formatted, 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to write %s", filePath)
	}

	goimportsCmd := exec.Command("goimports", filename)
	out, err := goimportsCmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "failed to run goimports on  %s", filePath)
	}
	err = ioutil.WriteFile(filePath, out, 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to write %s", filePath)
	}

	fmt.Printf("Wrote `%sFactory` to `%s`\n", name, getRelativePath(filePath))

	return nil
}

func getRelativePath(path string) string {
	_, caller, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(caller)
	return strings.Replace(path, basePath+"/", "", 1)
}

func generateTemplateParams(dir string, name string, pkgName string) (*FactoryTemplateParams, error) {
	packages, err := loadPackages(dir)
	if err != nil {
		return nil, err
	}

	str, file, err := getStruct(packages[pkgName], name)
	if err != nil {
		return nil, err
	}

	fields, err := ParseFields(str, file)
	if err != nil {
		return nil, err
	}
	var idField, valueField *Field
	var hasForeignID bool

	stages := []string{}
	blockTypes := map[string]string{}
	factoryTypes := map[string]string{}
	evalFieldsCnt := 0
	for _, field := range fields {
		if field.Stage != "-" {
			if !util.StringSliceContains(stages, field.Stage) {
				stages = append(stages, field.Stage)
			}
			if !field.IsID && !field.IsBlock {
				evalFieldsCnt++
			}
		}
		switch {
		case field.IsID:
			idField = field
			hasForeignID = field.IsReference
		case field.IsValue:
			valueField = field
		case field.IsFactory:
			factoryTypes[field.Type] = field.Name
		case field.IsBlock:
			blockTypes[field.Type] = field.Name
		}
	}

	return &FactoryTemplateParams{
		Package:                pkgName,
		Name:                   name,
		Stages:                 stages,
		BlockTypes:             blockTypes,
		FactoryTypes:           factoryTypes,
		Fields:                 fields,
		IDField:                idField,
		ValueField:             valueField,
		HasForeignID:           hasForeignID,
		NodeValueFunctionNames: util.NodeValueFunctionNames,
		EvalFieldsCnt:          evalFieldsCnt,
	}, nil
}

func loadPackages(dir string) (map[string]*ast.Package, error) {
	return parser.ParseDir(token.NewFileSet(), dir, nil, parser.AllErrors)
}

func getStruct(pkg *ast.Package, name string) (*ast.StructType, *ast.File, error) {
	var file *ast.File
	var str *ast.StructType
	var err error

	for _, f := range pkg.Files {
		ast.Inspect(f, func(node ast.Node) bool {
			typeSpec, ok := node.(*ast.TypeSpec)
			if !ok || typeSpec.Name.Name != name {
				return true
			}

			switch t := typeSpec.Type.(type) {
			case *ast.StructType:
				str = t
				file = f
				return true
			default:
				err = fmt.Errorf("'%s' does not refer to a struct", name)
			}

			return false
		})

		if str != nil {
			break
		}
	}

	return str, file, err
}
