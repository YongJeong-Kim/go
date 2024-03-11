// Code generated by MockGen. DO NOT EDIT.
// Source: service/service.go
//
// Generated by this command:
//
//	mockgen -package svcmock -destination service/mock/service.go -source service/service.go Servicer
//

// Package svcmock is a generated GoMock package.
package svcmock

import (
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockAccountServicer is a mock of AccountServicer interface.
type MockAccountServicer struct {
	ctrl     *gomock.Controller
	recorder *MockAccountServicerMockRecorder
}

// MockAccountServicerMockRecorder is the mock recorder for MockAccountServicer.
type MockAccountServicerMockRecorder struct {
	mock *MockAccountServicer
}

// NewMockAccountServicer creates a new mock instance.
func NewMockAccountServicer(ctrl *gomock.Controller) *MockAccountServicer {
	mock := &MockAccountServicer{ctrl: ctrl}
	mock.recorder = &MockAccountServicerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountServicer) EXPECT() *MockAccountServicerMockRecorder {
	return m.recorder
}

// Login mocks base method.
func (m *MockAccountServicer) Login(username, password string, duration time.Duration) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", username, password, duration)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockAccountServicerMockRecorder) Login(username, password, duration any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAccountServicer)(nil).Login), username, password, duration)
}
