FROM golang:1.24

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
RUN go build -o smaugServer ./cmd/server

# Default command can be overridden from docker-compose
CMD ["./smaugServer", "-k", "fa4f22ddvr", "-d", "postgresql://postgres:5432/citizix_db?user=test1_user&password=qwerwrbyrbgwb"]
