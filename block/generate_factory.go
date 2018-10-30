package block

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/opsidian/ocl/ocl"
	"github.com/opsidian/ocl/util"
)

type Field struct {
	Name        string
	FieldName   string
	Type        string
	Required    bool
	Stage       string
	IsID        bool
	IsValue     bool
	IsReference bool
	IsBlock     bool
	IsFactory   bool
}

type FactoryTemplateParams struct {
	Package                string
	Type                   string
	Name                   string
	HasForeignID           bool
	Stages                 []string
	BlockTypes             map[string]string
	FactoryTypes           map[string]string
	IDField                *Field
	ValueField             *Field
	Fields                 []*Field
	NodeValueFunctionNames map[string]string
}

const factoryTemplate = `
// Code generated by ocl generate. DO NOT EDIT.
package {{.Package}}

import (
	"errors"
	"fmt"

	"github.com/opsidian/ocl/ocl"
	"github.com/opsidian/ocl/util"
	"github.com/opsidian/parsley/parsley"
)

{{ $root := .}}
// New{{$root.Name}}Factory creates a new {{$root.Name}} block factory
func New{{$root.Name}}Factory(
	typeNode parsley.Node,
	idNode parsley.Node,
	paramNodes map[string]parsley.Node,
	blockNodes []parsley.Node,
) (ocl.BlockFactory, parsley.Error) {
	return &{{.Name}}Factory{
		typeNode:   typeNode,
		idNode:     idNode,
		paramNodes: paramNodes,
		blockNodes: blockNodes,
	}, nil
}

// {{$root.Name}}Factory will create and evaluate a {{$root.Name}} block
type {{$root.Name}}Factory struct {
	typeNode   parsley.Node
	idNode     parsley.Node
	paramNodes map[string]parsley.Node
	blockNodes []parsley.Node
	shortFormat bool
}

// CreateBlock creates a new {{$root.Name}} block
func (f *{{$root.Name}}Factory) CreateBlock(parentCtx interface{}) (ocl.Block, interface{}, parsley.Error) {
	var err parsley.Error

	block := &{{$root.Name}}{}

	if block.{{.IDField.FieldName}}, err = util.NodeStringValue(f.idNode, parentCtx); err != nil {
		return nil, nil, err
	}

	{{ if .HasForeignID }}
	idRegistry := parentCtx.(ocl.IDRegistryAware).GetIDRegistry()
	if !idRegistry.IDExists(id) {
		return nil, nil, parsley.NewErrorf(f.idNode.Pos(), "%q does not exist", id)
	}
	{{ end }}

	ctx := block.Context(parentCtx)

	{{ if .ValueField }}
	if valueNode, ok := f.paramNodes["_value"]; ok {
		f.shortFormat = true

		if block.{{.ValueField.FieldName}}, err = util.{{index $root.NodeValueFunctionNames .ValueField.Type}}(valueNode, ctx); err != nil {
			return nil, nil, err
		}
	}
	{{ end }}

	if len(f.blockNodes) > 0 {
		var childBlockFactory interface{}
		for _, childBlock := range f.blockNodes {
			if childBlockFactory, err = childBlock.Value(ctx); err != nil {
				return nil, nil, err
			}
			{{- if len .FactoryTypes}}
			switch b := childBlockFactory.(type) {
			{{- range $type, $fieldName := .FactoryTypes -}}
			case {{$type}}:
				block.{{$fieldName}} = append(block.{{$fieldName}}, b)
			{{- end}}
			default:
				panic(fmt.Sprintf("block type %T is not supported in {{.Name}}, please open a bug ticket", childBlockFactory))
			}
			{{ else }}
			panic(fmt.Sprintf("block type %T is not supported in {{.Name}}, please open a bug ticket", childBlockFactory))
			{{- end }}
		}
	}

	return block, ctx, nil
}

// EvalBlock evaluates all fields belonging to the given stage on a {{$root.Name}} block
func (f *{{$root.Name}}Factory) EvalBlock(ctx interface{}, stage string, res ocl.Block) parsley.Error {
	var err parsley.Error

	block, ok := res.(*{{$root.Name}})
	if !ok {
		panic("result must be a type of *{{$root.Name}}")
	}

	if preInterpreter, ok := res.(ocl.BlockPreInterpreter); ok {
		if err = preInterpreter.PreEval(ctx, stage); err != nil {
			return err
		}
	}

	if !f.shortFormat {
		switch stage {
		{{- range $stage := $root.Stages}}
		case "{{$stage}}":
			{{- range $root.Fields -}}{{- if and (eq .Stage $stage) (not .IsID) (not .IsFactory) (not .IsBlock)}}
			if valueNode, ok := f.paramNodes["{{.Name}}"]; ok {
				if block.{{.FieldName}}, err = util.{{index $root.NodeValueFunctionNames .Type}}(valueNode, ctx); err != nil {
					return err
				}
			}{{ if .Required }} else {
				return parsley.NewError(f.typeNode.Pos(), errors.New("{{.Name}} parameter is required"))
			}{{ end }}
			{{ end }}{{ end -}}
		{{end -}}
		default:
			panic(fmt.Sprintf("unknown stage: %s", stage))
		}

		switch stage {
			{{- range $stage := $root.Stages}}
			case "{{$stage}}":
				{{- range $root.Fields }}{{ if and (eq .Stage $stage) .IsFactory}}
				var childBlock ocl.Block
				var childBlockCtx interface{}
				for _, childBlockFactory := range block.{{.FieldName}} {
					if childBlock, childBlockCtx, err = childBlockFactory.CreateBlock(ctx); err != nil {
						return err
					}

					if err = childBlockFactory.EvalBlock(childBlockCtx, stage, childBlock); err != nil {
						return err
					}

					{{ if $root.BlockTypes -}}
					switch b := childBlock.(type) {
					{{- range $type, $fieldName := $root.BlockTypes -}}
					case {{$type}}:
						block.{{$fieldName}} = append(block.{{$fieldName}}, b)
					{{- end}}
					default:
						panic(fmt.Sprintf("block type %T is not supported in {{.Name}}, please open a bug ticket", childBlock))
					}
					{{ else }}
					panic(fmt.Sprintf("block type %T is not supported in {{.Name}}, please open a bug ticket", childBlock))
					{{- end }}
				}
				{{ end }}{{ end -}}
		{{end -}}
		default:
			panic(fmt.Sprintf("unknown stage: %s", stage))
		}
	}

	if postInterpreter, ok := res.(ocl.BlockPostInterpreter); ok {
		if err = postInterpreter.PostEval(ctx, stage); err != nil {
			return err
		}
	}

	return nil
}

// HasForeignID returns true if the block ID is referencing an other block id
func (f *{{$root.Name}}Factory) HasForeignID() bool {
	return {{.HasForeignID}}
}

// HasShortFormat returns true if the block can be defined in the short block format
func (f *{{$root.Name}}Factory) HasShortFormat() bool {
	return {{ if .ValueField }}true{{ else }}false{{ end }}
}
`

