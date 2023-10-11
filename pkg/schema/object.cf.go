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
			ID: "github.com/conflowio/conflow/pkg/schema.Object",
		},
		FieldNames:     map[string]string{"$id": "ID", "const": "Const", "default": "Default", "dependentRequired": "DependentRequired", "deprecated": "Deprecated", "description": "Description", "enum": "Enum", "examples": "Examples", "maxProperties": "MaxProperties", "minProperties": "MinProperties", "nullable": "Nullable", "properties": "Properties", "readOnly": "ReadOnly", "required": "Required", "title": "Title", "writeOnly": "WriteOnly", "x-annotations": "Annotations", "x-conflow-fields": "FieldNames", "x-conflow-parameters": "ParameterNames"},
		ParameterNames: map[string]string{"$id": "id", "dependentRequired": "dependent_required", "maxProperties": "max_properties", "minProperties": "min_properties", "properties": "property", "readOnly": "read_only", "writeOnly": "write_only", "x-annotations": "annotations", "x-conflow-fields": "field_names", "x-conflow-parameters": "parameter_names"},
		Properties: map[string]Schema{
			"$id": &String{},
			"const": &Map{
				AdditionalProperties: &Any{},
			},
			"default": &Map{
				AdditionalProperties: &Any{},
			},
			"dependentRequired": &Map{
				AdditionalProperties: &Array{
					Items: &String{},
				},
			},
			"deprecated":  &Boolean{},
			"description": &String{},
			"enum": &Array{
				Items: &Map{
					AdditionalProperties: &Any{},
				},
			},
			"examples": &Array{
				Items: &Any{},
			},
			"maxProperties": &Integer{
				Nullable: true,
			},
			"minProperties": &Integer{},
			"nullable":      &Boolean{},
			"properties": &Map{
				AdditionalProperties: &Reference{
					Ref: "github.com/conflowio/conflow/pkg/schema.Schema",
				},
			},
			"readOnly": &Boolean{},
			"required": &Array{
				Metadata: Metadata{
					Description: "It will contain the required parameter names",
				},
				Items: &String{},
			},
			"title":     &String{},
			"writeOnly": &Boolean{},
			"x-annotations": &Map{
				AdditionalProperties: &String{},
			},
			"x-conflow-fields": &Map{
				Metadata: Metadata{
					Description: "It will contain the json property name -> field name mapping, if they are different",
				},
				AdditionalProperties: &String{},
			},
			"x-conflow-parameters": &Map{
				Metadata: Metadata{
					Description: "It will contain the json property name -> parameter name mapping, if they are different",
				},
				AdditionalProperties: &String{},
			},
		},
	})
}

// NewObjectWithDefaults creates a new Object instance with default values
func NewObjectWithDefaults() *Object {
	b := &Object{}
	b.Properties = map[string]Schema{}
	return b
}