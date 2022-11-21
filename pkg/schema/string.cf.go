// Code generated by Conflow. DO NOT EDIT.

package schema

import (
	"github.com/conflowio/conflow/pkg/conflow/annotations"
)

func init() {
	Register(&Object{
		Metadata: Metadata{
			Annotations: map[string]string{
				annotations.Type: "configuration",
			},
			ID: "github.com/conflowio/conflow/pkg/schema.String",
		},
		FieldNames:     map[string]string{"$id": "ID", "const": "Const", "default": "Default", "deprecated": "Deprecated", "description": "Description", "enum": "Enum", "examples": "Examples", "format": "Format", "maxLength": "MaxLength", "minLength": "MinLength", "nullable": "Nullable", "pattern": "Pattern", "readOnly": "ReadOnly", "title": "Title", "writeOnly": "WriteOnly", "x-annotations": "Annotations"},
		ParameterNames: map[string]string{"$id": "id", "maxLength": "max_length", "minLength": "min_length", "readOnly": "read_only", "writeOnly": "write_only", "x-annotations": "annotations"},
		Properties: map[string]Schema{
			"$id": &String{},
			"const": &String{
				Nullable: true,
			},
			"default": &String{
				Nullable: true,
			},
			"deprecated":  &Boolean{},
			"description": &String{},
			"enum": &Array{
				Items: &String{},
			},
			"examples": &Array{
				Items: &Any{},
			},
			"format": &String{},
			"maxLength": &Integer{
				Nullable: true,
			},
			"minLength": &Integer{
				Minimum: Pointer(int64(0)),
			},
			"nullable": &Boolean{},
			"pattern": &String{
				Format:   "regex",
				Nullable: true,
			},
			"readOnly":  &Boolean{},
			"title":     &String{},
			"writeOnly": &Boolean{},
			"x-annotations": &Map{
				AdditionalProperties: &String{},
			},
		},
	})
}

// NewStringWithDefaults creates a new String instance with default values
func NewStringWithDefaults() *String {
	b := &String{}
	return b
}
