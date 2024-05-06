package migration

import (
	"github.com/angryronald/currency-service/cmd/currency-service/di"
	"github.com/angryronald/currency-service/internal/currency/infrastructure/repository/model"
)

func RunMigration() {
	// register all repository models
	di.AllDependencies.SQLDB.AutoMigrate(&model.CurrencyRepositoryModel{})
}
