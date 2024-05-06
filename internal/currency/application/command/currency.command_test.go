package command

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"gotest.tools/assert"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/angryronald/currency-service/internal/currency/application/model"
	"github.com/angryronald/currency-service/internal/currency/domain/currency"
	domainModel "github.com/angryronald/currency-service/internal/currency/domain/model"
)

var TEST_ID = uuid.MustParse("25ebc2e2-7fd2-47a0-bd2f-a8ec6e4d0163")

func TestCurrencyCommand_Add(t *testing.T) {
	// Initialize mocks
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock repository with dummy data or a test database
	mockService := currency.NewMockCurrencyServiceInterface(ctrl)

	command := NewCurrencyCommand(mockService)

	tests := []struct {
		name                       string
		currencyInput              *model.CurrencyApplicationModel
		expectedReturn             *model.CurrencyApplicationModel
		expectedErr                error
		mockReturn                 *domainModel.CurrencyDomainModel
		mockErr                    error
		serviceExpectedCalledCount int
	}{
		{
			name: "Success add currency",
			currencyInput: &model.CurrencyApplicationModel{
				Code: "USD",
				Name: "United States Dollar",
			},
			expectedReturn: &model.CurrencyApplicationModel{
				ID:   TEST_ID,
				Code: "USD",
				Name: "United States Dollar",
			},
			expectedErr: nil,
			mockReturn: &domainModel.CurrencyDomainModel{
				ID:   TEST_ID,
				Code: "USD",
				Name: "United States Dollar",
			},
			mockErr:                    nil,
			serviceExpectedCalledCount: 1,
		},
		{
			name: "Got error when adding currency data",
			currencyInput: &model.CurrencyApplicationModel{
				Code: "USD",
				Name: "United States Dollar",
			},
			expectedReturn:             nil,
			expectedErr:                errors.New("something went wrong"),
			mockReturn:                 nil,
			mockErr:                    errors.New("something went wrong"),
			serviceExpectedCalledCount: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockService.EXPECT().Add(gomock.Any(), test.currencyInput.ToDomainModel()).Return(test.mockReturn, test.mockErr).Times(test.serviceExpectedCalledCount)

			currencies, err := command.Add(context.TODO(), test.currencyInput)

			if !reflect.DeepEqual(test.expectedErr, err) {
				t.Errorf("Expected error: %v, got: %v", test.expectedErr, err)
			}
			assert.DeepEqual(t, test.expectedReturn, currencies)
		})
	}
}

func TestCurrencyCommand_MultipleAdd(t *testing.T) {
	// Initialize mocks
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock repository with dummy data or a test database
	mockService := currency.NewMockCurrencyServiceInterface(ctrl)

	command := NewCurrencyCommand(mockService)

	tests := []struct {
		name                       string
		currenciesInput            []*model.CurrencyApplicationModel
		expectedReturn             []*model.CurrencyApplicationModel
		expectedErr                error
		mockReturn                 []*domainModel.CurrencyDomainModel
		mockErr                    error
		serviceExpectedCalledCount int
	}{
		{
			name: "Success multiple add currencies",
			currenciesInput: []*model.CurrencyApplicationModel{
				{
					Code: "USD",
					Name: "United States Dollar",
				},
				{
					Code: "EUR",
					Name: "Euro",
				},
				{
					Code: "THB",
					Name: "Thailand Baht",
				},
				{
					Code: "IDR",
					Name: "Indonesian Rupiah",
				},
			},
			expectedReturn: []*model.CurrencyApplicationModel{
				{
					ID:   TEST_ID,
					Code: "USD",
					Name: "United States Dollar",
				},
				{
					ID:   TEST_ID,
					Code: "EUR",
					Name: "Euro",
				},
				{
					ID:   TEST_ID,
					Code: "THB",
					Name: "Thailand Baht",
				},
				{
					ID:   TEST_ID,
					Code: "IDR",
					Name: "Indonesian Rupiah",
				},
			},
			expectedErr: nil,
			mockReturn: []*domainModel.CurrencyDomainModel{
				{
					ID:   TEST_ID,
					Code: "USD",
					Name: "United States Dollar",
				},
				{
					ID:   TEST_ID,
					Code: "EUR",
					Name: "Euro",
				},
				{
					ID:   TEST_ID,
					Code: "THB",
					Name: "Thailand Baht",
				},
				{
					ID:   TEST_ID,
					Code: "IDR",
					Name: "Indonesian Rupiah",
				},
			},
			mockErr:                    nil,
			serviceExpectedCalledCount: 1,
		},
		{
			name: "Got error when multiple adding currencies",
			currenciesInput: []*model.CurrencyApplicationModel{
				{
					Code: "USD",
					Name: "United States Dollar",
				},
				{
					Code: "EUR",
					Name: "Euro",
				},
				{
					Code: "THB",
					Name: "Thailand Baht",
				},
				{
					Code: "IDR",
					Name: "Indonesian Rupiah",
				},
			},
			expectedReturn:             nil,
			expectedErr:                errors.New("something went wrong"),
			mockReturn:                 nil,
			mockErr:                    errors.New("something went wrong"),
			serviceExpectedCalledCount: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockService.EXPECT().MultipleAddOrUpdate(gomock.Any(), gomock.Any()).Return(test.mockReturn, test.mockErr).Times(test.serviceExpectedCalledCount)

			currencies, err := command.MultipleAddOrUpdate(context.TODO(), test.currenciesInput)

			if !reflect.DeepEqual(test.expectedErr, err) {
				t.Errorf("Expected error: %v, got: %v", test.expectedErr, err)
			}
			assert.DeepEqual(t, test.expectedReturn, currencies)
		})
	}
}
