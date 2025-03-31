FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates build-base

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/sentinel cmd/sentinel/main.go

# Create a minimal runtime image
FROM alpine:3.18

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy the binary from the builder stage
COPY --from=builder /app/bin/sentinel /app/bin/sentinel

# Set the PATH to include our binary
ENV PATH="/app/bin:${PATH}"

# Create a volume for agent data
VOLUME ["/app/data"]

# Set the entrypoint
ENTRYPOINT ["sentinel"]
CMD ["--help"]
