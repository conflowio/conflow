// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generator_test

import (
	"testing"

	"github.com/conflowio/conflow/pkg/test"
)

func TestSetParamBindArrayField(t *testing.T) {
	block := test.NewBlockWithDefaults()
	interpreter := test.BlockInterpreter{}
	upstream := []interface{}{"shared"}

	if err := interpreter.SetParam(block, "field_array", upstream); err != nil {
		t.Fatalf("SetParam: %v", err)
	}

	if block.FieldArray == nil {
		t.Fatal("FieldArray is nil")
	}
	if &block.FieldArray[0] == &upstream[0] {
		t.Fatal("FieldArray aliases upstream slice backing array")
	}
	if block.FieldArray[0] != "shared" {
		t.Fatalf("FieldArray[0] = %q, want %q", block.FieldArray[0], "shared")
	}

	upstream[0] = "mutated"
	if block.FieldArray[0] != "shared" {
		t.Fatalf("FieldArray mutated after upstream change")
	}
}

func TestSetParamBindMapField(t *testing.T) {
	block := test.NewBlockWithDefaults()
	interpreter := test.BlockInterpreter{}
	upstream := map[string]interface{}{"key": "shared"}

	if err := interpreter.SetParam(block, "field_map", upstream); err != nil {
		t.Fatalf("SetParam: %v", err)
	}

	if block.FieldMap == nil {
		t.Fatal("FieldMap is nil")
	}
	if _, ok := block.FieldMap["key"]; !ok {
		t.Fatal("FieldMap missing key")
	}
	if block.FieldMap["key"] != "shared" {
		t.Fatalf(`FieldMap["key"] = %q, want %q`, block.FieldMap["key"], "shared")
	}

	upstream["key"] = "mutated"
	if block.FieldMap["key"] != "shared" {
		t.Fatalf("FieldMap mutated after upstream change")
	}
}

func TestSetParamBindValueFieldCollection(t *testing.T) {
	block := test.NewBlockWithDefaults()
	interpreter := test.BlockInterpreter{}
	upstream := []interface{}{"value-item"}

	if err := interpreter.SetParam(block, "value", upstream); err != nil {
		t.Fatalf("SetParam: %v", err)
	}

	slice, ok := block.Value.([]interface{})
	if !ok {
		t.Fatalf("Value type = %T, want []interface{}", block.Value)
	}
	if &slice[0] == &upstream[0] {
		t.Fatal("Value slice aliases upstream backing array")
	}

	upstream[0] = "mutated"
	if slice[0] != "value-item" {
		t.Fatalf("Value slice mutated after upstream change")
	}
}

func TestSetParamBindScalarUnchanged(t *testing.T) {
	block := test.NewBlockWithDefaults()
	interpreter := test.BlockInterpreter{}

	if err := interpreter.SetParam(block, "field_string", "hello"); err != nil {
		t.Fatalf("SetParam: %v", err)
	}
	if block.FieldString != "hello" {
		t.Fatalf("FieldString = %q, want %q", block.FieldString, "hello")
	}
}
