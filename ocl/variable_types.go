package ocl

// Variable types
const (
	TypeInt          = "int64"
	TypeFloat        = "float64"
	TypeString       = "string"
	TypeBool         = "bool"
	TypeTimeDuration = "time.Duration"
	TypeArray        = "[]interface{}"
	TypeMap          = "map[string]interface{}"
)

// Valid variable types with descriptions
var VariableTypes = map[string]string{
	TypeInt:          "integer",
	TypeFloat:        "float",
	TypeString:       "string",
	TypeBool:         "boolean",
	TypeTimeDuration: "time duration",
	TypeArray:        "array",
	TypeMap:          "map",
}
