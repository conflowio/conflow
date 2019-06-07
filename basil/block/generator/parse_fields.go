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

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/util"
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

	field := &Field{
		Name:        name,
		ParamName:   paramName,
		IsRequired:  tags.GetBool(basil.BlockTagRequired),
		Stage:       tags.GetWithDefault(basil.BlockTagStage, "main"),
		IsID:        isID,
		IsValue:     tags.GetBool(basil.BlockTagValue),
		IsReference: tags.GetBool(basil.BlockTagReference),
		IsBlock:     tags.GetBool(basil.BlockTagBlock),
		IsOutput:    tags.GetBool(basil.BlockTagOutput),
	}

	setFieldType(astField.Type, field)

	if field.IsChannel {
		field.IsBlock = true
		field.IsRequired = true
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
	case *ast.ChanType:
		field.IsChannel = true
		setFieldType(t.Value, field)
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

func generateParamName(name string) string {
	re := regexp.MustCompile("[A-Z][a-z0-9]*")
	name = re.ReplaceAllStringFunc(name, func(str string) string {
		return "_" + strings.ToLower(str)
	})
	return strings.TrimLeft(name, "_")
}
