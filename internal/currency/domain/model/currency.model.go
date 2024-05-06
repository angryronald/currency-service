package model

import (
	"github.com/google/uuid"

	"github.com/angryronald/currency-service/internal/currency/infrastructure/repository/model"
)

type CurrencyDomainModel struct {
	ID   uuid.UUID
	Name string
	Code string
}

func (m *CurrencyDomainModel) FromRepositoryModel(repo *model.CurrencyRepositoryModel) {
	m.ID = repo.ID
	m.Name = repo.Name
	m.Code = repo.Code
}

func (m *CurrencyDomainModel) ToRepositoryModel() *model.CurrencyRepositoryModel {
	return &model.CurrencyRepositoryModel{
		ID:   m.ID,
		Name: m.Name,
		Code: m.Code,
	}
}

func NewCurrencyDomainModel(repo *model.CurrencyRepositoryModel) *CurrencyDomainModel {
	result := &CurrencyDomainModel{}
	result.FromRepositoryModel(repo)
	return result
}
