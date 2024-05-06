package model

import (
	"time"

	"github.com/google/uuid"
)

type CurrencyRepositoryModel struct {
	ID        uuid.UUID `gorm:"primaryKey,not null"`
	Name      string    `gorm:"not null"`
	Code      string    `gorm:"index:code_currencies_idx,unique,not null"`
	CreatedAt time.Time `gorm:"index:createdat_currencies_idx,not null"`
	CreatedBy uuid.UUID
}
