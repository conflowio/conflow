// Code generated by counterfeiter. DO NOT EDIT.
package jobfakes

import (
	"sync"

	"github.com/opsidian/basil/basil"
)

type FakePendingJob struct {
	JobIDStub        func() int
	jobIDMutex       sync.RWMutex
	jobIDArgsForCall []struct {
	}
	jobIDReturns struct {
		result1 int
	}
	jobIDReturnsOnCall map[int]struct {
		result1 int
	}
	JobNameStub        func() basil.ID
	jobNameMutex       sync.RWMutex
	jobNameArgsForCall []struct {
	}
	jobNameReturns struct {
		result1 basil.ID
	}
	jobNameReturnsOnCall map[int]struct {
		result1 basil.ID
	}
	LightweightStub        func() bool
	lightweightMutex       sync.RWMutex
	lightweightArgsForCall []struct {
	}
	lightweightReturns struct {
		result1 bool
	}
	lightweightReturnsOnCall map[int]struct {
		result1 bool
	}
	PendingStub        func() bool
	pendingMutex       sync.RWMutex
	pendingArgsForCall []struct {
	}
	pendingReturns struct {
		result1 bool
	}
	pendingReturnsOnCall map[int]struct {
		result1 bool
	}
	RunStub        func()
	runMutex       sync.RWMutex
	runArgsForCall []struct {
	}
	SetJobIDStub        func(int)
	setJobIDMutex       sync.RWMutex
	setJobIDArgsForCall []struct {
		arg1 int
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakePendingJob) JobID() int {
	fake.jobIDMutex.Lock()
	ret, specificReturn := fake.jobIDReturnsOnCall[len(fake.jobIDArgsForCall)]
	fake.jobIDArgsForCall = append(fake.jobIDArgsForCall, struct {
	}{})
	stub := fake.JobIDStub
	fakeReturns := fake.jobIDReturns
	fake.recordInvocation("JobID", []interface{}{})
	fake.jobIDMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakePendingJob) JobIDCallCount() int {
	fake.jobIDMutex.RLock()
	defer fake.jobIDMutex.RUnlock()
	return len(fake.jobIDArgsForCall)
}

func (fake *FakePendingJob) JobIDCalls(stub func() int) {
	fake.jobIDMutex.Lock()
	defer fake.jobIDMutex.Unlock()
	fake.JobIDStub = stub
}

func (fake *FakePendingJob) JobIDReturns(result1 int) {
	fake.jobIDMutex.Lock()
	defer fake.jobIDMutex.Unlock()
	fake.JobIDStub = nil
	fake.jobIDReturns = struct {
		result1 int
	}{result1}
}

func (fake *FakePendingJob) JobIDReturnsOnCall(i int, result1 int) {
	fake.jobIDMutex.Lock()
	defer fake.jobIDMutex.Unlock()
	fake.JobIDStub = nil
	if fake.jobIDReturnsOnCall == nil {
		fake.jobIDReturnsOnCall = make(map[int]struct {
			result1 int
		})
	}
	fake.jobIDReturnsOnCall[i] = struct {
		result1 int
	}{result1}
}

func (fake *FakePendingJob) JobName() basil.ID {
	fake.jobNameMutex.Lock()
	ret, specificReturn := fake.jobNameReturnsOnCall[len(fake.jobNameArgsForCall)]
	fake.jobNameArgsForCall = append(fake.jobNameArgsForCall, struct {
	}{})
	stub := fake.JobNameStub
	fakeReturns := fake.jobNameReturns
	fake.recordInvocation("JobName", []interface{}{})
	fake.jobNameMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakePendingJob) JobNameCallCount() int {
	fake.jobNameMutex.RLock()
	defer fake.jobNameMutex.RUnlock()
	return len(fake.jobNameArgsForCall)
}

func (fake *FakePendingJob) JobNameCalls(stub func() basil.ID) {
	fake.jobNameMutex.Lock()
	defer fake.jobNameMutex.Unlock()
	fake.JobNameStub = stub
}

func (fake *FakePendingJob) JobNameReturns(result1 basil.ID) {
	fake.jobNameMutex.Lock()
	defer fake.jobNameMutex.Unlock()
	fake.JobNameStub = nil
	fake.jobNameReturns = struct {
		result1 basil.ID
	}{result1}
}

func (fake *FakePendingJob) JobNameReturnsOnCall(i int, result1 basil.ID) {
	fake.jobNameMutex.Lock()
	defer fake.jobNameMutex.Unlock()
	fake.JobNameStub = nil
	if fake.jobNameReturnsOnCall == nil {
		fake.jobNameReturnsOnCall = make(map[int]struct {
			result1 basil.ID
		})
	}
	fake.jobNameReturnsOnCall[i] = struct {
		result1 basil.ID
	}{result1}
}

func (fake *FakePendingJob) Lightweight() bool {
	fake.lightweightMutex.Lock()
	ret, specificReturn := fake.lightweightReturnsOnCall[len(fake.lightweightArgsForCall)]
	fake.lightweightArgsForCall = append(fake.lightweightArgsForCall, struct {
	}{})
	stub := fake.LightweightStub
	fakeReturns := fake.lightweightReturns
	fake.recordInvocation("Lightweight", []interface{}{})
	fake.lightweightMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakePendingJob) LightweightCallCount() int {
	fake.lightweightMutex.RLock()
	defer fake.lightweightMutex.RUnlock()
	return len(fake.lightweightArgsForCall)
}

func (fake *FakePendingJob) LightweightCalls(stub func() bool) {
	fake.lightweightMutex.Lock()
	defer fake.lightweightMutex.Unlock()
	fake.LightweightStub = stub
}

func (fake *FakePendingJob) LightweightReturns(result1 bool) {
	fake.lightweightMutex.Lock()
	defer fake.lightweightMutex.Unlock()
	fake.LightweightStub = nil
	fake.lightweightReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakePendingJob) LightweightReturnsOnCall(i int, result1 bool) {
	fake.lightweightMutex.Lock()
	defer fake.lightweightMutex.Unlock()
	fake.LightweightStub = nil
	if fake.lightweightReturnsOnCall == nil {
		fake.lightweightReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.lightweightReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *FakePendingJob) Pending() bool {
	fake.pendingMutex.Lock()
	ret, specificReturn := fake.pendingReturnsOnCall[len(fake.pendingArgsForCall)]
	fake.pendingArgsForCall = append(fake.pendingArgsForCall, struct {
	}{})
	stub := fake.PendingStub
	fakeReturns := fake.pendingReturns
	fake.recordInvocation("Pending", []interface{}{})
	fake.pendingMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakePendingJob) PendingCallCount() int {
	fake.pendingMutex.RLock()
	defer fake.pendingMutex.RUnlock()
	return len(fake.pendingArgsForCall)
}

func (fake *FakePendingJob) PendingCalls(stub func() bool) {
	fake.pendingMutex.Lock()
	defer fake.pendingMutex.Unlock()
	fake.PendingStub = stub
}

func (fake *FakePendingJob) PendingReturns(result1 bool) {
	fake.pendingMutex.Lock()
	defer fake.pendingMutex.Unlock()
	fake.PendingStub = nil
	fake.pendingReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakePendingJob) PendingReturnsOnCall(i int, result1 bool) {
	fake.pendingMutex.Lock()
	defer fake.pendingMutex.Unlock()
	fake.PendingStub = nil
	if fake.pendingReturnsOnCall == nil {
		fake.pendingReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.pendingReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *FakePendingJob) Run() {
	fake.runMutex.Lock()
	fake.runArgsForCall = append(fake.runArgsForCall, struct {
	}{})
	stub := fake.RunStub
	fake.recordInvocation("Run", []interface{}{})
	fake.runMutex.Unlock()
	if stub != nil {
		fake.RunStub()
	}
}

func (fake *FakePendingJob) RunCallCount() int {
	fake.runMutex.RLock()
	defer fake.runMutex.RUnlock()
	return len(fake.runArgsForCall)
}

func (fake *FakePendingJob) RunCalls(stub func()) {
	fake.runMutex.Lock()
	defer fake.runMutex.Unlock()
	fake.RunStub = stub
}

func (fake *FakePendingJob) SetJobID(arg1 int) {
	fake.setJobIDMutex.Lock()
	fake.setJobIDArgsForCall = append(fake.setJobIDArgsForCall, struct {
		arg1 int
	}{arg1})
	stub := fake.SetJobIDStub
	fake.recordInvocation("SetJobID", []interface{}{arg1})
	fake.setJobIDMutex.Unlock()
	if stub != nil {
		fake.SetJobIDStub(arg1)
	}
}

func (fake *FakePendingJob) SetJobIDCallCount() int {
	fake.setJobIDMutex.RLock()
	defer fake.setJobIDMutex.RUnlock()
	return len(fake.setJobIDArgsForCall)
}

func (fake *FakePendingJob) SetJobIDCalls(stub func(int)) {
	fake.setJobIDMutex.Lock()
	defer fake.setJobIDMutex.Unlock()
	fake.SetJobIDStub = stub
}

func (fake *FakePendingJob) SetJobIDArgsForCall(i int) int {
	fake.setJobIDMutex.RLock()
	defer fake.setJobIDMutex.RUnlock()
	argsForCall := fake.setJobIDArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakePendingJob) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.jobIDMutex.RLock()
	defer fake.jobIDMutex.RUnlock()
	fake.jobNameMutex.RLock()
	defer fake.jobNameMutex.RUnlock()
	fake.lightweightMutex.RLock()
	defer fake.lightweightMutex.RUnlock()
	fake.pendingMutex.RLock()
	defer fake.pendingMutex.RUnlock()
	fake.runMutex.RLock()
	defer fake.runMutex.RUnlock()
	fake.setJobIDMutex.RLock()
	defer fake.setJobIDMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakePendingJob) recordInvocation(key string, args []interface{}) {
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
