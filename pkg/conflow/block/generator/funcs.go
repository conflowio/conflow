// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generator

import (
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
