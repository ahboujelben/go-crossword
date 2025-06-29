# Variables
CLI_NAME := go-crossword-cli
CLI_MAIN := ./cli
CLI_BINARY := $(CLI_NAME)

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

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@cd modules && go test ./...
	@cd cli && go test ./...

# Docker targets
.PHONY: docker-build-cli
docker-build-cli:
	@echo "Building CLI Docker image..."
	@docker build -t $(CLI_NAME) -f cli/Dockerfile .

.PHONY: docker-run-cli
docker-run-cli: docker-build-cli
	@echo "Running CLI Docker container..."
	@docker run --rm $(CLI_NAME) $(ARGS)

# Docker compose targets
.PHONY: docker-compose-up
docker-compose-up:
	@echo "Starting services with Docker Compose..."
	@docker-compose up -d

.PHONY: docker-compose-down
docker-compose-down:
	@echo "Stopping services with Docker Compose..."
	@docker-compose down

.PHONY: docker-compose-cli
docker-compose-cli:
	@echo "Running CLI with Docker Compose..."
	@docker-compose --profile cli run --rm -T cli $(ARGS)
