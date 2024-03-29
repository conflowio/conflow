// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package block

import (
	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/schema"
)

func IsBlockSchema(s schema.Schema) bool {
	if _, ok := s.(*schema.Reference); ok {
		return true
	}

	if a, ok := s.(*schema.Array); ok {
		if _, ok := a.GetItems().(*schema.Reference); ok {
			return true
		}
	}

	if a, ok := s.(*schema.Map); ok {
		if _, ok := a.GetAdditionalProperties().(*schema.Reference); ok {
			return true
		}
	}

	return false
}

func getNameSchemaForChildBlock(s *schema.Object, node conflow.BlockNode) (conflow.ID, schema.Schema) {
	if p, ok := s.PropertyByParameterName(string(node.ID())); ok {
		return node.ID(), p
	}

	if p, ok := s.PropertyByParameterName(string(node.ParameterName())); ok {
		return node.ParameterName(), p
	}

	return node.ParameterName(), nil
}