func GenerateFactory(dir string, name string, pkgName string) ([]byte, error) {
	params, err := generateTemplateParams(dir, name, pkgName)
	if err != nil {
		return nil, err
	}

	tmpl := template.Must(template.New("block_factory").Parse(factoryTemplate))

	res := &bytes.Buffer{}
	err = tmpl.Execute(res, params)
	if err != nil {
		return nil, err
	}

	formatted, err := format.Source(res.Bytes())
	if err != nil {
		return nil, err
	}

	return formatted, nil
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

	fields, err := getFields(str, file)
	if err != nil {
		return nil, err
	}
	var idField, valueField *Field
	var hasForeignID bool

	stages := []string{}
	blockTypes := map[string]string{}
	factoryTypes := map[string]string{}
	for _, field := range fields {
		if !util.StringSliceContains(stages, field.Stage) {
			stages = append(stages, field.Stage)
		}
		switch {
		case field.IsID:
			idField = field
			hasForeignID = field.IsReference
		case field.IsValue:
			valueField = field
		case field.IsFactory:
			factoryTypes[field.Type] = field.FieldName
		case field.IsBlock:
			blockTypes[field.Type] = field.FieldName
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

func getFields(str *ast.StructType, file *ast.File) ([]*Field, error) {
	fields := make([]*Field, 0, len(str.Fields.List))

	var idField string
	var valueField string
	for _, field := range str.Fields.List {
		var tag string
		fieldName := field.Names[0].String()
		if field.Tag != nil {
			var err error
			tag, err = strconv.Unquote(field.Tag.Value)
			if err != nil {
				return nil, fmt.Errorf("tag is invalid for %s", fieldName)
			}
		}
		ftype, valid := getFieldType(field.Type)

		tags := util.ParseFieldTag(reflect.StructTag(tag), "ocl", fieldName)

		for _, key := range tags.Keys() {
			if _, validTag := ocl.BlockTags[strings.ToLower(key)]; !validTag {
				return nil, fmt.Errorf("invalid tag %s on %s", key, fieldName)
			}
		}

		if tags.GetBool(ocl.BlockTagIgnore) {
			continue
		}

		isFactory := tags.GetBool(ocl.BlockTagFactory)
		isBlock := tags.GetBool(ocl.BlockTagBlock)
		if isBlock || isFactory {
			ftype = strings.TrimPrefix(ftype, "[]")
		}

		if !valid && !isBlock && !isFactory {
			return nil, fmt.Errorf("invalid field type: %s, use valid type or use ignore tag", fieldName)
		}

		isID := tags.GetBool(ocl.BlockTagID)
		if isID {
			if idField != "" {
				return nil, fmt.Errorf("multiple id fields were found: %s, %s", idField, fieldName)
			}
			idField = fieldName
		}

		isReference := tags.GetBool(ocl.BlockTagReference)
		if isReference && !isID {
			return nil, errors.New("reference tag can only be set on the id field")
		}

		isValue := tags.GetBool(ocl.BlockTagValue)
		if isValue {
			if valueField != "" {
				return nil, fmt.Errorf("multiple value fields were found: %s, %s", valueField, fieldName)
			}
			valueField = fieldName
		}
		if isValue && isID {
			return nil, errors.New("the value field can not be the id field")
		}

		name := tags.GetWithDefault(ocl.BlockTagName, generateName(fieldName))
		if name == "" {
			return nil, fmt.Errorf("name can not be empty: %s", fieldName)
		}

		stage := tags.GetWithDefault(ocl.BlockTagStage, "default")
		if stage == "" {
			return nil, fmt.Errorf("stage can not be empty: %s", fieldName)
		}

		fields = append(fields, &Field{
			Name:        name,
			FieldName:   fieldName,
			Required:    tags.GetBool(ocl.BlockTagRequired),
			Type:        ftype,
			Stage:       stage,
			IsID:        isID,
			IsValue:     isValue,
			IsReference: isReference,
			IsBlock:     isBlock,
			IsFactory:   isFactory,
		})
	}

	if valueField != "" {
		for _, field := range fields {
			if !field.IsValue && field.Required {
				return nil, errors.New("when setting a value field then other fields can not be required")
			}
		}
	}

	return fields, nil
}

func getFieldType(typeNode ast.Expr) (string, bool) {
	switch t := typeNode.(type) {
	case *ast.Ident:
		ftype := t.String()
		_, valid := ocl.VariableTypes[ftype]
		return ftype, valid
	default:
		b := &bytes.Buffer{}
		format.Node(b, token.NewFileSet(), t)
		_, valid := ocl.VariableTypes[b.String()]
		return b.String(), valid
	}

}

func generateName(str string) string {
	re := regexp.MustCompile("[A-Z][a-z0-9_]*")
	str = re.ReplaceAllStringFunc(str, func(str string) string {
		return "_" + strings.ToLower(str)
	})
	return strings.TrimLeft(str, "_")
}
