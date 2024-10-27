SHELL := /bin/bash # Use bash syntax

# Optional colors to beautify output
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

## Swagger UI
run-openapi-ui: ## runs swagger ui: http://localhost:4000/
	docker rm -f quizapp-swagger 2>/dev/null || true && \
    docker run --name quizapp-swagger -p 4000:8080 \
    --restart unless-stopped \
    -e SWAGGER_JSON=/docs/openapi/quiz.yaml \
    -v $(PWD)/docs/openapi:/docs/openapi \
    swaggerapi/swagger-ui:v5.17.14

## Quality
check-quality: ## runs code quality checks
	make lint
	make fmt
	make vet

# Append || true below if blocking local developement
lint: ## go linting. Update and use specific lint tool and options
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
	golangci-lint run -c ./.golangci.yml

lint-fix: ## go linting. Update and use specific lint tool and options
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
	golangci-lint run -c ./.golangci.yml --fix

vet: ## go vet
	go vet ./...

fmt: ## runs go formatter
	go fmt ./...

update-mocks:
	go install github.com/golang/mock/mockgen@v1.6.0
	go generate ./...

tidy: ## run go mod tidy
	go mod tidy
## Test
test-unit:
	make tidy
	go test ./... --tags=unit -coverprofile=coverage.out

coverage: ## displays test coverage report in html mode
	make test-unit
	go tool cover -html=coverage.out

.PHONY: all
## All
all: ## quality checks and tests
	make check-quality
	make test-unit

.PHONY: help
## Help
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)