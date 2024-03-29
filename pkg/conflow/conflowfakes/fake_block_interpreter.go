// Code generated by counterfeiter. DO NOT EDIT.
package conflowfakes

import (
	"sync"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/schema"
)

type FakeBlockInterpreter struct {
	CreateBlockStub        func(conflow.ID, *conflow.BlockContext) conflow.Block
	createBlockMutex       sync.RWMutex
	createBlockArgsForCall []struct {
		arg1 conflow.ID
		arg2 *conflow.BlockContext
	}
	createBlockReturns struct {
		result1 conflow.Block
	}
	createBlockReturnsOnCall map[int]struct {
		result1 conflow.Block
	}
	ParamStub        func(conflow.Block, conflow.ID) interface{}
	paramMutex       sync.RWMutex
	paramArgsForCall []struct {
		arg1 conflow.Block
		arg2 conflow.ID
	}
	paramReturns struct {
		result1 interface{}
	}
	paramReturnsOnCall map[int]struct {
		result1 interface{}
	}
	ParseContextStub        func(*conflow.ParseContext) *conflow.ParseContext
	parseContextMutex       sync.RWMutex
	parseContextArgsForCall []struct {
		arg1 *conflow.ParseContext
	}
	parseContextReturns struct {
		result1 *conflow.ParseContext
	}
	parseContextReturnsOnCall map[int]struct {
		result1 *conflow.ParseContext
	}
	SchemaStub        func() schema.Schema
	schemaMutex       sync.RWMutex
	schemaArgsForCall []struct {
	}
	schemaReturns struct {
		result1 schema.Schema
	}
	schemaReturnsOnCall map[int]struct {
		result1 schema.Schema
	}
	SetBlockStub        func(conflow.Block, conflow.ID, string, interface{}) error
	setBlockMutex       sync.RWMutex
	setBlockArgsForCall []struct {
		arg1 conflow.Block
		arg2 conflow.ID
		arg3 string
		arg4 interface{}
	}
	setBlockReturns struct {
		result1 error
	}
	setBlockReturnsOnCall map[int]struct {
		result1 error
	}
	SetParamStub        func(conflow.Block, conflow.ID, interface{}) error
	setParamMutex       sync.RWMutex
	setParamArgsForCall []struct {
		arg1 conflow.Block
		arg2 conflow.ID
		arg3 interface{}
	}
	setParamReturns struct {
		result1 error
	}
	setParamReturnsOnCall map[int]struct {
		result1 error
	}
	ValueParamNameStub        func() conflow.ID
	valueParamNameMutex       sync.RWMutex
	valueParamNameArgsForCall []struct {
	}
	valueParamNameReturns struct {
		result1 conflow.ID
	}
	valueParamNameReturnsOnCall map[int]struct {
		result1 conflow.ID
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeBlockInterpreter) CreateBlock(arg1 conflow.ID, arg2 *conflow.BlockContext) conflow.Block {
	fake.createBlockMutex.Lock()
	ret, specificReturn := fake.createBlockReturnsOnCall[len(fake.createBlockArgsForCall)]
	fake.createBlockArgsForCall = append(fake.createBlockArgsForCall, struct {
		arg1 conflow.ID
		arg2 *conflow.BlockContext
	}{arg1, arg2})
	stub := fake.CreateBlockStub
	fakeReturns := fake.createBlockReturns
	fake.recordInvocation("CreateBlock", []interface{}{arg1, arg2})
	fake.createBlockMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeBlockInterpreter) CreateBlockCallCount() int {
	fake.createBlockMutex.RLock()
	defer fake.createBlockMutex.RUnlock()
	return len(fake.createBlockArgsForCall)
}

func (fake *FakeBlockInterpreter) CreateBlockCalls(stub func(conflow.ID, *conflow.BlockContext) conflow.Block) {
	fake.createBlockMutex.Lock()
	defer fake.createBlockMutex.Unlock()
	fake.CreateBlockStub = stub
}

func (fake *FakeBlockInterpreter) CreateBlockArgsForCall(i int) (conflow.ID, *conflow.BlockContext) {
	fake.createBlockMutex.RLock()
	defer fake.createBlockMutex.RUnlock()
	argsForCall := fake.createBlockArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeBlockInterpreter) CreateBlockReturns(result1 conflow.Block) {
	fake.createBlockMutex.Lock()
	defer fake.createBlockMutex.Unlock()
	fake.CreateBlockStub = nil
	fake.createBlockReturns = struct {
		result1 conflow.Block
	}{result1}
}

func (fake *FakeBlockInterpreter) CreateBlockReturnsOnCall(i int, result1 conflow.Block) {
	fake.createBlockMutex.Lock()
	defer fake.createBlockMutex.Unlock()
	fake.CreateBlockStub = nil
	if fake.createBlockReturnsOnCall == nil {
		fake.createBlockReturnsOnCall = make(map[int]struct {
			result1 conflow.Block
		})
	}
	fake.createBlockReturnsOnCall[i] = struct {
		result1 conflow.Block
	}{result1}
}

func (fake *FakeBlockInterpreter) Param(arg1 conflow.Block, arg2 conflow.ID) interface{} {
	fake.paramMutex.Lock()
	ret, specificReturn := fake.paramReturnsOnCall[len(fake.paramArgsForCall)]
	fake.paramArgsForCall = append(fake.paramArgsForCall, struct {
		arg1 conflow.Block
		arg2 conflow.ID
	}{arg1, arg2})
	stub := fake.ParamStub
	fakeReturns := fake.paramReturns
	fake.recordInvocation("Param", []interface{}{arg1, arg2})
	fake.paramMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeBlockInterpreter) ParamCallCount() int {
	fake.paramMutex.RLock()
	defer fake.paramMutex.RUnlock()
	return len(fake.paramArgsForCall)
}

func (fake *FakeBlockInterpreter) ParamCalls(stub func(conflow.Block, conflow.ID) interface{}) {
	fake.paramMutex.Lock()
	defer fake.paramMutex.Unlock()
	fake.ParamStub = stub
}

func (fake *FakeBlockInterpreter) ParamArgsForCall(i int) (conflow.Block, conflow.ID) {
	fake.paramMutex.RLock()
	defer fake.paramMutex.RUnlock()
	argsForCall := fake.paramArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeBlockInterpreter) ParamReturns(result1 interface{}) {
	fake.paramMutex.Lock()
	defer fake.paramMutex.Unlock()
	fake.ParamStub = nil
	fake.paramReturns = struct {
		result1 interface{}
	}{result1}
}

func (fake *FakeBlockInterpreter) ParamReturnsOnCall(i int, result1 interface{}) {
	fake.paramMutex.Lock()
	defer fake.paramMutex.Unlock()
	fake.ParamStub = nil
	if fake.paramReturnsOnCall == nil {
		fake.paramReturnsOnCall = make(map[int]struct {
			result1 interface{}
		})
	}
	fake.paramReturnsOnCall[i] = struct {
		result1 interface{}
	}{result1}
}

func (fake *FakeBlockInterpreter) ParseContext(arg1 *conflow.ParseContext) *conflow.ParseContext {
	fake.parseContextMutex.Lock()
	ret, specificReturn := fake.parseContextReturnsOnCall[len(fake.parseContextArgsForCall)]
	fake.parseContextArgsForCall = append(fake.parseContextArgsForCall, struct {
		arg1 *conflow.ParseContext
	}{arg1})
	stub := fake.ParseContextStub
	fakeReturns := fake.parseContextReturns
	fake.recordInvocation("ParseContext", []interface{}{arg1})
	fake.parseContextMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeBlockInterpreter) ParseContextCallCount() int {
	fake.parseContextMutex.RLock()
	defer fake.parseContextMutex.RUnlock()
	return len(fake.parseContextArgsForCall)
}

func (fake *FakeBlockInterpreter) ParseContextCalls(stub func(*conflow.ParseContext) *conflow.ParseContext) {
	fake.parseContextMutex.Lock()
	defer fake.parseContextMutex.Unlock()
	fake.ParseContextStub = stub
}

func (fake *FakeBlockInterpreter) ParseContextArgsForCall(i int) *conflow.ParseContext {
	fake.parseContextMutex.RLock()
	defer fake.parseContextMutex.RUnlock()
	argsForCall := fake.parseContextArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeBlockInterpreter) ParseContextReturns(result1 *conflow.ParseContext) {
	fake.parseContextMutex.Lock()
	defer fake.parseContextMutex.Unlock()
	fake.ParseContextStub = nil
	fake.parseContextReturns = struct {
		result1 *conflow.ParseContext
	}{result1}
}

func (fake *FakeBlockInterpreter) ParseContextReturnsOnCall(i int, result1 *conflow.ParseContext) {
	fake.parseContextMutex.Lock()
	defer fake.parseContextMutex.Unlock()
	fake.ParseContextStub = nil
	if fake.parseContextReturnsOnCall == nil {
		fake.parseContextReturnsOnCall = make(map[int]struct {
			result1 *conflow.ParseContext
		})
	}
	fake.parseContextReturnsOnCall[i] = struct {
		result1 *conflow.ParseContext
	}{result1}
}

func (fake *FakeBlockInterpreter) Schema() schema.Schema {
	fake.schemaMutex.Lock()
	ret, specificReturn := fake.schemaReturnsOnCall[len(fake.schemaArgsForCall)]
	fake.schemaArgsForCall = append(fake.schemaArgsForCall, struct {
	}{})
	stub := fake.SchemaStub
	fakeReturns := fake.schemaReturns
	fake.recordInvocation("Schema", []interface{}{})
	fake.schemaMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeBlockInterpreter) SchemaCallCount() int {
	fake.schemaMutex.RLock()
	defer fake.schemaMutex.RUnlock()
	return len(fake.schemaArgsForCall)
}

func (fake *FakeBlockInterpreter) SchemaCalls(stub func() schema.Schema) {
	fake.schemaMutex.Lock()
	defer fake.schemaMutex.Unlock()
	fake.SchemaStub = stub
}

func (fake *FakeBlockInterpreter) SchemaReturns(result1 schema.Schema) {
	fake.schemaMutex.Lock()
	defer fake.schemaMutex.Unlock()
	fake.SchemaStub = nil
	fake.schemaReturns = struct {
		result1 schema.Schema
	}{result1}
}

func (fake *FakeBlockInterpreter) SchemaReturnsOnCall(i int, result1 schema.Schema) {
	fake.schemaMutex.Lock()
	defer fake.schemaMutex.Unlock()
	fake.SchemaStub = nil
	if fake.schemaReturnsOnCall == nil {
		fake.schemaReturnsOnCall = make(map[int]struct {
			result1 schema.Schema
		})
	}
	fake.schemaReturnsOnCall[i] = struct {
		result1 schema.Schema
	}{result1}
}

func (fake *FakeBlockInterpreter) SetBlock(arg1 conflow.Block, arg2 conflow.ID, arg3 string, arg4 interface{}) error {
	fake.setBlockMutex.Lock()
	ret, specificReturn := fake.setBlockReturnsOnCall[len(fake.setBlockArgsForCall)]
	fake.setBlockArgsForCall = append(fake.setBlockArgsForCall, struct {
		arg1 conflow.Block
		arg2 conflow.ID
		arg3 string
		arg4 interface{}
	}{arg1, arg2, arg3, arg4})
	stub := fake.SetBlockStub
	fakeReturns := fake.setBlockReturns
	fake.recordInvocation("SetBlock", []interface{}{arg1, arg2, arg3, arg4})
	fake.setBlockMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeBlockInterpreter) SetBlockCallCount() int {
	fake.setBlockMutex.RLock()
	defer fake.setBlockMutex.RUnlock()
	return len(fake.setBlockArgsForCall)
}

func (fake *FakeBlockInterpreter) SetBlockCalls(stub func(conflow.Block, conflow.ID, string, interface{}) error) {
	fake.setBlockMutex.Lock()
	defer fake.setBlockMutex.Unlock()
	fake.SetBlockStub = stub
}

func (fake *FakeBlockInterpreter) SetBlockArgsForCall(i int) (conflow.Block, conflow.ID, string, interface{}) {
	fake.setBlockMutex.RLock()
	defer fake.setBlockMutex.RUnlock()
	argsForCall := fake.setBlockArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeBlockInterpreter) SetBlockReturns(result1 error) {
	fake.setBlockMutex.Lock()
	defer fake.setBlockMutex.Unlock()
	fake.SetBlockStub = nil
	fake.setBlockReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBlockInterpreter) SetBlockReturnsOnCall(i int, result1 error) {
	fake.setBlockMutex.Lock()
	defer fake.setBlockMutex.Unlock()
	fake.SetBlockStub = nil
	if fake.setBlockReturnsOnCall == nil {
		fake.setBlockReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setBlockReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeBlockInterpreter) SetParam(arg1 conflow.Block, arg2 conflow.ID, arg3 interface{}) error {
	fake.setParamMutex.Lock()
	ret, specificReturn := fake.setParamReturnsOnCall[len(fake.setParamArgsForCall)]
	fake.setParamArgsForCall = append(fake.setParamArgsForCall, struct {
		arg1 conflow.Block
		arg2 conflow.ID
		arg3 interface{}
	}{arg1, arg2, arg3})
	stub := fake.SetParamStub
	fakeReturns := fake.setParamReturns
	fake.recordInvocation("SetParam", []interface{}{arg1, arg2, arg3})
	fake.setParamMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeBlockInterpreter) SetParamCallCount() int {
	fake.setParamMutex.RLock()
	defer fake.setParamMutex.RUnlock()
	return len(fake.setParamArgsForCall)
}

func (fake *FakeBlockInterpreter) SetParamCalls(stub func(conflow.Block, conflow.ID, interface{}) error) {
	fake.setParamMutex.Lock()
	defer fake.setParamMutex.Unlock()
	fake.SetParamStub = stub
}

func (fake *FakeBlockInterpreter) SetParamArgsForCall(i int) (conflow.Block, conflow.ID, interface{}) {
	fake.setParamMutex.RLock()
	defer fake.setParamMutex.RUnlock()
	argsForCall := fake.setParamArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeBlockInterpreter) SetParamReturns(result1 error) {
	fake.setParamMutex.Lock()
	defer fake.setParamMutex.Unlock()
	fake.SetParamStub = nil
	fake.setParamReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBlockInterpreter) SetParamReturnsOnCall(i int, result1 error) {
	fake.setParamMutex.Lock()
	defer fake.setParamMutex.Unlock()
	fake.SetParamStub = nil
	if fake.setParamReturnsOnCall == nil {
		fake.setParamReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setParamReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeBlockInterpreter) ValueParamName() conflow.ID {
	fake.valueParamNameMutex.Lock()
	ret, specificReturn := fake.valueParamNameReturnsOnCall[len(fake.valueParamNameArgsForCall)]
	fake.valueParamNameArgsForCall = append(fake.valueParamNameArgsForCall, struct {
	}{})
	stub := fake.ValueParamNameStub
	fakeReturns := fake.valueParamNameReturns
	fake.recordInvocation("ValueParamName", []interface{}{})
	fake.valueParamNameMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeBlockInterpreter) ValueParamNameCallCount() int {
	fake.valueParamNameMutex.RLock()
	defer fake.valueParamNameMutex.RUnlock()
	return len(fake.valueParamNameArgsForCall)
}

func (fake *FakeBlockInterpreter) ValueParamNameCalls(stub func() conflow.ID) {
	fake.valueParamNameMutex.Lock()
	defer fake.valueParamNameMutex.Unlock()
	fake.ValueParamNameStub = stub
}

func (fake *FakeBlockInterpreter) ValueParamNameReturns(result1 conflow.ID) {
	fake.valueParamNameMutex.Lock()
	defer fake.valueParamNameMutex.Unlock()
	fake.ValueParamNameStub = nil
	fake.valueParamNameReturns = struct {
		result1 conflow.ID
	}{result1}
}

func (fake *FakeBlockInterpreter) ValueParamNameReturnsOnCall(i int, result1 conflow.ID) {
	fake.valueParamNameMutex.Lock()
	defer fake.valueParamNameMutex.Unlock()
	fake.ValueParamNameStub = nil
	if fake.valueParamNameReturnsOnCall == nil {
		fake.valueParamNameReturnsOnCall = make(map[int]struct {
			result1 conflow.ID
		})
	}
	fake.valueParamNameReturnsOnCall[i] = struct {
		result1 conflow.ID
	}{result1}
}

func (fake *FakeBlockInterpreter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createBlockMutex.RLock()
	defer fake.createBlockMutex.RUnlock()
	fake.paramMutex.RLock()
	defer fake.paramMutex.RUnlock()
	fake.parseContextMutex.RLock()
	defer fake.parseContextMutex.RUnlock()
	fake.schemaMutex.RLock()
	defer fake.schemaMutex.RUnlock()
	fake.setBlockMutex.RLock()
	defer fake.setBlockMutex.RUnlock()
	fake.setParamMutex.RLock()
	defer fake.setParamMutex.RUnlock()
	fake.valueParamNameMutex.RLock()
	defer fake.valueParamNameMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeBlockInterpreter) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ conflow.BlockInterpreter = new(FakeBlockInterpreter)
