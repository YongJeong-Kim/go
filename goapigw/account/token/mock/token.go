// Code generated by MockGen. DO NOT EDIT.
// Source: token/token.go
//
// Generated by this command:
//
//	mockgen -package tkmock -destination token/mock/token.go -source token/token.go TokenMaker
//

// Package tkmock is a generated GoMock package.
package tkmock

import (
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockTokenMaker is a mock of TokenMaker interface.
type MockTokenMaker struct {
	ctrl     *gomock.Controller
	recorder *MockTokenMakerMockRecorder
}

// MockTokenMakerMockRecorder is the mock recorder for MockTokenMaker.
type MockTokenMakerMockRecorder struct {
	mock *MockTokenMaker
}

// NewMockTokenMaker creates a new mock instance.
func NewMockTokenMaker(ctrl *gomock.Controller) *MockTokenMaker {
	mock := &MockTokenMaker{ctrl: ctrl}
	mock.recorder = &MockTokenMakerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenMaker) EXPECT() *MockTokenMakerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockTokenMaker) Create(userID string, duration time.Duration) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", userID, duration)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockTokenMakerMockRecorder) Create(userID, duration any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTokenMaker)(nil).Create), userID, duration)
}
