// Code generated by counterfeiter. DO NOT EDIT.
package basilfakes

import (
	"sync"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/parsley"
)

type FakeFunctionNode struct {
	TokenStub        func() string
	tokenMutex       sync.RWMutex
	tokenArgsForCall []struct{}
	tokenReturns     struct {
		result1 string
	}
	tokenReturnsOnCall map[int]struct {
		result1 string
	}
	TypeStub        func() string
	typeMutex       sync.RWMutex
	typeArgsForCall []struct{}
	typeReturns     struct {
		result1 string
	}
	typeReturnsOnCall map[int]struct {
		result1 string
	}
	ValueStub        func(userCtx interface{}) (interface{}, parsley.Error)
	valueMutex       sync.RWMutex
	valueArgsForCall []struct {
		userCtx interface{}
	}
	valueReturns struct {
		result1 interface{}
		result2 parsley.Error
	}
	valueReturnsOnCall map[int]struct {
		result1 interface{}
		result2 parsley.Error
	}
	PosStub        func() parsley.Pos
	posMutex       sync.RWMutex
	posArgsForCall []struct{}
	posReturns     struct {
		result1 parsley.Pos
	}
	posReturnsOnCall map[int]struct {
		result1 parsley.Pos
	}
	ReaderPosStub        func() parsley.Pos
	readerPosMutex       sync.RWMutex
	readerPosArgsForCall []struct{}
	readerPosReturns     struct {
		result1 parsley.Pos
	}
	readerPosReturnsOnCall map[int]struct {
		result1 parsley.Pos
	}
	StaticCheckStub        func(userCtx interface{}) parsley.Error
	staticCheckMutex       sync.RWMutex
	staticCheckArgsForCall []struct {
		userCtx interface{}
	}
	staticCheckReturns struct {
		result1 parsley.Error
	}
	staticCheckReturnsOnCall map[int]struct {
		result1 parsley.Error
	}
	IDStub        func() basil.ID
	iDMutex       sync.RWMutex
	iDArgsForCall []struct{}
	iDReturns     struct {
		result1 basil.ID
	}
	iDReturnsOnCall map[int]struct {
		result1 basil.ID
	}
	ArgumentNodesStub        func() []parsley.Node
	argumentNodesMutex       sync.RWMutex
	argumentNodesArgsForCall []struct{}
	argumentNodesReturns     struct {
		result1 []parsley.Node
	}
	argumentNodesReturnsOnCall map[int]struct {
		result1 []parsley.Node
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeFunctionNode) Token() string {
	fake.tokenMutex.Lock()
	ret, specificReturn := fake.tokenReturnsOnCall[len(fake.tokenArgsForCall)]
	fake.tokenArgsForCall = append(fake.tokenArgsForCall, struct{}{})
	fake.recordInvocation("Token", []interface{}{})
	fake.tokenMutex.Unlock()
	if fake.TokenStub != nil {
		return fake.TokenStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.tokenReturns.result1
}

func (fake *FakeFunctionNode) TokenCallCount() int {
	fake.tokenMutex.RLock()
	defer fake.tokenMutex.RUnlock()
	return len(fake.tokenArgsForCall)
}

func (fake *FakeFunctionNode) TokenReturns(result1 string) {
	fake.TokenStub = nil
	fake.tokenReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeFunctionNode) TokenReturnsOnCall(i int, result1 string) {
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
	fake.typeArgsForCall = append(fake.typeArgsForCall, struct{}{})
	fake.recordInvocation("Type", []interface{}{})
	fake.typeMutex.Unlock()
	if fake.TypeStub != nil {
		return fake.TypeStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.typeReturns.result1
}

func (fake *FakeFunctionNode) TypeCallCount() int {
	fake.typeMutex.RLock()
	defer fake.typeMutex.RUnlock()
	return len(fake.typeArgsForCall)
}

func (fake *FakeFunctionNode) TypeReturns(result1 string) {
	fake.TypeStub = nil
	fake.typeReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeFunctionNode) TypeReturnsOnCall(i int, result1 string) {
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

func (fake *FakeFunctionNode) Value(userCtx interface{}) (interface{}, parsley.Error) {
	fake.valueMutex.Lock()
	ret, specificReturn := fake.valueReturnsOnCall[len(fake.valueArgsForCall)]
	fake.valueArgsForCall = append(fake.valueArgsForCall, struct {
		userCtx interface{}
	}{userCtx})
	fake.recordInvocation("Value", []interface{}{userCtx})
	fake.valueMutex.Unlock()
	if fake.ValueStub != nil {
		return fake.ValueStub(userCtx)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.valueReturns.result1, fake.valueReturns.result2
}

func (fake *FakeFunctionNode) ValueCallCount() int {
	fake.valueMutex.RLock()
	defer fake.valueMutex.RUnlock()
	return len(fake.valueArgsForCall)
}

func (fake *FakeFunctionNode) ValueArgsForCall(i int) interface{} {
	fake.valueMutex.RLock()
	defer fake.valueMutex.RUnlock()
	return fake.valueArgsForCall[i].userCtx
}

func (fake *FakeFunctionNode) ValueReturns(result1 interface{}, result2 parsley.Error) {
	fake.ValueStub = nil
	fake.valueReturns = struct {
		result1 interface{}
		result2 parsley.Error
	}{result1, result2}
}

func (fake *FakeFunctionNode) ValueReturnsOnCall(i int, result1 interface{}, result2 parsley.Error) {
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

func (fake *FakeFunctionNode) Pos() parsley.Pos {
	fake.posMutex.Lock()
	ret, specificReturn := fake.posReturnsOnCall[len(fake.posArgsForCall)]
	fake.posArgsForCall = append(fake.posArgsForCall, struct{}{})
	fake.recordInvocation("Pos", []interface{}{})
	fake.posMutex.Unlock()
	if fake.PosStub != nil {
		return fake.PosStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.posReturns.result1
}

func (fake *FakeFunctionNode) PosCallCount() int {
	fake.posMutex.RLock()
	defer fake.posMutex.RUnlock()
	return len(fake.posArgsForCall)
}

func (fake *FakeFunctionNode) PosReturns(result1 parsley.Pos) {
	fake.PosStub = nil
	fake.posReturns = struct {
		result1 parsley.Pos
	}{result1}
}

func (fake *FakeFunctionNode) PosReturnsOnCall(i int, result1 parsley.Pos) {
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
	fake.readerPosArgsForCall = append(fake.readerPosArgsForCall, struct{}{})
	fake.recordInvocation("ReaderPos", []interface{}{})
	fake.readerPosMutex.Unlock()
	if fake.ReaderPosStub != nil {
		return fake.ReaderPosStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.readerPosReturns.result1
}

func (fake *FakeFunctionNode) ReaderPosCallCount() int {
	fake.readerPosMutex.RLock()
	defer fake.readerPosMutex.RUnlock()
	return len(fake.readerPosArgsForCall)
}

func (fake *FakeFunctionNode) ReaderPosReturns(result1 parsley.Pos) {
	fake.ReaderPosStub = nil
	fake.readerPosReturns = struct {
		result1 parsley.Pos
	}{result1}
}

func (fake *FakeFunctionNode) ReaderPosReturnsOnCall(i int, result1 parsley.Pos) {
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

func (fake *FakeFunctionNode) StaticCheck(userCtx interface{}) parsley.Error {
	fake.staticCheckMutex.Lock()
	ret, specificReturn := fake.staticCheckReturnsOnCall[len(fake.staticCheckArgsForCall)]
	fake.staticCheckArgsForCall = append(fake.staticCheckArgsForCall, struct {
		userCtx interface{}
	}{userCtx})
	fake.recordInvocation("StaticCheck", []interface{}{userCtx})
	fake.staticCheckMutex.Unlock()
	if fake.StaticCheckStub != nil {
		return fake.StaticCheckStub(userCtx)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.staticCheckReturns.result1
}

func (fake *FakeFunctionNode) StaticCheckCallCount() int {
	fake.staticCheckMutex.RLock()
	defer fake.staticCheckMutex.RUnlock()
	return len(fake.staticCheckArgsForCall)
}

func (fake *FakeFunctionNode) StaticCheckArgsForCall(i int) interface{} {
	fake.staticCheckMutex.RLock()
	defer fake.staticCheckMutex.RUnlock()
	return fake.staticCheckArgsForCall[i].userCtx
}

func (fake *FakeFunctionNode) StaticCheckReturns(result1 parsley.Error) {
	fake.StaticCheckStub = nil
	fake.staticCheckReturns = struct {
		result1 parsley.Error
	}{result1}
}

func (fake *FakeFunctionNode) StaticCheckReturnsOnCall(i int, result1 parsley.Error) {
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

func (fake *FakeFunctionNode) ID() basil.ID {
	fake.iDMutex.Lock()
	ret, specificReturn := fake.iDReturnsOnCall[len(fake.iDArgsForCall)]
	fake.iDArgsForCall = append(fake.iDArgsForCall, struct{}{})
	fake.recordInvocation("ID", []interface{}{})
	fake.iDMutex.Unlock()
	if fake.IDStub != nil {
		return fake.IDStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.iDReturns.result1
}

func (fake *FakeFunctionNode) IDCallCount() int {
	fake.iDMutex.RLock()
	defer fake.iDMutex.RUnlock()
	return len(fake.iDArgsForCall)
}

func (fake *FakeFunctionNode) IDReturns(result1 basil.ID) {
	fake.IDStub = nil
	fake.iDReturns = struct {
		result1 basil.ID
	}{result1}
}

func (fake *FakeFunctionNode) IDReturnsOnCall(i int, result1 basil.ID) {
	fake.IDStub = nil
	if fake.iDReturnsOnCall == nil {
		fake.iDReturnsOnCall = make(map[int]struct {
			result1 basil.ID
		})
	}
	fake.iDReturnsOnCall[i] = struct {
		result1 basil.ID
	}{result1}
}

func (fake *FakeFunctionNode) ArgumentNodes() []parsley.Node {
	fake.argumentNodesMutex.Lock()
	ret, specificReturn := fake.argumentNodesReturnsOnCall[len(fake.argumentNodesArgsForCall)]
	fake.argumentNodesArgsForCall = append(fake.argumentNodesArgsForCall, struct{}{})
	fake.recordInvocation("ArgumentNodes", []interface{}{})
	fake.argumentNodesMutex.Unlock()
	if fake.ArgumentNodesStub != nil {
		return fake.ArgumentNodesStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.argumentNodesReturns.result1
}

func (fake *FakeFunctionNode) ArgumentNodesCallCount() int {
	fake.argumentNodesMutex.RLock()
	defer fake.argumentNodesMutex.RUnlock()
	return len(fake.argumentNodesArgsForCall)
}

func (fake *FakeFunctionNode) ArgumentNodesReturns(result1 []parsley.Node) {
	fake.ArgumentNodesStub = nil
	fake.argumentNodesReturns = struct {
		result1 []parsley.Node
	}{result1}
}

func (fake *FakeFunctionNode) ArgumentNodesReturnsOnCall(i int, result1 []parsley.Node) {
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

func (fake *FakeFunctionNode) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.tokenMutex.RLock()
	defer fake.tokenMutex.RUnlock()
	fake.typeMutex.RLock()
	defer fake.typeMutex.RUnlock()
	fake.valueMutex.RLock()
	defer fake.valueMutex.RUnlock()
	fake.posMutex.RLock()
	defer fake.posMutex.RUnlock()
	fake.readerPosMutex.RLock()
	defer fake.readerPosMutex.RUnlock()
	fake.staticCheckMutex.RLock()
	defer fake.staticCheckMutex.RUnlock()
	fake.iDMutex.RLock()
	defer fake.iDMutex.RUnlock()
	fake.argumentNodesMutex.RLock()
	defer fake.argumentNodesMutex.RUnlock()
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
