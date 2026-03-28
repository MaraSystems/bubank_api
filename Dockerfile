FROM golang:1.26-alpine3.23 AS builder
WORKDIR /app
COPY . .

RUN go mod download
RUN go build -v -o main main.go

RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.19.1/migrate.linux-amd64.tar.gz \
 | tar xvz

FROM alpine:3.23
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY db/migrations ./migrations

EXPOSE 8080
CMD [ "/app/main" ]
