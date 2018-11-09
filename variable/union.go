package variable

// Union is an interface for variables which can have multiple types
type Union interface {
	GetTypes() []string
	Value() interface{}
	Type() string
}
