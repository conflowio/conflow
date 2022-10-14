// Code generated by Conflow. DO NOT EDIT.

package blocks

import (
	"fmt"
	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/schema"
)

func init() {
	schema.Register(&schema.Object{
		Metadata: schema.Metadata{
			Annotations: map[string]string{"block.conflow.io/type": "task"},
			ID:          "github.com/conflowio/conflow/src/blocks.Exec",
		},
		JSONPropertyNames: map[string]string{"exit_code": "exitCode"},
		Name:              "Exec",
		Parameters: map[string]schema.Schema{
			"cmd": &schema.String{},
			"dir": &schema.String{},
			"env": &schema.Array{
				Items: &schema.String{},
			},
			"exit_code": &schema.Integer{
				Metadata: schema.Metadata{
					Annotations: map[string]string{"block.conflow.io/eval_stage": "close"},
					ReadOnly:    true,
				},
			},
			"id": &schema.String{
				Metadata: schema.Metadata{
					Annotations: map[string]string{"block.conflow.io/id": "true"},
					ReadOnly:    true,
				},
				Format: "conflow.ID",
			},
			"params": &schema.Array{
				Items: &schema.String{},
			},
			"stderr": &schema.Reference{
				Metadata: schema.Metadata{
					Annotations: map[string]string{"block.conflow.io/eval_stage": "init", "block.conflow.io/generated": "true"},
				},
				Nullable: true,
				Ref:      "github.com/conflowio/conflow/src/blocks.Stream",
			},
			"stdout": &schema.Reference{
				Metadata: schema.Metadata{
					Annotations: map[string]string{"block.conflow.io/eval_stage": "init", "block.conflow.io/generated": "true"},
				},
				Nullable: true,
				Ref:      "github.com/conflowio/conflow/src/blocks.Stream",
			},
		},
		Required: []string{"cmd", "stdout", "stderr"},
	})
}

// ExecInterpreter is the Conflow interpreter for the Exec block
type ExecInterpreter struct {
}

func (i ExecInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/src/blocks.Exec")
	return s
}

// Create creates a new Exec block
func (i ExecInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := &Exec{}
	b.id = id
	b.blockPublisher = blockCtx.BlockPublisher()
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i ExecInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i ExecInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Exec
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i ExecInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "cmd":
		return b.(*Exec).cmd
	case "dir":
		return b.(*Exec).dir
	case "env":
		return b.(*Exec).env
	case "exit_code":
		return b.(*Exec).exitCode
	case "id":
		return b.(*Exec).id
	case "params":
		return b.(*Exec).params
	default:
		panic(fmt.Errorf("unexpected parameter %q in Exec", name))
	}
}

func (i ExecInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Exec)
	switch name {
	case "cmd":
		b.cmd = value.(string)
	case "dir":
		b.dir = value.(string)
	case "env":
		b.env = make([]string, len(value.([]interface{})))
		for valuek, valuev := range value.([]interface{}) {
			b.env[valuek] = valuev.(string)
		}
	case "params":
		b.params = make([]string, len(value.([]interface{})))
		for valuek, valuev := range value.([]interface{}) {
			b.params[valuek] = valuev.(string)
		}
	}
	return nil
}

func (i ExecInterpreter) SetBlock(block conflow.Block, name conflow.ID, key string, value interface{}) error {
	b := block.(*Exec)
	switch name {
	case "stderr":
		b.stderr = value.(*Stream)
	case "stdout":
		b.stdout = value.(*Stream)
	}
	return nil
}
