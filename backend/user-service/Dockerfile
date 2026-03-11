FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -installsuffix cgo -ldflags="-w -s" \
    -o /app/bin/user-service ./cmd/server

FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

RUN adduser -D -g '' appuser

COPY --from=builder /app/bin/user-service /usr/local/bin/user-service

RUN mkdir -p /var/log/user-service && \
    chown -R appuser:appuser /var/log/user-service

USER appuser

EXPOSE 50051 8080

ENTRYPOINT ["/usr/local/bin/user-service"]
