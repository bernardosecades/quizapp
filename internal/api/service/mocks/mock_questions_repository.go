// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/quizapp/internal/api/service (interfaces: QuestionsRepository)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/quizapp/internal/api/model"
)

// MockQuestionsRepository is a mock of QuestionsRepository interface.
type MockQuestionsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockQuestionsRepositoryMockRecorder
}

// MockQuestionsRepositoryMockRecorder is the mock recorder for MockQuestionsRepository.
type MockQuestionsRepositoryMockRecorder struct {
	mock *MockQuestionsRepository
}

// NewMockQuestionsRepository creates a new mock instance.
func NewMockQuestionsRepository(ctrl *gomock.Controller) *MockQuestionsRepository {
	mock := &MockQuestionsRepository{ctrl: ctrl}
	mock.recorder = &MockQuestionsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQuestionsRepository) EXPECT() *MockQuestionsRepositoryMockRecorder {
	return m.recorder
}

// GetQuestions mocks base method.
func (m *MockQuestionsRepository) GetQuestions(arg0 context.Context) ([]*model.Question, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQuestions", arg0)
	ret0, _ := ret[0].([]*model.Question)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQuestions indicates an expected call of GetQuestions.
func (mr *MockQuestionsRepositoryMockRecorder) GetQuestions(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQuestions", reflect.TypeOf((*MockQuestionsRepository)(nil).GetQuestions), arg0)
}