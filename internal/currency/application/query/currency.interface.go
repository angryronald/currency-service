package query

import (
	"context"

	"github.com/angryronald/currency-service/internal/currency/application/model"
)

type CurrencyQueryInterface interface {
	List(ctx context.Context) ([]*model.CurrencyApplicationModel, error)
	GetByCode(ctx context.Context, currencyCode string) (*model.CurrencyApplicationModel, error)
}
