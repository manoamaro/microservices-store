// Code generated by mockery v2.23.4. DO NOT EDIT.

package mocks

import (
	service "github.com/manoamaro/microservices-store/order_service/internal/service"
	mock "github.com/stretchr/testify/mock"
)

// ProductService is an autogenerated mock type for the ProductService type
type ProductService struct {
	mock.Mock
}

// Get provides a mock function with given fields: productId
func (_m *ProductService) Get(productId string) (service.ProductDTO, error) {
	ret := _m.Called(productId)

	var r0 service.ProductDTO
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (service.ProductDTO, error)); ok {
		return rf(productId)
	}
	if rf, ok := ret.Get(0).(func(string) service.ProductDTO); ok {
		r0 = rf(productId)
	} else {
		r0 = ret.Get(0).(service.ProductDTO)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(productId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewProductService creates a new instance of ProductService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProductService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProductService {
	mock := &ProductService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
