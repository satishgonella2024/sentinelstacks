FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o sentinel cmd/sentinel/main.go
RUN go build -o api-server cmd/api/main.go

FROM alpine:latest
COPY --from=builder /app/sentinel /usr/local/bin/
COPY --from=builder /app/api-server /usr/local/bin/
COPY README.md /usr/local/share/doc/
COPY LICENSE /usr/local/share/doc/
COPY examples /usr/local/share/sentinel/examples/

# Create necessary directories
RUN mkdir -p /root/.sentinel/registry

# Set environment variables
ENV SENTINEL_HOME=/root/.sentinel
ENV PATH="/usr/local/bin:${PATH}"

# Default command
CMD ["api-server", "--port", "8080"] 