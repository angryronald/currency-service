package main

import (
	"context"
	"fmt"
	"log"
	netHttp "net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/angryronald/currency-service/cmd/currency-service/di"
	"github.com/angryronald/currency-service/cmd/currency-service/http"
	"github.com/angryronald/currency-service/cmd/currency-service/migration"
	"github.com/angryronald/currency-service/config"
	"github.com/angryronald/currency-service/internal/currency/infrastructure/repository"
	"github.com/angryronald/currency-service/internal/currency/infrastructure/repository/sync"
	"github.com/angryronald/currency-service/lib/file"
)

func init() {
	if config.GetValue(config.ENV) == "" {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if err = godotenv.Load(fmt.Sprintf("%s/default.env", dir)); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}
}

func main() {
	log := config.GetLogger()
	di.CollectDependencies(log)

	log.Println("starting ...")

	migration.RunMigration()
	runSeeder(log)
	runWorker(log)

	go runHTTP(log)
	go runHTTPProfiler(log)
	go runSubscribers(log)

	log.Println("Currency service is up")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	log.Println("shutting down currency service")
}

func runSeeder(log *logrus.Logger) {
	projectDir, err := file.GetProjectDir(config.GetValue(config.SERVICE_NAME))
	if err != nil {
		log.Fatalf("failed project directory is not found: %v", err)
	}

	seedFilePath := fmt.Sprintf("%s/seed/currency_code.json", projectDir)
	if err := repository.Seeding(
		di.AllDependencies.CurrencyMemcachedRepository,
		di.AllDependencies.CurrencySQLRepository,
		seedFilePath,
	); err != nil {
		log.Fatalf("failed seeding data: %v", err)
	}
}

func runWorker(log *logrus.Logger) {
	defaultPeriodInSec := 5
	periodInSecInString := config.GetValue(config.WORKER_PERIOD_IN_SEC)
	periodInSec, err := strconv.Atoi(periodInSecInString)
	if err != nil || periodInSec == 0 {
		periodInSec = defaultPeriodInSec
	}

	go sync.SynchronizeReadAndWriteData(
		di.AllDependencies.CurrencyMemcachedRepository,
		di.AllDependencies.CurrencySQLRepository,
		periodInSec,
		log,
	)

	log.Println("all workers are running")
}

func runSubscribers(log *logrus.Logger) {
	go di.AllDependencies.CurrencySubscriber.Run(context.Background())

	log.Println("all subscribers is listening")
}

func runHTTP(log *logrus.Logger) {
	port := config.GetValue(config.HTTP_PORT)

	if len(port) < 1 {
		panic(fmt.Sprintf("Environment Missing!\n*%s* is required", port))
	}

	var router *chi.Mux
	router = chi.NewRouter()

	router.Mount("/api", http.CompileRoute(router))

	server := &netHttp.Server{
		Addr:    port,
		Handler: router,
	}

	log.Println("HTTP transport run at ", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalln("cannot serve HTTP transport at port ", port, ": ", err)
	}
}

func runHTTPProfiler(log *logrus.Logger) {
	profilerPort := config.GetValue(config.HTTP_PROFILER_PORT)

	if len(profilerPort) < 1 {
		panic(fmt.Sprintf("Environment Missing!\n*%s* is required", profilerPort))
	}

	var router *chi.Mux
	router = chi.NewRouter()

	router.Mount("/profiler", http.CompileProfilingRoute(router))

	server := &netHttp.Server{
		Addr:    profilerPort,
		Handler: router,
	}

	log.Println("HTTP Profiler transport run at ", profilerPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalln("cannot serve HTTP profiler at port ", profilerPort, ": ", err)
	}
}
