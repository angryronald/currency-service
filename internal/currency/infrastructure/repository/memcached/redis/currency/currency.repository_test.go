package currency

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/angryronald/currency-service/internal/currency/constant"
	"github.com/angryronald/currency-service/internal/currency/infrastructure/repository/model"
	"github.com/angryronald/currency-service/lib/cast"
)

var TEST_ID = uuid.MustParse("25ebc2e2-7fd2-47a0-bd2f-a8ec6e4d0163")

func TestCurrencyRepository_FindAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewCurrencyRepository(redisClient, 300)

	tests := []struct {
		name           string
		currenciesData []*model.CurrencyRepositoryModel
		expectedErr    error
	}{
		{
			name: "FindAll success return data",
			currenciesData: []*model.CurrencyRepositoryModel{
				{ID: uuid.New() /* other fields */},
				{ID: uuid.New() /* other fields */},
			},
			expectedErr: nil,
		},
		{
			name:           "FindAll with no results",
			currenciesData: nil,
			expectedErr:    constant.ErrNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if len(test.currenciesData) > 0 {
				currenciesInBytes, err := cast.ToBytes(test.currenciesData)
				if err != nil {
					fmt.Printf("error on casting object to bytes: %+v", err)
				}

				if err = redisClient.Set(
					context.TODO(),
					constant.CURRENCY_CACHE_KEY,
					currenciesInBytes,
					time.Duration(300*time.Second)).Err(); err != nil {
					fmt.Printf("error on inserting to cache: %+v", err)
				}
			}

			currencies, err := repo.FindAll(context.TODO())

			if !reflect.DeepEqual(err, test.expectedErr) {
				t.Errorf("Expected error: %v, got: %v", test.expectedErr, err)
			}

			if test.expectedErr == nil {
				if len(currencies) != len(test.currenciesData) {
					t.Errorf("Expected %d currencies, got %d", len(test.currenciesData), len(currencies))
				}
			}

			if len(test.currenciesData) > 0 {
				if err = redisClient.Del(
					context.TODO(),
					constant.CURRENCY_CACHE_KEY).Err(); err != nil {
					fmt.Printf("error on deleting from cache: %+v", err)
				}
			}
		})
	}
}

func TestCurrencyRepository_FindByCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewCurrencyRepository(redisClient, 300)

	tests := []struct {
		name         string
		code         string
		currencyData *model.CurrencyRepositoryModel
		expectedErr  error
	}{
		{
			name:         "FindByCode with valid code",
			code:         "IDR",
			currencyData: &model.CurrencyRepositoryModel{ID: uuid.New(), Code: "IDR", Name: "Indonesian Rupiah"},
			expectedErr:  nil,
		},
		{
			name:         "FindByCode with no results",
			code:         "THB",
			currencyData: nil,
			expectedErr:  constant.ErrNotFound,
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.currencyData != nil {
				currencyInBytes, err := cast.ToBytes(test.currencyData)
				if err != nil {
					fmt.Printf("error on casting object to bytes: %+v", err)
				}

				if err = redisClient.Set(
					context.TODO(),
					fmt.Sprintf("%s.CODE.%s", constant.CURRENCY_CACHE_KEY, test.currencyData.Code),
					currencyInBytes,
					time.Duration(300*time.Second)).Err(); err != nil {
					fmt.Printf("error on inserting to cache: %+v", err)
				}
			}

			currency, err := repo.FindByCode(context.TODO(), test.code)
			if !reflect.DeepEqual(err, test.expectedErr) {
				t.Errorf("Expected error: %v, got: %v", test.expectedErr, err)
			}
			if !reflect.DeepEqual(currency, test.currencyData) {
				t.Errorf("Expected: %v, got: %v", test.currencyData, currency)
			}

			if test.currencyData != nil {
				if err = redisClient.Del(
					context.TODO(),
					fmt.Sprintf("%s.CODE.%s", constant.CURRENCY_CACHE_KEY, test.currencyData.Code)).Err(); err != nil {
					fmt.Printf("error on deleting from cache: %+v", err)
				}
			}
		})
	}
}

func TestCurrencyRepository_FindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewCurrencyRepository(redisClient, 300)

	tests := []struct {
		name          string
		currencyID    uuid.UUID
		currencyData  *model.CurrencyRepositoryModel
		expectedErr   error
		expectedCount int
	}{
		{
			name:          "FindByID with valid currency ID",
			currencyID:    TEST_ID,
			currencyData:  &model.CurrencyRepositoryModel{ID: TEST_ID /* other fields */},
			expectedErr:   nil,
			expectedCount: 1,
		},
		{
			name:          "FindByID with non-existent currency ID",
			currencyID:    uuid.New(),
			currencyData:  nil,
			expectedErr:   constant.ErrNotFound,
			expectedCount: 0,
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.currencyData != nil {
				currencyInBytes, err := cast.ToBytes(test.currencyData)
				if err != nil {
					fmt.Printf("error on casting object to bytes: %+v", err)
				}

				if err = redisClient.Set(
					context.TODO(),
					fmt.Sprintf("%s.ID.%s", constant.CURRENCY_CACHE_KEY, test.currencyData.ID),
					currencyInBytes,
					time.Duration(300*time.Second)).Err(); err != nil {
					fmt.Printf("error on inserting to cache: %+v", err)
				}
			}

			currency, err := repo.FindByID(context.TODO(), test.currencyID)

			if err != test.expectedErr {
				t.Errorf("Expected error: %v, got: %v", test.expectedErr, err)
			}

			if !reflect.DeepEqual(currency, test.currencyData) {
				t.Errorf("Expected: %v, got: %v", test.currencyData, currency)
			}

			if test.currencyData != nil {
				if err = redisClient.Del(
					context.TODO(),
					fmt.Sprintf("%s.ID.%s", constant.CURRENCY_CACHE_KEY, test.currencyData.ID)).Err(); err != nil {
					fmt.Printf("error on deleting from cache: %+v", err)
				}
			}
		})
	}
}

