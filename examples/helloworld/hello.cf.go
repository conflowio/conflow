// Code generated by Conflow. DO NOT EDIT.

package main

import (
	"fmt"
	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/schema"
)

// HelloInterpreter is the conflow interpreter for the Hello block
type HelloInterpreter struct {
	s schema.Schema
}

func (i HelloInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Metadata: schema.Metadata{
				Description: "It is capable to print some greetings",
			},
			Name: "Hello",
			Properties: map[string]schema.Schema{
				"greeting": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"eval_stage": "close"},
						ReadOnly:    true,
					},
				},
				"id": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"id": "true"},
						ReadOnly:    true,
					},
					Format: "conflow.ID",
				},
				"to": &schema.String{},
			},
			Required: []string{"to"},
		}
	}
	return i.s
}

// Create creates a new Hello block
func (i HelloInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	return &Hello{
		id: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i HelloInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i HelloInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Hello
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i HelloInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "greeting":
		return b.(*Hello).greeting
	case "id":
		return b.(*Hello).id
	case "to":
		return b.(*Hello).to
	default:
		panic(fmt.Errorf("unexpected parameter %q in Hello", name))
	}
}

func (i HelloInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Hello)
	switch name {
	case "to":
		b.to = value.(string)
	}
	return nil
}

func (i HelloInterpreter) SetBlock(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}