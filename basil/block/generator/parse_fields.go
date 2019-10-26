// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

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

	"github.com/opsidian/basil/basil/variable"

	"github.com/opsidian/basil/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"

	"github.com/opsidian/basil/basil"
)

// ParseFields parses all fields of a given go struct
func ParseFields(str *ast.StructType, file *ast.File) (Fields, error) {
	fields := make(Fields, 0, len(str.Fields.List))

	var idField string
	var valueField string

	for _, field := range str.Fields.List {
		field, err := parseField(field)
		if err != nil {
			return nil, err
		}
		if field == nil {
			continue
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

		fields = append(fields, field)
	}

	if idField == "" {
		return nil, fmt.Errorf("no fields were found with the \"id\" tag")
	}

	if valueField != "" {
		for _, field := range fields {
			if !field.IsValue && field.IsRequired {
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

	params, err := parseFieldTag(reflect.StructTag(tag))
	if err != nil {
		return nil, err
	}

	for k := range params {
		if _, validTag := basil.BlockTags[k]; !validTag {
			return nil, fmt.Errorf("invalid tag %q on field %q", k, name)
		}
	}

	if params[basil.BlockTagIgnore] == true {
		return nil, nil
	}

	var paramName basil.ID
	if params[basil.BlockTagName] == nil {
		paramName = generateParamName(name)
	} else {
		paramName, err = variable.IdentifierValue(params[basil.BlockTagName])
		if err != nil {
			return nil, fmt.Errorf("was expecting string value for tag %q on field %q", basil.BlockTagName, name)
		}
	}
	var stage basil.ID
	if params[basil.BlockTagStage] == nil {
		stage = "main"
	} else {
		stage, err = variable.IdentifierValue(params[basil.BlockTagStage])
		if err != nil {
			return nil, fmt.Errorf("was expecting string value for tag %q on field %q", basil.BlockTagStage, name)
		}
	}

	isID := params[basil.BlockTagID] == true
	if isID {
		paramName = "id"
	}

	field := &Field{
		Name:        name,
		ParamName:   string(paramName),
		IsRequired:  params[basil.BlockTagRequired] == true,
		Stage:       string(stage),
		Default:     params[basil.BlockTagDefault],
		IsID:        isID,
		IsValue:     params[basil.BlockTagValue] == true,
		IsReference: params[basil.BlockTagReference] == true,
		IsBlock:     params[basil.BlockTagBlock] == true,
		IsOutput:    params[basil.BlockTagOutput] == true,
		IsGenerated: params[basil.BlockTagGenerated] == true,
	}

	setFieldType(astField.Type, field)

	if field.IsGenerated {
		field.IsBlock = true
		field.IsRequired = true
		field.Stage = "init"
	}

	if field.IsOutput {
		field.Stage = "close"
	}

	return field, nil
}

func setFieldType(typeNode ast.Expr, field *Field) {
	switch t := typeNode.(type) {
	case *ast.Ident:
		field.Type = t.String()
		return
	case *ast.ArrayType:
		if field.IsBlock {
			field.IsMany = true
			setFieldType(t.Elt, field)
			return
		}
	}

	b := &bytes.Buffer{}
	if err := format.Node(b, token.NewFileSet(), typeNode); err != nil {
		panic(err)
	}
	field.Type = b.String()
}

func parseFieldTag(tag reflect.StructTag) (map[basil.ID]interface{}, error) {
	value := strings.TrimSpace(tag.Get("basil"))
	if value == "" {
		return nil, nil
	}
	f := text.NewFile("", []byte(value))
	fs := parsley.NewFileSet(f)
	ctx := parsley.NewContext(fs, text.NewReader(f))
	val, err := parsley.Evaluate(ctx, parser.StructTag())
	if err != nil {
		return nil, err
	}
	return val.(map[basil.ID]interface{}), nil
}

func generateParamName(name string) basil.ID {
	re := regexp.MustCompile("[A-Z][a-z0-9]*")
	name = re.ReplaceAllStringFunc(name, func(str string) string {
		return "_" + strings.ToLower(str)
	})
	return basil.ID(strings.TrimLeft(name, "_"))
}
