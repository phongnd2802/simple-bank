FROM golang:1.24.1-alpine3.20 AS builder

WORKDIR /app

COPY . .

RUN go build -o main main.go

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/main .
COPY config/production.yaml config/production.yaml

EXPOSE 8002

ENTRYPOINT [ "/app/main" ]

