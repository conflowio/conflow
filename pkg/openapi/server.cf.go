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
			ID: "github.com/conflowio/conflow/pkg/openapi.Server",
		},
		FieldNames:     map[string]string{"description": "Description", "url": "URL", "variables": "Variables"},
		ParameterNames: map[string]string{"variables": "variable"},
		Properties: map[string]schema.Schema{
			"description": &schema.String{},
			"url": &schema.String{
				Format: "uri",
			},
			"variables": &schema.Map{
				AdditionalProperties: &schema.Reference{
					Nullable: true,
					Ref:      "github.com/conflowio/conflow/pkg/openapi.ServerVariable",
				},
			},
		},
	})
}

// NewServerWithDefaults creates a new Server instance with default values
func NewServerWithDefaults() *Server {
	b := &Server{}
	b.Variables = map[string]*ServerVariable{}
	return b
}

// ServerInterpreter is the Conflow interpreter for the Server block
type ServerInterpreter struct {
}

func (i ServerInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/openapi.Server")
	return s
}

// Create creates a new Server block
func (i ServerInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := NewServerWithDefaults()
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i ServerInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i ServerInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Server
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i ServerInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "description":
		return b.(*Server).Description
	case "url":
		return b.(*Server).URL
	default:
		panic(fmt.Errorf("unexpected parameter %q in Server", name))
	}
}

func (i ServerInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Server)
	switch name {
	case "description":
		b.Description = schema.Value[string](value)
	case "url":
		b.URL = schema.Value[types.URL](value)
	}
	return nil
}

func (i ServerInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	b := block.(*Server)
	switch name {
	case "variable":
		b.Variables[key] = value.(*ServerVariable)
	}
	return nil
}
