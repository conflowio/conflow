// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

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
