# Build stage
FROM golang:1.20-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /go-service .

# Final stage
FROM alpine:3.18
COPY --from=build /go-service /usr/local/bin/go-service
ENTRYPOINT ["/usr/local/bin/go-service"]
