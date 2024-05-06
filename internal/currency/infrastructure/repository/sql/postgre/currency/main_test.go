package currency

import (
	"os"
	"testing"

	"gorm.io/gorm"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/sirupsen/logrus"

	"github.com/angryronald/currency-service/internal/currency/infrastructure/repository/model"
	"github.com/angryronald/currency-service/lib/test/docker/postgres"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		logrus.Fatalf("Could not connect to Docker: %s\n", err)
	}

	var postgresResource *dockertest.Resource
	db, postgresResource = postgres.GenerateInstance(pool)
	if db != nil {
		logrus.Debugf("Success generate postgres instance\n")
	}

	db.AutoMigrate(&model.CurrencyRepositoryModel{})

	code := m.Run() // execute all the tests

	// Delete the Docker container
	if err := pool.Purge(postgresResource); err != nil {
		logrus.Fatalf("Could not purge postgres resource: %s\n", err)
	}

	defer os.Exit(code)
}
