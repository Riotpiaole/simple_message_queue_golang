FROM golang:1.20.4-alpine3.17 AS builder

# Install the C lib for kafka
RUN apk add --no-progress --no-cache gcc musl-dev

# Set the working directory
WORKDIR /app

# Copy the source code
COPY main.go ./
COPY * ./

# Download Go module dependencies
COPY ../go.mod ./
COPY ../go.sum ./
RUN go mod download 

# Build the Go application
RUN go build -tags musl -ldflags '-extldflags "-static"' -o producer .

CMD ["./producer"]
EXPOSE 3000