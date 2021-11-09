// Code generated by Conflow. DO NOT EDIT.

package blocks

import (
	"fmt"
	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/schema"
	"time"
)

// TickerInterpreter is the conflow interpreter for the Ticker block
type TickerInterpreter struct {
	s schema.Schema
}

func (i TickerInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Name: "Ticker",
			Properties: map[string]schema.Schema{
				"id": &schema.String{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"block.conflow.io/id": "true"},
						ReadOnly:    true,
					},
					Format: "conflow.ID",
				},
				"interval": &schema.TimeDuration{},
				"tick": &schema.Reference{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"block.conflow.io/eval_stage": "init", "block.conflow.io/generated": "true"},
						Pointer:     true,
					},
					Ref: "http://conflow.schema/github.com/conflowio/conflow/blocks.Tick",
				},
			},
			Required: []string{"interval", "tick"},
		}
	}
	return i.s
}

// Create creates a new Ticker block
func (i TickerInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	return &Ticker{
		id:             id,
		blockPublisher: blockCtx.BlockPublisher(),
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i TickerInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i TickerInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Ticker
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i TickerInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "id":
		return b.(*Ticker).id
	case "interval":
		return b.(*Ticker).interval
	default:
		panic(fmt.Errorf("unexpected parameter %q in Ticker", name))
	}
}

func (i TickerInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Ticker)
	switch name {
	case "interval":
		b.interval = value.(time.Duration)
	}
	return nil
}

func (i TickerInterpreter) SetBlock(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Ticker)
	switch name {
	case "tick":
		b.tick = value.(*Tick)
	}
	return nil
}
