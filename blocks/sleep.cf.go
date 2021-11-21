// Code generated by Conflow. DO NOT EDIT.

package blocks

import (
	"fmt"
	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/schema"
	"time"
)

// SleepInterpreter is the conflow interpreter for the Sleep block
type SleepInterpreter struct {
	s schema.Schema
}

func (i SleepInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Metadata: schema.Metadata{
				Annotations: map[string]string{"block.conflow.io/type": "task"},
			},
			Name: "Sleep",
			Parameters: map[string]schema.Schema{
				"duration": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"block.conflow.io/value": "true"},
					},
					Format: "duration-go",
				},
				"id": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"block.conflow.io/id": "true"},
						ReadOnly:    true,
					},
					Format: "conflow.ID",
				},
			},
			Required: []string{"duration"},
		}
	}
	return i.s
}

// Create creates a new Sleep block
func (i SleepInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	return &Sleep{
		id: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i SleepInterpreter) ValueParamName() conflow.ID {
	return "duration"
}

// ParseContext returns with the parse context for the block
func (i SleepInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Sleep
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i SleepInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "duration":
		return b.(*Sleep).duration
	case "id":
		return b.(*Sleep).id
	default:
		panic(fmt.Errorf("unexpected parameter %q in Sleep", name))
	}
}

func (i SleepInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Sleep)
	switch name {
	case "duration":
		b.duration = value.(time.Duration)
	}
	return nil
}

func (i SleepInterpreter) SetBlock(block conflow.Block, name conflow.ID, value interface{}) error {
	return nil
}
