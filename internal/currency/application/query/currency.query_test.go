package query

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"gotest.tools/assert"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/angryronald/currency-service/internal/currency/application/model"
	"github.com/angryronald/currency-service/internal/currency/constant"
	"github.com/angryronald/currency-service/internal/currency/domain/currency"
	domainModel "github.com/angryronald/currency-service/internal/currency/domain/model"
)

var TEST_ID = uuid.MustParse("25ebc2e2-7fd2-47a0-bd2f-a8ec6e4d0163")

func TestCurrencyService_List(t *testing.T) {
	// Initialize mocks
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock repository with dummy data or a test database
	mockService := currency.NewMockCurrencyServiceInterface(ctrl)
	mockFallbackService := currency.NewMockCurrencyServiceInterface(ctrl)
	log := logrus.New()

	query := NewCurrencyQuery(mockService, mockFallbackService, log)

	tests := []struct {
		name                                     string
		expectedReturn                           []*model.CurrencyApplicationModel
		expectedErr                              error
		mockReturn                               []*domainModel.CurrencyDomainModel
		mockErr                                  error
		serviceExpectedCalledCount               int
		fallbackMockReturn                       []*domainModel.CurrencyDomainModel
		fallbackMockErr                          error
		fallbackServiceExpectedCalledCount       int
		handleFallbackMockParam                  []*domainModel.CurrencyDomainModel
		handleFallbackMockErr                    error
		handleFallbackServiceExpectedCalledCount int
	}{
		{
			name: "Success get all of currencies",
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
			mockErr:                                  nil,
			serviceExpectedCalledCount:               1,
			fallbackMockReturn:                       nil,
			fallbackMockErr:                          nil,
			fallbackServiceExpectedCalledCount:       0,
			handleFallbackMockParam:                  nil,
			handleFallbackMockErr:                    nil,
			handleFallbackServiceExpectedCalledCount: 0,
		},
		{
			name:                                     "Got error when get currency data",
			expectedReturn:                           nil,
			expectedErr:                              errors.New("something went wrong"),
			mockReturn:                               nil,
			mockErr:                                  errors.New("something went wrong"),
			serviceExpectedCalledCount:               1,
			fallbackMockReturn:                       nil,
			fallbackMockErr:                          nil,
			fallbackServiceExpectedCalledCount:       0,
			handleFallbackMockParam:                  nil,
			handleFallbackMockErr:                    nil,
			handleFallbackServiceExpectedCalledCount: 0,
		},
		{
			name: "Currencies are not exists in memcached, success collect currencies from database",
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
			expectedErr:                nil,
			mockReturn:                 nil,
			mockErr:                    constant.ErrNotFound,
			serviceExpectedCalledCount: 1,
			fallbackMockReturn: []*domainModel.CurrencyDomainModel{
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
			fallbackMockErr:                    nil,
			fallbackServiceExpectedCalledCount: 1,
			handleFallbackMockParam: []*domainModel.CurrencyDomainModel{
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
			handleFallbackMockErr:                    nil,
			handleFallbackServiceExpectedCalledCount: 1,
		},
		{
			name:                                     "Currencies are not exists in memcached and database",
			expectedReturn:                           nil,
			expectedErr:                              constant.ErrNotFound,
			mockReturn:                               nil,
			mockErr:                                  constant.ErrNotFound,
			serviceExpectedCalledCount:               1,
			fallbackMockReturn:                       nil,
			fallbackMockErr:                          constant.ErrNotFound,
			fallbackServiceExpectedCalledCount:       1,
			handleFallbackMockParam:                  nil,
			handleFallbackMockErr:                    nil,
			handleFallbackServiceExpectedCalledCount: 0,
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockService.EXPECT().List(gomock.Any()).Return(test.mockReturn, test.mockErr).Times(test.serviceExpectedCalledCount)

			mockFallbackService.EXPECT().List(gomock.Any()).Return(test.fallbackMockReturn, test.fallbackMockErr).Times(test.fallbackServiceExpectedCalledCount)

			mockService.EXPECT().MultipleAddOrUpdate(gomock.Any(), test.handleFallbackMockParam).Return(test.handleFallbackMockParam, test.handleFallbackMockErr).Times(test.handleFallbackServiceExpectedCalledCount)

			currencies, err := query.List(context.TODO())

			if !reflect.DeepEqual(test.expectedErr, err) {
				t.Errorf("Expected error: %v, got: %v", test.expectedErr, err)
			}
			assert.DeepEqual(t, test.expectedReturn, currencies)
		})
	}
}

func TestCurrencyService_GetByCode(t *testing.T) {
	// Initialize mocks
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock repository with dummy data or a test database
	mockService := currency.NewMockCurrencyServiceInterface(ctrl)
	mockFallbackService := currency.NewMockCurrencyServiceInterface(ctrl)
	log := logrus.New()

	query := NewCurrencyQuery(mockService, mockFallbackService, log)

	tests := []struct {
		name                                     string
		currencyCode                             string
		expectedReturn                           *model.CurrencyApplicationModel
		expectedErr                              error
		mockReturn                               *domainModel.CurrencyDomainModel
		mockErr                                  error
		serviceExpectedCalledCount               int
		fallbackMockReturn                       *domainModel.CurrencyDomainModel
		fallbackMockErr                          error
		fallbackServiceExpectedCalledCount       int
		handleFallbackMockParam                  *domainModel.CurrencyDomainModel
		handleFallbackMockErr                    error
		handleFallbackServiceExpectedCalledCount int
	}{
		{
			name:         "Success get currency",
			currencyCode: "USD",
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
			mockErr:                                  nil,
			serviceExpectedCalledCount:               1,
			fallbackMockReturn:                       nil,
			fallbackMockErr:                          nil,
			fallbackServiceExpectedCalledCount:       0,
			handleFallbackMockParam:                  nil,
			handleFallbackMockErr:                    nil,
			handleFallbackServiceExpectedCalledCount: 0,
		},
		{
			name:                                     "Got error when get currency data",
			currencyCode:                             "USD",
			expectedReturn:                           nil,
			expectedErr:                              errors.New("something went wrong"),
			mockReturn:                               nil,
			mockErr:                                  errors.New("something went wrong"),
			serviceExpectedCalledCount:               1,
			fallbackMockReturn:                       nil,
			fallbackMockErr:                          nil,
			fallbackServiceExpectedCalledCount:       0,
			handleFallbackMockParam:                  nil,
			handleFallbackMockErr:                    nil,
			handleFallbackServiceExpectedCalledCount: 0,
		},
		{
			name:         "Currency is not exists in memcached, success collect from database",
			currencyCode: "USD",
			expectedReturn: &model.CurrencyApplicationModel{
				ID:   TEST_ID,
				Code: "USD",
				Name: "United States Dollar",
			},
			expectedErr:                nil,
			mockReturn:                 nil,
			mockErr:                    constant.ErrNotFound,
			serviceExpectedCalledCount: 1,
			fallbackMockReturn: &domainModel.CurrencyDomainModel{
				ID:   TEST_ID,
				Code: "USD",
				Name: "United States Dollar",
			},
			fallbackMockErr:                    nil,
			fallbackServiceExpectedCalledCount: 1,
			handleFallbackMockParam: &domainModel.CurrencyDomainModel{
				ID:   TEST_ID,
				Code: "USD",
				Name: "United States Dollar",
			},
			handleFallbackMockErr:                    nil,
			handleFallbackServiceExpectedCalledCount: 1,
		},
		{
			name:                                     "Currency is not exists in memcached and database",
			currencyCode:                             "USD",
			expectedReturn:                           nil,
			expectedErr:                              constant.ErrNotFound,
			mockReturn:                               nil,
			mockErr:                                  constant.ErrNotFound,
			serviceExpectedCalledCount:               1,
			fallbackMockReturn:                       nil,
			fallbackMockErr:                          constant.ErrNotFound,
			fallbackServiceExpectedCalledCount:       1,
			handleFallbackMockParam:                  nil,
			handleFallbackMockErr:                    nil,
			handleFallbackServiceExpectedCalledCount: 0,
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockService.EXPECT().GetByCode(gomock.Any(), test.currencyCode).Return(test.mockReturn, test.mockErr).Times(test.serviceExpectedCalledCount)

			mockFallbackService.EXPECT().GetByCode(gomock.Any(), test.currencyCode).Return(test.fallbackMockReturn, test.fallbackMockErr).Times(test.fallbackServiceExpectedCalledCount)

			mockService.EXPECT().Add(gomock.Any(), test.handleFallbackMockParam).Return(test.handleFallbackMockParam, test.handleFallbackMockErr).Times(test.handleFallbackServiceExpectedCalledCount)

			currencies, err := query.GetByCode(context.TODO(), test.currencyCode)

			if !reflect.DeepEqual(test.expectedErr, err) {
				t.Errorf("Expected error: %v, got: %v", test.expectedErr, err)
			}
			assert.DeepEqual(t, test.expectedReturn, currencies)
		})
	}
}
