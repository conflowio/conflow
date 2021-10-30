// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"fmt"

	"github.com/opsidian/conflow/basil/schema"
	"github.com/opsidian/conflow/conflow"
)

// Block is the directive for marking structs as basil blocks
//
// @block
type Block struct {
	// @id
	id   conflow.ID
	Path string
	// @enum ["ignore", "init", "parse", "resolve"]
	EvalStage string
}

func (b *Block) ID() conflow.ID {
	return b.id
}

func (b *Block) ApplyToSchema(s schema.Schema) error {
	if _, ok := s.(*schema.Object); !ok {
		return fmt.Errorf("@block can only be used on a struct")
	}

	if b.EvalStage != "" {
		s.(*schema.Object).SetAnnotation("eval_stage", b.EvalStage)
	}

	return nil
}
