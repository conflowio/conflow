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
			ID: "github.com/conflowio/conflow/pkg/schema.Boolean",
		},
		FieldNames:     map[string]string{"$id": "ID", "const": "Const", "default": "Default", "deprecated": "Deprecated", "description": "Description", "enum": "Enum", "examples": "Examples", "nullable": "Nullable", "readOnly": "ReadOnly", "title": "Title", "writeOnly": "WriteOnly", "x-annotations": "Annotations"},
		ParameterNames: map[string]string{"$id": "id", "readOnly": "read_only", "writeOnly": "write_only", "x-annotations": "annotations"},
		Properties: map[string]Schema{
			"$id": &String{},
			"const": &Boolean{
				Nullable: true,
			},
			"default": &Boolean{
				Nullable: true,
			},
			"deprecated":  &Boolean{},
			"description": &String{},
			"enum": &Array{
				Items: &Boolean{},
			},
			"examples": &Array{
				Items: &Any{},
			},
			"nullable":  &Boolean{},
			"readOnly":  &Boolean{},
			"title":     &String{},
			"writeOnly": &Boolean{},
			"x-annotations": &Map{
				AdditionalProperties: &String{},
			},
		},
	})
}

// NewBooleanWithDefaults creates a new Boolean instance with default values
func NewBooleanWithDefaults() *Boolean {
	b := &Boolean{}
	return b
}
