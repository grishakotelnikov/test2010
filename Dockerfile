FROM golang:1.23-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o /app/main ./cmd

FROM alpine:latest


COPY --from=builder /app/main /main
COPY .env /app/.env
COPY ./migrations /app/migrations


CMD ["/main"]
