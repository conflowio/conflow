package generator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/variable"
)

// Field contains a metadata for a block field
type Field struct {
	Name        string
	ParamName   string
	Type        string
	Stage       string
	IsRequired  bool
	IsParam     bool
	IsID        bool
	IsValue     bool
	IsReference bool
	IsBlock     bool
	IsNode      bool
	IsOutput    bool
	IsChannel   bool
}

// Validate validates the field tags
func (f *Field) Validate() error {
	_, validType := variable.Types[f.Type]
	if !validType && !f.IsBlock && !f.IsNode {
		return fmt.Errorf("invalid field type on field %q, use valid type or use ignore tag", f.Name)
	}

	if f.hasMultipleKinds() {
		return fmt.Errorf("field %q must only have one tag of: id, value, block or node", f.Name)
	}

	if !basil.IDRegExp.MatchString(f.ParamName) {
		return fmt.Errorf("\"name\" tag is invalid on field %q, it must be a valid identifier", f.Name)
	}

	if f.IsID && f.Type != variable.TypeIdentifier {
		return fmt.Errorf("field %q must be defined as %s", f.Name, variable.TypeIdentifier)
	}

	if f.IsReference && !f.IsID {
		return errors.New("the \"reference\" tag can only be set on the id field")
	}

	if f.IsBlock || f.IsNode {
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
	if f.IsNode {
		typeCnt++
	}
	return typeCnt > 1
}