func TestCurrencyRepository_Insert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewCurrencyRepository(redisClient, 300)

	tests := []struct {
		name           string
		currencyData   *model.CurrencyRepositoryModel
		expectedErr    error
		expectedCount  int
		expectedExists bool
		expectedReturn *model.CurrencyRepositoryModel
	}{
		{
			name: "Insert new currency",
			currencyData: &model.CurrencyRepositoryModel{
				ID:   TEST_ID,
				Code: "USD",
				Name: "United States Dollar",
				// Set other fields as needed
			},
			expectedReturn: &model.CurrencyRepositoryModel{
				ID:   TEST_ID,
				Code: "USD",
				Name: "United States Dollar",
				// Set other fields as needed
			},
			expectedErr:    nil,
			expectedCount:  1,
			expectedExists: true,
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			insertedCurrency, err := repo.Insert(context.TODO(), test.currencyData)

			if !reflect.DeepEqual(err, test.expectedErr) {
				t.Errorf("Expected error: %v, got: %v", test.expectedErr, err)
			}

			if !reflect.DeepEqual(insertedCurrency, test.expectedReturn) {
				t.Errorf("Expected: %v, got: %v", test.expectedReturn, insertedCurrency)
			}

			resultRaw, err := redisClient.Get(context.TODO(), fmt.Sprintf("%s.ID.%s", constant.CURRENCY_CACHE_KEY, test.currencyData.ID)).Result()
			if err != nil {
				t.Errorf("Expected: %v, got: %v", nil, err)
			}

			currency := &model.CurrencyRepositoryModel{}
			if err = cast.FromBytes([]byte(resultRaw), currency); err != nil {
				fmt.Printf("error on casting object from bytes: %+v", err)
			}

			if !reflect.DeepEqual(insertedCurrency, currency) {
				t.Errorf("Expected: %v, got: %v", insertedCurrency, currency)
			}

			if test.currencyData != nil {
				if err = redisClient.Del(
					context.TODO(),
					fmt.Sprintf("%s.ID.%s", constant.CURRENCY_CACHE_KEY, test.currencyData.ID)).Err(); err != nil {
					fmt.Printf("error on deleting from cache: %+v", err)
				}
			}
		})
	}
}

func TestCurrencyRepository_BulkInsert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewCurrencyRepository(redisClient, 300)

	tests := []struct {
		name           string
		currencyData   []*model.CurrencyRepositoryModel
		expectedErr    error
		expectedCount  int
		expectedExists bool
		expectedReturn []*model.CurrencyRepositoryModel
	}{
		{
			name: "Insert new currency",
			currencyData: []*model.CurrencyRepositoryModel{
				{
					ID:   TEST_ID,
					Code: "USD",
					Name: "United States Dollar",
				},
				// Set other fields as needed
			},
			expectedReturn: []*model.CurrencyRepositoryModel{
				{
					ID:   TEST_ID,
					Code: "USD",
					Name: "United States Dollar",
				},
				// Set other fields as needed
			},
			expectedErr:    nil,
			expectedCount:  1,
			expectedExists: true,
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := redisClient.Del(
				context.TODO(),
				constant.CURRENCY_CACHE_KEY).Err(); err != nil {
				fmt.Printf("error on deleting from cache: %+v", err)
			}

			insertedCurrency, err := repo.BulkInsert(context.TODO(), test.currencyData)

			if !reflect.DeepEqual(err, test.expectedErr) {
				t.Errorf("Expected error: %v, got: %v", test.expectedErr, err)
			}

			if !reflect.DeepEqual(insertedCurrency, test.expectedReturn) {
				t.Errorf("Expected: %v, got: %v", test.expectedReturn, insertedCurrency)
			}

			resultRaw, err := redisClient.Get(context.TODO(), constant.CURRENCY_CACHE_KEY).Result()
			if err != nil {
				t.Errorf("Expected: %v, got: %v", nil, err)
			}

			currencies := []*model.CurrencyRepositoryModel{}
			if err = cast.FromBytes([]byte(resultRaw), &currencies); err != nil {
				t.Errorf("Expected: %v, got: %v", nil, err)
			}

			if test.expectedErr == nil {
				if len(currencies) != len(test.expectedReturn) {
					t.Errorf("Expected %d currencies, got %d", len(currencies), len(test.expectedReturn))
				}
			}

			if len(currencies) > 0 {
				if err = redisClient.Del(
					context.TODO(),
					constant.CURRENCY_CACHE_KEY).Err(); err != nil {
					fmt.Printf("error on deleting from cache: %+v", err)
				}

				for _, currency := range currencies {
					if err = redisClient.Del(
						context.TODO(),
						fmt.Sprintf("%s.ID.%s", constant.CURRENCY_CACHE_KEY, currency.ID)).Err(); err != nil {
						fmt.Printf("error on deleting from cache: %+v", err)
					}

					if err = redisClient.Del(
						context.TODO(),
						fmt.Sprintf("%s.CODE.%s", constant.CURRENCY_CACHE_KEY, currency.ID)).Err(); err != nil {
						fmt.Printf("error on deleting from cache: %+v", err)
					}
				}
			}
		})
	}
}
