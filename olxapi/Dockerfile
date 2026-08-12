FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd

FROM alpine:3.24.1

RUN adduser -D -g '' gouser

WORKDIR /app

COPY --from=builder /app/main .

USER gouser

EXPOSE 8080

CMD ["./main"]