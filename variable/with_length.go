package variable

import (
	"fmt"

	"github.com/opsidian/basil/basil"
)

// WithLength represents variables which have a length
type WithLength struct {
	value interface{}
}

// NewWithLength creates a new variable with a length
func NewWithLength(value interface{}) *WithLength {
	if unionType, ok := value.(Union); ok {
		value = unionType.Value()
	}
	return &WithLength{
		value: value,
	}
}

// Value returns with the contained value
func (w *WithLength) Value() interface{} {
	return w.value
}

// Type returns the type of the value
func (w *WithLength) Type() string {
	switch w.value.(type) {
	case string:
		return TypeString
	case basil.ID:
		return TypeIdentifier
	case []interface{}:
		return TypeArray
	case []string:
		return TypeStringArray
	case map[string]interface{}:
		return TypeMap
	default:
		panic(fmt.Sprintf("unexpected type: %T", w.value))
	}
}

// IsWithLengthType returns true if the given value is a type with a length
func IsWithLengthType(val interface{}) bool {
	switch val.(type) {
	case string, basil.ID, []interface{}, []string, map[string]interface{}:
		return true
	default:
		return false
	}
}
