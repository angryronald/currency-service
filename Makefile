APP_DIRECTORY := ./cmd/currency-service
MAIN_FILE := main.go

.PHONY: run
run:
	@go run $(APP_DIRECTORY)/$(MAIN_FILE)

.PHONY: docker
docker: Dockerfile
	echo "building the currency-service container..."
	docker build --label "version=1.0" -t currency-service:1.0 .

.PHONY: test
test: 
	go test -v -p 1 -timeout 99999s $(shell go list ./... | grep -v '/event/')

.PHONY: test-coverage
test-coverage:
	go test -v -p 1 -timeout 99999s -count=1 -race -coverpkg=./... -coverprofile=coverage.out -covermode=atomic $(shell go list ./... | grep -v '/event/')

.PHONY: spin-up-dependent-containers-locally
spin-up-dependent-containers-locally:
	docker-compose -f ./docker/docker-compose.local.yaml up -d 

.PHONY: spin-down-dependent-containers-locally
spin-down-dependent-containers-locally:
	docker-compose -f ./docker/docker-compose.local.yaml down 