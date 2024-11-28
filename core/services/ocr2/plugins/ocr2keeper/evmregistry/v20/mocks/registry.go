// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	big "math/big"

	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"

	generated "github.com/goplugin/pluginv3.0/v2/core/gethwrappers/generated"

	keeper_registry_wrapper2_0 "github.com/goplugin/pluginv3.0/v2/core/gethwrappers/generated/keeper_registry_wrapper2_0"

	mock "github.com/stretchr/testify/mock"

	types "github.com/ethereum/go-ethereum/core/types"
)

// Registry is an autogenerated mock type for the Registry type
type Registry struct {
	mock.Mock
}

// GetActiveUpkeepIDs provides a mock function with given fields: opts, startIndex, maxCount
func (_m *Registry) GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error) {
	ret := _m.Called(opts, startIndex, maxCount)

	if len(ret) == 0 {
		panic("no return value specified for GetActiveUpkeepIDs")
	}

	var r0 []*big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts, *big.Int, *big.Int) ([]*big.Int, error)); ok {
		return rf(opts, startIndex, maxCount)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts, *big.Int, *big.Int) []*big.Int); ok {
		r0 = rf(opts, startIndex, maxCount)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts, *big.Int, *big.Int) error); ok {
		r1 = rf(opts, startIndex, maxCount)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetState provides a mock function with given fields: opts
func (_m *Registry) GetState(opts *bind.CallOpts) (keeper_registry_wrapper2_0.GetState, error) {
	ret := _m.Called(opts)

	if len(ret) == 0 {
		panic("no return value specified for GetState")
	}

	var r0 keeper_registry_wrapper2_0.GetState
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) (keeper_registry_wrapper2_0.GetState, error)); ok {
		return rf(opts)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) keeper_registry_wrapper2_0.GetState); ok {
		r0 = rf(opts)
	} else {
		r0 = ret.Get(0).(keeper_registry_wrapper2_0.GetState)
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts) error); ok {
		r1 = rf(opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUpkeep provides a mock function with given fields: opts, id
func (_m *Registry) GetUpkeep(opts *bind.CallOpts, id *big.Int) (keeper_registry_wrapper2_0.UpkeepInfo, error) {
	ret := _m.Called(opts, id)

	if len(ret) == 0 {
		panic("no return value specified for GetUpkeep")
	}

	var r0 keeper_registry_wrapper2_0.UpkeepInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts, *big.Int) (keeper_registry_wrapper2_0.UpkeepInfo, error)); ok {
		return rf(opts, id)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts, *big.Int) keeper_registry_wrapper2_0.UpkeepInfo); ok {
		r0 = rf(opts, id)
	} else {
		r0 = ret.Get(0).(keeper_registry_wrapper2_0.UpkeepInfo)
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts, *big.Int) error); ok {
		r1 = rf(opts, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ParseLog provides a mock function with given fields: log
func (_m *Registry) ParseLog(log types.Log) (generated.AbigenLog, error) {
	ret := _m.Called(log)

	if len(ret) == 0 {
		panic("no return value specified for ParseLog")
	}

	var r0 generated.AbigenLog
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Log) (generated.AbigenLog, error)); ok {
		return rf(log)
	}
	if rf, ok := ret.Get(0).(func(types.Log) generated.AbigenLog); ok {
		r0 = rf(log)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(generated.AbigenLog)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Log) error); ok {
		r1 = rf(log)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRegistry creates a new instance of Registry. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRegistry(t interface {
	mock.TestingT
	Cleanup(func())
}) *Registry {
	mock := &Registry{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
