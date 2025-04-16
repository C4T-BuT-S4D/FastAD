FROM golang:1.24.2-alpine

WORKDIR /app
COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o /app/service /app/main.go

CMD ["/app/service"]
