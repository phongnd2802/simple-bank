FROM golang:1.24.1-alpine3.20 AS builder

WORKDIR /app

COPY . .

RUN go build -o main main.go
RUN apk add curl
RUN curl -fsSL \
https://raw.githubusercontent.com/pressly/goose/master/install.sh |\
GOOSE_INSTALL=/app sh -s v3.5.0


FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/bin/goose ./goose
COPY config/production.yaml config/production.yaml

COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration

EXPOSE 8002

CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]

