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

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . BlockWithInit
type BlockWithInit interface {
	basil.Block
	basil.BlockInitialiser
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . BlockWithMain
type BlockWithMain interface {
	basil.Block
	basil.BlockRunner
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . BlockWithClose
type BlockWithClose interface {
	basil.Block
	basil.BlockCloser
}

//go:generate basil generate
type Block struct {
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

func (b *Block) ID() basil.ID {
	return b.IDField
}

func (b *Block) ParseContextOverride() basil.ParseContextOverride {
	return basil.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"testblock": BlockInterpreter{},
		},
	}
}

func (b *Block) Compare(b2 *Block, input string) {
	interpreter := BlockInterpreter{}
	compareBlocks(b, b2, interpreter, input)
	Expect(len(b.Blocks)).To(Equal(len(b2.Blocks)), "child block count does not match, input: %s", input)
	for i, c := range b2.Blocks {
		compareBlocks(c, b2.Blocks[i], interpreter, input)
	}
}

func compareBlocks(i1, i2 interface{}, interpreter basil.BlockInterpreter, input string) {
	b1 := i1.(basil.Block)
	b2 := i2.(basil.Block)

	Expect(b1.ID()).To(Equal(b2.ID()), "id does not match, input: %s", input)

	for paramName, _ := range interpreter.Params() {
		v1 := interpreter.Param(b1, paramName)
		v2 := interpreter.Param(b2, paramName)
		if v2 != nil {
			Expect(v1).To(Equal(v2), "%s does not match, input: %s", paramName, input)
		} else {
			Expect(v1).To(BeNil(), "%s does not match, input: %s", paramName, input)
		}
	}
}
