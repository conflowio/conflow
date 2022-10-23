// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package test

import (
	"fmt"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/conflow/block"
	"github.com/conflowio/conflow/pkg/conflow/types"
	"github.com/conflowio/conflow/pkg/schema"
)

//counterfeiter:generate . BlockWithInit
type BlockWithInit interface {
	conflow.Block
	conflow.BlockInitialiser
}

//counterfeiter:generate . BlockWithRun
type BlockWithRun interface {
	conflow.Block
	conflow.BlockRunner
}

//counterfeiter:generate . BlockWithClose
type BlockWithClose interface {
	conflow.Block
	conflow.BlockCloser
}

// @block "configuration"
type Block struct {
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
	FieldTimeDuration types.Duration
	// @name "custom_field"
	FieldCustomName string

	// @name "testblock"
	BlockArray []*Block

	// @name "testblockmap"
	BlockMap map[string]*Block
}

func (b *Block) ID() conflow.ID {
	return b.IDField
}

func (b *Block) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"testblock": BlockInterpreter{},
		},
	}
}

func (b *Block) Compare(b2 *Block, input string) {
	interpreter := BlockInterpreter{}
	compareBlocks(b, b2, interpreter, input)
	Expect(len(b.BlockArray)).To(Equal(len(b2.BlockArray)), "BlockArray count does not match, input: %s", input)
	for i2 := range b2.BlockArray {
		found := false
		for i := range b.BlockArray {
			if b.BlockArray[i].ID() == b2.BlockArray[i2].ID() {
				compareBlocks(b.BlockArray[i], b2.BlockArray[i2], interpreter, input)
				found = true
				break
			}
		}
		if !found {
			ginkgo.Fail(fmt.Sprintf("block not found with id %s", b2.BlockArray[i2].ID()))
		}
	}

	Expect(len(b.BlockMap)).To(Equal(len(b2.BlockMap)), "BlockMap count does not match, input: %s", input)
	for k := range b2.BlockMap {
		if _, ok := b.BlockMap[k]; !ok {
			ginkgo.Fail(fmt.Sprintf("block not found with key %s", k))
		}
		compareBlocks(b.BlockMap[k], b2.BlockMap[k], interpreter, input)
	}
}

func compareBlocks(b1, b2 conflow.Identifiable, interpreter conflow.BlockInterpreter, input string) {
	Expect(b1.ID()).To(Equal(b2.ID()), "id does not match, input: %s", input)

	o := interpreter.Schema().(*schema.Object)
	for jsonPropertyName, p := range o.Properties {
		if block.IsBlockSchema(p) {
			continue
		}
		parameterName := o.ParameterName(jsonPropertyName)

		v1 := interpreter.Param(b1, conflow.ID(parameterName))
		v2 := interpreter.Param(b2, conflow.ID(parameterName))
		if v2 != nil {
			Expect(v1).To(Equal(v2), "%s does not match, input: %s", parameterName, input)
		} else {
			Expect(v1).To(BeNil(), "%s does not match, input: %s", parameterName, input)
		}
	}
}
