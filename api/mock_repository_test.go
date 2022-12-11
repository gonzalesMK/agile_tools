// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package main is a generated GoMock package.
package main

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRepoInterface is a mock of RepoInterface interface.
type MockRepoInterface struct {
	ctrl     *gomock.Controller
	recorder *MockRepoInterfaceMockRecorder
}

// MockRepoInterfaceMockRecorder is the mock recorder for MockRepoInterface.
type MockRepoInterfaceMockRecorder struct {
	mock *MockRepoInterface
}

// NewMockRepoInterface creates a new mock instance.
func NewMockRepoInterface(ctrl *gomock.Controller) *MockRepoInterface {
	mock := &MockRepoInterface{ctrl: ctrl}
	mock.recorder = &MockRepoInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepoInterface) EXPECT() *MockRepoInterfaceMockRecorder {
	return m.recorder
}

// ClearPlayerStatusInRoom mocks base method.
func (m *MockRepoInterface) ClearPlayerStatusInRoom(roomID uint, statusID int8) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClearPlayerStatusInRoom", roomID, statusID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ClearPlayerStatusInRoom indicates an expected call of ClearPlayerStatusInRoom.
func (mr *MockRepoInterfaceMockRecorder) ClearPlayerStatusInRoom(roomID, statusID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearPlayerStatusInRoom", reflect.TypeOf((*MockRepoInterface)(nil).ClearPlayerStatusInRoom), roomID, statusID)
}

// DeleteById mocks base method.
func (m *MockRepoInterface) DeleteById(model interface{}, id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteById", model, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteById indicates an expected call of DeleteById.
func (mr *MockRepoInterfaceMockRecorder) DeleteById(model, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteById", reflect.TypeOf((*MockRepoInterface)(nil).DeleteById), model, id)
}

// GetOneById mocks base method.
func (m *MockRepoInterface) GetOneById(model interface{}, id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOneById", model, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetOneById indicates an expected call of GetOneById.
func (mr *MockRepoInterfaceMockRecorder) GetOneById(model, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOneById", reflect.TypeOf((*MockRepoInterface)(nil).GetOneById), model, id)
}

// GetPlayersFromRoom mocks base method.
func (m *MockRepoInterface) GetPlayersFromRoom(arg0 uint) ([]Users, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPlayersFromRoom", arg0)
	ret0, _ := ret[0].([]Users)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPlayersFromRoom indicates an expected call of GetPlayersFromRoom.
func (mr *MockRepoInterfaceMockRecorder) GetPlayersFromRoom(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPlayersFromRoom", reflect.TypeOf((*MockRepoInterface)(nil).GetPlayersFromRoom), arg0)
}

// Save mocks base method.
func (m *MockRepoInterface) Save(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockRepoInterfaceMockRecorder) Save(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockRepoInterface)(nil).Save), arg0)
}

// UpdateFieldById mocks base method.
func (m *MockRepoInterface) UpdateFieldById(id uint, content interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFieldById", id, content)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateFieldById indicates an expected call of UpdateFieldById.
func (mr *MockRepoInterfaceMockRecorder) UpdateFieldById(id, content interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFieldById", reflect.TypeOf((*MockRepoInterface)(nil).UpdateFieldById), id, content)
}
