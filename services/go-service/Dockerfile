# Build stage
FROM golang:1.20-alpine AS build

WORKDIR /src

# Copy entire project
COPY . .

# Download all dependencies automatically
RUN go mod tidy

# Build the Go service
RUN go build -o /go-service .

# Final stage
FROM alpine:3.18
COPY --from=build /go-service /usr/local/bin/go-service
ENTRYPOINT ["/usr/local/bin/go-service"]
