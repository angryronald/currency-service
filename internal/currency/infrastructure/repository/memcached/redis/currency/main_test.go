package currency

import (
	"fmt"
	"log"
	"os"
	"sync"
	"testing"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/redis/go-redis/v9"

	redisInDocker "github.com/angryronald/currency-service/lib/test/docker/redis"
)

var redisClient *redis.Client
var poolResourceMap map[*dockertest.Pool]*dockertest.Resource
var lock sync.Mutex

func init() {
	poolResourceMap = map[*dockertest.Pool]*dockertest.Resource{}
}

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s\n", err)
	}

	redisClientTemp, redisResource := redisInDocker.GenerateInstance(pool)
	if redisClientTemp != nil {
		redisClient = redisClientTemp
		fmt.Printf("Success generate redis instance\n")
	}

	code := m.Run() // execute all the tests

	// Delete the Docker container
	for pool, instance := range poolResourceMap {
		if err := pool.Purge(instance); err != nil {
			log.Fatalf("Could not purge postgres resource: %s\n", err)
		}
	}

	// Delete the Docker container
	if err := pool.Purge(redisResource); err != nil {
		log.Fatalf("Could not purge redis resource: %s\n", err)
	}

	defer os.Exit(code)
}
