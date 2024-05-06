package postgres

import (
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"

	internalDocker "github.com/angryronald/currency-service/lib/test/docker"
)

// ref: https://jonnylangefeld.com/blog/how-to-write-a-go-api-part-3-testing-with-dockertest

var dbInstanceLock sync.Mutex

func GenerateInstance(pool *dockertest.Pool) (*gorm.DB, *dockertest.Resource) {
	dbInstanceLock.Lock()
	defer dbInstanceLock.Unlock()

	var db *gorm.DB
	port := internalDocker.GetAvailablePort(5432)

	// Pull an image, create a container based on it and set all necessary parameters
	opts := dockertest.RunOptions{
		Repository:   "mdillon/postgis",
		Tag:          "latest",
		Env:          []string{"POSTGRES_PASSWORD=password"},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	// Run the Docker container
	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// Exponential retry to connect to database while it is booting
	if err := pool.Retry(func() error {
		databaseConnStr := fmt.Sprintf("host=localhost port=%s user=postgres dbname=postgres password=password sslmode=disable", port)
		db, err = gorm.Open(postgres.Open(databaseConnStr), &gorm.Config{})
		if err != nil {
			log.Println("Database not ready yet (it is booting up, wait for a few tries)...")
			return err
		}

		// Tests if database is reachable
		dbinstance, err := db.DB()
		if err != nil {
			log.Println("Database instance cannot created")
			return err
		}
		return dbinstance.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	return db, resource
}
