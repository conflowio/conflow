package variable

// Variable contains a variable's name, type and value
type Variable struct {
	name    string
	varType string
	value   interface{}
}

// NewVariable creates a new variable
func NewVariable(name string, varType string, value interface{}) *Variable {
	return &Variable{
		name:    name,
		varType: varType,
		value:   value,
	}
}

// Name returns with the variable's name
func (v *Variable) Name() string {
	return v.name
}

// Type returns with the variable's type
func (v *Variable) Type() string {
	return v.varType
}

// Value returns with the variable's value
func (v *Variable) Value() interface{} {
	return v.value
}

// SetValue sets the variable's value
func (v *Variable) SetValue(value interface{}) {
	v.value = value
}
