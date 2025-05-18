# Stage 1: Build the Go application
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download dependencies. If you have private modules, you might need to configure git for HTTPS or SSH.
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
# Using CGO_ENABLED=0 to ensure static linking for alpine base image if no CGO is used.
# The output binary is specified as gok-proxy-proxy in your README build instructions.
RUN CGO_ENABLED=0 go build -v -o /gok-proxy-proxy ./cmd/proxy/

# Stage 2: Create the final lightweight image
FROM alpine:latest

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /gok-proxy-proxy /app/gok-proxy-proxy

# Copy the configuration file
# Ensure config.yaml is in the root of your project when building the image
COPY config.yaml /app/config.yaml

# Expose the port the proxy listens on (default is 8080 from your config)
EXPOSE 8080

# Set the entrypoint for the container
ENTRYPOINT ["/app/gok-proxy-proxy"]