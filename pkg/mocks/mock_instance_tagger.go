// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ouzi-dev/node-tagger/pkg/aws (interfaces: NodeTagger)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1 "k8s.io/api/core/v1"
)

// MockNodeTagger is a mock of NodeTagger interface
type MockNodeTagger struct {
	ctrl     *gomock.Controller
	recorder *MockNodeTaggerMockRecorder
}

// MockNodeTaggerMockRecorder is the mock recorder for MockNodeTagger
type MockNodeTaggerMockRecorder struct {
	mock *MockNodeTagger
}

// NewMockNodeTagger creates a new mock instance
func NewMockNodeTagger(ctrl *gomock.Controller) *MockNodeTagger {
	mock := &MockNodeTagger{ctrl: ctrl}
	mock.recorder = &MockNodeTaggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNodeTagger) EXPECT() *MockNodeTaggerMockRecorder {
	return m.recorder
}

// EnsureInstanceNodeHasTags mocks base method
func (m *MockNodeTagger) EnsureInstanceNodeHasTags(arg0 *v1.Node, arg1 map[string]string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnsureInstanceNodeHasTags", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// EnsureInstanceNodeHasTags indicates an expected call of EnsureInstanceNodeHasTags
func (mr *MockNodeTaggerMockRecorder) EnsureInstanceNodeHasTags(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnsureInstanceNodeHasTags", reflect.TypeOf((*MockNodeTagger)(nil).EnsureInstanceNodeHasTags), arg0, arg1)
}