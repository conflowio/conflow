package block

import (
	"errors"
	"fmt"
	"strings"

	"github.com/opsidian/ocl/identifier"
	"github.com/opsidian/ocl/ocl"
)

// Field contains a metadata for a block field
type Field struct {
	Name        string
	ParamName   string
	Type        string
	Required    bool
	Stage       string
	IsID        bool
	IsValue     bool
	IsReference bool
	IsBlock     bool
	IsFactory   bool
}

// Validate validates the field tags
func (f *Field) Validate() error {
	_, validType := ocl.VariableTypes[f.Type]
	if !validType && !f.IsBlock && !f.IsFactory {
		return fmt.Errorf("invalid field type on field %q, use valid type or use ignore tag", f.Name)
	}

	if f.hasMultipleKinds() {
		return fmt.Errorf("field %q must only have one tag of: id, value, block or factory", f.Name)
	}

	if !identifier.RegExp.MatchString(f.ParamName) {
		return fmt.Errorf("\"name\" tag is invalid on field %q, it must be a valid identifier", f.Name)
	}

	if f.IsID && f.Type != "string" {
		return fmt.Errorf("field %q must be defined as string", f.Name)
	}

	if f.IsReference && !f.IsID {
		return errors.New("the \"reference\" tag can only be set on the id field")
	}

	if f.IsBlock || f.IsFactory {
		if !strings.HasPrefix(f.Type, "[]") {
			return fmt.Errorf("field %q must be an array", f.Name)
		}
	}

	if f.Stage == "" {
		return fmt.Errorf("\"stage\" can not be empty on field %q", f.Name)
	}

	return nil
}

func (f *Field) hasMultipleKinds() bool {
	typeCnt := 0
	if f.IsID {
		typeCnt++
	}
	if f.IsValue {
		typeCnt++
	}
	if f.IsBlock {
		typeCnt++
	}
	if f.IsFactory {
		typeCnt++
	}
	return typeCnt > 1
}
