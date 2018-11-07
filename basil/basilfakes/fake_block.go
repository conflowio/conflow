// Code generated by counterfeiter. DO NOT EDIT.
package basilfakes

import (
	"sync"

	"github.com/opsidian/basil/basil"
)

type FakeBlock struct {
	IDStub        func() basil.ID
	iDMutex       sync.RWMutex
	iDArgsForCall []struct{}
	iDReturns     struct {
		result1 basil.ID
	}
	iDReturnsOnCall map[int]struct {
		result1 basil.ID
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
	ContextStub        func(parentCtx interface{}) interface{}
	contextMutex       sync.RWMutex
	contextArgsForCall []struct {
		parentCtx interface{}
	}
	contextReturns struct {
		result1 interface{}
	}
	contextReturnsOnCall map[int]struct {
		result1 interface{}
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeBlock) ID() basil.ID {
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

func (fake *FakeBlock) IDCallCount() int {
	fake.iDMutex.RLock()
	defer fake.iDMutex.RUnlock()
	return len(fake.iDArgsForCall)
}

func (fake *FakeBlock) IDReturns(result1 basil.ID) {
	fake.IDStub = nil
	fake.iDReturns = struct {
		result1 basil.ID
	}{result1}
}

func (fake *FakeBlock) IDReturnsOnCall(i int, result1 basil.ID) {
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

func (fake *FakeBlock) Type() string {
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

func (fake *FakeBlock) TypeCallCount() int {
	fake.typeMutex.RLock()
	defer fake.typeMutex.RUnlock()
	return len(fake.typeArgsForCall)
}

func (fake *FakeBlock) TypeReturns(result1 string) {
	fake.TypeStub = nil
	fake.typeReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeBlock) TypeReturnsOnCall(i int, result1 string) {
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

func (fake *FakeBlock) Context(parentCtx interface{}) interface{} {
	fake.contextMutex.Lock()
	ret, specificReturn := fake.contextReturnsOnCall[len(fake.contextArgsForCall)]
	fake.contextArgsForCall = append(fake.contextArgsForCall, struct {
		parentCtx interface{}
	}{parentCtx})
	fake.recordInvocation("Context", []interface{}{parentCtx})
	fake.contextMutex.Unlock()
	if fake.ContextStub != nil {
		return fake.ContextStub(parentCtx)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.contextReturns.result1
}

func (fake *FakeBlock) ContextCallCount() int {
	fake.contextMutex.RLock()
	defer fake.contextMutex.RUnlock()
	return len(fake.contextArgsForCall)
}

func (fake *FakeBlock) ContextArgsForCall(i int) interface{} {
	fake.contextMutex.RLock()
	defer fake.contextMutex.RUnlock()
	return fake.contextArgsForCall[i].parentCtx
}

func (fake *FakeBlock) ContextReturns(result1 interface{}) {
	fake.ContextStub = nil
	fake.contextReturns = struct {
		result1 interface{}
	}{result1}
}

func (fake *FakeBlock) ContextReturnsOnCall(i int, result1 interface{}) {
	fake.ContextStub = nil
	if fake.contextReturnsOnCall == nil {
		fake.contextReturnsOnCall = make(map[int]struct {
			result1 interface{}
		})
	}
	fake.contextReturnsOnCall[i] = struct {
		result1 interface{}
	}{result1}
}

func (fake *FakeBlock) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.iDMutex.RLock()
	defer fake.iDMutex.RUnlock()
	fake.typeMutex.RLock()
	defer fake.typeMutex.RUnlock()
	fake.contextMutex.RLock()
	defer fake.contextMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeBlock) recordInvocation(key string, args []interface{}) {
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

var _ basil.Block = new(FakeBlock)
