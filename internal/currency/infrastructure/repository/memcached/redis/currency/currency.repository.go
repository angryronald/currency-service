package currency

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/angryronald/currency-service/internal/currency/constant"
	"github.com/angryronald/currency-service/internal/currency/infrastructure/repository"
	"github.com/angryronald/currency-service/internal/currency/infrastructure/repository/model"
	"github.com/angryronald/currency-service/lib/cast"
)

type CurrencyRepository struct {
	client                  *redis.Client
	key                     string
	expirationDurationInSec int
}

func (r *CurrencyRepository) FindAll(ctx context.Context) ([]*model.CurrencyRepositoryModel, error) {
	resultRaw, err := r.client.Get(ctx, r.key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, constant.ErrNotFound
		}
		return nil, err
	}

	currencies := []*model.CurrencyRepositoryModel{}
	if err = cast.FromBytes([]byte(resultRaw), &currencies); err != nil {
		return nil, err
	}

	if len(currencies) == 0 {
		return nil, constant.ErrNotFound
	}

	return currencies, nil
}

func (r *CurrencyRepository) FindByCode(ctx context.Context, code string) (*model.CurrencyRepositoryModel, error) {
	resultRaw, err := r.client.Get(ctx, fmt.Sprintf("%s.CODE.%s", r.key, code)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, constant.ErrNotFound
		}
		return nil, err
	}

	currency := &model.CurrencyRepositoryModel{}
	if err = cast.FromBytes([]byte(resultRaw), currency); err != nil {
		return nil, err
	}

	if currency == nil {
		return nil, constant.ErrNotFound
	}

	return currency, nil
}

func (r *CurrencyRepository) FindByID(ctx context.Context, ID uuid.UUID) (*model.CurrencyRepositoryModel, error) {
	resultRaw, err := r.client.Get(ctx, fmt.Sprintf("%s.ID.%s", r.key, ID)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, constant.ErrNotFound
		}
		return nil, err
	}

	currency := &model.CurrencyRepositoryModel{}
	if err = cast.FromBytes([]byte(resultRaw), currency); err != nil {
		return nil, err
	}

	if currency == nil {
		return nil, constant.ErrNotFound
	}

	return currency, nil
}

func (r *CurrencyRepository) Insert(ctx context.Context, currency *model.CurrencyRepositoryModel) (*model.CurrencyRepositoryModel, error) {
	existingCurrencies, err := r.FindAll(ctx)
	if err != nil {
		if err != constant.ErrNotFound {
			return nil, err
		}
		existingCurrencies = []*model.CurrencyRepositoryModel{}
	}

	existingCurrencies = append(existingCurrencies, currency)

	existingCurrenciesInBytes, err := cast.ToBytes(existingCurrencies)
	if err != nil {
		return nil, err
	}

	if err = r.client.Set(
		ctx,
		r.key,
		existingCurrenciesInBytes,
		time.Duration(r.expirationDurationInSec*int(time.Second))).Err(); err != nil {
		return nil, err
	}

	currencyInBytes, err := cast.ToBytes(currency)
	if err != nil {
		return nil, err
	}

	if err = r.client.Set(
		ctx,
		fmt.Sprintf("%s.ID.%s", r.key, currency.ID),
		currencyInBytes,
		time.Duration(r.expirationDurationInSec*int(time.Second))).Err(); err != nil {
		return nil, err
	}

	if err = r.client.Set(
		ctx,
		fmt.Sprintf("%s.CODE.%s", r.key, currency.Code),
		currencyInBytes,
		time.Duration(r.expirationDurationInSec*int(time.Second))).Err(); err != nil {
		return nil, err
	}

	return currency, nil
}

func (r *CurrencyRepository) BulkUpsert(ctx context.Context, currencies []*model.CurrencyRepositoryModel) ([]*model.CurrencyRepositoryModel, error) {
	existingCurrencies, err := r.FindAll(ctx)
	if err != nil {
		if err != constant.ErrNotFound {
			return nil, err
		}
		existingCurrencies = []*model.CurrencyRepositoryModel{}
	}

	existingCurrencies = append(existingCurrencies, currencies...)

	existingCurrenciesInBytes, err := cast.ToBytes(existingCurrencies)
	if err != nil {
		return nil, err
	}

	if err = r.client.Set(
		ctx,
		r.key,
		existingCurrenciesInBytes,
		time.Duration(r.expirationDurationInSec*int(time.Second))).Err(); err != nil {
		return nil, err
	}

	for _, currency := range currencies {
		currencyInBytes, err := cast.ToBytes(currency)
		if err != nil {
			return nil, err
		}

		if err = r.client.Set(
			ctx,
			fmt.Sprintf("%s.ID.%s", r.key, currency.ID),
			currencyInBytes,
			time.Duration(r.expirationDurationInSec*int(time.Second))).Err(); err != nil {
			return nil, err
		}

		if err = r.client.Set(
			ctx,
			fmt.Sprintf("%s.CODE.%s", r.key, currency.Code),
			currencyInBytes,
			time.Duration(r.expirationDurationInSec*int(time.Second))).Err(); err != nil {
			return nil, err
		}
	}

	return currencies, nil
}

// NewCurrencyRepository create redis currency repository, expirationDuration need to be set in redisClient when doing dependency injection on GenericRepo
func NewCurrencyRepository(
	client *redis.Client,
	expirationDurationInSec int,
) repository.CurrencyRepository {
	return &CurrencyRepository{
		client:                  client,
		key:                     constant.CURRENCY_CACHE_KEY,
		expirationDurationInSec: expirationDurationInSec,
	}
}
