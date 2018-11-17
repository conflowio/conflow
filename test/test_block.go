package test

import (
	"time"

	"github.com/opsidian/basil/block"

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

	BlockNodes []basil.BlockNode `basil:"node"`
	Blocks     []*TestBlock      `basil:"block"`
}

func (t *TestBlock) ID() basil.ID {
	return t.IDField
}

func (t *TestBlock) Type() string {
	return "testblock"
}

func (t *TestBlock) Context(ctx interface{}) interface{} {
	return ctx
}

func (t *TestBlock) BlockRegistry() block.Registry {
	return block.Registry{
		"testblock": TestBlockInterpreter{},
	}
}
