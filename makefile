MODULE = $(shell go list -m)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || echo "1.0.0")
PACKAGES := $(shell go list ./... | grep -v /vendor/)
LDFLAGS := -ldflags "-X main.Version=${VERSION}"

CONFIG_FILE ?= ./config/local.yml
DB_SOURCE ?= "$(shell sed -n 's/^dbsource:[[:space:]]*"\(.*\)"/\1/p' $(CONFIG_FILE))"
MIGRATE := docker run -v $(shell pwd)/migrations:/migrations --network host migrate/migrate:v4.10.0 -path=/migrations/ -database $(DB_SOURCE)

PID_FILE := './.pid'
FSWATCH_FILE := './fswatch.cfg'

.PHONY: default
default: help

# Read More https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## run unit tests
	@echo "mode: count" > coverage-all.out
	@$(foreach pkg,$(PACKAGES), \
		go test -p=1 -cover -covermode=count -coverprofile=coverage.out ${pkg}; \
		tail -n +2 coverage.out >> coverage-all.out;)

.PHONY: test-cover
test-cover: test ## run unit tests and show test coverage information
	go tool cover -html=coverage-all.out

.PHONY: run
run: ## run the API server
	go run ${LDFLAGS} cmd/server/main.go

.PHONY: start
start: ## start docker compose
	docker-compose up -d --remove-orphans --build 

.PHONY: stop
stop: ## stop docker compose
	docker-compose down --remove-orphans --volumes 

.PHONY: run-restart
run-restart: ## restart the API server
	@pkill -P `cat $(PID_FILE)` || true
	@printf '%*s\n' "80" '' | tr ' ' -
	@echo "Source file changed. Restarting server..."
	@go run ${LDFLAGS} cmd/server/main.go & echo $$! > $(PID_FILE)
	@printf '%*s\n' "80" '' | tr ' ' -

run-live: ## run the API server with live reload support (requires fswatch)
	@go run ${LDFLAGS} cmd/server/main.go & echo $$! > $(PID_FILE)
	@fswatch -x -o --event Created --event Updated --event Renamed -r internal pkg cmd config | xargs -n1 -I {} make run-restart

.PHONY: build
build:  ## build the API server binary
	CGO_ENABLED=0 go build ${LDFLAGS} -a -o server $(MODULE)/cmd/server

.PHONY: build-docker
build-docker: ## build the API server as a docker image
	docker build -f cmd/server/Dockerfile -t server .

.PHONY: clean
clean: ## remove temporary files
	rm -rf server coverage.out coverage-all.out

.PHONY: version
version: ## display the version of the API server
	@echo $(VERSION)

.PHONY: db-start
db-start: ## start the database server
	docker run --name ev_webservice_mongodb  -d mongo:latest 

.PHONY: db-stop
db-stop: ## stop the database server
	docker stop ev_webservice_mongodb

.PHONY: lint
lint: ## run golint on all Go package
	@golint $(PACKAGES)


.PHONY: fmt
fmt: ## run "go fmt" on all Go packages
	@go fmt $(PACKAGES)
