// Code generated by Conflow. DO NOT EDIT.

package main

import (
	"fmt"
	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/schema"
)

// LicensifyInterpreter is the conflow interpreter for the Licensify block
type LicensifyInterpreter struct {
	s schema.Schema
}

func (i LicensifyInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Name: "Licensify",
			Properties: map[string]schema.Schema{
				"id": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"id": "true"},
						ReadOnly:    true,
					},
					Format: "conflow.ID",
				},
				"license": &schema.String{},
				"path":    &schema.String{},
			},
			Required: []string{"path", "license"},
		}
	}
	return i.s
}

// Create creates a new Licensify block
func (i LicensifyInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	return &Licensify{
		id: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i LicensifyInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i LicensifyInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Licensify
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i LicensifyInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*Licensify).id
	case "license":
		return b.(*Licensify).license
	case "path":
		return b.(*Licensify).path
	default:
		panic(fmt.Errorf("unexpected parameter %q in Licensify", name))
	}
}

func (i LicensifyInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Licensify)
	switch name {
	case "license":
		b.license = value.(string)
	case "path":
		b.path = value.(string)
	}
	return nil
}

func (i LicensifyInterpreter) SetBlock(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}