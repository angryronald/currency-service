package query

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/angryronald/currency-service/internal/currency/application/model"
	"github.com/angryronald/currency-service/internal/currency/constant"
	"github.com/angryronald/currency-service/internal/currency/domain/currency"
	domainModel "github.com/angryronald/currency-service/internal/currency/domain/model"
)

type CurrencyQuery struct {
	currencyService         currency.CurrencyServiceInterface
	fallbackCurrencyService currency.CurrencyServiceInterface
	log                     *logrus.Logger
}

func (q *CurrencyQuery) fallbackList(ctx context.Context) ([]*domainModel.CurrencyDomainModel, error) {
	currencies, err := q.fallbackCurrencyService.List(ctx)
	if err != nil {
		return nil, err
	}

	// Note: to keep the services availability high, when it is failing to write to memcached, should not breaking request
	if _, err := q.currencyService.MultipleAdd(ctx, currencies); err != nil {
		q.log.Warnf("failing to write on memcached as a fallback scenario (List): %v (%v)", err, currencies)
	}

	return currencies, nil
}

func (q *CurrencyQuery) List(ctx context.Context) ([]*model.CurrencyApplicationModel, error) {
	currenciesDomain, err := q.currencyService.List(ctx)
	if err != nil {
		if err != constant.ErrNotFound {
			return nil, err
		}
		if currenciesDomain, err = q.fallbackList(ctx); err != nil {
			return nil, err
		}
	}

	result := []*model.CurrencyApplicationModel{}
	for _, currency := range currenciesDomain {
		result = append(result, model.NewCurrencyApplicationModel(currency))
	}

	return result, nil
}

func (q *CurrencyQuery) fallbackGetByCode(ctx context.Context, currencyCode string) (*domainModel.CurrencyDomainModel, error) {
	currency, err := q.fallbackCurrencyService.GetByCode(ctx, currencyCode)
	if err != nil {
		return nil, err
	}

	// Note: to keep the services availability high, when it is failing to write to memcached, should not breaking request
	if _, err := q.currencyService.Add(ctx, currency); err != nil {
		q.log.Warnf("failing to write on memcached as a fallback scenario (GetByCode): %v (%v)", err, currency)
	}

	return currency, nil
}

func (q *CurrencyQuery) GetByCode(ctx context.Context, currencyCode string) (*model.CurrencyApplicationModel, error) {
	currency, err := q.currencyService.GetByCode(ctx, currencyCode)
	if err != nil {
		if err != constant.ErrNotFound {
			return nil, err
		}
		if currency, err = q.fallbackGetByCode(ctx, currencyCode); err != nil {
			return nil, err
		}
	}

	return model.NewCurrencyApplicationModel(currency), nil
}

func NewCurrencyQuery(
	currencyService currency.CurrencyServiceInterface,
	fallbackCurrencyService currency.CurrencyServiceInterface,
	log *logrus.Logger,
) CurrencyQueryInterface {
	return &CurrencyQuery{
		currencyService:         currencyService,
		fallbackCurrencyService: fallbackCurrencyService,
		log:                     log,
	}
}
