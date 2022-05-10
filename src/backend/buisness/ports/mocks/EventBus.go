// Code generated by mockery v2.12.0. DO NOT EDIT.

package mocks

import (
	bus "iu7-2022-sd-labs/buisness/ports/bus"

	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// EventBus is an autogenerated mock type for the EventBus type
type EventBus struct {
	mock.Mock
}

// Notify provides a mock function with given fields: event
func (_m *EventBus) Notify(event bus.Event) {
	_m.Called(event)
}

// NewEventBus creates a new instance of EventBus. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewEventBus(t testing.TB) *EventBus {
	mock := &EventBus{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}