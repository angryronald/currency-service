package currency

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/angryronald/currency-service/internal/currency/constant"
	"github.com/angryronald/currency-service/internal/currency/domain/model"
	"github.com/angryronald/currency-service/internal/currency/infrastructure/event/publisher"
	"github.com/angryronald/currency-service/internal/currency/infrastructure/repository"
	repositoryModel "github.com/angryronald/currency-service/internal/currency/infrastructure/repository/model"
	"github.com/angryronald/currency-service/lib/cast"
)

type CurrencyService struct {
	publisher  publisher.Publisher
	repository repository.CurrencyRepository
	log        *logrus.Logger
}

func (s *CurrencyService) List(ctx context.Context) ([]*model.CurrencyDomainModel, error) {
	currenciesRepo, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	result := []*model.CurrencyDomainModel{}
	for _, currency := range currenciesRepo {
		result = append(result, model.NewCurrencyDomainModel(currency))
	}

	return result, nil
}

func (s *CurrencyService) GetByCode(ctx context.Context, currencyCode string) (*model.CurrencyDomainModel, error) {
	currencyRepo, err := s.repository.FindByCode(ctx, currencyCode)
	if err != nil {
		return nil, err
	}

	return model.NewCurrencyDomainModel(currencyRepo), nil
}

func (s *CurrencyService) Add(ctx context.Context, currency *model.CurrencyDomainModel) (*model.CurrencyDomainModel, error) {
	currencyRepo, err := s.repository.Insert(ctx, currency.ToRepositoryModel())
	if err != nil {
		return nil, err
	}

	currency.FromRepositoryModel(currencyRepo)

	currencyInBytes, err := cast.ToBytes(currency)
	if err != nil {
		return nil, err
	}

	if err = s.publisher.Publish(ctx, string(constant.CURRENCY_ADDED_EVENT), currency.ID.String(), currencyInBytes); err != nil {
		s.log.Warnf(`[CurrencyService.Add] error: %s StackTrace: %v`, err, errors.WithStack(err))
	}

	return currency, nil
}

func (s *CurrencyService) MultipleAdd(ctx context.Context, currencies []*model.CurrencyDomainModel) ([]*model.CurrencyDomainModel, error) {
	var err error

	currenciesRepo := []*repositoryModel.CurrencyRepositoryModel{}
	for _, currency := range currencies {
		currenciesRepo = append(currenciesRepo, currency.ToRepositoryModel())
	}

	if currenciesRepo, err = s.repository.BulkInsert(ctx, currenciesRepo); err != nil {
		return nil, err
	}

	result := []*model.CurrencyDomainModel{}
	for _, currency := range currenciesRepo {
		result = append(result, model.NewCurrencyDomainModel(currency))
	}

	return result, nil
}

func NewCurrencyService(
	publisher publisher.Publisher,
	repository repository.CurrencyRepository,
	log *logrus.Logger,
) CurrencyServiceInterface {
	return &CurrencyService{
		publisher:  publisher,
		repository: repository,
		log:        log,
	}
}
