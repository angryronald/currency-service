package sync

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/angryronald/currency-service/internal/currency/infrastructure/repository"
	"github.com/angryronald/currency-service/internal/currency/infrastructure/repository/model"
	"github.com/angryronald/currency-service/lib/worker"
)

func SynchronizeReadAndWriteData(
	currencyMemcachedRepository repository.CurrencyRepository,
	currencySQLRepository repository.CurrencyRepository,
	periodInSec int,
	log *logrus.Logger,
) {
	worker.RunFuncEveryGivenPeriod(
		func(ctx context.Context) error {
			var err error
			var currencies []*model.CurrencyRepositoryModel

			if currencies, err = currencySQLRepository.FindAll(ctx); err != nil {
				return err
			}

			if _, err = currencyMemcachedRepository.BulkUpsert(ctx, currencies); err != nil {
				return err
			}

			return nil
		}, periodInSec, log,
	)
}
