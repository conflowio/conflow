// Code generated by Conflow. DO NOT EDIT.

package main

import (
	"fmt"
	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/schema"
	"time"
)

func init() {
	schema.Register(&schema.Object{
		Metadata: schema.Metadata{
			Annotations: map[string]string{"block.conflow.io/type": "task"},
			ID:          "github.com/conflowio/conflow/examples/benchmark.Benchmark",
		},
		Name: "Benchmark",
		Parameters: map[string]schema.Schema{
			"counter": &schema.Integer{
				Metadata: schema.Metadata{
					Annotations: map[string]string{"block.conflow.io/eval_stage": "close"},
					ReadOnly:    true,
				},
			},
			"duration": &schema.String{
				Format: "duration-go",
			},
			"elapsed": &schema.String{
				Format: "duration-go",
			},
			"id": &schema.String{
				Metadata: schema.Metadata{
					Annotations: map[string]string{"block.conflow.io/id": "true"},
					ReadOnly:    true,
				},
				Format: "conflow.ID",
			},
			"run": &schema.Reference{
				Metadata: schema.Metadata{
					Annotations: map[string]string{"block.conflow.io/eval_stage": "init", "block.conflow.io/generated": "true"},
				},
				Nullable: true,
				Ref:      "github.com/conflowio/conflow/examples/benchmark.BenchmarkRun",
			},
		},
		Required: []string{"duration", "run"},
	})
}

// BenchmarkInterpreter is the Conflow interpreter for the Benchmark block
type BenchmarkInterpreter struct {
}

func (i BenchmarkInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/examples/benchmark.Benchmark")
	return s
}

// Create creates a new Benchmark block
func (i BenchmarkInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	b := &Benchmark{}
	b.id = id
	b.blockPublisher = blockCtx.BlockPublisher()
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i BenchmarkInterpreter) ValueParamName() conflow.ID {
	return ""
}

// ParseContext returns with the parse context for the block
func (i BenchmarkInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	var nilBlock *Benchmark
	if b, ok := conflow.Block(nilBlock).(conflow.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i BenchmarkInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	switch name {
	case "counter":
		return b.(*Benchmark).counter
	case "duration":
		return b.(*Benchmark).duration
	case "elapsed":
		return b.(*Benchmark).elapsed
	case "id":
		return b.(*Benchmark).id
	default:
		panic(fmt.Errorf("unexpected parameter %q in Benchmark", name))
	}
}

func (i BenchmarkInterpreter) SetParam(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Benchmark)
	switch name {
	case "duration":
		b.duration = value.(time.Duration)
	case "elapsed":
		b.elapsed = value.(time.Duration)
	}
	return nil
}

func (i BenchmarkInterpreter) SetBlock(block conflow.Block, name conflow.ID, value interface{}) error {
	b := block.(*Benchmark)
	switch name {
	case "run":
		b.run = value.(*BenchmarkRun)
	}
	return nil
}
