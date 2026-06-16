// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generator

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/conflowio/conflow/pkg/conflow/annotations"
	"github.com/conflowio/conflow/pkg/conflow/block"
	"github.com/conflowio/conflow/pkg/internal/utils"
	"github.com/conflowio/conflow/pkg/schema"
)

func commonFuncs(blockSchema schema.Schema, imports map[string]string) map[string]interface{} {
	return map[string]interface{}{
		"assignValue": func(s schema.Schema, valueName, resultName string) string {
			return s.AssignValue(imports, valueName, resultName)
		},
		"bindAndAssignValue": func(s schema.Schema, paramName, valueName, resultName string) string {
			schemaSel := utils.EnsureUniqueGoPackageSelector(imports, "github.com/conflowio/conflow/pkg/schema")
			bindSel := utils.EnsureUniqueGoPackageSelector(imports, "github.com/conflowio/conflow/pkg/conflow/bind")

			bindPrefix := fmt.Sprintf(`propSchema, _ := i.Schema().(*%sObject).PropertyByParameterName(%q)
bound, err := %sBindValue(propSchema, %s)
if err != nil {
	return err
}
`, schemaSel, paramName, bindSel, valueName)

			switch s.Type() {
			case schema.TypeArray:
				valuesSel := utils.EnsureUniqueGoPackageSelector(imports, "github.com/conflowio/conflow/pkg/values")
				a := s.(*schema.Array)
				normalize := fmt.Sprintf(`slice, err := %sAsInterfaceSlice(bound)
if err != nil {
	return err
}
`, valuesSel)
				if a.Items.Type() == schema.TypeAny {
					return bindPrefix + normalize + fmt.Sprintf("%s = slice\n", resultName)
				}
				return bindPrefix + normalize + fmt.Sprintf(`%s = make(%s, len(slice))
for slicek, slicev := range slice {
	%s
}`, resultName, a.GoType(imports), indent(a.Items.AssignValue(imports, "slicev", fmt.Sprintf("%s[slicek]", resultName))))
			case schema.TypeMap:
				valuesSel := utils.EnsureUniqueGoPackageSelector(imports, "github.com/conflowio/conflow/pkg/values")
				m := s.(*schema.Map)
				normalize := fmt.Sprintf(`goMap, err := %sAsStringInterfaceMap(bound)
if err != nil {
	return err
}
`, valuesSel)
				if m.AdditionalProperties.Type() == schema.TypeAny {
					return bindPrefix + normalize + fmt.Sprintf("%s = goMap\n", resultName)
				}
				return bindPrefix + normalize + fmt.Sprintf(`%s = make(map[string]%s, len(goMap))
for goMapk, goMapv := range goMap {
	%s
}`, resultName, m.AdditionalProperties.GoType(imports), indent(m.AdditionalProperties.AssignValue(imports, "goMapv", fmt.Sprintf("%s[goMapk]", resultName))))
			case schema.TypeAny:
				valuesSel := utils.EnsureUniqueGoPackageSelector(imports, "github.com/conflowio/conflow/pkg/values")
				return bindPrefix + fmt.Sprintf(`if slice, err := %sAsInterfaceSlice(bound); err == nil {
	%s = slice
} else if goMap, err := %sAsStringInterfaceMap(bound); err == nil {
	%s = goMap
} else {
	%s = bound
}`, valuesSel, resultName, valuesSel, resultName, resultName)
			default:
				return bindPrefix + s.AssignValue(imports, "bound", resultName)
			}
		},
		"filterInputs": func(props map[string]schema.Schema) map[string]schema.Schema {
			return filterSchemaProperties(props, func(s schema.Schema) bool {
				return !s.GetReadOnly()
			})
		},
		"filterParams": func(props map[string]schema.Schema) map[string]schema.Schema {
			return filterSchemaProperties(props, func(s schema.Schema) bool {
				return !block.IsBlockSchema(s)
			})
		},
		"filterBlocks": func(props map[string]schema.Schema) map[string]schema.Schema {
			return filterSchemaProperties(props, block.IsBlockSchema)
		},
		"filterNonID": func(props map[string]schema.Schema) map[string]schema.Schema {
			return filterSchemaProperties(props, func(s schema.Schema) bool {
				return s.GetAnnotation(annotations.ID) != "true"
			})
		},
		"filterDefaults": func(props map[string]schema.Schema) map[string]schema.Schema {
			return filterSchemaProperties(props, func(s schema.Schema) bool {
				return s.DefaultValue() != nil
			})
		},
		"getParameterName": func(name string) string {
			return blockSchema.(*schema.Object).ParameterName(name)
		},
		"getFieldName": func(name string) string {
			return blockSchema.(*schema.Object).FieldName(name)
		},
		"getType": func(s schema.Schema) string {
			return s.GoType(imports)
		},
		"isArray": func(s schema.Schema) bool {
			_, ok := s.(*schema.Array)
			return ok
		},
		"isMap": func(s schema.Schema) bool {
			_, ok := s.(*schema.Map)
			return ok
		},
		"title": func(s string) string {
			return cases.Title(language.English, cases.NoLower).String(s)
		},
		"ensureUniqueGoPackageSelector": utils.EnsureUniqueGoPackageSelector,
	}
}

func indent(s string) string {
	return strings.ReplaceAll(s, "\n", "\n\t")
}
