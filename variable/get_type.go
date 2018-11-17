package variable

import (
	"time"

	"github.com/opsidian/basil/basil"
)

// GetType returns with the type of the given value
func GetType(value interface{}) string {
	switch value.(type) {
	case []interface{}:
		return TypeArray
	case *Basic:
		return TypeBasic
	case bool:
		return TypeBool
	case float64:
		return TypeFloat
	case basil.ID:
		return TypeIdentifier
	case int64:
		return TypeInteger
	case map[string]interface{}:
		return TypeMap
	case *Number:
		return TypeNumber
	case string:
		return TypeString
	case []string:
		return TypeStringArray
	case time.Duration:
		return TypeTimeDuration
	case *WithLength:
		return TypeWithLength
	default:
		return TypeUnknown
	}
}
