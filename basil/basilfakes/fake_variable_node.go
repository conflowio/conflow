// Code generated by counterfeiter. DO NOT EDIT.
package basilfakes

import (
	"sync"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/parsley"
)

type FakeVariableNode struct {
	IDStub        func() basil.ID
	iDMutex       sync.RWMutex
	iDArgsForCall []struct {
	}
	iDReturns struct {
		result1 basil.ID
	}
	iDReturnsOnCall map[int]struct {
		result1 basil.ID
	}
	ParamNameStub        func() basil.ID
	paramNameMutex       sync.RWMutex
	paramNameArgsForCall []struct {
	}
	paramNameReturns struct {
		result1 basil.ID
	}
	paramNameReturnsOnCall map[int]struct {
		result1 basil.ID
	}
	ParentIDStub        func() basil.ID
	parentIDMutex       sync.RWMutex
	parentIDArgsForCall []struct {
	}
	parentIDReturns struct {
		result1 basil.ID
	}
	parentIDReturnsOnCall map[int]struct {
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
	SchemaStub        func() interface{}
	schemaMutex       sync.RWMutex
	schemaArgsForCall []struct {
	}
	schemaReturns struct {
		result1 interface{}
	}
	schemaReturnsOnCall map[int]struct {
		result1 interface{}
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
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeVariableNode) ID() basil.ID {
	fake.iDMutex.Lock()
	ret, specificReturn := fake.iDReturnsOnCall[len(fake.iDArgsForCall)]
	fake.iDArgsForCall = append(fake.iDArgsForCall, struct {
	}{})
	stub := fake.IDStub
	fakeReturns := fake.iDReturns
	fake.recordInvocation("ID", []interface{}{})
	fake.iDMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeVariableNode) IDCallCount() int {
	fake.iDMutex.RLock()
	defer fake.iDMutex.RUnlock()
	return len(fake.iDArgsForCall)
}

func (fake *FakeVariableNode) IDCalls(stub func() basil.ID) {
	fake.iDMutex.Lock()
	defer fake.iDMutex.Unlock()
	fake.IDStub = stub
}

func (fake *FakeVariableNode) IDReturns(result1 basil.ID) {
	fake.iDMutex.Lock()
	defer fake.iDMutex.Unlock()
	fake.IDStub = nil
	fake.iDReturns = struct {
		result1 basil.ID
	}{result1}
}

func (fake *FakeVariableNode) IDReturnsOnCall(i int, result1 basil.ID) {
	fake.iDMutex.Lock()
	defer fake.iDMutex.Unlock()
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

func (fake *FakeVariableNode) ParamName() basil.ID {
	fake.paramNameMutex.Lock()
	ret, specificReturn := fake.paramNameReturnsOnCall[len(fake.paramNameArgsForCall)]
	fake.paramNameArgsForCall = append(fake.paramNameArgsForCall, struct {
	}{})
	stub := fake.ParamNameStub
	fakeReturns := fake.paramNameReturns
	fake.recordInvocation("ParamName", []interface{}{})
	fake.paramNameMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeVariableNode) ParamNameCallCount() int {
	fake.paramNameMutex.RLock()
	defer fake.paramNameMutex.RUnlock()
	return len(fake.paramNameArgsForCall)
}

func (fake *FakeVariableNode) ParamNameCalls(stub func() basil.ID) {
	fake.paramNameMutex.Lock()
	defer fake.paramNameMutex.Unlock()
	fake.ParamNameStub = stub
}

func (fake *FakeVariableNode) ParamNameReturns(result1 basil.ID) {
	fake.paramNameMutex.Lock()
	defer fake.paramNameMutex.Unlock()
	fake.ParamNameStub = nil
	fake.paramNameReturns = struct {
		result1 basil.ID
	}{result1}
}

func (fake *FakeVariableNode) ParamNameReturnsOnCall(i int, result1 basil.ID) {
	fake.paramNameMutex.Lock()
	defer fake.paramNameMutex.Unlock()
	fake.ParamNameStub = nil
	if fake.paramNameReturnsOnCall == nil {
		fake.paramNameReturnsOnCall = make(map[int]struct {
			result1 basil.ID
		})
	}
	fake.paramNameReturnsOnCall[i] = struct {
		result1 basil.ID
	}{result1}
}

func (fake *FakeVariableNode) ParentID() basil.ID {
	fake.parentIDMutex.Lock()
	ret, specificReturn := fake.parentIDReturnsOnCall[len(fake.parentIDArgsForCall)]
	fake.parentIDArgsForCall = append(fake.parentIDArgsForCall, struct {
	}{})
	stub := fake.ParentIDStub
	fakeReturns := fake.parentIDReturns
	fake.recordInvocation("ParentID", []interface{}{})
	fake.parentIDMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeVariableNode) ParentIDCallCount() int {
	fake.parentIDMutex.RLock()
	defer fake.parentIDMutex.RUnlock()
	return len(fake.parentIDArgsForCall)
}

func (fake *FakeVariableNode) ParentIDCalls(stub func() basil.ID) {
	fake.parentIDMutex.Lock()
	defer fake.parentIDMutex.Unlock()
	fake.ParentIDStub = stub
}

func (fake *FakeVariableNode) ParentIDReturns(result1 basil.ID) {
	fake.parentIDMutex.Lock()
	defer fake.parentIDMutex.Unlock()
	fake.ParentIDStub = nil
	fake.parentIDReturns = struct {
		result1 basil.ID
	}{result1}
}

func (fake *FakeVariableNode) ParentIDReturnsOnCall(i int, result1 basil.ID) {
	fake.parentIDMutex.Lock()
	defer fake.parentIDMutex.Unlock()
	fake.ParentIDStub = nil
	if fake.parentIDReturnsOnCall == nil {
		fake.parentIDReturnsOnCall = make(map[int]struct {
			result1 basil.ID
		})
	}
	fake.parentIDReturnsOnCall[i] = struct {
		result1 basil.ID
	}{result1}
}

func (fake *FakeVariableNode) Pos() parsley.Pos {
	fake.posMutex.Lock()
	ret, specificReturn := fake.posReturnsOnCall[len(fake.posArgsForCall)]
	fake.posArgsForCall = append(fake.posArgsForCall, struct {
	}{})
	stub := fake.PosStub
	fakeReturns := fake.posReturns
	fake.recordInvocation("Pos", []interface{}{})
	fake.posMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeVariableNode) PosCallCount() int {
	fake.posMutex.RLock()
	defer fake.posMutex.RUnlock()
	return len(fake.posArgsForCall)
}

func (fake *FakeVariableNode) PosCalls(stub func() parsley.Pos) {
	fake.posMutex.Lock()
	defer fake.posMutex.Unlock()
	fake.PosStub = stub
}

func (fake *FakeVariableNode) PosReturns(result1 parsley.Pos) {
	fake.posMutex.Lock()
	defer fake.posMutex.Unlock()
	fake.PosStub = nil
	fake.posReturns = struct {
		result1 parsley.Pos
	}{result1}
}

func (fake *FakeVariableNode) PosReturnsOnCall(i int, result1 parsley.Pos) {
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

func (fake *FakeVariableNode) ReaderPos() parsley.Pos {
	fake.readerPosMutex.Lock()
	ret, specificReturn := fake.readerPosReturnsOnCall[len(fake.readerPosArgsForCall)]
	fake.readerPosArgsForCall = append(fake.readerPosArgsForCall, struct {
	}{})
	stub := fake.ReaderPosStub
	fakeReturns := fake.readerPosReturns
	fake.recordInvocation("ReaderPos", []interface{}{})
	fake.readerPosMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeVariableNode) ReaderPosCallCount() int {
	fake.readerPosMutex.RLock()
	defer fake.readerPosMutex.RUnlock()
	return len(fake.readerPosArgsForCall)
}

func (fake *FakeVariableNode) ReaderPosCalls(stub func() parsley.Pos) {
	fake.readerPosMutex.Lock()
	defer fake.readerPosMutex.Unlock()
	fake.ReaderPosStub = stub
}

func (fake *FakeVariableNode) ReaderPosReturns(result1 parsley.Pos) {
	fake.readerPosMutex.Lock()
	defer fake.readerPosMutex.Unlock()
	fake.ReaderPosStub = nil
	fake.readerPosReturns = struct {
		result1 parsley.Pos
	}{result1}
}

func (fake *FakeVariableNode) ReaderPosReturnsOnCall(i int, result1 parsley.Pos) {
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

func (fake *FakeVariableNode) Schema() interface{} {
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

func (fake *FakeVariableNode) SchemaCallCount() int {
	fake.schemaMutex.RLock()
	defer fake.schemaMutex.RUnlock()
	return len(fake.schemaArgsForCall)
}

func (fake *FakeVariableNode) SchemaCalls(stub func() interface{}) {
	fake.schemaMutex.Lock()
	defer fake.schemaMutex.Unlock()
	fake.SchemaStub = stub
}

func (fake *FakeVariableNode) SchemaReturns(result1 interface{}) {
	fake.schemaMutex.Lock()
	defer fake.schemaMutex.Unlock()
	fake.SchemaStub = nil
	fake.schemaReturns = struct {
		result1 interface{}
	}{result1}
}

func (fake *FakeVariableNode) SchemaReturnsOnCall(i int, result1 interface{}) {
	fake.schemaMutex.Lock()
	defer fake.schemaMutex.Unlock()
	fake.SchemaStub = nil
	if fake.schemaReturnsOnCall == nil {
		fake.schemaReturnsOnCall = make(map[int]struct {
			result1 interface{}
		})
	}
	fake.schemaReturnsOnCall[i] = struct {
		result1 interface{}
	}{result1}
}

func (fake *FakeVariableNode) Token() string {
	fake.tokenMutex.Lock()
	ret, specificReturn := fake.tokenReturnsOnCall[len(fake.tokenArgsForCall)]
	fake.tokenArgsForCall = append(fake.tokenArgsForCall, struct {
	}{})
	stub := fake.TokenStub
	fakeReturns := fake.tokenReturns
	fake.recordInvocation("Token", []interface{}{})
	fake.tokenMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeVariableNode) TokenCallCount() int {
	fake.tokenMutex.RLock()
	defer fake.tokenMutex.RUnlock()
	return len(fake.tokenArgsForCall)
}

func (fake *FakeVariableNode) TokenCalls(stub func() string) {
	fake.tokenMutex.Lock()
	defer fake.tokenMutex.Unlock()
	fake.TokenStub = stub
}

func (fake *FakeVariableNode) TokenReturns(result1 string) {
	fake.tokenMutex.Lock()
	defer fake.tokenMutex.Unlock()
	fake.TokenStub = nil
	fake.tokenReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeVariableNode) TokenReturnsOnCall(i int, result1 string) {
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

func (fake *FakeVariableNode) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.iDMutex.RLock()
	defer fake.iDMutex.RUnlock()
	fake.paramNameMutex.RLock()
	defer fake.paramNameMutex.RUnlock()
	fake.parentIDMutex.RLock()
	defer fake.parentIDMutex.RUnlock()
	fake.posMutex.RLock()
	defer fake.posMutex.RUnlock()
	fake.readerPosMutex.RLock()
	defer fake.readerPosMutex.RUnlock()
	fake.schemaMutex.RLock()
	defer fake.schemaMutex.RUnlock()
	fake.tokenMutex.RLock()
	defer fake.tokenMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeVariableNode) recordInvocation(key string, args []interface{}) {
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

var _ basil.VariableNode = new(FakeVariableNode)
