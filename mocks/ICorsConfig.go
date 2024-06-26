// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	gin "github.com/gin-gonic/gin"

	mock "github.com/stretchr/testify/mock"
)

// ICorsConfig is an autogenerated mock type for the ICorsConfig type
type ICorsConfig struct {
	mock.Mock
}

// CorsConfig provides a mock function with given fields:
func (_m *ICorsConfig) CorsConfig() gin.HandlerFunc {
	ret := _m.Called()

	var r0 gin.HandlerFunc
	if rf, ok := ret.Get(0).(func() gin.HandlerFunc); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gin.HandlerFunc)
		}
	}

	return r0
}

// NewICorsConfig creates a new instance of ICorsConfig. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewICorsConfig(t interface {
	mock.TestingT
	Cleanup(func())
}) *ICorsConfig {
	mock := &ICorsConfig{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
