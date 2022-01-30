// Code generated by counterfeiter. DO NOT EDIT.
package conflowfakes

import (
	"sync"

	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/schema"
	"github.com/conflowio/parsley/parsley"
)

type FakeParameterNode struct {
	CreateContainerStub        func(*conflow.EvalContext, conflow.RuntimeConfig, conflow.BlockContainer, interface{}, []conflow.WaitGroup, bool) conflow.JobContainer
	createContainerMutex       sync.RWMutex
	createContainerArgsForCall []struct {
		arg1 *conflow.EvalContext
		arg2 conflow.RuntimeConfig
		arg3 conflow.BlockContainer
		arg4 interface{}
		arg5 []conflow.WaitGroup
		arg6 bool
	}
	createContainerReturns struct {
		result1 conflow.JobContainer
	}
	createContainerReturnsOnCall map[int]struct {
		result1 conflow.JobContainer
	}
	DependenciesStub        func() conflow.Dependencies
	dependenciesMutex       sync.RWMutex
	dependenciesArgsForCall []struct {
	}
	dependenciesReturns struct {
		result1 conflow.Dependencies
	}
	dependenciesReturnsOnCall map[int]struct {
		result1 conflow.Dependencies
	}
	DirectivesStub        func() []conflow.BlockNode
	directivesMutex       sync.RWMutex
	directivesArgsForCall []struct {
	}
	directivesReturns struct {
		result1 []conflow.BlockNode
	}
	directivesReturnsOnCall map[int]struct {
		result1 []conflow.BlockNode
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
	GeneratesStub        func() []conflow.ID
	generatesMutex       sync.RWMutex
	generatesArgsForCall []struct {
	}
	generatesReturns struct {
		result1 []conflow.ID
	}
	generatesReturnsOnCall map[int]struct {
		result1 []conflow.ID
	}
	IDStub        func() conflow.ID
	iDMutex       sync.RWMutex
	iDArgsForCall []struct {
	}
	iDReturns struct {
		result1 conflow.ID
	}
	iDReturnsOnCall map[int]struct {
		result1 conflow.ID
	}
	IsDeclarationStub        func() bool
	isDeclarationMutex       sync.RWMutex
	isDeclarationArgsForCall []struct {
	}
	isDeclarationReturns struct {
		result1 bool
	}
	isDeclarationReturnsOnCall map[int]struct {
		result1 bool
	}
	NameStub        func() conflow.ID
	nameMutex       sync.RWMutex
	nameArgsForCall []struct {
	}
	nameReturns struct {
		result1 conflow.ID
	}
	nameReturnsOnCall map[int]struct {
		result1 conflow.ID
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
	ProvidesStub        func() []conflow.ID
	providesMutex       sync.RWMutex
	providesArgsForCall []struct {
	}
	providesReturns struct {
		result1 []conflow.ID
	}
	providesReturnsOnCall map[int]struct {
		result1 []conflow.ID
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
	SetSchemaStub        func(schema.Schema)
	setSchemaMutex       sync.RWMutex
	setSchemaArgsForCall []struct {
		arg1 schema.Schema
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
	ValueNodeStub        func() parsley.Node
	valueNodeMutex       sync.RWMutex
	valueNodeArgsForCall []struct {
	}
	valueNodeReturns struct {
		result1 parsley.Node
	}
	valueNodeReturnsOnCall map[int]struct {
		result1 parsley.Node
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeParameterNode) CreateContainer(arg1 *conflow.EvalContext, arg2 conflow.RuntimeConfig, arg3 conflow.BlockContainer, arg4 interface{}, arg5 []conflow.WaitGroup, arg6 bool) conflow.JobContainer {
	var arg5Copy []conflow.WaitGroup
	if arg5 != nil {
		arg5Copy = make([]conflow.WaitGroup, len(arg5))
		copy(arg5Copy, arg5)
	}
	fake.createContainerMutex.Lock()
	ret, specificReturn := fake.createContainerReturnsOnCall[len(fake.createContainerArgsForCall)]
	fake.createContainerArgsForCall = append(fake.createContainerArgsForCall, struct {
		arg1 *conflow.EvalContext
		arg2 conflow.RuntimeConfig
		arg3 conflow.BlockContainer
		arg4 interface{}
		arg5 []conflow.WaitGroup
		arg6 bool
	}{arg1, arg2, arg3, arg4, arg5Copy, arg6})
	stub := fake.CreateContainerStub
	fakeReturns := fake.createContainerReturns
	fake.recordInvocation("CreateContainer", []interface{}{arg1, arg2, arg3, arg4, arg5Copy, arg6})
	fake.createContainerMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5, arg6)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeParameterNode) CreateContainerCallCount() int {
	fake.createContainerMutex.RLock()
	defer fake.createContainerMutex.RUnlock()
	return len(fake.createContainerArgsForCall)
}

func (fake *FakeParameterNode) CreateContainerCalls(stub func(*conflow.EvalContext, conflow.RuntimeConfig, conflow.BlockContainer, interface{}, []conflow.WaitGroup, bool) conflow.JobContainer) {
	fake.createContainerMutex.Lock()
	defer fake.createContainerMutex.Unlock()
	fake.CreateContainerStub = stub
}

func (fake *FakeParameterNode) CreateContainerArgsForCall(i int) (*conflow.EvalContext, conflow.RuntimeConfig, conflow.BlockContainer, interface{}, []conflow.WaitGroup, bool) {
	fake.createContainerMutex.RLock()
	defer fake.createContainerMutex.RUnlock()
	argsForCall := fake.createContainerArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5, argsForCall.arg6
}

func (fake *FakeParameterNode) CreateContainerReturns(result1 conflow.JobContainer) {
	fake.createContainerMutex.Lock()
	defer fake.createContainerMutex.Unlock()
	fake.CreateContainerStub = nil
	fake.createContainerReturns = struct {
		result1 conflow.JobContainer
	}{result1}
}

func (fake *FakeParameterNode) CreateContainerReturnsOnCall(i int, result1 conflow.JobContainer) {
	fake.createContainerMutex.Lock()
	defer fake.createContainerMutex.Unlock()
	fake.CreateContainerStub = nil
	if fake.createContainerReturnsOnCall == nil {
		fake.createContainerReturnsOnCall = make(map[int]struct {
			result1 conflow.JobContainer
		})
	}
	fake.createContainerReturnsOnCall[i] = struct {
		result1 conflow.JobContainer
	}{result1}
}

func (fake *FakeParameterNode) Dependencies() conflow.Dependencies {
	fake.dependenciesMutex.Lock()
	ret, specificReturn := fake.dependenciesReturnsOnCall[len(fake.dependenciesArgsForCall)]
	fake.dependenciesArgsForCall = append(fake.dependenciesArgsForCall, struct {
	}{})
	stub := fake.DependenciesStub
	fakeReturns := fake.dependenciesReturns
	fake.recordInvocation("Dependencies", []interface{}{})
	fake.dependenciesMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeParameterNode) DependenciesCallCount() int {
	fake.dependenciesMutex.RLock()
	defer fake.dependenciesMutex.RUnlock()
	return len(fake.dependenciesArgsForCall)
}

func (fake *FakeParameterNode) DependenciesCalls(stub func() conflow.Dependencies) {
	fake.dependenciesMutex.Lock()
	defer fake.dependenciesMutex.Unlock()
	fake.DependenciesStub = stub
}

func (fake *FakeParameterNode) DependenciesReturns(result1 conflow.Dependencies) {
	fake.dependenciesMutex.Lock()
	defer fake.dependenciesMutex.Unlock()
	fake.DependenciesStub = nil
	fake.dependenciesReturns = struct {
		result1 conflow.Dependencies
	}{result1}
}

func (fake *FakeParameterNode) DependenciesReturnsOnCall(i int, result1 conflow.Dependencies) {
	fake.dependenciesMutex.Lock()
	defer fake.dependenciesMutex.Unlock()
	fake.DependenciesStub = nil
	if fake.dependenciesReturnsOnCall == nil {
		fake.dependenciesReturnsOnCall = make(map[int]struct {
			result1 conflow.Dependencies
		})
	}
	fake.dependenciesReturnsOnCall[i] = struct {
		result1 conflow.Dependencies
	}{result1}
}

func (fake *FakeParameterNode) Directives() []conflow.BlockNode {
	fake.directivesMutex.Lock()
	ret, specificReturn := fake.directivesReturnsOnCall[len(fake.directivesArgsForCall)]
	fake.directivesArgsForCall = append(fake.directivesArgsForCall, struct {
	}{})
	stub := fake.DirectivesStub
	fakeReturns := fake.directivesReturns
	fake.recordInvocation("Directives", []interface{}{})
	fake.directivesMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeParameterNode) DirectivesCallCount() int {
	fake.directivesMutex.RLock()
	defer fake.directivesMutex.RUnlock()
	return len(fake.directivesArgsForCall)
}

func (fake *FakeParameterNode) DirectivesCalls(stub func() []conflow.BlockNode) {
	fake.directivesMutex.Lock()
	defer fake.directivesMutex.Unlock()
	fake.DirectivesStub = stub
}

func (fake *FakeParameterNode) DirectivesReturns(result1 []conflow.BlockNode) {
	fake.directivesMutex.Lock()
	defer fake.directivesMutex.Unlock()
	fake.DirectivesStub = nil
	fake.directivesReturns = struct {
		result1 []conflow.BlockNode
	}{result1}
}

func (fake *FakeParameterNode) DirectivesReturnsOnCall(i int, result1 []conflow.BlockNode) {
	fake.directivesMutex.Lock()
	defer fake.directivesMutex.Unlock()
	fake.DirectivesStub = nil
	if fake.directivesReturnsOnCall == nil {
		fake.directivesReturnsOnCall = make(map[int]struct {
			result1 []conflow.BlockNode
		})
	}
	fake.directivesReturnsOnCall[i] = struct {
		result1 []conflow.BlockNode
	}{result1}
}

func (fake *FakeParameterNode) EvalStage() conflow.EvalStage {
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

func (fake *FakeParameterNode) EvalStageCallCount() int {
	fake.evalStageMutex.RLock()
	defer fake.evalStageMutex.RUnlock()
	return len(fake.evalStageArgsForCall)
}

func (fake *FakeParameterNode) EvalStageCalls(stub func() conflow.EvalStage) {
	fake.evalStageMutex.Lock()
	defer fake.evalStageMutex.Unlock()
	fake.EvalStageStub = stub
}

func (fake *FakeParameterNode) EvalStageReturns(result1 conflow.EvalStage) {
	fake.evalStageMutex.Lock()
	defer fake.evalStageMutex.Unlock()
	fake.EvalStageStub = nil
	fake.evalStageReturns = struct {
		result1 conflow.EvalStage
	}{result1}
}

func (fake *FakeParameterNode) EvalStageReturnsOnCall(i int, result1 conflow.EvalStage) {
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

func (fake *FakeParameterNode) Generates() []conflow.ID {
	fake.generatesMutex.Lock()
	ret, specificReturn := fake.generatesReturnsOnCall[len(fake.generatesArgsForCall)]
	fake.generatesArgsForCall = append(fake.generatesArgsForCall, struct {
	}{})
	stub := fake.GeneratesStub
	fakeReturns := fake.generatesReturns
	fake.recordInvocation("Generates", []interface{}{})
	fake.generatesMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeParameterNode) GeneratesCallCount() int {
	fake.generatesMutex.RLock()
	defer fake.generatesMutex.RUnlock()
	return len(fake.generatesArgsForCall)
}

func (fake *FakeParameterNode) GeneratesCalls(stub func() []conflow.ID) {
	fake.generatesMutex.Lock()
	defer fake.generatesMutex.Unlock()
	fake.GeneratesStub = stub
}

func (fake *FakeParameterNode) GeneratesReturns(result1 []conflow.ID) {
	fake.generatesMutex.Lock()
	defer fake.generatesMutex.Unlock()
	fake.GeneratesStub = nil
	fake.generatesReturns = struct {
		result1 []conflow.ID
	}{result1}
}

func (fake *FakeParameterNode) GeneratesReturnsOnCall(i int, result1 []conflow.ID) {
	fake.generatesMutex.Lock()
	defer fake.generatesMutex.Unlock()
	fake.GeneratesStub = nil
	if fake.generatesReturnsOnCall == nil {
		fake.generatesReturnsOnCall = make(map[int]struct {
			result1 []conflow.ID
		})
	}
	fake.generatesReturnsOnCall[i] = struct {
		result1 []conflow.ID
	}{result1}
}

func (fake *FakeParameterNode) ID() conflow.ID {
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

func (fake *FakeParameterNode) IDCallCount() int {
	fake.iDMutex.RLock()
	defer fake.iDMutex.RUnlock()
	return len(fake.iDArgsForCall)
}

func (fake *FakeParameterNode) IDCalls(stub func() conflow.ID) {
	fake.iDMutex.Lock()
	defer fake.iDMutex.Unlock()
	fake.IDStub = stub
}

func (fake *FakeParameterNode) IDReturns(result1 conflow.ID) {
	fake.iDMutex.Lock()
	defer fake.iDMutex.Unlock()
	fake.IDStub = nil
	fake.iDReturns = struct {
		result1 conflow.ID
	}{result1}
}

func (fake *FakeParameterNode) IDReturnsOnCall(i int, result1 conflow.ID) {
	fake.iDMutex.Lock()
	defer fake.iDMutex.Unlock()
	fake.IDStub = nil
	if fake.iDReturnsOnCall == nil {
		fake.iDReturnsOnCall = make(map[int]struct {
			result1 conflow.ID
		})
	}
	fake.iDReturnsOnCall[i] = struct {
		result1 conflow.ID
	}{result1}
}

func (fake *FakeParameterNode) IsDeclaration() bool {
	fake.isDeclarationMutex.Lock()
	ret, specificReturn := fake.isDeclarationReturnsOnCall[len(fake.isDeclarationArgsForCall)]
	fake.isDeclarationArgsForCall = append(fake.isDeclarationArgsForCall, struct {
	}{})
	stub := fake.IsDeclarationStub
	fakeReturns := fake.isDeclarationReturns
	fake.recordInvocation("IsDeclaration", []interface{}{})
	fake.isDeclarationMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeParameterNode) IsDeclarationCallCount() int {
	fake.isDeclarationMutex.RLock()
	defer fake.isDeclarationMutex.RUnlock()
	return len(fake.isDeclarationArgsForCall)
}

func (fake *FakeParameterNode) IsDeclarationCalls(stub func() bool) {
	fake.isDeclarationMutex.Lock()
	defer fake.isDeclarationMutex.Unlock()
	fake.IsDeclarationStub = stub
}

func (fake *FakeParameterNode) IsDeclarationReturns(result1 bool) {
	fake.isDeclarationMutex.Lock()
	defer fake.isDeclarationMutex.Unlock()
	fake.IsDeclarationStub = nil
	fake.isDeclarationReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakeParameterNode) IsDeclarationReturnsOnCall(i int, result1 bool) {
	fake.isDeclarationMutex.Lock()
	defer fake.isDeclarationMutex.Unlock()
	fake.IsDeclarationStub = nil
	if fake.isDeclarationReturnsOnCall == nil {
		fake.isDeclarationReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.isDeclarationReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *FakeParameterNode) Name() conflow.ID {
	fake.nameMutex.Lock()
	ret, specificReturn := fake.nameReturnsOnCall[len(fake.nameArgsForCall)]
	fake.nameArgsForCall = append(fake.nameArgsForCall, struct {
	}{})
	stub := fake.NameStub
	fakeReturns := fake.nameReturns
	fake.recordInvocation("Name", []interface{}{})
	fake.nameMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeParameterNode) NameCallCount() int {
	fake.nameMutex.RLock()
	defer fake.nameMutex.RUnlock()
	return len(fake.nameArgsForCall)
}

func (fake *FakeParameterNode) NameCalls(stub func() conflow.ID) {
	fake.nameMutex.Lock()
	defer fake.nameMutex.Unlock()
	fake.NameStub = stub
}

func (fake *FakeParameterNode) NameReturns(result1 conflow.ID) {
	fake.nameMutex.Lock()
	defer fake.nameMutex.Unlock()
	fake.NameStub = nil
	fake.nameReturns = struct {
		result1 conflow.ID
	}{result1}
}

func (fake *FakeParameterNode) NameReturnsOnCall(i int, result1 conflow.ID) {
	fake.nameMutex.Lock()
	defer fake.nameMutex.Unlock()
	fake.NameStub = nil
	if fake.nameReturnsOnCall == nil {
		fake.nameReturnsOnCall = make(map[int]struct {
			result1 conflow.ID
		})
	}
	fake.nameReturnsOnCall[i] = struct {
		result1 conflow.ID
	}{result1}
}

func (fake *FakeParameterNode) Pos() parsley.Pos {
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

func (fake *FakeParameterNode) PosCallCount() int {
	fake.posMutex.RLock()
	defer fake.posMutex.RUnlock()
	return len(fake.posArgsForCall)
}

func (fake *FakeParameterNode) PosCalls(stub func() parsley.Pos) {
	fake.posMutex.Lock()
	defer fake.posMutex.Unlock()
	fake.PosStub = stub
}

func (fake *FakeParameterNode) PosReturns(result1 parsley.Pos) {
	fake.posMutex.Lock()
	defer fake.posMutex.Unlock()
	fake.PosStub = nil
	fake.posReturns = struct {
		result1 parsley.Pos
	}{result1}
}

func (fake *FakeParameterNode) PosReturnsOnCall(i int, result1 parsley.Pos) {
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

func (fake *FakeParameterNode) Provides() []conflow.ID {
	fake.providesMutex.Lock()
	ret, specificReturn := fake.providesReturnsOnCall[len(fake.providesArgsForCall)]
	fake.providesArgsForCall = append(fake.providesArgsForCall, struct {
	}{})
	stub := fake.ProvidesStub
	fakeReturns := fake.providesReturns
	fake.recordInvocation("Provides", []interface{}{})
	fake.providesMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeParameterNode) ProvidesCallCount() int {
	fake.providesMutex.RLock()
	defer fake.providesMutex.RUnlock()
	return len(fake.providesArgsForCall)
}

func (fake *FakeParameterNode) ProvidesCalls(stub func() []conflow.ID) {
	fake.providesMutex.Lock()
	defer fake.providesMutex.Unlock()
	fake.ProvidesStub = stub
}

func (fake *FakeParameterNode) ProvidesReturns(result1 []conflow.ID) {
	fake.providesMutex.Lock()
	defer fake.providesMutex.Unlock()
	fake.ProvidesStub = nil
	fake.providesReturns = struct {
		result1 []conflow.ID
	}{result1}
}

func (fake *FakeParameterNode) ProvidesReturnsOnCall(i int, result1 []conflow.ID) {
	fake.providesMutex.Lock()
	defer fake.providesMutex.Unlock()
	fake.ProvidesStub = nil
	if fake.providesReturnsOnCall == nil {
		fake.providesReturnsOnCall = make(map[int]struct {
			result1 []conflow.ID
		})
	}
	fake.providesReturnsOnCall[i] = struct {
		result1 []conflow.ID
	}{result1}
}

func (fake *FakeParameterNode) ReaderPos() parsley.Pos {
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

func (fake *FakeParameterNode) ReaderPosCallCount() int {
	fake.readerPosMutex.RLock()
	defer fake.readerPosMutex.RUnlock()
	return len(fake.readerPosArgsForCall)
}

func (fake *FakeParameterNode) ReaderPosCalls(stub func() parsley.Pos) {
	fake.readerPosMutex.Lock()
	defer fake.readerPosMutex.Unlock()
	fake.ReaderPosStub = stub
}

func (fake *FakeParameterNode) ReaderPosReturns(result1 parsley.Pos) {
	fake.readerPosMutex.Lock()
	defer fake.readerPosMutex.Unlock()
	fake.ReaderPosStub = nil
	fake.readerPosReturns = struct {
		result1 parsley.Pos
	}{result1}
}

func (fake *FakeParameterNode) ReaderPosReturnsOnCall(i int, result1 parsley.Pos) {
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

func (fake *FakeParameterNode) Schema() interface{} {
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

func (fake *FakeParameterNode) SchemaCallCount() int {
	fake.schemaMutex.RLock()
	defer fake.schemaMutex.RUnlock()
	return len(fake.schemaArgsForCall)
}

func (fake *FakeParameterNode) SchemaCalls(stub func() interface{}) {
	fake.schemaMutex.Lock()
	defer fake.schemaMutex.Unlock()
	fake.SchemaStub = stub
}

func (fake *FakeParameterNode) SchemaReturns(result1 interface{}) {
	fake.schemaMutex.Lock()
	defer fake.schemaMutex.Unlock()
	fake.SchemaStub = nil
	fake.schemaReturns = struct {
		result1 interface{}
	}{result1}
}

func (fake *FakeParameterNode) SchemaReturnsOnCall(i int, result1 interface{}) {
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

func (fake *FakeParameterNode) SetSchema(arg1 schema.Schema) {
	fake.setSchemaMutex.Lock()
	fake.setSchemaArgsForCall = append(fake.setSchemaArgsForCall, struct {
		arg1 schema.Schema
	}{arg1})
	stub := fake.SetSchemaStub
	fake.recordInvocation("SetSchema", []interface{}{arg1})
	fake.setSchemaMutex.Unlock()
	if stub != nil {
		fake.SetSchemaStub(arg1)
	}
}

func (fake *FakeParameterNode) SetSchemaCallCount() int {
	fake.setSchemaMutex.RLock()
	defer fake.setSchemaMutex.RUnlock()
	return len(fake.setSchemaArgsForCall)
}

func (fake *FakeParameterNode) SetSchemaCalls(stub func(schema.Schema)) {
	fake.setSchemaMutex.Lock()
	defer fake.setSchemaMutex.Unlock()
	fake.SetSchemaStub = stub
}

func (fake *FakeParameterNode) SetSchemaArgsForCall(i int) schema.Schema {
	fake.setSchemaMutex.RLock()
	defer fake.setSchemaMutex.RUnlock()
	argsForCall := fake.setSchemaArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeParameterNode) StaticCheck(arg1 interface{}) parsley.Error {
	fake.staticCheckMutex.Lock()
	ret, specificReturn := fake.staticCheckReturnsOnCall[len(fake.staticCheckArgsForCall)]
	fake.staticCheckArgsForCall = append(fake.staticCheckArgsForCall, struct {
		arg1 interface{}
	}{arg1})
	stub := fake.StaticCheckStub
	fakeReturns := fake.staticCheckReturns
	fake.recordInvocation("StaticCheck", []interface{}{arg1})
	fake.staticCheckMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeParameterNode) StaticCheckCallCount() int {
	fake.staticCheckMutex.RLock()
	defer fake.staticCheckMutex.RUnlock()
	return len(fake.staticCheckArgsForCall)
}

func (fake *FakeParameterNode) StaticCheckCalls(stub func(interface{}) parsley.Error) {
	fake.staticCheckMutex.Lock()
	defer fake.staticCheckMutex.Unlock()
	fake.StaticCheckStub = stub
}

func (fake *FakeParameterNode) StaticCheckArgsForCall(i int) interface{} {
	fake.staticCheckMutex.RLock()
	defer fake.staticCheckMutex.RUnlock()
	argsForCall := fake.staticCheckArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeParameterNode) StaticCheckReturns(result1 parsley.Error) {
	fake.staticCheckMutex.Lock()
	defer fake.staticCheckMutex.Unlock()
	fake.StaticCheckStub = nil
	fake.staticCheckReturns = struct {
		result1 parsley.Error
	}{result1}
}

func (fake *FakeParameterNode) StaticCheckReturnsOnCall(i int, result1 parsley.Error) {
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

func (fake *FakeParameterNode) Token() string {
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

func (fake *FakeParameterNode) TokenCallCount() int {
	fake.tokenMutex.RLock()
	defer fake.tokenMutex.RUnlock()
	return len(fake.tokenArgsForCall)
}

func (fake *FakeParameterNode) TokenCalls(stub func() string) {
	fake.tokenMutex.Lock()
	defer fake.tokenMutex.Unlock()
	fake.TokenStub = stub
}

func (fake *FakeParameterNode) TokenReturns(result1 string) {
	fake.tokenMutex.Lock()
	defer fake.tokenMutex.Unlock()
	fake.TokenStub = nil
	fake.tokenReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeParameterNode) TokenReturnsOnCall(i int, result1 string) {
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

func (fake *FakeParameterNode) Value(arg1 interface{}) (interface{}, parsley.Error) {
	fake.valueMutex.Lock()
	ret, specificReturn := fake.valueReturnsOnCall[len(fake.valueArgsForCall)]
	fake.valueArgsForCall = append(fake.valueArgsForCall, struct {
		arg1 interface{}
	}{arg1})
	stub := fake.ValueStub
	fakeReturns := fake.valueReturns
	fake.recordInvocation("Value", []interface{}{arg1})
	fake.valueMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeParameterNode) ValueCallCount() int {
	fake.valueMutex.RLock()
	defer fake.valueMutex.RUnlock()
	return len(fake.valueArgsForCall)
}

func (fake *FakeParameterNode) ValueCalls(stub func(interface{}) (interface{}, parsley.Error)) {
	fake.valueMutex.Lock()
	defer fake.valueMutex.Unlock()
	fake.ValueStub = stub
}

func (fake *FakeParameterNode) ValueArgsForCall(i int) interface{} {
	fake.valueMutex.RLock()
	defer fake.valueMutex.RUnlock()
	argsForCall := fake.valueArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeParameterNode) ValueReturns(result1 interface{}, result2 parsley.Error) {
	fake.valueMutex.Lock()
	defer fake.valueMutex.Unlock()
	fake.ValueStub = nil
	fake.valueReturns = struct {
		result1 interface{}
		result2 parsley.Error
	}{result1, result2}
}

func (fake *FakeParameterNode) ValueReturnsOnCall(i int, result1 interface{}, result2 parsley.Error) {
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

func (fake *FakeParameterNode) ValueNode() parsley.Node {
	fake.valueNodeMutex.Lock()
	ret, specificReturn := fake.valueNodeReturnsOnCall[len(fake.valueNodeArgsForCall)]
	fake.valueNodeArgsForCall = append(fake.valueNodeArgsForCall, struct {
	}{})
	stub := fake.ValueNodeStub
	fakeReturns := fake.valueNodeReturns
	fake.recordInvocation("ValueNode", []interface{}{})
	fake.valueNodeMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeParameterNode) ValueNodeCallCount() int {
	fake.valueNodeMutex.RLock()
	defer fake.valueNodeMutex.RUnlock()
	return len(fake.valueNodeArgsForCall)
}

func (fake *FakeParameterNode) ValueNodeCalls(stub func() parsley.Node) {
	fake.valueNodeMutex.Lock()
	defer fake.valueNodeMutex.Unlock()
	fake.ValueNodeStub = stub
}

func (fake *FakeParameterNode) ValueNodeReturns(result1 parsley.Node) {
	fake.valueNodeMutex.Lock()
	defer fake.valueNodeMutex.Unlock()
	fake.ValueNodeStub = nil
	fake.valueNodeReturns = struct {
		result1 parsley.Node
	}{result1}
}

func (fake *FakeParameterNode) ValueNodeReturnsOnCall(i int, result1 parsley.Node) {
	fake.valueNodeMutex.Lock()
	defer fake.valueNodeMutex.Unlock()
	fake.ValueNodeStub = nil
	if fake.valueNodeReturnsOnCall == nil {
		fake.valueNodeReturnsOnCall = make(map[int]struct {
			result1 parsley.Node
		})
	}
	fake.valueNodeReturnsOnCall[i] = struct {
		result1 parsley.Node
	}{result1}
}

func (fake *FakeParameterNode) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createContainerMutex.RLock()
	defer fake.createContainerMutex.RUnlock()
	fake.dependenciesMutex.RLock()
	defer fake.dependenciesMutex.RUnlock()
	fake.directivesMutex.RLock()
	defer fake.directivesMutex.RUnlock()
	fake.evalStageMutex.RLock()
	defer fake.evalStageMutex.RUnlock()
	fake.generatesMutex.RLock()
	defer fake.generatesMutex.RUnlock()
	fake.iDMutex.RLock()
	defer fake.iDMutex.RUnlock()
	fake.isDeclarationMutex.RLock()
	defer fake.isDeclarationMutex.RUnlock()
	fake.nameMutex.RLock()
	defer fake.nameMutex.RUnlock()
	fake.posMutex.RLock()
	defer fake.posMutex.RUnlock()
	fake.providesMutex.RLock()
	defer fake.providesMutex.RUnlock()
	fake.readerPosMutex.RLock()
	defer fake.readerPosMutex.RUnlock()
	fake.schemaMutex.RLock()
	defer fake.schemaMutex.RUnlock()
	fake.setSchemaMutex.RLock()
	defer fake.setSchemaMutex.RUnlock()
	fake.staticCheckMutex.RLock()
	defer fake.staticCheckMutex.RUnlock()
	fake.tokenMutex.RLock()
	defer fake.tokenMutex.RUnlock()
	fake.valueMutex.RLock()
	defer fake.valueMutex.RUnlock()
	fake.valueNodeMutex.RLock()
	defer fake.valueNodeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeParameterNode) recordInvocation(key string, args []interface{}) {
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

var _ conflow.ParameterNode = new(FakeParameterNode)
