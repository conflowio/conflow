// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"fmt"

	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/schema"
)

// Block is the directive for marking structs as conflow blocks
//
// @block "directive"
type Block struct {
	// @id
	id conflow.ID
	// @enum ["ignore", "init", "parse", "resolve"]
	EvalStage string
	Path      string
	// @enum ["configuration", "directive", "generator", "main", "task"]
	// @value
	// @required
	Type string
}

func (b *Block) ID() conflow.ID {
	return b.id
}

func (b *Block) ApplyToSchema(s schema.Schema) error {
	if _, ok := s.(*schema.Object); !ok {
		return fmt.Errorf("@block can only be used on a struct")
	}

	if b.EvalStage != "" {
		s.(*schema.Object).SetAnnotation(conflow.AnnotationEvalStage, b.EvalStage)
	}

	return nil
}
