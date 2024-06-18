// Code generated by MockGen. DO NOT EDIT.
// Source: .\service\service.go
//
// Generated by this command:
//
//	mockgen -source .\service\service.go -destination service/svcmock/service_mock.go -package svcmock Servicer
//

// Package svcmock is a generated GoMock package.
package svcmock

import (
	repository "gounread/repository"
	service "gounread/service"
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockMessage is a mock of Message interface.
type MockMessage struct {
	ctrl     *gomock.Controller
	recorder *MockMessageMockRecorder
}

// MockMessageMockRecorder is the mock recorder for MockMessage.
type MockMessageMockRecorder struct {
	mock *MockMessage
}

// NewMockMessage creates a new mock instance.
func NewMockMessage(ctrl *gomock.Controller) *MockMessage {
	mock := &MockMessage{ctrl: ctrl}
	mock.recorder = &MockMessageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMessage) EXPECT() *MockMessageMockRecorder {
	return m.recorder
}

// CreateMessage mocks base method.
func (m *MockMessage) CreateMessage(param *repository.CreateMessageParam) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMessage", param)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateMessage indicates an expected call of CreateMessage.
func (mr *MockMessageMockRecorder) CreateMessage(param any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMessage", reflect.TypeOf((*MockMessage)(nil).CreateMessage), param)
}

// GetAllRoomsReadMessageTime mocks base method.
func (m *MockMessage) GetAllRoomsReadMessageTime(userID string) []*repository.GetAllRoomsReadMessageTimeResult {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllRoomsReadMessageTime", userID)
	ret0, _ := ret[0].([]*repository.GetAllRoomsReadMessageTimeResult)
	return ret0
}

// GetAllRoomsReadMessageTime indicates an expected call of GetAllRoomsReadMessageTime.
func (mr *MockMessageMockRecorder) GetAllRoomsReadMessageTime(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllRoomsReadMessageTime", reflect.TypeOf((*MockMessage)(nil).GetAllRoomsReadMessageTime), userID)
}

// GetMessageReadTime mocks base method.
func (m *MockMessage) GetMessageReadTime(roomID, userID string) (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessageReadTime", roomID, userID)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessageReadTime indicates an expected call of GetMessageReadTime.
func (mr *MockMessageMockRecorder) GetMessageReadTime(roomID, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessageReadTime", reflect.TypeOf((*MockMessage)(nil).GetMessageReadTime), roomID, userID)
}

// GetRecentMessageByRoomID mocks base method.
func (m *MockMessage) GetRecentMessageByRoomID(roomID string) (*repository.GetRecentMessageByRoomIDResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecentMessageByRoomID", roomID)
	ret0, _ := ret[0].(*repository.GetRecentMessageByRoomIDResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecentMessageByRoomID indicates an expected call of GetRecentMessageByRoomID.
func (mr *MockMessageMockRecorder) GetRecentMessageByRoomID(roomID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecentMessageByRoomID", reflect.TypeOf((*MockMessage)(nil).GetRecentMessageByRoomID), roomID)
}

// GetRecentMessages mocks base method.
func (m *MockMessage) GetRecentMessages(roomID string, limit int) []*repository.GetRecentMessagesResult {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecentMessages", roomID, limit)
	ret0, _ := ret[0].([]*repository.GetRecentMessagesResult)
	return ret0
}

// GetRecentMessages indicates an expected call of GetRecentMessages.
func (mr *MockMessageMockRecorder) GetRecentMessages(roomID, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecentMessages", reflect.TypeOf((*MockMessage)(nil).GetRecentMessages), roomID, limit)
}

// GetRoomsUnreadMessageCount mocks base method.
func (m *MockMessage) GetRoomsUnreadMessageCount(times []*repository.GetAllRoomsReadMessageTimeResult) ([]*service.GetRoomsUnreadMessageCountResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRoomsUnreadMessageCount", times)
	ret0, _ := ret[0].([]*service.GetRoomsUnreadMessageCountResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRoomsUnreadMessageCount indicates an expected call of GetRoomsUnreadMessageCount.
func (mr *MockMessageMockRecorder) GetRoomsUnreadMessageCount(times any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRoomsUnreadMessageCount", reflect.TypeOf((*MockMessage)(nil).GetRoomsUnreadMessageCount), times)
}

// GetUnreadMessageCount mocks base method.
func (m *MockMessage) GetUnreadMessageCount(roomID string, t time.Time) (*int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnreadMessageCount", roomID, t)
	ret0, _ := ret[0].(*int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnreadMessageCount indicates an expected call of GetUnreadMessageCount.
func (mr *MockMessageMockRecorder) GetUnreadMessageCount(roomID, t any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnreadMessageCount", reflect.TypeOf((*MockMessage)(nil).GetUnreadMessageCount), roomID, t)
}

// GetUnreadMessages mocks base method.
func (m *MockMessage) GetUnreadMessages(roomID string, start, end time.Time) []*repository.GetMessagesByRoomIDAndTimeResult {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnreadMessages", roomID, start, end)
	ret0, _ := ret[0].([]*repository.GetMessagesByRoomIDAndTimeResult)
	return ret0
}

// GetUnreadMessages indicates an expected call of GetUnreadMessages.
func (mr *MockMessageMockRecorder) GetUnreadMessages(roomID, start, end any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnreadMessages", reflect.TypeOf((*MockMessage)(nil).GetUnreadMessages), roomID, start, end)
}

// ReadMessage mocks base method.
func (m *MockMessage) ReadMessage(roomID, userID string) (time.Time, time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadMessage", roomID, userID)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(time.Time)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ReadMessage indicates an expected call of ReadMessage.
func (mr *MockMessageMockRecorder) ReadMessage(roomID, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadMessage", reflect.TypeOf((*MockMessage)(nil).ReadMessage), roomID, userID)
}

// UpdateMessageReadTime mocks base method.
func (m *MockMessage) UpdateMessageReadTime(roomID, userID string, t time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMessageReadTime", roomID, userID, t)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMessageReadTime indicates an expected call of UpdateMessageReadTime.
func (mr *MockMessageMockRecorder) UpdateMessageReadTime(roomID, userID, t any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMessageReadTime", reflect.TypeOf((*MockMessage)(nil).UpdateMessageReadTime), roomID, userID, t)
}

// UpdateRecentMessage mocks base method.
func (m *MockMessage) UpdateRecentMessage(roomID, recentMessage string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRecentMessage", roomID, recentMessage)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRecentMessage indicates an expected call of UpdateRecentMessage.
func (mr *MockMessageMockRecorder) UpdateRecentMessage(roomID, recentMessage any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRecentMessage", reflect.TypeOf((*MockMessage)(nil).UpdateRecentMessage), roomID, recentMessage)
}

// MockRoom is a mock of Room interface.
type MockRoom struct {
	ctrl     *gomock.Controller
	recorder *MockRoomMockRecorder
}

// MockRoomMockRecorder is the mock recorder for MockRoom.
type MockRoomMockRecorder struct {
	mock *MockRoom
}

// NewMockRoom creates a new mock instance.
func NewMockRoom(ctrl *gomock.Controller) *MockRoom {
	mock := &MockRoom{ctrl: ctrl}
	mock.recorder = &MockRoomMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRoom) EXPECT() *MockRoomMockRecorder {
	return m.recorder
}

// CreateRoom mocks base method.
func (m *MockRoom) CreateRoom(users []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRoom", users)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateRoom indicates an expected call of CreateRoom.
func (mr *MockRoomMockRecorder) CreateRoom(users any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRoom", reflect.TypeOf((*MockRoom)(nil).CreateRoom), users)
}

// GetRoomsByUserID mocks base method.
func (m *MockRoom) GetRoomsByUserID(userID string) ([]*repository.GetRoomsByUserIDResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRoomsByUserID", userID)
	ret0, _ := ret[0].([]*repository.GetRoomsByUserIDResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRoomsByUserID indicates an expected call of GetRoomsByUserID.
func (mr *MockRoomMockRecorder) GetRoomsByUserID(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRoomsByUserID", reflect.TypeOf((*MockRoom)(nil).GetRoomsByUserID), userID)
}

// GetUsersByRoomID mocks base method.
func (m *MockRoom) GetUsersByRoomID(roomID string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersByRoomID", roomID)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersByRoomID indicates an expected call of GetUsersByRoomID.
func (mr *MockRoomMockRecorder) GetUsersByRoomID(roomID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersByRoomID", reflect.TypeOf((*MockRoom)(nil).GetUsersByRoomID), roomID)
}

// MockServicer is a mock of Servicer interface.
type MockServicer struct {
	ctrl     *gomock.Controller
	recorder *MockServicerMockRecorder
}

// MockServicerMockRecorder is the mock recorder for MockServicer.
type MockServicerMockRecorder struct {
	mock *MockServicer
}

// NewMockServicer creates a new mock instance.
func NewMockServicer(ctrl *gomock.Controller) *MockServicer {
	mock := &MockServicer{ctrl: ctrl}
	mock.recorder = &MockServicerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServicer) EXPECT() *MockServicerMockRecorder {
	return m.recorder
}

// CreateMessage mocks base method.
func (m *MockServicer) CreateMessage(param *repository.CreateMessageParam) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMessage", param)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateMessage indicates an expected call of CreateMessage.
func (mr *MockServicerMockRecorder) CreateMessage(param any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMessage", reflect.TypeOf((*MockServicer)(nil).CreateMessage), param)
}

// CreateRoom mocks base method.
func (m *MockServicer) CreateRoom(users []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRoom", users)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateRoom indicates an expected call of CreateRoom.
func (mr *MockServicerMockRecorder) CreateRoom(users any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRoom", reflect.TypeOf((*MockServicer)(nil).CreateRoom), users)
}

// GetAllRoomsReadMessageTime mocks base method.
func (m *MockServicer) GetAllRoomsReadMessageTime(userID string) []*repository.GetAllRoomsReadMessageTimeResult {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllRoomsReadMessageTime", userID)
	ret0, _ := ret[0].([]*repository.GetAllRoomsReadMessageTimeResult)
	return ret0
}

// GetAllRoomsReadMessageTime indicates an expected call of GetAllRoomsReadMessageTime.
func (mr *MockServicerMockRecorder) GetAllRoomsReadMessageTime(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllRoomsReadMessageTime", reflect.TypeOf((*MockServicer)(nil).GetAllRoomsReadMessageTime), userID)
}

// GetMessageReadTime mocks base method.
func (m *MockServicer) GetMessageReadTime(roomID, userID string) (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessageReadTime", roomID, userID)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessageReadTime indicates an expected call of GetMessageReadTime.
func (mr *MockServicerMockRecorder) GetMessageReadTime(roomID, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessageReadTime", reflect.TypeOf((*MockServicer)(nil).GetMessageReadTime), roomID, userID)
}

// GetRecentMessageByRoomID mocks base method.
func (m *MockServicer) GetRecentMessageByRoomID(roomID string) (*repository.GetRecentMessageByRoomIDResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecentMessageByRoomID", roomID)
	ret0, _ := ret[0].(*repository.GetRecentMessageByRoomIDResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecentMessageByRoomID indicates an expected call of GetRecentMessageByRoomID.
func (mr *MockServicerMockRecorder) GetRecentMessageByRoomID(roomID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecentMessageByRoomID", reflect.TypeOf((*MockServicer)(nil).GetRecentMessageByRoomID), roomID)
}

// GetRecentMessages mocks base method.
func (m *MockServicer) GetRecentMessages(roomID string, limit int) []*repository.GetRecentMessagesResult {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecentMessages", roomID, limit)
	ret0, _ := ret[0].([]*repository.GetRecentMessagesResult)
	return ret0
}

// GetRecentMessages indicates an expected call of GetRecentMessages.
func (mr *MockServicerMockRecorder) GetRecentMessages(roomID, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecentMessages", reflect.TypeOf((*MockServicer)(nil).GetRecentMessages), roomID, limit)
}

// GetRoomsByUserID mocks base method.
func (m *MockServicer) GetRoomsByUserID(userID string) ([]*repository.GetRoomsByUserIDResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRoomsByUserID", userID)
	ret0, _ := ret[0].([]*repository.GetRoomsByUserIDResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRoomsByUserID indicates an expected call of GetRoomsByUserID.
func (mr *MockServicerMockRecorder) GetRoomsByUserID(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRoomsByUserID", reflect.TypeOf((*MockServicer)(nil).GetRoomsByUserID), userID)
}

// GetRoomsUnreadMessageCount mocks base method.
func (m *MockServicer) GetRoomsUnreadMessageCount(times []*repository.GetAllRoomsReadMessageTimeResult) ([]*service.GetRoomsUnreadMessageCountResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRoomsUnreadMessageCount", times)
	ret0, _ := ret[0].([]*service.GetRoomsUnreadMessageCountResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRoomsUnreadMessageCount indicates an expected call of GetRoomsUnreadMessageCount.
func (mr *MockServicerMockRecorder) GetRoomsUnreadMessageCount(times any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRoomsUnreadMessageCount", reflect.TypeOf((*MockServicer)(nil).GetRoomsUnreadMessageCount), times)
}

// GetUnreadMessageCount mocks base method.
func (m *MockServicer) GetUnreadMessageCount(roomID string, t time.Time) (*int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnreadMessageCount", roomID, t)
	ret0, _ := ret[0].(*int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnreadMessageCount indicates an expected call of GetUnreadMessageCount.
func (mr *MockServicerMockRecorder) GetUnreadMessageCount(roomID, t any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnreadMessageCount", reflect.TypeOf((*MockServicer)(nil).GetUnreadMessageCount), roomID, t)
}

// GetUnreadMessages mocks base method.
func (m *MockServicer) GetUnreadMessages(roomID string, start, end time.Time) []*repository.GetMessagesByRoomIDAndTimeResult {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnreadMessages", roomID, start, end)
	ret0, _ := ret[0].([]*repository.GetMessagesByRoomIDAndTimeResult)
	return ret0
}

// GetUnreadMessages indicates an expected call of GetUnreadMessages.
func (mr *MockServicerMockRecorder) GetUnreadMessages(roomID, start, end any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnreadMessages", reflect.TypeOf((*MockServicer)(nil).GetUnreadMessages), roomID, start, end)
}

// GetUsersByRoomID mocks base method.
func (m *MockServicer) GetUsersByRoomID(roomID string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersByRoomID", roomID)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersByRoomID indicates an expected call of GetUsersByRoomID.
func (mr *MockServicerMockRecorder) GetUsersByRoomID(roomID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersByRoomID", reflect.TypeOf((*MockServicer)(nil).GetUsersByRoomID), roomID)
}

// ReadMessage mocks base method.
func (m *MockServicer) ReadMessage(roomID, userID string) (time.Time, time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadMessage", roomID, userID)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(time.Time)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ReadMessage indicates an expected call of ReadMessage.
func (mr *MockServicerMockRecorder) ReadMessage(roomID, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadMessage", reflect.TypeOf((*MockServicer)(nil).ReadMessage), roomID, userID)
}

// UpdateMessageReadTime mocks base method.
func (m *MockServicer) UpdateMessageReadTime(roomID, userID string, t time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMessageReadTime", roomID, userID, t)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMessageReadTime indicates an expected call of UpdateMessageReadTime.
func (mr *MockServicerMockRecorder) UpdateMessageReadTime(roomID, userID, t any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMessageReadTime", reflect.TypeOf((*MockServicer)(nil).UpdateMessageReadTime), roomID, userID, t)
}

// UpdateRecentMessage mocks base method.
func (m *MockServicer) UpdateRecentMessage(roomID, recentMessage string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRecentMessage", roomID, recentMessage)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRecentMessage indicates an expected call of UpdateRecentMessage.
func (mr *MockServicerMockRecorder) UpdateRecentMessage(roomID, recentMessage any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRecentMessage", reflect.TypeOf((*MockServicer)(nil).UpdateRecentMessage), roomID, recentMessage)
}
