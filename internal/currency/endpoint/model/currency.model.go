package model

import "github.com/angryronald/currency-service/internal/currency/application/model"

type CurrencyResponse struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func (m *CurrencyResponse) FromApplicationModel(currency *model.CurrencyApplicationModel) {
	m.Name = currency.Name
	m.Code = currency.Code
}

func NewCurrencyResponse(currency *model.CurrencyApplicationModel) *CurrencyResponse {
	result := &CurrencyResponse{}
	result.FromApplicationModel(currency)
	return result
}

type AddCurrencyRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func (m *AddCurrencyRequest) ToApplicationModel() *model.CurrencyApplicationModel {
	return &model.CurrencyApplicationModel{
		Name: m.Name,
		Code: m.Code,
	}
}
