// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	context "context"

	types "github.com/goplugin/pluginv3.0/v2/common/types"
	mock "github.com/stretchr/testify/mock"
)

// HeadTrackable is an autogenerated mock type for the HeadTrackable type
type HeadTrackable[H types.Head[BLOCK_HASH], BLOCK_HASH types.Hashable] struct {
	mock.Mock
}

// OnNewLongestChain provides a mock function with given fields: ctx, head
func (_m *HeadTrackable[H, BLOCK_HASH]) OnNewLongestChain(ctx context.Context, head H) {
	_m.Called(ctx, head)
}

// NewHeadTrackable creates a new instance of HeadTrackable. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHeadTrackable[H types.Head[BLOCK_HASH], BLOCK_HASH types.Hashable](t interface {
	mock.TestingT
	Cleanup(func())
}) *HeadTrackable[H, BLOCK_HASH] {
	mock := &HeadTrackable[H, BLOCK_HASH]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}