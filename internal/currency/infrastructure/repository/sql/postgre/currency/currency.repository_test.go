package currency

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/angryronald/currency-service/internal/currency/constant"
	"github.com/angryronald/currency-service/internal/currency/infrastructure/repository/model"
)

var TEST_ID = uuid.MustParse("25ebc2e2-7fd2-47a0-bd2f-a8ec6e4d0163")

func TestCurrencyRepository_FindAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		preseed        []*model.CurrencyRepositoryModel
		name           string
		currenciesData []*model.CurrencyRepositoryModel
		expectedErr    error
	}{
		{
			name: "FindAll success return data",
			preseed: []*model.CurrencyRepositoryModel{
				{ID: uuid.New(), Code: "ABC"},
				{ID: uuid.New(), Code: "DEF"},
			},
			currenciesData: []*model.CurrencyRepositoryModel{
				{ID: uuid.New(), Code: "ABC"},
				{ID: uuid.New(), Code: "DEF"},
			},
			expectedErr: nil,
		},
		{
			name:           "FindAll with no results",
			preseed:        nil,
			currenciesData: nil,
			expectedErr:    constant.ErrNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.preseed != nil {
				fmt.Println("error inserting preseed data: ", db.Create(test.preseed).Error)
			}

			repo := NewCurrencyRepository(db)

			currencies, err := repo.FindAll(context.TODO())

			if err != test.expectedErr {
				t.Errorf("Expected error: %v, got: %v", test.expectedErr, err)
			}

			if len(currencies) != len(test.currenciesData) {
				t.Errorf("Expected %d currencies, got %d", len(test.currenciesData), len(currencies))
			}

			db.Delete(test.preseed)
		})
	}
}

func TestCurrencyRepository_FindByCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewCurrencyRepository(db)

	ID := uuid.New()
	tests := []struct {
		preseed       *model.CurrencyRepositoryModel
		name          string
		code          string
		currencyData  *model.CurrencyRepositoryModel
		expectedErr   error
		expectedCount int
	}{
		{
			name:          "FindByCode with valid code",
			preseed:       &model.CurrencyRepositoryModel{ID: ID, Code: "IDR", Name: "Indonesian Rupiah"},
			code:          "IDR",
			currencyData:  &model.CurrencyRepositoryModel{ID: ID, Code: "IDR", Name: "Indonesian Rupiah"},
			expectedErr:   nil,
			expectedCount: 1,
		},
		{
			name:          "FindByCode with no results",
			preseed:       nil,
			code:          "THB",
			currencyData:  nil,
			expectedErr:   constant.ErrNotFound,
			expectedCount: 0,
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.preseed != nil {
				fmt.Println("error inserting preseed data: ", db.Create(test.preseed).Error)
			}

			currency, err := repo.FindByCode(context.TODO(), test.code)

			if err != test.expectedErr {
				t.Errorf("Expected error: %v, got: %v", test.expectedErr, err)
			}

			if currency != nil {
				test.currencyData.CreatedAt = currency.CreatedAt
			}
			if !reflect.DeepEqual(currency, test.currencyData) {
				t.Errorf("Expected %v, got %v", test.currencyData, currency)
			}

			db.Delete(test.preseed)
		})
	}
}

func TestCurrencyRepository_FindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewCurrencyRepository(db)

	ID := uuid.New()
	tests := []struct {
		name          string
		preseed       *model.CurrencyRepositoryModel
		currencyID    uuid.UUID
		currencyData  *model.CurrencyRepositoryModel
		expectedErr   error
		expectedCount int
	}{
		{
			name:          "FindByID with valid currency ID",
			preseed:       &model.CurrencyRepositoryModel{ID: ID},
			currencyID:    ID,
			currencyData:  &model.CurrencyRepositoryModel{ID: ID},
			expectedErr:   nil,
			expectedCount: 1,
		},
		{
			name:          "FindByID with non-existent currency ID",
			preseed:       nil,
			currencyID:    uuid.New(),
			currencyData:  nil,
			expectedErr:   constant.ErrNotFound,
			expectedCount: 0,
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.preseed != nil {
				fmt.Println("error inserting preseed data: ", db.Create(test.preseed).Error)
			}

			currency, err := repo.FindByID(context.TODO(), test.currencyID)

			if err != test.expectedErr {
				t.Errorf("Expected error: %v, got: %v", test.expectedErr, err)
			}

			if currency != nil {
				test.currencyData.CreatedAt = currency.CreatedAt
			}
			if !reflect.DeepEqual(currency, test.currencyData) {
				t.Errorf("Expected: %v, got: %v", test.currencyData, currency)
			}

			db.Delete(test.preseed)
		})
	}
}

