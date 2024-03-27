FROM golang:1.22-alpine3.18 AS builder

ENV CGO_ENABLED=1

RUN apk add --no-cache gcc musl-dev

COPY . /kinopoisk-telegram-bot
WORKDIR /kinopoisk-telegram-bot

RUN go mod download
RUN go build -o ./bot cmd/bot/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /kinopoisk-telegram-bot/bot .
COPY --from=0 /kinopoisk-telegram-bot/configs configs/

EXPOSE 8080

CMD ["./bot"]

