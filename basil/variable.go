package basil

// VariableNode stores a variable reference. It always refers to a named block's parameter.
//go:generate counterfeiter . VariableNode
type VariableNode interface {
	IdentifiableNode
	ParamName() ID
}
