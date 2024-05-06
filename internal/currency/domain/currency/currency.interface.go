package currency

import (
	"context"

	"github.com/angryronald/currency-service/internal/currency/domain/model"
)

type CurrencyServiceInterface interface {
	List(ctx context.Context) ([]*model.CurrencyDomainModel, error)
	GetByCode(ctx context.Context, currencyCode string) (*model.CurrencyDomainModel, error)
	Add(ctx context.Context, currency *model.CurrencyDomainModel) (*model.CurrencyDomainModel, error)
	MultipleAddOrUpdate(ctx context.Context, currencies []*model.CurrencyDomainModel) ([]*model.CurrencyDomainModel, error)
}
