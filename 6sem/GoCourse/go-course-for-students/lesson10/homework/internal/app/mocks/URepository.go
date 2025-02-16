// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	user "homework10/internal/entities/user"

	mock "github.com/stretchr/testify/mock"
)

// URepository is an autogenerated mock type for the URepository type
type URepository struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: _a0
func (_m *URepository) CreateUser(_a0 user.User) (user.User, error) {
	ret := _m.Called(_a0)

	var r0 user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(user.User) (user.User, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(user.User) user.User); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(user.User)
	}

	if rf, ok := ret.Get(1).(func(user.User) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteUser provides a mock function with given fields: _a0
func (_m *URepository) DeleteUser(_a0 user.User) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(user.User) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetCurrentUser provides a mock function with given fields: id
func (_m *URepository) GetCurrentUser(id int64) (user.User, error) {
	ret := _m.Called(id)

	var r0 user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(int64) (user.User, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int64) user.User); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(user.User)
	}

	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUser provides a mock function with given fields: _a0
func (_m *URepository) UpdateUser(_a0 user.User) (user.User, error) {
	ret := _m.Called(_a0)

	var r0 user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(user.User) (user.User, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(user.User) user.User); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(user.User)
	}

	if rf, ok := ret.Get(1).(func(user.User) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewURepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewURepository creates a new instance of URepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewURepository(t mockConstructorTestingTNewURepository) *URepository {
	mock := &URepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
