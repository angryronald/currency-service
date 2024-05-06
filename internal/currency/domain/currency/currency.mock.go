// Code generated by MockGen. DO NOT EDIT.
// Source: internal/currency/domain/currency/currency.interface.go

// Package currency is a generated GoMock package.
package currency

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	model "github.com/angryronald/currency-service/internal/currency/domain/model"
)

// MockCurrencyServiceInterface is a mock of CurrencyServiceInterface interface.
type MockCurrencyServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockCurrencyServiceInterfaceMockRecorder
}

// MockCurrencyServiceInterfaceMockRecorder is the mock recorder for MockCurrencyServiceInterface.
type MockCurrencyServiceInterfaceMockRecorder struct {
	mock *MockCurrencyServiceInterface
}

// NewMockCurrencyServiceInterface creates a new mock instance.
func NewMockCurrencyServiceInterface(ctrl *gomock.Controller) *MockCurrencyServiceInterface {
	mock := &MockCurrencyServiceInterface{ctrl: ctrl}
	mock.recorder = &MockCurrencyServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCurrencyServiceInterface) EXPECT() *MockCurrencyServiceInterfaceMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockCurrencyServiceInterface) Add(ctx context.Context, currency *model.CurrencyDomainModel) (*model.CurrencyDomainModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, currency)
	ret0, _ := ret[0].(*model.CurrencyDomainModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockCurrencyServiceInterfaceMockRecorder) Add(ctx, currency interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockCurrencyServiceInterface)(nil).Add), ctx, currency)
}

// GetByCode mocks base method.
func (m *MockCurrencyServiceInterface) GetByCode(ctx context.Context, currencyCode string) (*model.CurrencyDomainModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByCode", ctx, currencyCode)
	ret0, _ := ret[0].(*model.CurrencyDomainModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByCode indicates an expected call of GetByCode.
func (mr *MockCurrencyServiceInterfaceMockRecorder) GetByCode(ctx, currencyCode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByCode", reflect.TypeOf((*MockCurrencyServiceInterface)(nil).GetByCode), ctx, currencyCode)
}

// List mocks base method.
func (m *MockCurrencyServiceInterface) List(ctx context.Context) ([]*model.CurrencyDomainModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx)
	ret0, _ := ret[0].([]*model.CurrencyDomainModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockCurrencyServiceInterfaceMockRecorder) List(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockCurrencyServiceInterface)(nil).List), ctx)
}

// MultipleAdd mocks base method.
func (m *MockCurrencyServiceInterface) MultipleAdd(ctx context.Context, currencies []*model.CurrencyDomainModel) ([]*model.CurrencyDomainModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MultipleAdd", ctx, currencies)
	ret0, _ := ret[0].([]*model.CurrencyDomainModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MultipleAdd indicates an expected call of MultipleAdd.
func (mr *MockCurrencyServiceInterfaceMockRecorder) MultipleAdd(ctx, currencies interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MultipleAdd", reflect.TypeOf((*MockCurrencyServiceInterface)(nil).MultipleAdd), ctx, currencies)
}
