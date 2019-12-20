package block

import "github.com/opsidian/basil/basil"

func newMainInterpreter(
	interpreter basil.BlockInterpreter, moduleParams map[basil.ID]basil.ParameterDescriptor,
) basil.BlockInterpreter {
	params := make(map[basil.ID]basil.ParameterDescriptor, len(interpreter.Params())+len(moduleParams))
	for k, v := range interpreter.Params() {
		params[k] = v
	}
	for k, v := range moduleParams {
		params[k] = v
	}

	return &mainInterpreter{
		BlockInterpreter: interpreter,
		params:           params,
	}
}

type mainInterpreter struct {
	basil.BlockInterpreter
	params map[basil.ID]basil.ParameterDescriptor
}

func (m *mainInterpreter) Params() map[basil.ID]basil.ParameterDescriptor {
	return m.params
}
