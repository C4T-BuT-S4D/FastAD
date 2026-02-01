FROM golang:1.25-alpine AS builder

ARG SERVICE_NAME
ENV CGO_ENABLED=0

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
            -o "/service" \
            "./cmd/${SERVICE_NAME}/main.go"

ARG COMPRESS_BINARIES="false"
ENV COMPRESS_BINARIES=$COMPRESS_BINARIES
RUN if [ "${COMPRESS_BINARIES}" = "true" ]; then upx --lzma -9 "/service"; fi

FROM alpine:3.21

COPY --from=builder /service /service

CMD ["/service"]
