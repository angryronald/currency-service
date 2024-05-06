package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/angryronald/currency-service/internal/currency/infrastructure/repository/model"
)

type CurrencyRepository interface {
	FindAll(ctx context.Context) ([]*model.CurrencyRepositoryModel, error)
	FindByCode(ctx context.Context, currencyCode string) (*model.CurrencyRepositoryModel, error)
	FindByID(ctx context.Context, ID uuid.UUID) (*model.CurrencyRepositoryModel, error)
	Insert(ctx context.Context, currency *model.CurrencyRepositoryModel) (*model.CurrencyRepositoryModel, error)
	BulkInsert(ctx context.Context, currencies []*model.CurrencyRepositoryModel) ([]*model.CurrencyRepositoryModel, error)
}
