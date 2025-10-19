# Stage 1: Build the application
FROM golang:1.24-alpine AS build

WORKDIR /app

# Copy the entire project with workspace structure
COPY . .

# Download dependencies for the workspace
WORKDIR /app
RUN go work sync

# Build the MCP server
WORKDIR /app/mcp
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-crossword-mcp .

# Stage 2: Create the final, minimal image
FROM scratch

WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /go-crossword-mcp .

# Set the entrypoint for the container to run the MCP server.
ENTRYPOINT ["/app/go-crossword-mcp"]
