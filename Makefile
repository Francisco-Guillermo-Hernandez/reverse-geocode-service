GO := go
GOPATH := $(shell $(GO) env GOPATH)
GOROOT := $(shell $(GO) env GOROOT)

SERVICE_NAME := reverse-geocode-service
FUNCTION_NAME := $(SERVICE_NAME)
BUILD_DIR := build
BIN_NAME := $(BUILD_DIR)/bootstrap
ZIP_FILE := bootstrap.zip
SERVER_DIR := ./cmd/server/main.go
LAMBDA_DIR := ./cmd/lambda/main.go
DATA_DIR := ./data

start-infra:
	@echo "Starting infrastructure..."
	docker compose up -d

stop-infra:
	@echo "Stopping infrastructure..."
	docker compose down -v --remove-orphans

build-function:
	@echo "Building Go Lambda function $(FUNCTION_NAME)..."
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 \
	$(GO) build $(GOFLAGS) -o $(BIN_NAME) $(LAMBDA_DIR)
	@mkdir -p $(BUILD_DIR)/data
	@cp $(DATA_DIR)/geocode.json $(BUILD_DIR)/data/geocode.json
	@cp $(DATA_DIR)/countries.csv $(BUILD_DIR)/data/countries.csv
	@zip $(BUILD_DIR)/$(ZIP_FILE) -j $(BIN_NAME)
	@echo "Build completed successfully!"

build-local:
	$(GO) build $(GOFLAGS) -o $(BIN_NAME) $(SERVER_DIR)
	@echo "Build completed successfully!"

quality:
	sonar-scanner

live-reload:
	air --build.cmd "go build -o bin/$(SERVICE_NAME) $(SERVER_DIR)" --build.entrypoint "./bin/$(SERVICE_NAME)"

serve:
	@echo "Starting local development..."
	@make live-reload
	
deploy: 
	@cd cloud
	@cdk bootstrap --profile Developer
	@cdk deploy --profile Developer