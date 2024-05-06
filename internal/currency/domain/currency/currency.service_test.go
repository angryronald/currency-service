package currency

import (
	"context"
	"errors"
	"testing"

	"github.com/angryronald/go-kit/cast"
	"github.com/angryronald/go-kit/publisher"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/angryronald/currency-service/internal/currency/constant"
	"github.com/angryronald/currency-service/internal/currency/domain/model"
	"github.com/angryronald/currency-service/internal/currency/infrastructure/repository"
	repositoryModel "github.com/angryronald/currency-service/internal/currency/infrastructure/repository/model"
)

var TEST_ID = uuid.MustParse("25ebc2e2-7fd2-47a0-bd2f-a8ec6e4d0163")

func TestCurrencyService_List(t *testing.T) {
	// Initialize mocks
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock repository with dummy data or a test database
	mockRepo := repository.NewMockCurrencyRepository(ctrl)
	mockPublisher := publisher.NewMockPublisher(ctrl)

	service := NewCurrencyService(mockPublisher, mockRepo, logrus.New())

	tests := []struct {
		name                         string
		expectedReturn               []*model.CurrencyDomainModel
		expectedErr                  error
		mockReturn                   []*repositoryModel.CurrencyRepositoryModel
		mockErr                      error
		repoExpectedCalledCount      int
		publisherExpectedCalledCount int
	}{
		{
			name: "Success get all of currencies",
			expectedReturn: []*model.CurrencyDomainModel{
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
			mockReturn: []*repositoryModel.CurrencyRepositoryModel{
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
			mockErr:                      nil,
			repoExpectedCalledCount:      1,
			publisherExpectedCalledCount: 0,
		},
		{
			name:                         "Got error when fetching currency data from repository",
			expectedReturn:               nil,
			expectedErr:                  errors.New("something went wrong"),
			mockReturn:                   nil,
			mockErr:                      errors.New("something went wrong"),
			repoExpectedCalledCount:      1,
			publisherExpectedCalledCount: 0,
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo.EXPECT().FindAll(gomock.Any()).Return(test.mockReturn, test.mockErr).Times(test.repoExpectedCalledCount)

			currencies, err := service.List(context.Background())

			assert.Equal(t, err, test.expectedErr)
			assert.Equal(t, test.expectedReturn, currencies)
		})
	}
}

func TestCurrencyService_GetByCode(t *testing.T) {
	// Initialize mocks
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock repository with dummy data or a test database
	mockRepo := repository.NewMockCurrencyRepository(ctrl)
	mockPublisher := publisher.NewMockPublisher(ctrl)

	service := NewCurrencyService(mockPublisher, mockRepo, logrus.New())

	tests := []struct {
		name                         string
		currencyCode                 string
		expectedReturn               *model.CurrencyDomainModel
		expectedErr                  error
		mockReturn                   *repositoryModel.CurrencyRepositoryModel
		mockErr                      error
		repoExpectedCalledCount      int
		publisherExpectedCalledCount int
	}{
		{
			name:         "Success get currency",
			currencyCode: "USD",
			expectedReturn: &model.CurrencyDomainModel{
				ID:   TEST_ID,
				Code: "USD",
				Name: "United States Dollar",
			},
			expectedErr: nil,
			mockReturn: &repositoryModel.CurrencyRepositoryModel{
				ID:   TEST_ID,
				Code: "USD",
				Name: "United States Dollar",
			},
			mockErr:                      nil,
			repoExpectedCalledCount:      1,
			publisherExpectedCalledCount: 0,
		},
		{
			name:                         "Got error when fetching currency data from repository",
			currencyCode:                 "USD",
			expectedReturn:               nil,
			expectedErr:                  errors.New("something went wrong"),
			mockReturn:                   nil,
			mockErr:                      errors.New("something went wrong"),
			repoExpectedCalledCount:      1,
			publisherExpectedCalledCount: 0,
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo.EXPECT().FindByCode(gomock.Any(), test.currencyCode).Return(test.mockReturn, test.mockErr).Times(test.repoExpectedCalledCount)

			currencies, err := service.GetByCode(context.Background(), test.currencyCode)

			assert.Equal(t, err, test.expectedErr)
			assert.Equal(t, test.expectedReturn, currencies)
		})
	}
}

