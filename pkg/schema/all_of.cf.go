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
			ID: "github.com/conflowio/conflow/pkg/schema.AllOf",
		},
		FieldNames:     map[string]string{"allOf": "Schemas", "const": "Const", "default": "Default", "enum": "Enum", "nullable": "Nullable"},
		ParameterNames: map[string]string{"allOf": "schema"},
		Properties: map[string]Schema{
			"allOf": &Array{
				Items: &Reference{
					Ref: "github.com/conflowio/conflow/pkg/schema.Schema",
				},
				MinItems: 1,
			},
			"const":   &Any{},
			"default": &Any{},
			"enum": &Array{
				Items: &Any{},
			},
			"nullable": &Boolean{},
		},
		Required: []string{"schema"},
	})
}

// NewAllOfWithDefaults creates a new AllOf instance with default values
func NewAllOfWithDefaults() *AllOf {
	b := &AllOf{}
	return b
}
