FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN go build -o /app/bin/app ./cmd/main.go

FROM alpine:3.18 AS runner

WORKDIR /app

COPY --from=builder /app/bin/app /app/app

COPY config/local.yaml /app/config/local.yaml

ENV CONFIG_PATH /app/config/local.yaml

COPY migrations /app/migrations

CMD ["sh", "-c", "goose -dir /app/migrations postgres 'user=postgres password=postgres dbname=postgres sslmode=disable' up && ./app"]
