// Code generated by Basil. DO NOT EDIT.

package blocks

import (
	"fmt"
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/schema"
)

// ExecInterpreter is the basil interpreter for the Exec block
type ExecInterpreter struct {
	s schema.Schema
}

func (i ExecInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Object{
			Properties: map[string]schema.Schema{
				"cmd": &schema.String{},
				"dir": &schema.String{},
				"env": &schema.Array{
					Items: &schema.String{},
				},
				"exit_code": &schema.Integer{
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
					Format: "basil.ID",
				},
				"params": &schema.Array{
					Items: &schema.String{},
				},
				"stderr": &schema.Reference{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"eval_stage": "init", "generated": "true"},
						Pointer:     true,
					},
					Ref: "http://basil.schema/github.com/opsidian/basil/blocks/Stream",
				},
				"stdout": &schema.Reference{
					Metadata: schema.Metadata{
						Annotations: map[string]string{"eval_stage": "init", "generated": "true"},
						Pointer:     true,
					},
					Ref: "http://basil.schema/github.com/opsidian/basil/blocks/Stream",
				},
			},
			Required:         []string{"cmd", "stdout", "stderr"},
			StructProperties: map[string]string{"exit_code": "exitCode"},
		}
	}
	return i.s
}

// Create creates a new Exec block
func (i ExecInterpreter) CreateBlock(id basil.ID) basil.Block {
	return &Exec{
		id: id,
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i ExecInterpreter) ValueParamName() basil.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i ExecInterpreter) ParseContext(ctx *basil.ParseContext) *basil.ParseContext {
	var nilBlock *Exec
	if b, ok := basil.Block(nilBlock).(basil.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i ExecInterpreter) Param(b basil.Block, name basil.ID) interface{} {
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

func (i ExecInterpreter) SetParam(block basil.Block, name basil.ID, value interface{}) error {
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

func (i ExecInterpreter) SetBlock(block basil.Block, name basil.ID, value interface{}) error {
	b := block.(*Exec)
	switch name {
	case "stderr":
		b.stderr = value.(*Stream)
	case "stdout":
		b.stdout = value.(*Stream)
	}
	return nil
}
