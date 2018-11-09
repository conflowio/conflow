package variable

// Provider is an interface for looking up variables
//go:generate counterfeiter . Provider
type Provider interface {
	GetVar(name ID) (interface{}, bool)
	LookupVar(lookup LookUp) (interface{}, error)
}

// ProviderAware defines a function to access a variable provider
type ProviderAware interface {
	VariableProvider() Provider
}

// LookUp is a variable lookup function
type LookUp func(provider Provider) (interface{}, error)
