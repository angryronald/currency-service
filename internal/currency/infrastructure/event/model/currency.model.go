package model

import (
	"github.com/google/uuid"

	"github.com/angryronald/currency-service/internal/currency/infrastructure/repository/model"
)

type CurrencyEventModel struct {
	ID   uuid.UUID
	Name string
	Code string
}

func (m *CurrencyEventModel) ToRepositoryModel() *model.CurrencyRepositoryModel {
	return &model.CurrencyRepositoryModel{
		ID:   m.ID,
		Name: m.Name,
		Code: m.Code,
	}
}
