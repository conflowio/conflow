// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package test

import (
	"time"

	. "github.com/onsi/gomega"
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/block"
)

//go:generate basil generate
type Directive struct {
	IDField           basil.ID    `basil:"id"`
	Value             interface{} `basil:"value"`
	FieldString       string
	FieldInt          int64
	FieldFloat        float64
	FieldBool         bool
	FieldArray        []interface{}
	FieldMap          map[string]interface{}
	FieldTimeDuration time.Duration
	FieldCustomName   string `basil:"name=custom_field"`

	Blocks []*Block `basil:"block,name=testblock"`
}

func (d *Directive) ID() basil.ID {
	return d.IDField
}

func (d *Directive) EvalStage() basil.EvalStage {
	return basil.EvalStageInit
}

func (d *Directive) RuntimeConfig() basil.RuntimeConfig {
	return basil.RuntimeConfig{}
}

func (d *Directive) ParseContextOverride() basil.ParseContextOverride {
	return basil.ParseContextOverride{
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
