// Code generated by MockGen. DO NOT EDIT.
// Source: service/friend.go
//
// Generated by this command:
//
//	mockgen -destination service/mock/friend/friend_mock.go -package mockfriendsvc -source service/friend.go Friender
//

// Package mockfriendsvc is a generated GoMock package.
package mockfriendsvc

import (
	context "context"
	repository "gorelationship/repository"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockFriender is a mock of Friender interface.
type MockFriender struct {
	ctrl     *gomock.Controller
	recorder *MockFrienderMockRecorder
}

// MockFrienderMockRecorder is the mock recorder for MockFriender.
type MockFrienderMockRecorder struct {
	mock *MockFriender
}

// NewMockFriender creates a new mock instance.
func NewMockFriender(ctrl *gomock.Controller) *MockFriender {
	mock := &MockFriender{ctrl: ctrl}
	mock.recorder = &MockFrienderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFriender) EXPECT() *MockFrienderMockRecorder {
	return m.recorder
}

// Accept mocks base method.
func (m *MockFriender) Accept(ctx context.Context, requestUserID, acceptUserID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Accept", ctx, requestUserID, acceptUserID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Accept indicates an expected call of Accept.
func (mr *MockFrienderMockRecorder) Accept(ctx, requestUserID, acceptUserID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Accept", reflect.TypeOf((*MockFriender)(nil).Accept), ctx, requestUserID, acceptUserID)
}

// Count mocks base method.
func (m *MockFriender) Count(ctx context.Context, userID string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, userID)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockFrienderMockRecorder) Count(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockFriender)(nil).Count), ctx, userID)
}

// FromRequestCount mocks base method.
func (m *MockFriender) FromRequestCount(ctx context.Context, userID string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FromRequestCount", ctx, userID)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FromRequestCount indicates an expected call of FromRequestCount.
func (mr *MockFrienderMockRecorder) FromRequestCount(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FromRequestCount", reflect.TypeOf((*MockFriender)(nil).FromRequestCount), ctx, userID)
}

// List mocks base method.
func (m *MockFriender) List(ctx context.Context, userID string) ([]repository.ListResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, userID)
	ret0, _ := ret[0].([]repository.ListResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockFrienderMockRecorder) List(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockFriender)(nil).List), ctx, userID)
}

// ListFromRequests mocks base method.
func (m *MockFriender) ListFromRequests(ctx context.Context, userID string) ([]repository.ListFromRequestsResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFromRequests", ctx, userID)
	ret0, _ := ret[0].([]repository.ListFromRequestsResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFromRequests indicates an expected call of ListFromRequests.
func (mr *MockFrienderMockRecorder) ListFromRequests(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFromRequests", reflect.TypeOf((*MockFriender)(nil).ListFromRequests), ctx, userID)
}

// ListMutuals mocks base method.
func (m *MockFriender) ListMutuals(ctx context.Context, userID1, userID2 string) ([]repository.ListMutualsResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListMutuals", ctx, userID1, userID2)
	ret0, _ := ret[0].([]repository.ListMutualsResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListMutuals indicates an expected call of ListMutuals.
func (mr *MockFrienderMockRecorder) ListMutuals(ctx, userID1, userID2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListMutuals", reflect.TypeOf((*MockFriender)(nil).ListMutuals), ctx, userID1, userID2)
}

// MutualCount mocks base method.
func (m *MockFriender) MutualCount(ctx context.Context, userID1, userID2 string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MutualCount", ctx, userID1, userID2)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MutualCount indicates an expected call of MutualCount.
func (mr *MockFrienderMockRecorder) MutualCount(ctx, userID1, userID2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MutualCount", reflect.TypeOf((*MockFriender)(nil).MutualCount), ctx, userID1, userID2)
}

// Request mocks base method.
func (m *MockFriender) Request(ctx context.Context, requestUserID, acceptUserID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Request", ctx, requestUserID, acceptUserID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Request indicates an expected call of Request.
func (mr *MockFrienderMockRecorder) Request(ctx, requestUserID, acceptUserID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Request", reflect.TypeOf((*MockFriender)(nil).Request), ctx, requestUserID, acceptUserID)
}
