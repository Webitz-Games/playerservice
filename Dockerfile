FROM golang:1.21-alpine as builder
WORKDIR /build
COPY . .
RUN go get -d -v ./...
RUN go build -o player-service-api ./cmd/player
EXPOSE 8080

FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/player-service-api player-service-api
EXPOSE 8080
CMD ["/app/player-service-api"]