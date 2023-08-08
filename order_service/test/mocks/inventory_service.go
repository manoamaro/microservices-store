// Code generated by mockery v2.23.4. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// InventoryService is an autogenerated mock type for the InventoryService type
type InventoryService struct {
	mock.Mock
}

// Add provides a mock function with given fields: productId, amount
func (_m *InventoryService) Add(productId string, amount uint) (uint, error) {
	ret := _m.Called(productId, amount)

	var r0 uint
	var r1 error
	if rf, ok := ret.Get(0).(func(string, uint) (uint, error)); ok {
		return rf(productId, amount)
	}
	if rf, ok := ret.Get(0).(func(string, uint) uint); ok {
		r0 = rf(productId, amount)
	} else {
		r0 = ret.Get(0).(uint)
	}

	if rf, ok := ret.Get(1).(func(string, uint) error); ok {
		r1 = rf(productId, amount)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: productId
func (_m *InventoryService) Get(productId string) (uint, error) {
	ret := _m.Called(productId)

	var r0 uint
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (uint, error)); ok {
		return rf(productId)
	}
	if rf, ok := ret.Get(0).(func(string) uint); ok {
		r0 = rf(productId)
	} else {
		r0 = ret.Get(0).(uint)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(productId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Reserve provides a mock function with given fields: cartId, productId, amount
func (_m *InventoryService) Reserve(cartId string, productId string, amount uint) (uint, error) {
	ret := _m.Called(cartId, productId, amount)

	var r0 uint
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, uint) (uint, error)); ok {
		return rf(cartId, productId, amount)
	}
	if rf, ok := ret.Get(0).(func(string, string, uint) uint); ok {
		r0 = rf(cartId, productId, amount)
	} else {
		r0 = ret.Get(0).(uint)
	}

	if rf, ok := ret.Get(1).(func(string, string, uint) error); ok {
		r1 = rf(cartId, productId, amount)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Subtract provides a mock function with given fields: productId, amount
func (_m *InventoryService) Subtract(productId string, amount uint) (uint, error) {
	ret := _m.Called(productId, amount)

	var r0 uint
	var r1 error
	if rf, ok := ret.Get(0).(func(string, uint) (uint, error)); ok {
		return rf(productId, amount)
	}
	if rf, ok := ret.Get(0).(func(string, uint) uint); ok {
		r0 = rf(productId, amount)
	} else {
		r0 = ret.Get(0).(uint)
	}

	if rf, ok := ret.Get(1).(func(string, uint) error); ok {
		r1 = rf(productId, amount)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewInventoryService creates a new instance of InventoryService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewInventoryService(t interface {
	mock.TestingT
	Cleanup(func())
}) *InventoryService {
	mock := &InventoryService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
