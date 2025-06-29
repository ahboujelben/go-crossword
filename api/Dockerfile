# Build stage
FROM golang:1.24-alpine AS build

WORKDIR /app

# Copy the entire project with workspace structure
COPY . .

# Download dependencies for the workspace
WORKDIR /app
RUN go work sync

# Build the API application
WORKDIR /app/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-crossword-api .

# Final stage
FROM alpine:3.19

WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /go-crossword-api .

# Expose the default API port
EXPOSE 8080

# Set environment variables with defaults
ENV PORT=8080
ENV OLLAMA_URL=http://ollama:11434
ENV OLLAMA_MODEL=llama3:8b

# Set the entry point
ENTRYPOINT ["/app/go-crossword-api"]
