# Build Stage
FROM golang:1.25-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/gin-crud-demo cmd/gin-crud-demo/main.go

# Final Stage
FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /app/bin/gin-crud-demo /app/gin-crud-demo
COPY --from=builder /app/config /app/config
COPY --from=builder /app/internal/database/migrations /app/internal/database/migrations

EXPOSE 8080

ENTRYPOINT ["/app/gin-crud-demo"]
