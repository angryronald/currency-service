# currency-service
Handle all responsibility related to currency. This is a sample of implementation DDD with Layered Architecture, CQRS, SOLID Principle, and DI. 

## Dependencies
1. RabbitMQ
2. Redis
3. Postgres

## How to run
1. Spinning up all the dependencies
```bash
    make spin-up-dependent-containers-locally
```
2. Start up the currency service
```bash
    go run cmd/currency-service/main.go
```

## How to test
1. Without Coverage
```bash
    make test
```
2. With coverage
```bash
    make test-coverage
```

### Notes on testing
Due to issue with dockertest and rabbitmq, the testing on publisher and subscriber which using rabbitmq need to be done manually. Since the time to creating and destroy the container is overlapping with each test.
The testing provided only on unit test and integration test for now.

## Description
1. **User Interface** is reflected by endpoints which responsible for handling the requests and generate responses
2. **Application Layer** is reflected by application package which contains use cases
3. **Domain Layer** is reflected by domain package which contains the business logic for specific domain
4. **Infrastructure Layer** is reflected by infrastructure package which contains **event** for domain events, **repository** for repository layer, and **external** for external layer

![Layered Architecture - Eric Evans, 2003](https://github.com/angryronald/currency-service/blob/main/docs/DDD-Layered-Architecture.jpg)

## Documentation
API Documentation can be found inside /docs
