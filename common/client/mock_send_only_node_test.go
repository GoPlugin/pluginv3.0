// Code generated by mockery v2.43.2. DO NOT EDIT.

package client

import (
	context "context"

	types "github.com/goplugin/pluginv3.0/v2/common/types"
	mock "github.com/stretchr/testify/mock"
)

// mockSendOnlyNode is an autogenerated mock type for the SendOnlyNode type
type mockSendOnlyNode[CHAIN_ID types.ID, RPC interface{}] struct {
	mock.Mock
}

type mockSendOnlyNode_Expecter[CHAIN_ID types.ID, RPC interface{}] struct {
	mock *mock.Mock
}

func (_m *mockSendOnlyNode[CHAIN_ID, RPC]) EXPECT() *mockSendOnlyNode_Expecter[CHAIN_ID, RPC] {
	return &mockSendOnlyNode_Expecter[CHAIN_ID, RPC]{mock: &_m.Mock}
}

// Close provides a mock function with given fields:
func (_m *mockSendOnlyNode[CHAIN_ID, RPC]) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockSendOnlyNode_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type mockSendOnlyNode_Close_Call[CHAIN_ID types.ID, RPC interface{}] struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *mockSendOnlyNode_Expecter[CHAIN_ID, RPC]) Close() *mockSendOnlyNode_Close_Call[CHAIN_ID, RPC] {
	return &mockSendOnlyNode_Close_Call[CHAIN_ID, RPC]{Call: _e.mock.On("Close")}
}

func (_c *mockSendOnlyNode_Close_Call[CHAIN_ID, RPC]) Run(run func()) *mockSendOnlyNode_Close_Call[CHAIN_ID, RPC] {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockSendOnlyNode_Close_Call[CHAIN_ID, RPC]) Return(_a0 error) *mockSendOnlyNode_Close_Call[CHAIN_ID, RPC] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockSendOnlyNode_Close_Call[CHAIN_ID, RPC]) RunAndReturn(run func() error) *mockSendOnlyNode_Close_Call[CHAIN_ID, RPC] {
	_c.Call.Return(run)
	return _c
}

// ConfiguredChainID provides a mock function with given fields:
func (_m *mockSendOnlyNode[CHAIN_ID, RPC]) ConfiguredChainID() CHAIN_ID {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ConfiguredChainID")
	}

	var r0 CHAIN_ID
	if rf, ok := ret.Get(0).(func() CHAIN_ID); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(CHAIN_ID)
	}

	return r0
}

// mockSendOnlyNode_ConfiguredChainID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ConfiguredChainID'
type mockSendOnlyNode_ConfiguredChainID_Call[CHAIN_ID types.ID, RPC interface{}] struct {
	*mock.Call
}

// ConfiguredChainID is a helper method to define mock.On call
func (_e *mockSendOnlyNode_Expecter[CHAIN_ID, RPC]) ConfiguredChainID() *mockSendOnlyNode_ConfiguredChainID_Call[CHAIN_ID, RPC] {
	return &mockSendOnlyNode_ConfiguredChainID_Call[CHAIN_ID, RPC]{Call: _e.mock.On("ConfiguredChainID")}
}

func (_c *mockSendOnlyNode_ConfiguredChainID_Call[CHAIN_ID, RPC]) Run(run func()) *mockSendOnlyNode_ConfiguredChainID_Call[CHAIN_ID, RPC] {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockSendOnlyNode_ConfiguredChainID_Call[CHAIN_ID, RPC]) Return(_a0 CHAIN_ID) *mockSendOnlyNode_ConfiguredChainID_Call[CHAIN_ID, RPC] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockSendOnlyNode_ConfiguredChainID_Call[CHAIN_ID, RPC]) RunAndReturn(run func() CHAIN_ID) *mockSendOnlyNode_ConfiguredChainID_Call[CHAIN_ID, RPC] {
	_c.Call.Return(run)
	return _c
}

// Name provides a mock function with given fields:
func (_m *mockSendOnlyNode[CHAIN_ID, RPC]) Name() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Name")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// mockSendOnlyNode_Name_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Name'
type mockSendOnlyNode_Name_Call[CHAIN_ID types.ID, RPC interface{}] struct {
	*mock.Call
}

// Name is a helper method to define mock.On call
func (_e *mockSendOnlyNode_Expecter[CHAIN_ID, RPC]) Name() *mockSendOnlyNode_Name_Call[CHAIN_ID, RPC] {
	return &mockSendOnlyNode_Name_Call[CHAIN_ID, RPC]{Call: _e.mock.On("Name")}
}

func (_c *mockSendOnlyNode_Name_Call[CHAIN_ID, RPC]) Run(run func()) *mockSendOnlyNode_Name_Call[CHAIN_ID, RPC] {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockSendOnlyNode_Name_Call[CHAIN_ID, RPC]) Return(_a0 string) *mockSendOnlyNode_Name_Call[CHAIN_ID, RPC] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockSendOnlyNode_Name_Call[CHAIN_ID, RPC]) RunAndReturn(run func() string) *mockSendOnlyNode_Name_Call[CHAIN_ID, RPC] {
	_c.Call.Return(run)
	return _c
}

// RPC provides a mock function with given fields:
func (_m *mockSendOnlyNode[CHAIN_ID, RPC]) RPC() RPC {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for RPC")
	}

	var r0 RPC
	if rf, ok := ret.Get(0).(func() RPC); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(RPC)
	}

	return r0
}

