FROM golang:1.21-alpine

ENV CGO_ENABLED=0

WORKDIR /app
COPY cmd ./cmd
COPY pkg ./pkg
COPY internal ./internal
COPY go.* ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o /data_service ./cmd/data_service/main.go

CMD "/data_service"