func TestCurrencyRepository_Insert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewCurrencyRepository(db)

	tests := []struct {
		name           string
		preseed        *model.CurrencyRepositoryModel
		currencyData   *model.CurrencyRepositoryModel
		expectedErr    error
		expectedCount  int
		expectedExists bool
	}{
		{
			name:    "Insert new currency",
			preseed: nil,
			currencyData: &model.CurrencyRepositoryModel{
				ID:   TEST_ID,
				Name: "Euro",
				Code: "EUR",
				// Set other fields as needed
			},
			expectedErr:    nil,
			expectedCount:  1,
			expectedExists: true,
		},
		{
			name: "Insert new currency with existing Code",
			preseed: &model.CurrencyRepositoryModel{
				ID:   uuid.New(),
				Code: "EUR",
				Name: "Following the previous inserted data",
				// Set other fields as needed
			},
			currencyData: &model.CurrencyRepositoryModel{
				ID:   uuid.New(),
				Code: "EUR",
				Name: "Following the previous inserted data",
				// Set other fields as needed
			},
			expectedErr:    constant.ErrConflict,
			expectedCount:  1,
			expectedExists: false,
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.preseed != nil {
				fmt.Println("error inserting preseed data: ", db.Create(test.preseed).Error)
			}

			insertedCurrency, err := repo.Insert(context.TODO(), test.currencyData)

			if !reflect.DeepEqual(test.expectedErr, err) {
				t.Errorf("Expected error: %v, got: %v", test.expectedErr, err)
			}

			if insertedCurrency != nil {
				actual := &model.CurrencyRepositoryModel{}
				err = db.Where("id = ?", insertedCurrency.ID).First(&actual).Error

				actual.CreatedAt = insertedCurrency.CreatedAt
				if !reflect.DeepEqual(insertedCurrency, actual) {
					t.Errorf("Expected: %v, got: %v", insertedCurrency, actual)
				}
			}

			db.Delete(insertedCurrency)
			if test.preseed != nil {
				db.Delete(test.preseed)
			}
		})
	}
}

func TestCurrencyRepository_BulkInsert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewCurrencyRepository(db)

	tests := []struct {
		name           string
		currencyData   []*model.CurrencyRepositoryModel
		expectedErr    error
		expectedCount  int
		expectedExists bool
		mockReturn     []*model.CurrencyRepositoryModel
		mockErr        error
	}{
		{
			name: "BulkInsert new currencies",
			currencyData: []*model.CurrencyRepositoryModel{
				{
					ID:   TEST_ID,
					Name: "Euro",
					Code: "EUR",
					// Set other fields as needed
				},
			},
			expectedErr:    nil,
			expectedCount:  1,
			expectedExists: true,
			mockReturn: []*model.CurrencyRepositoryModel{
				{
					ID:   TEST_ID,
					Name: "Euro",
					Code: "EUR",
					// Set other fields as needed
				},
			},
			mockErr: nil,
		},
		{
			name: "BulkInsert new currencies with existing Code",
			currencyData: []*model.CurrencyRepositoryModel{
				{
					ID:   uuid.New(),
					Code: "USD",
					Name: "United States Dollar",
					// Set other fields as needed
				},
				{
					ID:   uuid.New(),
					Code: "USD",
					Name: "United States Dollar",
					// Set other fields as needed
				},
			},
			expectedErr:    constant.ErrConflict,
			expectedCount:  1,
			expectedExists: false,
			mockReturn:     nil,
			mockErr:        constant.ErrConflict,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			insertedCurrencies, err := repo.BulkInsert(context.TODO(), test.currencyData)

			if !reflect.DeepEqual(test.expectedErr, err) {
				t.Errorf("Expected error: %v, got: %v", test.expectedErr, err)
			}

			for _, insertedCurrency := range insertedCurrencies {
				if insertedCurrency != nil {
					actual := &model.CurrencyRepositoryModel{}
					err = db.Where("id = ?", insertedCurrency.ID).First(&actual).Error

					actual.CreatedAt = insertedCurrency.CreatedAt
					if !reflect.DeepEqual(insertedCurrency, actual) {
						t.Errorf("Expected: %v, got: %v", insertedCurrency, actual)
					}
				}
			}

			db.Delete(insertedCurrencies)
		})
	}
}
