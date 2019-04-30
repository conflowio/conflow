package generator

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/util"
)

// ParseFields parses all fields of a given go struct
func ParseFields(str *ast.StructType, file *ast.File) ([]*Field, error) {
	fields := make([]*Field, 0, len(str.Fields.List))

	var idField string
	var valueField string

	for _, field := range str.Fields.List {
		field, err := parseField(field)
		if err != nil {
			return nil, err
		}

		if field.IsID {
			if idField != "" {
				return nil, fmt.Errorf("multiple id fields were found: %s, %s", idField, field.Name)
			}
			idField = field.Name
		}

		if field.IsValue {
			if valueField != "" {
				return nil, fmt.Errorf("multiple value fields were found: %s, %s", valueField, field.Name)
			}
			valueField = field.Name
		}

		if field != nil {
			fields = append(fields, field)
		}
	}

	if idField == "" {
		return nil, fmt.Errorf("no fields were found with the \"id\" tag")
	}

	if valueField != "" {
		for _, field := range fields {
			if !field.IsValue && field.Required {
				return nil, errors.New("when setting a value field then no other fields can be required")
			}
		}
	}

	return fields, nil
}

func parseField(astField *ast.Field) (*Field, error) {
	var tag string
	name := astField.Names[0].String()
	if astField.Tag != nil {
		var err error
		tag, err = strconv.Unquote(astField.Tag.Value)
		if err != nil {
			return nil, fmt.Errorf("tag is invalid for %s", name)
		}
	}

	tags := util.ParseFieldTag(reflect.StructTag(tag), "basil", name)

	for _, key := range tags.Keys() {
		if _, validTag := basil.BlockTags[strings.ToLower(key)]; !validTag {
			return nil, fmt.Errorf("invalid tag %q on field %q", key, name)
		}
	}

	if tags.GetBool(basil.BlockTagIgnore) {
		return nil, nil
	}

	paramName := tags.GetWithDefault(basil.BlockTagName, generateParamName(name))

	isID := tags.GetBool(basil.BlockTagID)
	if isID {
		paramName = "id"
	}

	blockType, isBlock := tags.Get("block")
	if isBlock {
		paramName = blockType
	}

	field := &Field{
		Name:        name,
		ParamName:   paramName,
		Required:    tags.GetBool(basil.BlockTagRequired),
		Type:        getFieldType(astField.Type),
		Stage:       tags.GetWithDefault(basil.BlockTagStage, "default"),
		IsID:        isID,
		IsValue:     tags.GetBool(basil.BlockTagValue),
		IsReference: tags.GetBool(basil.BlockTagReference),
		IsBlock:     isBlock,
		IsNode:      tags.GetBool(basil.BlockTagNode),
	}

	if !field.IsID && !field.IsBlock && !field.IsNode {
		field.IsParam = true
	}

	return field, nil
}

func getFieldType(typeNode ast.Expr) string {
	switch t := typeNode.(type) {
	case *ast.Ident:
		return t.String()
	default:
		b := &bytes.Buffer{}
		if err := format.Node(b, token.NewFileSet(), t); err != nil {
			panic(err)
		}
		return b.String()
	}
}

func generateParamName(name string) string {
	re := regexp.MustCompile("[A-Z][a-z0-9]*")
	name = re.ReplaceAllStringFunc(name, func(str string) string {
		return "_" + strings.ToLower(str)
	})
	return strings.TrimLeft(name, "_")
}
