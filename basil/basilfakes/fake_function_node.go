// Code generated by counterfeiter. DO NOT EDIT.
package basilfakes

import (
	"sync"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/parsley"
)

type FakeFunctionNode struct {
	ArgumentNodesStub        func() []parsley.Node
	argumentNodesMutex       sync.RWMutex
	argumentNodesArgsForCall []struct {
	}
	argumentNodesReturns struct {
		result1 []parsley.Node
	}
	argumentNodesReturnsOnCall map[int]struct {
		result1 []parsley.Node
	}
	NameStub        func() basil.ID
	nameMutex       sync.RWMutex
	nameArgsForCall []struct {
	}
	nameReturns struct {
		result1 basil.ID
	}
	nameReturnsOnCall map[int]struct {
		result1 basil.ID
	}
	PosStub        func() parsley.Pos
	posMutex       sync.RWMutex
	posArgsForCall []struct {
	}
	posReturns struct {
		result1 parsley.Pos
	}
	posReturnsOnCall map[int]struct {
		result1 parsley.Pos
	}
	ReaderPosStub        func() parsley.Pos
	readerPosMutex       sync.RWMutex
	readerPosArgsForCall []struct {
	}
	readerPosReturns struct {
		result1 parsley.Pos
	}
	readerPosReturnsOnCall map[int]struct {
		result1 parsley.Pos
	}
	StaticCheckStub        func(interface{}) parsley.Error
	staticCheckMutex       sync.RWMutex
	staticCheckArgsForCall []struct {
		arg1 interface{}
	}
	staticCheckReturns struct {
		result1 parsley.Error
	}
	staticCheckReturnsOnCall map[int]struct {
		result1 parsley.Error
	}
	TokenStub        func() string
	tokenMutex       sync.RWMutex
	tokenArgsForCall []struct {
	}
	tokenReturns struct {
		result1 string
	}
	tokenReturnsOnCall map[int]struct {
		result1 string
	}
	TypeStub        func() string
	typeMutex       sync.RWMutex
	typeArgsForCall []struct {
	}
	typeReturns struct {
		result1 string
	}
	typeReturnsOnCall map[int]struct {
		result1 string
	}
	ValueStub        func(interface{}) (interface{}, parsley.Error)
	valueMutex       sync.RWMutex
	valueArgsForCall []struct {
		arg1 interface{}
	}
	valueReturns struct {
		result1 interface{}
		result2 parsley.Error
	}
	valueReturnsOnCall map[int]struct {
		result1 interface{}
		result2 parsley.Error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeFunctionNode) ArgumentNodes() []parsley.Node {
	fake.argumentNodesMutex.Lock()
	ret, specificReturn := fake.argumentNodesReturnsOnCall[len(fake.argumentNodesArgsForCall)]
	fake.argumentNodesArgsForCall = append(fake.argumentNodesArgsForCall, struct {
	}{})
	fake.recordInvocation("ArgumentNodes", []interface{}{})
	fake.argumentNodesMutex.Unlock()
	if fake.ArgumentNodesStub != nil {
		return fake.ArgumentNodesStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.argumentNodesReturns
	return fakeReturns.result1
}

func (fake *FakeFunctionNode) ArgumentNodesCallCount() int {
	fake.argumentNodesMutex.RLock()
	defer fake.argumentNodesMutex.RUnlock()
	return len(fake.argumentNodesArgsForCall)
}

func (fake *FakeFunctionNode) ArgumentNodesCalls(stub func() []parsley.Node) {
	fake.argumentNodesMutex.Lock()
	defer fake.argumentNodesMutex.Unlock()
	fake.ArgumentNodesStub = stub
}

func (fake *FakeFunctionNode) ArgumentNodesReturns(result1 []parsley.Node) {
	fake.argumentNodesMutex.Lock()
	defer fake.argumentNodesMutex.Unlock()
	fake.ArgumentNodesStub = nil
	fake.argumentNodesReturns = struct {
		result1 []parsley.Node
	}{result1}
}

func (fake *FakeFunctionNode) ArgumentNodesReturnsOnCall(i int, result1 []parsley.Node) {
	fake.argumentNodesMutex.Lock()
	defer fake.argumentNodesMutex.Unlock()
	fake.ArgumentNodesStub = nil
	if fake.argumentNodesReturnsOnCall == nil {
		fake.argumentNodesReturnsOnCall = make(map[int]struct {
			result1 []parsley.Node
		})
	}
	fake.argumentNodesReturnsOnCall[i] = struct {
		result1 []parsley.Node
	}{result1}
}

func (fake *FakeFunctionNode) Name() basil.ID {
	fake.nameMutex.Lock()
	ret, specificReturn := fake.nameReturnsOnCall[len(fake.nameArgsForCall)]
	fake.nameArgsForCall = append(fake.nameArgsForCall, struct {
	}{})
	fake.recordInvocation("Name", []interface{}{})
	fake.nameMutex.Unlock()
	if fake.NameStub != nil {
		return fake.NameStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.nameReturns
	return fakeReturns.result1
}

func (fake *FakeFunctionNode) NameCallCount() int {
	fake.nameMutex.RLock()
	defer fake.nameMutex.RUnlock()
	return len(fake.nameArgsForCall)
}

func (fake *FakeFunctionNode) NameCalls(stub func() basil.ID) {
	fake.nameMutex.Lock()
	defer fake.nameMutex.Unlock()
	fake.NameStub = stub
}

func (fake *FakeFunctionNode) NameReturns(result1 basil.ID) {
	fake.nameMutex.Lock()
	defer fake.nameMutex.Unlock()
	fake.NameStub = nil
	fake.nameReturns = struct {
		result1 basil.ID
	}{result1}
}

func (fake *FakeFunctionNode) NameReturnsOnCall(i int, result1 basil.ID) {
	fake.nameMutex.Lock()
	defer fake.nameMutex.Unlock()
	fake.NameStub = nil
	if fake.nameReturnsOnCall == nil {
		fake.nameReturnsOnCall = make(map[int]struct {
			result1 basil.ID
		})
	}
	fake.nameReturnsOnCall[i] = struct {
		result1 basil.ID
	}{result1}
}

func (fake *FakeFunctionNode) Pos() parsley.Pos {
	fake.posMutex.Lock()
	ret, specificReturn := fake.posReturnsOnCall[len(fake.posArgsForCall)]
	fake.posArgsForCall = append(fake.posArgsForCall, struct {
	}{})
	fake.recordInvocation("Pos", []interface{}{})
	fake.posMutex.Unlock()
	if fake.PosStub != nil {
		return fake.PosStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.posReturns
	return fakeReturns.result1
}

func (fake *FakeFunctionNode) PosCallCount() int {
	fake.posMutex.RLock()
	defer fake.posMutex.RUnlock()
	return len(fake.posArgsForCall)
}

func (fake *FakeFunctionNode) PosCalls(stub func() parsley.Pos) {
	fake.posMutex.Lock()
	defer fake.posMutex.Unlock()
	fake.PosStub = stub
}

func (fake *FakeFunctionNode) PosReturns(result1 parsley.Pos) {
	fake.posMutex.Lock()
	defer fake.posMutex.Unlock()
	fake.PosStub = nil
	fake.posReturns = struct {
		result1 parsley.Pos
	}{result1}
}

func (fake *FakeFunctionNode) PosReturnsOnCall(i int, result1 parsley.Pos) {
	fake.posMutex.Lock()
	defer fake.posMutex.Unlock()
	fake.PosStub = nil
	if fake.posReturnsOnCall == nil {
		fake.posReturnsOnCall = make(map[int]struct {
			result1 parsley.Pos
		})
	}
	fake.posReturnsOnCall[i] = struct {
		result1 parsley.Pos
	}{result1}
}

func (fake *FakeFunctionNode) ReaderPos() parsley.Pos {
	fake.readerPosMutex.Lock()
	ret, specificReturn := fake.readerPosReturnsOnCall[len(fake.readerPosArgsForCall)]
	fake.readerPosArgsForCall = append(fake.readerPosArgsForCall, struct {
	}{})
	fake.recordInvocation("ReaderPos", []interface{}{})
	fake.readerPosMutex.Unlock()
	if fake.ReaderPosStub != nil {
		return fake.ReaderPosStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.readerPosReturns
	return fakeReturns.result1
}

func (fake *FakeFunctionNode) ReaderPosCallCount() int {
	fake.readerPosMutex.RLock()
	defer fake.readerPosMutex.RUnlock()
	return len(fake.readerPosArgsForCall)
}

func (fake *FakeFunctionNode) ReaderPosCalls(stub func() parsley.Pos) {
	fake.readerPosMutex.Lock()
	defer fake.readerPosMutex.Unlock()
	fake.ReaderPosStub = stub
}

func (fake *FakeFunctionNode) ReaderPosReturns(result1 parsley.Pos) {
	fake.readerPosMutex.Lock()
	defer fake.readerPosMutex.Unlock()
	fake.ReaderPosStub = nil
	fake.readerPosReturns = struct {
		result1 parsley.Pos
	}{result1}
}

func (fake *FakeFunctionNode) ReaderPosReturnsOnCall(i int, result1 parsley.Pos) {
	fake.readerPosMutex.Lock()
	defer fake.readerPosMutex.Unlock()
	fake.ReaderPosStub = nil
	if fake.readerPosReturnsOnCall == nil {
		fake.readerPosReturnsOnCall = make(map[int]struct {
			result1 parsley.Pos
		})
	}
	fake.readerPosReturnsOnCall[i] = struct {
		result1 parsley.Pos
	}{result1}
}

func (fake *FakeFunctionNode) StaticCheck(arg1 interface{}) parsley.Error {
	fake.staticCheckMutex.Lock()
	ret, specificReturn := fake.staticCheckReturnsOnCall[len(fake.staticCheckArgsForCall)]
	fake.staticCheckArgsForCall = append(fake.staticCheckArgsForCall, struct {
		arg1 interface{}
	}{arg1})
	fake.recordInvocation("StaticCheck", []interface{}{arg1})
	fake.staticCheckMutex.Unlock()
	if fake.StaticCheckStub != nil {
		return fake.StaticCheckStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.staticCheckReturns
	return fakeReturns.result1
}

func (fake *FakeFunctionNode) StaticCheckCallCount() int {
	fake.staticCheckMutex.RLock()
	defer fake.staticCheckMutex.RUnlock()
	return len(fake.staticCheckArgsForCall)
}

func (fake *FakeFunctionNode) StaticCheckCalls(stub func(interface{}) parsley.Error) {
	fake.staticCheckMutex.Lock()
	defer fake.staticCheckMutex.Unlock()
	fake.StaticCheckStub = stub
}

func (fake *FakeFunctionNode) StaticCheckArgsForCall(i int) interface{} {
	fake.staticCheckMutex.RLock()
	defer fake.staticCheckMutex.RUnlock()
	argsForCall := fake.staticCheckArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeFunctionNode) StaticCheckReturns(result1 parsley.Error) {
	fake.staticCheckMutex.Lock()
	defer fake.staticCheckMutex.Unlock()
	fake.StaticCheckStub = nil
	fake.staticCheckReturns = struct {
		result1 parsley.Error
	}{result1}
}

func (fake *FakeFunctionNode) StaticCheckReturnsOnCall(i int, result1 parsley.Error) {
	fake.staticCheckMutex.Lock()
	defer fake.staticCheckMutex.Unlock()
	fake.StaticCheckStub = nil
	if fake.staticCheckReturnsOnCall == nil {
		fake.staticCheckReturnsOnCall = make(map[int]struct {
			result1 parsley.Error
		})
	}
	fake.staticCheckReturnsOnCall[i] = struct {
		result1 parsley.Error
	}{result1}
}

func (fake *FakeFunctionNode) Token() string {
	fake.tokenMutex.Lock()
	ret, specificReturn := fake.tokenReturnsOnCall[len(fake.tokenArgsForCall)]
	fake.tokenArgsForCall = append(fake.tokenArgsForCall, struct {
	}{})
	fake.recordInvocation("Token", []interface{}{})
	fake.tokenMutex.Unlock()
	if fake.TokenStub != nil {
		return fake.TokenStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.tokenReturns
	return fakeReturns.result1
}

func (fake *FakeFunctionNode) TokenCallCount() int {
	fake.tokenMutex.RLock()
	defer fake.tokenMutex.RUnlock()
	return len(fake.tokenArgsForCall)
}

func (fake *FakeFunctionNode) TokenCalls(stub func() string) {
	fake.tokenMutex.Lock()
	defer fake.tokenMutex.Unlock()
	fake.TokenStub = stub
}

func (fake *FakeFunctionNode) TokenReturns(result1 string) {
	fake.tokenMutex.Lock()
	defer fake.tokenMutex.Unlock()
	fake.TokenStub = nil
	fake.tokenReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeFunctionNode) TokenReturnsOnCall(i int, result1 string) {
	fake.tokenMutex.Lock()
	defer fake.tokenMutex.Unlock()
	fake.TokenStub = nil
	if fake.tokenReturnsOnCall == nil {
		fake.tokenReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.tokenReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeFunctionNode) Type() string {
	fake.typeMutex.Lock()
	ret, specificReturn := fake.typeReturnsOnCall[len(fake.typeArgsForCall)]
	fake.typeArgsForCall = append(fake.typeArgsForCall, struct {
	}{})
	fake.recordInvocation("Type", []interface{}{})
	fake.typeMutex.Unlock()
	if fake.TypeStub != nil {
		return fake.TypeStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.typeReturns
	return fakeReturns.result1
}

func (fake *FakeFunctionNode) TypeCallCount() int {
	fake.typeMutex.RLock()
	defer fake.typeMutex.RUnlock()
	return len(fake.typeArgsForCall)
}

func (fake *FakeFunctionNode) TypeCalls(stub func() string) {
	fake.typeMutex.Lock()
	defer fake.typeMutex.Unlock()
	fake.TypeStub = stub
}

func (fake *FakeFunctionNode) TypeReturns(result1 string) {
	fake.typeMutex.Lock()
	defer fake.typeMutex.Unlock()
	fake.TypeStub = nil
	fake.typeReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeFunctionNode) TypeReturnsOnCall(i int, result1 string) {
	fake.typeMutex.Lock()
	defer fake.typeMutex.Unlock()
	fake.TypeStub = nil
	if fake.typeReturnsOnCall == nil {
		fake.typeReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.typeReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeFunctionNode) Value(arg1 interface{}) (interface{}, parsley.Error) {
	fake.valueMutex.Lock()
	ret, specificReturn := fake.valueReturnsOnCall[len(fake.valueArgsForCall)]
	fake.valueArgsForCall = append(fake.valueArgsForCall, struct {
		arg1 interface{}
	}{arg1})
	fake.recordInvocation("Value", []interface{}{arg1})
	fake.valueMutex.Unlock()
	if fake.ValueStub != nil {
		return fake.ValueStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.valueReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeFunctionNode) ValueCallCount() int {
	fake.valueMutex.RLock()
	defer fake.valueMutex.RUnlock()
	return len(fake.valueArgsForCall)
}

func (fake *FakeFunctionNode) ValueCalls(stub func(interface{}) (interface{}, parsley.Error)) {
	fake.valueMutex.Lock()
	defer fake.valueMutex.Unlock()
	fake.ValueStub = stub
}

func (fake *FakeFunctionNode) ValueArgsForCall(i int) interface{} {
	fake.valueMutex.RLock()
	defer fake.valueMutex.RUnlock()
	argsForCall := fake.valueArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeFunctionNode) ValueReturns(result1 interface{}, result2 parsley.Error) {
	fake.valueMutex.Lock()
	defer fake.valueMutex.Unlock()
	fake.ValueStub = nil
	fake.valueReturns = struct {
		result1 interface{}
		result2 parsley.Error
	}{result1, result2}
}

func (fake *FakeFunctionNode) ValueReturnsOnCall(i int, result1 interface{}, result2 parsley.Error) {
	fake.valueMutex.Lock()
	defer fake.valueMutex.Unlock()
	fake.ValueStub = nil
	if fake.valueReturnsOnCall == nil {
		fake.valueReturnsOnCall = make(map[int]struct {
			result1 interface{}
			result2 parsley.Error
		})
	}
	fake.valueReturnsOnCall[i] = struct {
		result1 interface{}
		result2 parsley.Error
	}{result1, result2}
}

func (fake *FakeFunctionNode) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.argumentNodesMutex.RLock()
	defer fake.argumentNodesMutex.RUnlock()
	fake.nameMutex.RLock()
	defer fake.nameMutex.RUnlock()
	fake.posMutex.RLock()
	defer fake.posMutex.RUnlock()
	fake.readerPosMutex.RLock()
	defer fake.readerPosMutex.RUnlock()
	fake.staticCheckMutex.RLock()
	defer fake.staticCheckMutex.RUnlock()
	fake.tokenMutex.RLock()
	defer fake.tokenMutex.RUnlock()
	fake.typeMutex.RLock()
	defer fake.typeMutex.RUnlock()
	fake.valueMutex.RLock()
	defer fake.valueMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeFunctionNode) recordInvocation(key string, args []interface{}) {
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

var _ basil.FunctionNode = new(FakeFunctionNode)
