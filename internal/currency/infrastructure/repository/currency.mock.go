// Code generated by MockGen. DO NOT EDIT.
// Source: internal/currency/infrastructure/repository/currency.interface.go

// Package repository is a generated GoMock package.
package repository

import (
	context "context"
	reflect "reflect"

	model "github.com/angryronald/currency-service/internal/currency/infrastructure/repository/model"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockCurrencyRepository is a mock of CurrencyRepository interface.
type MockCurrencyRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCurrencyRepositoryMockRecorder
}

// MockCurrencyRepositoryMockRecorder is the mock recorder for MockCurrencyRepository.
type MockCurrencyRepositoryMockRecorder struct {
	mock *MockCurrencyRepository
}

// NewMockCurrencyRepository creates a new mock instance.
func NewMockCurrencyRepository(ctrl *gomock.Controller) *MockCurrencyRepository {
	mock := &MockCurrencyRepository{ctrl: ctrl}
	mock.recorder = &MockCurrencyRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCurrencyRepository) EXPECT() *MockCurrencyRepositoryMockRecorder {
	return m.recorder
}

// BulkInsert mocks base method.
func (m *MockCurrencyRepository) BulkInsert(ctx context.Context, currencies []*model.CurrencyRepositoryModel) ([]*model.CurrencyRepositoryModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BulkInsert", ctx, currencies)
	ret0, _ := ret[0].([]*model.CurrencyRepositoryModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BulkInsert indicates an expected call of BulkInsert.
func (mr *MockCurrencyRepositoryMockRecorder) BulkInsert(ctx, currencies interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BulkInsert", reflect.TypeOf((*MockCurrencyRepository)(nil).BulkInsert), ctx, currencies)
}

// FindAll mocks base method.
func (m *MockCurrencyRepository) FindAll(ctx context.Context) ([]*model.CurrencyRepositoryModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx)
	ret0, _ := ret[0].([]*model.CurrencyRepositoryModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockCurrencyRepositoryMockRecorder) FindAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockCurrencyRepository)(nil).FindAll), ctx)
}

// FindByCode mocks base method.
func (m *MockCurrencyRepository) FindByCode(ctx context.Context, currencyCode string) (*model.CurrencyRepositoryModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByCode", ctx, currencyCode)
	ret0, _ := ret[0].(*model.CurrencyRepositoryModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByCode indicates an expected call of FindByCode.
func (mr *MockCurrencyRepositoryMockRecorder) FindByCode(ctx, currencyCode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByCode", reflect.TypeOf((*MockCurrencyRepository)(nil).FindByCode), ctx, currencyCode)
}

// FindByID mocks base method.
func (m *MockCurrencyRepository) FindByID(ctx context.Context, ID uuid.UUID) (*model.CurrencyRepositoryModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, ID)
	ret0, _ := ret[0].(*model.CurrencyRepositoryModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockCurrencyRepositoryMockRecorder) FindByID(ctx, ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockCurrencyRepository)(nil).FindByID), ctx, ID)
}

// Insert mocks base method.
func (m *MockCurrencyRepository) Insert(ctx context.Context, currency *model.CurrencyRepositoryModel) (*model.CurrencyRepositoryModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, currency)
	ret0, _ := ret[0].(*model.CurrencyRepositoryModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockCurrencyRepositoryMockRecorder) Insert(ctx, currency interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockCurrencyRepository)(nil).Insert), ctx, currency)
}
