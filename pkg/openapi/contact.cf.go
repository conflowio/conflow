// Code generated by Conflow. DO NOT EDIT.

package openapi

import (
	"fmt"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/conflow/annotations"
	"github.com/conflowio/conflow/pkg/conflow/types"
	"github.com/conflowio/conflow/pkg/schema"
)

func init() {
	schema.Register(&schema.Object{
		Metadata: schema.Metadata{
			Annotations: map[string]string{
				annotations.Type: "configuration",
			},
			ID: "github.com/conflowio/conflow/pkg/openapi.Contact",
		},
		FieldNames: map[string]string{"email": "Email", "name": "Name", "url": "URL"},
		Properties: map[string]schema.Schema{
			"email": &schema.String{
				Format: "email",
			},
			"name": &schema.String{},
			"url": &schema.String{
				Format: "uri",
			},
		},
	})
}

// NewContactWithDefaults creates a new Contact instance with default values
func NewContactWithDefaults() *Contact {
	b := &Contact{}
	return b
}

// ContactInterpreter is the Conflow interpreter for the Contact block
type ContactInterpreter struct {
}

func (i ContactInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/openapi.Contact")
	return s
}

// Create creates a new Contact block
func (i ContactInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := NewContactWithDefaults()
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i ContactInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i ContactInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Contact
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i ContactInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "email":
		return b.(*Contact).Email
	case "name":
		return b.(*Contact).Name
	case "url":
		return b.(*Contact).URL
	default:
		panic(fmt.Errorf("unexpected parameter %q in Contact", name))
	}
}

func (i ContactInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Contact)
	switch name {
	case "email":
		b.Email = schema.Value[types.Email](value)
	case "name":
		b.Name = schema.Value[string](value)
	case "url":
		b.URL = schema.Value[types.URL](value)
	}
	return nil
}

func (i ContactInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	return nil
}
