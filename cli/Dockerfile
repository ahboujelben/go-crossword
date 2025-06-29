# Build stage
FROM golang:1.24-alpine AS build

WORKDIR /app

# Copy the entire project with workspace structure
COPY . .

# Download dependencies for the workspace
WORKDIR /app
RUN go work sync

# Build the CLI application
WORKDIR /app/cli
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-crossword-cli .

# Final stage
FROM alpine:3.19

WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /go-crossword-cli .

# Set the entry point
ENTRYPOINT ["/app/go-crossword-cli"]
