package model

import (
	"github.com/google/uuid"

	domainModel "github.com/angryronald/currency-service/internal/currency/domain/model"
)

type CurrencyApplicationModel struct {
	ID   uuid.UUID
	Name string
	Code string
}

func (m *CurrencyApplicationModel) FromDomainModel(currencyDomain *domainModel.CurrencyDomainModel) {
	m.ID = currencyDomain.ID
	m.Name = currencyDomain.Name
	m.Code = currencyDomain.Code
}

func (m *CurrencyApplicationModel) ToDomainModel() *domainModel.CurrencyDomainModel {
	return &domainModel.CurrencyDomainModel{
		ID:   m.ID,
		Name: m.Name,
		Code: m.Code,
	}
}

func NewCurrencyApplicationModel(currencyDomain *domainModel.CurrencyDomainModel) *CurrencyApplicationModel {
	result := &CurrencyApplicationModel{}
	result.FromDomainModel(currencyDomain)
	return result
}
