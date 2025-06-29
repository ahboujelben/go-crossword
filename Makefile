# Variables
CLI_NAME := go-crossword-cli
API_NAME := go-crossword-api
CLI_MAIN := ./cli
API_MAIN := ./api
CLI_BINARY := $(CLI_NAME)
API_BINARY := $(API_NAME)

# Build the CLI application
.PHONY: build-cli
build-cli:
	@echo "Building CLI application..."
	@cd $(CLI_MAIN) && go build -o ../$(CLI_BINARY) .
	@echo "CLI build complete: $(CLI_BINARY)"

# Run the CLI application
.PHONY: run-cli
run-cli: 
	@echo "Running CLI application..."
	go run ${CLI_MAIN}

# Build the API server
.PHONY: build-api
build-api:
	@echo "Building API server..."
	@cd $(API_MAIN) && go build -o ../$(API_BINARY) .
	@echo "API build complete: $(API_BINARY)"

# Run the API server
.PHONY: run-api
run-api: build-api
	@echo "Running API server..."
	go run ${API_MAIN}

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@cd modules && go test ./...
	@cd cli && go test ./...
	@cd api && go test ./...

# Docker targets
.PHONY: docker-build-cli
docker-build-cli:
	@echo "Building CLI Docker image..."
	@docker build -t $(CLI_NAME) -f cli/Dockerfile .

.PHONY: docker-run-cli
docker-run-cli: docker-build-cli
	@echo "Running CLI Docker container..."
	@docker run --rm $(CLI_NAME) $(ARGS)

.PHONY: docker-build-api
docker-build-api:
	@echo "Building API Docker image..."
	@docker build -t $(API_NAME) -f api/Dockerfile .

.PHONY: docker-run-api
docker-run-api: docker-build-api
	@echo "Running API Docker container..."
	@docker run --rm -p 8080:8080 $(API_NAME)

# Docker compose targets
.PHONY: docker-compose-up
docker-compose-up:
	@echo "Starting services with Docker Compose..."
	@docker-compose up -d

.PHONY: docker-compose-down
docker-compose-down:
	@echo "Stopping services with Docker Compose..."
	@docker-compose down
