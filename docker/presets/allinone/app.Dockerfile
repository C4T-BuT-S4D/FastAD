FROM golang:1.23-alpine as builder

RUN apk add --no-cache upx

WORKDIR /app

COPY cmd cmd
COPY pkg pkg
COPY internal internal
COPY go.* ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
        go build \
            -trimpath \
            -ldflags="-w -s" \
            -o "/allinone" \
            "./cmd/allinone/main.go" \
    && upx --lzma -9 "/allinone"

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
        go build \
            -trimpath \
            -ldflags="-w -s" \
            -o "/migrator" \
            "./cmd/migrator/main.go" \
    && upx --lzma -9 "/migrator"


FROM alpine:3.20

COPY --from=builder /allinone /allinone
COPY --from=builder /migrator /migrator

CMD ["/bin/sh", "-c", "/migrator init && /migrator migrate && /allinone"]
