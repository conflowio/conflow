// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package test

import (
	"time"

	. "github.com/onsi/gomega"

	"github.com/opsidian/conflow/conflow"
	"github.com/opsidian/conflow/conflow/block"
)

var _ conflow.BlockDirective = &Directive{}

// @block
type Directive struct {
	// @id
	IDField conflow.ID
	// @value
	Value             interface{}
	FieldString       string
	FieldInt          int64
	FieldFloat        float64
	FieldBool         bool
	FieldArray        []interface{}
	FieldMap          map[string]interface{}
	FieldTimeDuration time.Duration
	// @name "custom_field"
	FieldCustomName string

	// @name "testblock"
	Blocks []*Block
}

func (d *Directive) ID() conflow.ID {
	return d.IDField
}

func (d *Directive) EvalStage() conflow.EvalStage {
	return conflow.EvalStageInit
}

func (d *Directive) ApplyToRuntimeConfig(config *conflow.RuntimeConfig) {
}

func (d *Directive) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"testblock": BlockInterpreter{},
		},
	}
}

func (d *Directive) Compare(d2 *Directive, input string) {
	compareBlocks(d, d2, DirectiveInterpreter{}, input)
	Expect(len(d.Blocks)).To(Equal(len(d2.Blocks)), "child block count does not match, input: %s", input)
	for i, c := range d2.Blocks {
		compareBlocks(c, d2.Blocks[i], BlockInterpreter{}, input)
	}
}
