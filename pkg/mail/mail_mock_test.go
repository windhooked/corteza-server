// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/mail/mail.go

// Package mail is a generated GoMock package.
package mail

import (
	gomock "github.com/golang/mock/gomock"
	mail_v2 "gopkg.in/mail.v2"
	reflect "reflect"
)

// MockDialer is a mock of Dialer interface
type MockDialer struct {
	ctrl     *gomock.Controller
	recorder *MockDialerMockRecorder
}

// MockDialerMockRecorder is the mock recorder for MockDialer
type MockDialerMockRecorder struct {
	mock *MockDialer
}

// NewMockDialer creates a new mock instance
func NewMockDialer(ctrl *gomock.Controller) *MockDialer {
	mock := &MockDialer{ctrl: ctrl}
	mock.recorder = &MockDialerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDialer) EXPECT() *MockDialerMockRecorder {
	return m.recorder
}

// DialAndSend mocks base method
func (m *MockDialer) DialAndSend(arg0 ...*mail_v2.Message) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DialAndSend", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// DialAndSend indicates an expected call of DialAndSend
func (mr *MockDialerMockRecorder) DialAndSend(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DialAndSend", reflect.TypeOf((*MockDialer)(nil).DialAndSend), arg0...)
}
