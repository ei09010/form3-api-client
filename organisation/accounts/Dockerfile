FROM golang:1.17-alpine

# Set working directory
WORKDIR /test

ENV CGO_ENABLED=0

COPY go.* .
RUN go mod download

# Copy local code to the container image.
COPY . .

# Run tests
CMD go test -v ./...