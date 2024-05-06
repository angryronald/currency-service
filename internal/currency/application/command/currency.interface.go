package command

import (
	"context"

	"github.com/angryronald/currency-service/internal/currency/application/model"
)

type CurrencyCommandInterface interface {
	Add(ctx context.Context, currency *model.CurrencyApplicationModel) (*model.CurrencyApplicationModel, error)
	MultipleAdd(ctx context.Context, currencies []*model.CurrencyApplicationModel) ([]*model.CurrencyApplicationModel, error)
}
