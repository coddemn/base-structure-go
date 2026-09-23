FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /bin/app ./cmd/app
RUN go build -o /bin/migrate ./cmd/migrate

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /bin/app /bin/app
COPY --from=builder /bin/migrate /bin/migrate
COPY --from=builder /app/configs /app/configs
