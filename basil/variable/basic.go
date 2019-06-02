package variable

import (
	"fmt"
	"time"

	"github.com/opsidian/basil/basil"
)

// Basic represents all the basic types
type Basic struct {
	value interface{}
}

// NewBasic creates a new basic variable
func NewBasic(value interface{}) *Basic {
	if unionType, ok := value.(Union); ok {
		value = unionType.Value()
	}
	return &Basic{
		value: value,
	}
}

// Value returns with the contained value
func (b *Basic) Value() interface{} {
	return b.value
}

// Type returns the type of the value
func (b *Basic) Type() string {
	switch b.value.(type) {
	case bool:
		return TypeBool
	case float64:
		return TypeFloat
	case basil.ID:
		return TypeIdentifier
	case int64:
		return TypeInteger
	case Number:
		return TypeNumber
	case string:
		return TypeString
	case time.Time:
		return TypeTime
	case time.Duration:
		return TypeTimeDuration
	default:
		panic(fmt.Sprintf("unexpected type: %T", b.value))
	}
}

// IsBasicType returns true if the given value is a number type
func IsBasicType(val interface{}) bool {
	switch val.(type) {
	case bool, float64, basil.ID, int64, Number, string, time.Duration, time.Time:
		return true
	default:
		return false
	}
}
