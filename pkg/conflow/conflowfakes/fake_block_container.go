// Code generated by counterfeiter. DO NOT EDIT.
package conflowfakes

import (
	"sync"

	"github.com/conflowio/parsley/parsley"

	"github.com/conflowio/conflow/pkg/conflow"
)

type FakeBlockContainer struct {
	CloseStub        func()
	closeMutex       sync.RWMutex
	closeArgsForCall []struct {
	}
	EvalStageStub        func() conflow.EvalStage
	evalStageMutex       sync.RWMutex
	evalStageArgsForCall []struct {
	}
	evalStageReturns struct {
		result1 conflow.EvalStage
	}
	evalStageReturnsOnCall map[int]struct {
		result1 conflow.EvalStage
	}
	NodeStub        func() conflow.Node
	nodeMutex       sync.RWMutex
	nodeArgsForCall []struct {
	}
	nodeReturns struct {
		result1 conflow.Node
	}
	nodeReturnsOnCall map[int]struct {
		result1 conflow.Node
	}
	ParamStub        func(conflow.ID) interface{}
	paramMutex       sync.RWMutex
	paramArgsForCall []struct {
		arg1 conflow.ID
	}
	paramReturns struct {
		result1 interface{}
	}
	paramReturnsOnCall map[int]struct {
		result1 interface{}
	}
	SetChildStub        func(conflow.Container)
	setChildMutex       sync.RWMutex
	setChildArgsForCall []struct {
		arg1 conflow.Container
	}
	SetErrorStub        func(parsley.Error)
	setErrorMutex       sync.RWMutex
	setErrorArgsForCall []struct {
		arg1 parsley.Error
	}
	ValueStub        func() (interface{}, parsley.Error)
	valueMutex       sync.RWMutex
	valueArgsForCall []struct {
	}
	valueReturns struct {
		result1 interface{}
		result2 parsley.Error
	}
	valueReturnsOnCall map[int]struct {
		result1 interface{}
		result2 parsley.Error
	}
	WaitGroupsStub        func() []conflow.WaitGroup
	waitGroupsMutex       sync.RWMutex
	waitGroupsArgsForCall []struct {
	}
	waitGroupsReturns struct {
		result1 []conflow.WaitGroup
	}
	waitGroupsReturnsOnCall map[int]struct {
		result1 []conflow.WaitGroup
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeBlockContainer) Close() {
	fake.closeMutex.Lock()
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct {
	}{})
	stub := fake.CloseStub
	fake.recordInvocation("Close", []interface{}{})
	fake.closeMutex.Unlock()
	if stub != nil {
		fake.CloseStub()
	}
}

func (fake *FakeBlockContainer) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *FakeBlockContainer) CloseCalls(stub func()) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = stub
}

func (fake *FakeBlockContainer) EvalStage() conflow.EvalStage {
	fake.evalStageMutex.Lock()
	ret, specificReturn := fake.evalStageReturnsOnCall[len(fake.evalStageArgsForCall)]
	fake.evalStageArgsForCall = append(fake.evalStageArgsForCall, struct {
	}{})
	stub := fake.EvalStageStub
	fakeReturns := fake.evalStageReturns
	fake.recordInvocation("EvalStage", []interface{}{})
	fake.evalStageMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeBlockContainer) EvalStageCallCount() int {
	fake.evalStageMutex.RLock()
	defer fake.evalStageMutex.RUnlock()
	return len(fake.evalStageArgsForCall)
}

func (fake *FakeBlockContainer) EvalStageCalls(stub func() conflow.EvalStage) {
	fake.evalStageMutex.Lock()
	defer fake.evalStageMutex.Unlock()
	fake.EvalStageStub = stub
}

func (fake *FakeBlockContainer) EvalStageReturns(result1 conflow.EvalStage) {
	fake.evalStageMutex.Lock()
	defer fake.evalStageMutex.Unlock()
	fake.EvalStageStub = nil
	fake.evalStageReturns = struct {
		result1 conflow.EvalStage
	}{result1}
}

func (fake *FakeBlockContainer) EvalStageReturnsOnCall(i int, result1 conflow.EvalStage) {
	fake.evalStageMutex.Lock()
	defer fake.evalStageMutex.Unlock()
	fake.EvalStageStub = nil
	if fake.evalStageReturnsOnCall == nil {
		fake.evalStageReturnsOnCall = make(map[int]struct {
			result1 conflow.EvalStage
		})
	}
	fake.evalStageReturnsOnCall[i] = struct {
		result1 conflow.EvalStage
	}{result1}
}

