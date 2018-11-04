package test

import "time"

//go:generate basil generate TestBlock
type TestBlock struct {
	IDField           string      `basil:"id"`
	Value             interface{} `basil:"value"`
	FieldString       string
	FieldInt          int64
	FieldFloat        float64
	FieldBool         bool
	FieldArray        []interface{}
	FieldMap          map[string]interface{}
	FieldTimeDuration time.Duration
	FieldCustomName   string `basil:"name=custom_field"`

	BlockFactories []*TestBlockFactory `basil:"factory"`
	Blocks         []*TestBlock        `basil:"block"`
}

func (t *TestBlock) ID() string {
	return t.IDField
}

func (t *TestBlock) Type() string {
	return "testblock"
}

func (t *TestBlock) Context(ctx interface{}) interface{} {
	return ctx
}
