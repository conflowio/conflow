package test

import (
	"time"

	"github.com/opsidian/basil/basil/block"

	"github.com/opsidian/basil/basil"
)

//go:generate basil generate
type TestBlock struct {
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

	Blocks []*TestBlock `basil:"block=testblock"`
}

func (t *TestBlock) ParseContext(ctx *basil.ParseContext) *basil.ParseContext {
	return ctx.New(basil.ParseContextConfig{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"testblock": TestBlockInterpreter{},
		},
	})
}