func TestCurrencyService_Add(t *testing.T) {
	// Initialize mocks
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock repository with dummy data or a test database
	mockRepo := repository.NewMockCurrencyRepository(ctrl)
	mockPublisher := publisher.NewMockPublisher(ctrl)

	service := NewCurrencyService(mockPublisher, mockRepo, logrus.New())

	tests := []struct {
		name                         string
		currencyInput                *model.CurrencyDomainModel
		expectedReturn               *model.CurrencyDomainModel
		expectedErr                  error
		mockReturn                   *repositoryModel.CurrencyRepositoryModel
		mockErr                      error
		repoExpectedCalledCount      int
		publisherExpectedCalledCount int
	}{
		{
			name: "Success add currency",
			currencyInput: &model.CurrencyDomainModel{
				Code: "USD",
				Name: "United States Dollar",
			},
			expectedReturn: &model.CurrencyDomainModel{
				ID:   TEST_ID,
				Code: "USD",
				Name: "United States Dollar",
			},
			expectedErr: nil,
			mockReturn: &repositoryModel.CurrencyRepositoryModel{
				ID:   TEST_ID,
				Code: "USD",
				Name: "United States Dollar",
			},
			mockErr:                      nil,
			repoExpectedCalledCount:      1,
			publisherExpectedCalledCount: 1,
		},
		{
			name: "Got error when inserting currency data to repository",
			currencyInput: &model.CurrencyDomainModel{
				Code: "USD",
				Name: "United States Dollar",
			},
			expectedReturn:               nil,
			expectedErr:                  errors.New("something went wrong"),
			mockReturn:                   nil,
			mockErr:                      errors.New("something went wrong"),
			repoExpectedCalledCount:      1,
			publisherExpectedCalledCount: 0,
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo.EXPECT().Insert(gomock.Any(), test.currencyInput.ToRepositoryModel()).Return(test.mockReturn, test.mockErr).Times(test.repoExpectedCalledCount)
			currencyInBytes, _ := cast.ToBytes(test.expectedReturn)
			mockPublisher.EXPECT().Publish(gomock.Any(), string(constant.CURRENCY_ADDED_EVENT), gomock.Any(), currencyInBytes).Times(test.publisherExpectedCalledCount)

			currencies, err := service.Add(context.Background(), test.currencyInput)

			assert.Equal(t, err, test.expectedErr)
			assert.Equal(t, test.expectedReturn, currencies)
		})
	}
}

func TestCurrencyService_MultipleAdd(t *testing.T) {
	// Initialize mocks
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock repository with dummy data or a test database
	mockRepo := repository.NewMockCurrencyRepository(ctrl)
	mockPublisher := publisher.NewMockPublisher(ctrl)

	service := NewCurrencyService(mockPublisher, mockRepo, logrus.New())
	tests := []struct {
		name                    string
		currenciesInput         []*model.CurrencyDomainModel
		expectedReturn          []*model.CurrencyDomainModel
		expectedErr             error
		mockReturn              []*repositoryModel.CurrencyRepositoryModel
		mockErr                 error
		repoExpectedCalledCount int
	}{
		{
			name: "Success multiple add currencies",
			currenciesInput: []*model.CurrencyDomainModel{
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
			expectedReturn: []*model.CurrencyDomainModel{
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
			mockReturn: []*repositoryModel.CurrencyRepositoryModel{
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
			mockErr:                 nil,
			repoExpectedCalledCount: 1,
		},
		{
			name: "Got error when bulk insert currencies to repository",
			currenciesInput: []*model.CurrencyDomainModel{
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
			expectedReturn:          nil,
			expectedErr:             errors.New("something went wrong"),
			mockReturn:              nil,
			mockErr:                 errors.New("something went wrong"),
			repoExpectedCalledCount: 1,
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo.EXPECT().BulkInsert(gomock.Any(), gomock.Any()).Return(test.mockReturn, test.mockErr).Times(test.repoExpectedCalledCount)

			currencies, err := service.MultipleAdd(context.Background(), test.currenciesInput)

			assert.Equal(t, err, test.expectedErr)
			assert.Equal(t, test.expectedReturn, currencies)
		})
	}
}
