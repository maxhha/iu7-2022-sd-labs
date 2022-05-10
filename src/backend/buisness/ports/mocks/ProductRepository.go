// Code generated by mockery v2.12.0. DO NOT EDIT.

package mocks

import (
	entities "iu7-2022-sd-labs/buisness/entities"

	mock "github.com/stretchr/testify/mock"

	repositories "iu7-2022-sd-labs/buisness/ports/repositories"

	testing "testing"
)

// ProductRepository is an autogenerated mock type for the ProductRepository type
type ProductRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: room
func (_m *ProductRepository) Create(room *entities.Product) error {
	ret := _m.Called(room)

	var r0 error
	if rf, ok := ret.Get(0).(func(*entities.Product) error); ok {
		r0 = rf(room)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: id
func (_m *ProductRepository) Delete(id string) (entities.Product, error) {
	ret := _m.Called(id)

	var r0 entities.Product
	if rf, ok := ret.Get(0).(func(string) entities.Product); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(entities.Product)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Find provides a mock function with given fields: params
func (_m *ProductRepository) Find(params *repositories.ProductFindParams) ([]entities.Product, error) {
	ret := _m.Called(params)

	var r0 []entities.Product
	if rf, ok := ret.Get(0).(func(*repositories.ProductFindParams) []entities.Product); ok {
		r0 = rf(params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entities.Product)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*repositories.ProductFindParams) error); ok {
		r1 = rf(params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: id
func (_m *ProductRepository) Get(id string) (entities.Product, error) {
	ret := _m.Called(id)

	var r0 entities.Product
	if rf, ok := ret.Get(0).(func(string) entities.Product); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(entities.Product)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: id, updateFn
func (_m *ProductRepository) Update(id string, updateFn func(*entities.Product) error) (entities.Product, error) {
	ret := _m.Called(id, updateFn)

	var r0 entities.Product
	if rf, ok := ret.Get(0).(func(string, func(*entities.Product) error) entities.Product); ok {
		r0 = rf(id, updateFn)
	} else {
		r0 = ret.Get(0).(entities.Product)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, func(*entities.Product) error) error); ok {
		r1 = rf(id, updateFn)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewProductRepository creates a new instance of ProductRepository. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewProductRepository(t testing.TB) *ProductRepository {
	mock := &ProductRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