func (fake *FakeBlockContainer) Node() conflow.Node {
	fake.nodeMutex.Lock()
	ret, specificReturn := fake.nodeReturnsOnCall[len(fake.nodeArgsForCall)]
	fake.nodeArgsForCall = append(fake.nodeArgsForCall, struct {
	}{})
	stub := fake.NodeStub
	fakeReturns := fake.nodeReturns
	fake.recordInvocation("Node", []interface{}{})
	fake.nodeMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeBlockContainer) NodeCallCount() int {
	fake.nodeMutex.RLock()
	defer fake.nodeMutex.RUnlock()
	return len(fake.nodeArgsForCall)
}

func (fake *FakeBlockContainer) NodeCalls(stub func() conflow.Node) {
	fake.nodeMutex.Lock()
	defer fake.nodeMutex.Unlock()
	fake.NodeStub = stub
}

func (fake *FakeBlockContainer) NodeReturns(result1 conflow.Node) {
	fake.nodeMutex.Lock()
	defer fake.nodeMutex.Unlock()
	fake.NodeStub = nil
	fake.nodeReturns = struct {
		result1 conflow.Node
	}{result1}
}

func (fake *FakeBlockContainer) NodeReturnsOnCall(i int, result1 conflow.Node) {
	fake.nodeMutex.Lock()
	defer fake.nodeMutex.Unlock()
	fake.NodeStub = nil
	if fake.nodeReturnsOnCall == nil {
		fake.nodeReturnsOnCall = make(map[int]struct {
			result1 conflow.Node
		})
	}
	fake.nodeReturnsOnCall[i] = struct {
		result1 conflow.Node
	}{result1}
}

func (fake *FakeBlockContainer) Param(arg1 conflow.ID) interface{} {
	fake.paramMutex.Lock()
	ret, specificReturn := fake.paramReturnsOnCall[len(fake.paramArgsForCall)]
	fake.paramArgsForCall = append(fake.paramArgsForCall, struct {
		arg1 conflow.ID
	}{arg1})
	stub := fake.ParamStub
	fakeReturns := fake.paramReturns
	fake.recordInvocation("Param", []interface{}{arg1})
	fake.paramMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeBlockContainer) ParamCallCount() int {
	fake.paramMutex.RLock()
	defer fake.paramMutex.RUnlock()
	return len(fake.paramArgsForCall)
}

func (fake *FakeBlockContainer) ParamCalls(stub func(conflow.ID) interface{}) {
	fake.paramMutex.Lock()
	defer fake.paramMutex.Unlock()
	fake.ParamStub = stub
}

func (fake *FakeBlockContainer) ParamArgsForCall(i int) conflow.ID {
	fake.paramMutex.RLock()
	defer fake.paramMutex.RUnlock()
	argsForCall := fake.paramArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeBlockContainer) ParamReturns(result1 interface{}) {
	fake.paramMutex.Lock()
	defer fake.paramMutex.Unlock()
	fake.ParamStub = nil
	fake.paramReturns = struct {
		result1 interface{}
	}{result1}
}

