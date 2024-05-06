package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/angryronald/go-kit/cast"

	"github.com/angryronald/currency-service/internal/currency/infrastructure/repository/model"
)

func Seeding(
	memcachedRepository CurrencyRepository,
	databaseRepository CurrencyRepository,
	filepath string,
) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	var jsonData map[string]string
	if err = cast.FromBytes(data, &jsonData); err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
	}

	currencies := []*model.CurrencyRepositoryModel{}
	for code, name := range jsonData {
		currencies = append(currencies, &model.CurrencyRepositoryModel{
			Code: code,
			Name: name,
		})
	}

	if existingCurrencies, err := databaseRepository.FindAll(context.Background()); existingCurrencies == nil || err != nil {
		if currencies, err = databaseRepository.BulkInsert(context.Background(), currencies); err != nil {
			return fmt.Errorf("error seeding to database: %v", err)
		}
	}

	if existingCurrencies, err := memcachedRepository.FindAll(context.Background()); existingCurrencies == nil || err != nil {
		if _, err = memcachedRepository.BulkInsert(context.Background(), currencies); err != nil {
			return fmt.Errorf("error seeding to memcached: %v", err)
		}
	}

	return nil
}
