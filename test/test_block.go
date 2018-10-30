package test

import "time"

//go:generate ocl generate TestBlock
type TestBlock struct {
	IDField           string `ocl:"id"`
	Value             int64  `ocl:"value"`
	FieldString       string
	FieldInt          int64
	FieldFloat        float64
	FieldBool         bool
	FieldArray        []interface{}
	FieldMap          map[string]interface{}
	FieldTimeDuration time.Duration
	FieldCustomName   string `ocl:"name=custom_field"`

	BlockFactories []*TestBlockFactory `ocl:"factory"`
	Blocks         []*TestBlock        `ocl:"block"`
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
