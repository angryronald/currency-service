package currency

import (
	"context"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"

	"github.com/angryronald/currency-service/internal/currency/constant"
	"github.com/angryronald/currency-service/internal/currency/infrastructure/repository"
	"github.com/angryronald/currency-service/internal/currency/infrastructure/repository/model"
)

type CurrencyRepository struct {
	db *gorm.DB
}

func (r *CurrencyRepository) FindAll(ctx context.Context) ([]*model.CurrencyRepositoryModel, error) {
	var err error
	var result []*model.CurrencyRepositoryModel

	if err = r.db.Find(&result).Error; err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, constant.ErrNotFound
	}

	return result, nil
}

func (r *CurrencyRepository) FindByCode(ctx context.Context, code string) (*model.CurrencyRepositoryModel, error) {
	var err error
	result := &model.CurrencyRepositoryModel{}

	where := `code = ?`
	args := []interface{}{
		code,
	}
	if err = r.db.Where(
		where,
		args,
	).First(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, constant.ErrNotFound
		}
		return nil, err
	}

	if result == nil {
		return nil, constant.ErrNotFound
	}

	return result, nil
}

func (r *CurrencyRepository) FindByID(ctx context.Context, ID uuid.UUID) (*model.CurrencyRepositoryModel, error) {
	var err error
	result := &model.CurrencyRepositoryModel{}

	where := `"id" = ?`
	args := []interface{}{
		ID,
	}
	if err = r.db.Where(
		where,
		args,
	).First(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, constant.ErrNotFound
		}
		return nil, err
	}

	if result == nil {
		return nil, constant.ErrNotFound
	}

	return result, nil
}

func (r *CurrencyRepository) Insert(ctx context.Context, currency *model.CurrencyRepositoryModel) (*model.CurrencyRepositoryModel, error) {
	var err error
	currency.ID = uuid.New()
	currency.CreatedAt = time.Now().UTC()

	if err = r.db.Create(currency).Error; err != nil {
		if err == gorm.ErrDuplicatedKey || strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, constant.ErrConflict
		}
		return nil, err
	}

	return currency, nil
}

func (r *CurrencyRepository) BulkUpsert(ctx context.Context, currencies []*model.CurrencyRepositoryModel) ([]*model.CurrencyRepositoryModel, error) {
	var err error

	for _, currency := range currencies {
		currency.ID = uuid.New()
		currency.CreatedAt = time.Now().UTC()
	}

	if err = r.db.Save(currencies).Error; err != nil {
		if err == gorm.ErrDuplicatedKey || strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, constant.ErrConflict
		}
		return nil, err
	}

	return currencies, nil
}

// NewCurrencyRepository create postgre currency repository
func NewCurrencyRepository(
	db *gorm.DB,
) repository.CurrencyRepository {
	return &CurrencyRepository{
		db: db,
	}
}
