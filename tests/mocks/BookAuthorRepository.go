// Code generated by mockery v2.47.0. DO NOT EDIT.

package mocks

import (
	models "github.com/joaooliveira247/go_olist_challenge/src/models"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// BookAuthorRepository is an autogenerated mock type for the BookAuthorRepository type
type BookAuthorRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: relationship
func (_m *BookAuthorRepository) Create(relationship *models.BookAuthor) error {
	ret := _m.Called(relationship)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.BookAuthor) error); ok {
		r0 = rf(relationship)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: bookID
func (_m *BookAuthorRepository) Delete(bookID uuid.UUID) error {
	ret := _m.Called(bookID)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) error); ok {
		r0 = rf(bookID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewBookAuthorRepository creates a new instance of BookAuthorRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBookAuthorRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *BookAuthorRepository {
	mock := &BookAuthorRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
