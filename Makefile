# Variables
CLI_NAME := go-crossword-cli
CLI_MAIN := ./cli
CLI_BINARY := $(CLI_NAME)
MCP_NAME := go-crossword-mcp
MCP_MAIN := ./mcp
MCP_BINARY := $(MCP_NAME)

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

# Build the MCP server
.PHONY: build-mcp
build-mcp:
	@echo "Building MCP server..."
	@cd $(MCP_MAIN) && go build -o ../$(MCP_BINARY) .
	@echo "MCP server build complete: $(MCP_BINARY)"

# Run the MCP server
.PHONY: run-mcp
run-mcp:
	@echo "Running MCP server..."
	go run ${MCP_MAIN}

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@cd modules && go test ./...
	@cd cli && go test ./...
	@cd mcp && go test ./...

# Docker targets
.PHONY: docker-build-cli
docker-build-cli:
	@echo "Building CLI Docker image..."
	@docker build -t $(CLI_NAME) -f cli/Dockerfile .

.PHONY: docker-run-cli
docker-run-cli: docker-build-cli
	@echo "Running CLI Docker container..."
	@docker run --rm $(CLI_NAME) $(ARGS)

.PHONY: docker-build-mcp
docker-build-mcp:
	@echo "Building MCP Docker image..."
	@docker build -t $(MCP_NAME) -f mcp/Dockerfile .

.PHONY: docker-run-mcp
docker-run-mcp: docker-build-mcp
	@echo "Running MCP Docker container..."
	@docker run --rm -i $(MCP_NAME)
