package config

import "github.com/angryronald/go-kit/constant"

// noinspection ALL
const (
	ENV             = "ENV"
	ENV_DEVELOPMENT = constant.DEVELOPMENT
	ENV_LOCAL       = constant.LOCAL
	ENV_PRODUCTION  = constant.PRODUCTION

	HTTP_PORT          = "HTTP_PORT"
	HTTP_PROFILER_PORT = "HTTP_PROFILER_PORT"

	DATABASE_CONNECTION_STRING = "DATABASE_CONNECTION_STRING"

	DATABASE_DRIVER = "DATABASE_DRIVER"
	DATABASE_HOST   = "DATABASE_HOST"
	DATABASE_PORT   = "DATABASE_PORT"
	DATABASE_USER   = "DATABASE_USER"
	DATABASE_PASS   = "DATABASE_PASS"
	DATABASE_NAME   = "DATABASE_NAME"
	DATABASE_SSL    = "DATABASE_SSL"

	REDIS_ADDRESS            = "REDIS_ADDRESS"
	REDIS_PASSWORD           = "REDIS_PASSWORD"
	REDIS_DB                 = "REDIS_DB"
	REDIS_DEFAULT_EXPIRATION = "REDIS_DEFAULT_EXPIRATION"

	SERVICE_NAME = "SERVICE_NAME"
	SERVICE_ID   = "SERVICE_ID"

	RABBITMQ_HOST     = "RABBITMQ_HOST"
	RABBITMQ_PORT     = "RABBITMQ_PORT"
	RABBITMQ_USERNAME = "RABBITMQ_USERNAME"
	RABBITMQ_PASSWORD = "RABBITMQ_PASSWORD"

	ALLOWED_CLIENTS = "ALLOWED_CLIENTS"

	WORKER_PERIOD_IN_SEC = "WORKER_PERIOD_IN_SEC"
)
