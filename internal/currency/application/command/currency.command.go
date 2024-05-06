package command

import (
	"context"

	"github.com/angryronald/currency-service/internal/currency/application/model"
	"github.com/angryronald/currency-service/internal/currency/domain/currency"
	domainModel "github.com/angryronald/currency-service/internal/currency/domain/model"
)

type CurrencyCommand struct {
	currencyService currency.CurrencyServiceInterface
}

func (c *CurrencyCommand) Add(ctx context.Context, currency *model.CurrencyApplicationModel) (*model.CurrencyApplicationModel, error) {
	currencyDomain, err := c.currencyService.Add(ctx, currency.ToDomainModel())
	if err != nil {
		return nil, err
	}

	currency.FromDomainModel(currencyDomain)

	return currency, nil
}

func (c *CurrencyCommand) MultipleAdd(ctx context.Context, currencies []*model.CurrencyApplicationModel) ([]*model.CurrencyApplicationModel, error) {
	var err error

	currenciesDomain := []*domainModel.CurrencyDomainModel{}
	for _, currency := range currencies {
		currenciesDomain = append(currenciesDomain, currency.ToDomainModel())
	}

	if currenciesDomain, err = c.currencyService.MultipleAdd(ctx, currenciesDomain); err != nil {
		return nil, err
	}

	result := []*model.CurrencyApplicationModel{}
	for _, currency := range currenciesDomain {
		currencyApplication := &model.CurrencyApplicationModel{}
		currencyApplication.FromDomainModel(currency)
		result = append(result, currencyApplication)
	}

	return result, nil
}

func NewCurrencyCommand(
	currencyService currency.CurrencyServiceInterface,
) CurrencyCommandInterface {
	return &CurrencyCommand{
		currencyService: currencyService,
	}
}