func (fake *FakeBlockContainer) ParamReturnsOnCall(i int, result1 interface{}) {
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

func (fake *FakeBlockContainer) SetChild(arg1 conflow.Container) {
	fake.setChildMutex.Lock()
	fake.setChildArgsForCall = append(fake.setChildArgsForCall, struct {
		arg1 conflow.Container
	}{arg1})
	stub := fake.SetChildStub
	fake.recordInvocation("SetChild", []interface{}{arg1})
	fake.setChildMutex.Unlock()
	if stub != nil {
		fake.SetChildStub(arg1)
	}
}

func (fake *FakeBlockContainer) SetChildCallCount() int {
	fake.setChildMutex.RLock()
	defer fake.setChildMutex.RUnlock()
	return len(fake.setChildArgsForCall)
}

func (fake *FakeBlockContainer) SetChildCalls(stub func(conflow.Container)) {
	fake.setChildMutex.Lock()
	defer fake.setChildMutex.Unlock()
	fake.SetChildStub = stub
}

func (fake *FakeBlockContainer) SetChildArgsForCall(i int) conflow.Container {
	fake.setChildMutex.RLock()
	defer fake.setChildMutex.RUnlock()
	argsForCall := fake.setChildArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeBlockContainer) SetError(arg1 parsley.Error) {
	fake.setErrorMutex.Lock()
	fake.setErrorArgsForCall = append(fake.setErrorArgsForCall, struct {
		arg1 parsley.Error
	}{arg1})
	stub := fake.SetErrorStub
	fake.recordInvocation("SetError", []interface{}{arg1})
	fake.setErrorMutex.Unlock()
	if stub != nil {
		fake.SetErrorStub(arg1)
	}
}

func (fake *FakeBlockContainer) SetErrorCallCount() int {
	fake.setErrorMutex.RLock()
	defer fake.setErrorMutex.RUnlock()
	return len(fake.setErrorArgsForCall)
}

func (fake *FakeBlockContainer) SetErrorCalls(stub func(parsley.Error)) {
	fake.setErrorMutex.Lock()
	defer fake.setErrorMutex.Unlock()
	fake.SetErrorStub = stub
}

func (fake *FakeBlockContainer) SetErrorArgsForCall(i int) parsley.Error {
	fake.setErrorMutex.RLock()
	defer fake.setErrorMutex.RUnlock()
	argsForCall := fake.setErrorArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeBlockContainer) Value() (interface{}, parsley.Error) {
	fake.valueMutex.Lock()
	ret, specificReturn := fake.valueReturnsOnCall[len(fake.valueArgsForCall)]
	fake.valueArgsForCall = append(fake.valueArgsForCall, struct {
	}{})
	stub := fake.ValueStub
	fakeReturns := fake.valueReturns
	fake.recordInvocation("Value", []interface{}{})
	fake.valueMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeBlockContainer) ValueCallCount() int {
	fake.valueMutex.RLock()
	defer fake.valueMutex.RUnlock()
	return len(fake.valueArgsForCall)
}

func (fake *FakeBlockContainer) ValueCalls(stub func() (interface{}, parsley.Error)) {
	fake.valueMutex.Lock()
	defer fake.valueMutex.Unlock()
	fake.ValueStub = stub
}

func (fake *FakeBlockContainer) ValueReturns(result1 interface{}, result2 parsley.Error) {
	fake.valueMutex.Lock()
	defer fake.valueMutex.Unlock()
	fake.ValueStub = nil
	fake.valueReturns = struct {
		result1 interface{}
		result2 parsley.Error
	}{result1, result2}
}

func (fake *FakeBlockContainer) ValueReturnsOnCall(i int, result1 interface{}, result2 parsley.Error) {
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

func (fake *FakeBlockContainer) WaitGroups() []conflow.WaitGroup {
	fake.waitGroupsMutex.Lock()
	ret, specificReturn := fake.waitGroupsReturnsOnCall[len(fake.waitGroupsArgsForCall)]
	fake.waitGroupsArgsForCall = append(fake.waitGroupsArgsForCall, struct {
	}{})
	stub := fake.WaitGroupsStub
	fakeReturns := fake.waitGroupsReturns
	fake.recordInvocation("WaitGroups", []interface{}{})
	fake.waitGroupsMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeBlockContainer) WaitGroupsCallCount() int {
	fake.waitGroupsMutex.RLock()
	defer fake.waitGroupsMutex.RUnlock()
	return len(fake.waitGroupsArgsForCall)
}

func (fake *FakeBlockContainer) WaitGroupsCalls(stub func() []conflow.WaitGroup) {
	fake.waitGroupsMutex.Lock()
	defer fake.waitGroupsMutex.Unlock()
	fake.WaitGroupsStub = stub
}

func (fake *FakeBlockContainer) WaitGroupsReturns(result1 []conflow.WaitGroup) {
	fake.waitGroupsMutex.Lock()
	defer fake.waitGroupsMutex.Unlock()
	fake.WaitGroupsStub = nil
	fake.waitGroupsReturns = struct {
		result1 []conflow.WaitGroup
	}{result1}
}

func (fake *FakeBlockContainer) WaitGroupsReturnsOnCall(i int, result1 []conflow.WaitGroup) {
	fake.waitGroupsMutex.Lock()
	defer fake.waitGroupsMutex.Unlock()
	fake.WaitGroupsStub = nil
	if fake.waitGroupsReturnsOnCall == nil {
		fake.waitGroupsReturnsOnCall = make(map[int]struct {
			result1 []conflow.WaitGroup
		})
	}
	fake.waitGroupsReturnsOnCall[i] = struct {
		result1 []conflow.WaitGroup
	}{result1}
}

func (fake *FakeBlockContainer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	fake.evalStageMutex.RLock()
	defer fake.evalStageMutex.RUnlock()
	fake.nodeMutex.RLock()
	defer fake.nodeMutex.RUnlock()
	fake.paramMutex.RLock()
	defer fake.paramMutex.RUnlock()
	fake.setChildMutex.RLock()
	defer fake.setChildMutex.RUnlock()
	fake.setErrorMutex.RLock()
	defer fake.setErrorMutex.RUnlock()
	fake.valueMutex.RLock()
	defer fake.valueMutex.RUnlock()
	fake.waitGroupsMutex.RLock()
	defer fake.waitGroupsMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeBlockContainer) recordInvocation(key string, args []interface{}) {
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

var _ conflow.BlockContainer = new(FakeBlockContainer)