// mockSendOnlyNode_RPC_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RPC'
type mockSendOnlyNode_RPC_Call[CHAIN_ID types.ID, RPC interface{}] struct {
	*mock.Call
}

// RPC is a helper method to define mock.On call
func (_e *mockSendOnlyNode_Expecter[CHAIN_ID, RPC]) RPC() *mockSendOnlyNode_RPC_Call[CHAIN_ID, RPC] {
	return &mockSendOnlyNode_RPC_Call[CHAIN_ID, RPC]{Call: _e.mock.On("RPC")}
}

func (_c *mockSendOnlyNode_RPC_Call[CHAIN_ID, RPC]) Run(run func()) *mockSendOnlyNode_RPC_Call[CHAIN_ID, RPC] {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockSendOnlyNode_RPC_Call[CHAIN_ID, RPC]) Return(_a0 RPC) *mockSendOnlyNode_RPC_Call[CHAIN_ID, RPC] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockSendOnlyNode_RPC_Call[CHAIN_ID, RPC]) RunAndReturn(run func() RPC) *mockSendOnlyNode_RPC_Call[CHAIN_ID, RPC] {
	_c.Call.Return(run)
	return _c
}

// Start provides a mock function with given fields: _a0
func (_m *mockSendOnlyNode[CHAIN_ID, RPC]) Start(_a0 context.Context) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Start")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockSendOnlyNode_Start_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Start'
type mockSendOnlyNode_Start_Call[CHAIN_ID types.ID, RPC interface{}] struct {
	*mock.Call
}

// Start is a helper method to define mock.On call
//   - _a0 context.Context
func (_e *mockSendOnlyNode_Expecter[CHAIN_ID, RPC]) Start(_a0 interface{}) *mockSendOnlyNode_Start_Call[CHAIN_ID, RPC] {
	return &mockSendOnlyNode_Start_Call[CHAIN_ID, RPC]{Call: _e.mock.On("Start", _a0)}
}

func (_c *mockSendOnlyNode_Start_Call[CHAIN_ID, RPC]) Run(run func(_a0 context.Context)) *mockSendOnlyNode_Start_Call[CHAIN_ID, RPC] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *mockSendOnlyNode_Start_Call[CHAIN_ID, RPC]) Return(_a0 error) *mockSendOnlyNode_Start_Call[CHAIN_ID, RPC] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockSendOnlyNode_Start_Call[CHAIN_ID, RPC]) RunAndReturn(run func(context.Context) error) *mockSendOnlyNode_Start_Call[CHAIN_ID, RPC] {
	_c.Call.Return(run)
	return _c
}

// State provides a mock function with given fields:
func (_m *mockSendOnlyNode[CHAIN_ID, RPC]) State() nodeState {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for State")
	}

	var r0 nodeState
	if rf, ok := ret.Get(0).(func() nodeState); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(nodeState)
	}

	return r0
}

// mockSendOnlyNode_State_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'State'
type mockSendOnlyNode_State_Call[CHAIN_ID types.ID, RPC interface{}] struct {
	*mock.Call
}

// State is a helper method to define mock.On call
func (_e *mockSendOnlyNode_Expecter[CHAIN_ID, RPC]) State() *mockSendOnlyNode_State_Call[CHAIN_ID, RPC] {
	return &mockSendOnlyNode_State_Call[CHAIN_ID, RPC]{Call: _e.mock.On("State")}
}

func (_c *mockSendOnlyNode_State_Call[CHAIN_ID, RPC]) Run(run func()) *mockSendOnlyNode_State_Call[CHAIN_ID, RPC] {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockSendOnlyNode_State_Call[CHAIN_ID, RPC]) Return(_a0 nodeState) *mockSendOnlyNode_State_Call[CHAIN_ID, RPC] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockSendOnlyNode_State_Call[CHAIN_ID, RPC]) RunAndReturn(run func() nodeState) *mockSendOnlyNode_State_Call[CHAIN_ID, RPC] {
	_c.Call.Return(run)
	return _c
}

// String provides a mock function with given fields:
func (_m *mockSendOnlyNode[CHAIN_ID, RPC]) String() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for String")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// mockSendOnlyNode_String_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'String'
type mockSendOnlyNode_String_Call[CHAIN_ID types.ID, RPC interface{}] struct {
	*mock.Call
}

// String is a helper method to define mock.On call
func (_e *mockSendOnlyNode_Expecter[CHAIN_ID, RPC]) String() *mockSendOnlyNode_String_Call[CHAIN_ID, RPC] {
	return &mockSendOnlyNode_String_Call[CHAIN_ID, RPC]{Call: _e.mock.On("String")}
}

func (_c *mockSendOnlyNode_String_Call[CHAIN_ID, RPC]) Run(run func()) *mockSendOnlyNode_String_Call[CHAIN_ID, RPC] {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockSendOnlyNode_String_Call[CHAIN_ID, RPC]) Return(_a0 string) *mockSendOnlyNode_String_Call[CHAIN_ID, RPC] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockSendOnlyNode_String_Call[CHAIN_ID, RPC]) RunAndReturn(run func() string) *mockSendOnlyNode_String_Call[CHAIN_ID, RPC] {
	_c.Call.Return(run)
	return _c
}

// newMockSendOnlyNode creates a new instance of mockSendOnlyNode. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockSendOnlyNode[CHAIN_ID types.ID, RPC interface{}](t interface {
	mock.TestingT
	Cleanup(func())
}) *mockSendOnlyNode[CHAIN_ID, RPC] {
	mock := &mockSendOnlyNode[CHAIN_ID, RPC]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
