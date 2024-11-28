// Code generated by mockery v2.43.2. DO NOT EDIT.

package client

import (
	types "github.com/goplugin/pluginv3.0/v2/common/types"
	mock "github.com/stretchr/testify/mock"
)

// mockNodeSelector is an autogenerated mock type for the NodeSelector type
type mockNodeSelector[CHAIN_ID types.ID, RPC interface{}] struct {
	mock.Mock
}

type mockNodeSelector_Expecter[CHAIN_ID types.ID, RPC interface{}] struct {
	mock *mock.Mock
}

func (_m *mockNodeSelector[CHAIN_ID, RPC]) EXPECT() *mockNodeSelector_Expecter[CHAIN_ID, RPC] {
	return &mockNodeSelector_Expecter[CHAIN_ID, RPC]{mock: &_m.Mock}
}

// Name provides a mock function with given fields:
func (_m *mockNodeSelector[CHAIN_ID, RPC]) Name() string {
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

// mockNodeSelector_Name_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Name'
type mockNodeSelector_Name_Call[CHAIN_ID types.ID, RPC interface{}] struct {
	*mock.Call
}

// Name is a helper method to define mock.On call
func (_e *mockNodeSelector_Expecter[CHAIN_ID, RPC]) Name() *mockNodeSelector_Name_Call[CHAIN_ID, RPC] {
	return &mockNodeSelector_Name_Call[CHAIN_ID, RPC]{Call: _e.mock.On("Name")}
}

func (_c *mockNodeSelector_Name_Call[CHAIN_ID, RPC]) Run(run func()) *mockNodeSelector_Name_Call[CHAIN_ID, RPC] {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockNodeSelector_Name_Call[CHAIN_ID, RPC]) Return(_a0 string) *mockNodeSelector_Name_Call[CHAIN_ID, RPC] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockNodeSelector_Name_Call[CHAIN_ID, RPC]) RunAndReturn(run func() string) *mockNodeSelector_Name_Call[CHAIN_ID, RPC] {
	_c.Call.Return(run)
	return _c
}

// Select provides a mock function with given fields:
func (_m *mockNodeSelector[CHAIN_ID, RPC]) Select() Node[CHAIN_ID, RPC] {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Select")
	}

	var r0 Node[CHAIN_ID, RPC]
	if rf, ok := ret.Get(0).(func() Node[CHAIN_ID, RPC]); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Node[CHAIN_ID, RPC])
		}
	}

	return r0
}

// mockNodeSelector_Select_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Select'
type mockNodeSelector_Select_Call[CHAIN_ID types.ID, RPC interface{}] struct {
	*mock.Call
}

// Select is a helper method to define mock.On call
func (_e *mockNodeSelector_Expecter[CHAIN_ID, RPC]) Select() *mockNodeSelector_Select_Call[CHAIN_ID, RPC] {
	return &mockNodeSelector_Select_Call[CHAIN_ID, RPC]{Call: _e.mock.On("Select")}
}

func (_c *mockNodeSelector_Select_Call[CHAIN_ID, RPC]) Run(run func()) *mockNodeSelector_Select_Call[CHAIN_ID, RPC] {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockNodeSelector_Select_Call[CHAIN_ID, RPC]) Return(_a0 Node[CHAIN_ID, RPC]) *mockNodeSelector_Select_Call[CHAIN_ID, RPC] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockNodeSelector_Select_Call[CHAIN_ID, RPC]) RunAndReturn(run func() Node[CHAIN_ID, RPC]) *mockNodeSelector_Select_Call[CHAIN_ID, RPC] {
	_c.Call.Return(run)
	return _c
}

// newMockNodeSelector creates a new instance of mockNodeSelector. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockNodeSelector[CHAIN_ID types.ID, RPC interface{}](t interface {
	mock.TestingT
	Cleanup(func())
}) *mockNodeSelector[CHAIN_ID, RPC] {
	mock := &mockNodeSelector[CHAIN_ID, RPC]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
